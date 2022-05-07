// a Telegram Bot Service framework
package tgolf

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Project-Nessie/nessielight"
	"github.com/yanzay/tbot/v2"
)

var logger *log.Logger

// 描述了一个命令
type Command struct {
	// 继承 BotCommand
	tbot.BotCommand
	// Init 函数用于触发 starter 后的初始化和鉴权
	Init func(from *tbot.User, chat tbot.Chat) bool
	// 命令参数
	Param    []Parameter
	Callback func(argv []Argument, from *tbot.User, chatid string)
	Handler  func(*tbot.Message)
}

// tg bot server
type Server struct {
	Bot       *tbot.Server
	Client    *tbot.Client
	db        nessielight.KVDatabase
	commands  map[string]*Command
	callbacks map[string]func(*tbot.CallbackQuery)
}

// Send formatted message to a chat with html parsing
func (r *Server) Sendf(chatid string, format string, v ...interface{}) (*tbot.Message, error) {
	return r.Client.SendMessage(chatid, fmt.Sprintf(format, v...), tbot.OptParseModeHTML)
}

// Send message with inline button
func (r *Server) SendfWithBtn(chatid string, btnMatrix [][]tbot.InlineKeyboardButton, format string, v ...interface{}) (*tbot.Message, error) {
	btns := &tbot.InlineKeyboardMarkup{
		InlineKeyboard: btnMatrix,
	}
	return r.Client.SendMessage(chatid, fmt.Sprintf(format, v...),
		tbot.OptInlineKeyboardMarkup(btns), tbot.OptParseModeHTML)
}

// starter 为命令的触发字符串。若开头为 / 则会作为显示命令，否则为隐式命令。
// description 可选，用于描述命令。开头为 / 的命令会以 start - description 的形式打印到日志中，方便在
// Bot Father 那 setcommand。init 指在触发 start 后，获取参数前的检查（例如权限），返回值为
// true 则继续获取参数，否则终止。params 描述了参数列表，包含每个参数的描述，校验器，f 即回调函数
func (r *Server) Register(starter string, description string, init func(from *tbot.User, chat tbot.Chat) bool,
	params []Parameter, f func(argv []Argument, from *tbot.User, chatid string)) {

	logger.Printf("register command: start=%s, params: %v", starter, params)

	// 挂个名，命令本身没有内容
	if f == nil {
		return
	}

	handler := func(m *tbot.Message) {
		logger.Printf("handle command: %s", starter)
		from := m.From
		chat := m.Chat
		text := m.Text

		if init != nil && !init(from, chat) {
			return
		}

		argv := make([]Argument, len(params))
		for i, v := range params {
			argv[i] = Argument{Field: v.Field}
		}
		current := 0
		if from != nil && r.db.Get(fmt.Sprintf("user/%d", from.ID)) != nil && text != "/cancel" {
			r.Sendf(chat.ID, "You're currently doing another job, send /cancel to cancel it")
			return
		}

		logger.Printf("user %d invoke %s", from.ID, starter)

		if current == len(argv) {
			f(argv, from, chat.ID)
			return
		}
		r.Sendf(chat.ID, "Enter "+params[current].Description)

		// 每一次被调用，填充一个参数
		// return true if finish
		mhandler := func(m *tbot.Message) bool {
			logger.Printf("user %d invoke \"%s\" for argument \"%s\"", from.ID, starter, argv[current].Key)
			if m.Text == "/cancel" {
				r.Sendf(m.Chat.ID, "Operation canceled")
				return true
			}
			if argv[current].Validator == nil || argv[current].Validator(m.Text) {
				argv[current].Value = m.Text
				current = current + 1
				if current == len(argv) {
					f(argv, m.From, m.Chat.ID)
					return true
				}
				r.Sendf(m.Chat.ID, "Enter %s\nSend /cancel to stop current operation", params[current].Description)
				return false
			} else {
				r.Sendf(m.Chat.ID, "invalid argument %s. try again", argv[current].Key)
				return false
			}
		}

		r.db.Set(fmt.Sprintf("user/%d", from.ID), mhandler)
	}
	r.Bot.HandleMessage(starter, handler)

	r.commands[starter] = &Command{
		BotCommand: tbot.BotCommand{
			Command:     starter,
			Description: description,
		},
		Init:     init,
		Param:    params,
		Callback: f,
		Handler:  handler,
	}
}

// filter：决定是否使用本 callback
func (r *Server) RegisterInlineButton(data string, handler func(cq *tbot.CallbackQuery)) {
	r.callbacks[data] = handler
}

func (r *Server) HandleCallback(cq *tbot.CallbackQuery) {
	logger.Printf("HandleCallback: %s, message: %s", cq.Data, cq.Message.Text)
	data := cq.Data
	if r.callbacks[data] != nil {
		r.callbacks[data](cq)
	}
}

func (r *Server) HandleMessage(m *tbot.Message) {
	logger.Printf("receive message: %s \"%s\"", m.Chat.Title, m.Text)
	if m.From != nil {
		if handler := r.db.Get(fmt.Sprintf("user/%d", m.From.ID)); handler != nil {
			typedhandler := handler.(func(*tbot.Message) bool)
			if typedhandler(m) {
				r.db.Set(fmt.Sprintf("user/%d", m.From.ID), nil)
			}
		} else {
			r.Sendf(m.Chat.ID, "I can't understand >_<")
		}
	}
}

func (r *Server) Start() error {
	r.Bot.HandleMessage(".*", r.HandleMessage)
	r.Bot.HandleCallback(r.HandleCallback)
	commands := make([]tbot.BotCommand, 0, len(r.commands))
	for _, v := range r.commands {
		if v.Command[0] == '/' {
			commands = append(commands, tbot.BotCommand{
				Command:     v.Command[1:],
				Description: v.Description,
			})
		}
	}

	commandlist := "commands:\n"
	for _, v := range commands {
		commandlist += v.Command + " - " + v.Description + "\n"
	}
	logger.Print(commandlist)

	if err := r.Bot.Start(); err != nil {
		return err
	}
	if err := r.Client.SetMyCommands(commands); err != nil {
		return err
	}
	return nil
}

func (r *Server) StartCommand(starter string, from *tbot.User, chat tbot.Chat) error {
	if r.commands[starter] == nil {
		return fmt.Errorf("command not found: starter=%s", starter)
	}
	// NOTE: 要注意初始化 handleMessage 里要用到的部分
	trigger := tbot.Message{
		From: from,
		Chat: chat,
		Text: starter,
	}
	r.commands[starter].Handler(&trigger)
	return nil
}

func (r *Server) EditCallbackBtn(cq *tbot.CallbackQuery, btnMatrix [][]tbot.InlineKeyboardButton) (*tbot.Message, error) {
	chatid := cq.Message.Chat.ID
	msgid := cq.Message.MessageID
	return r.Client.EditMessageReplyMarkup(chatid, msgid,
		tbot.OptInlineKeyboardMarkup(&tbot.InlineKeyboardMarkup{InlineKeyboard: btnMatrix}))
}

func (r *Server) EditCallbackMsg(cq *tbot.CallbackQuery, format string, v ...interface{}) (*tbot.Message, error) {
	chatid := cq.Message.Chat.ID
	msgid := cq.Message.MessageID
	return r.Client.EditMessageText(chatid, msgid, fmt.Sprintf(format, v...),
		tbot.OptParseModeHTML, tbot.OptInlineKeyboardMarkup(cq.Message.ReplyMarkup))
}
func (r *Server) EditCallbackMsgWithBtn(cq *tbot.CallbackQuery, btnMatrix [][]tbot.InlineKeyboardButton,
	format string, v ...interface{}) (*tbot.Message, error) {
	chatid := cq.Message.Chat.ID
	msgid := cq.Message.MessageID
	return r.Client.EditMessageText(chatid, msgid, fmt.Sprintf(format, v...),
		tbot.OptParseModeHTML, tbot.OptInlineKeyboardMarkup(&tbot.InlineKeyboardMarkup{InlineKeyboard: btnMatrix}))
}

func NewServerFromTbot(bot *tbot.Server) Server {
	db := nessielight.NewMemoryDB()
	server := Server{
		Bot:       bot,
		db:        &db,
		Client:    bot.Client(),
		commands:  make(map[string]*Command),
		callbacks: make(map[string]func(*tbot.CallbackQuery)),
	}
	return server
}

func NewServer(botToken string, webhookUrl string, listen string) Server {
	bot := tbot.New(botToken, tbot.WithWebhook(webhookUrl, listen))
	bot.Use(func(h tbot.UpdateHandler) tbot.UpdateHandler {
		return func(u *tbot.Update) {
			start := time.Now()
			h(u)
			logger.Printf("Handle time: %v", time.Since(start))
		}
	})
	server := NewServerFromTbot(bot)
	return server
}

type Field struct {
	Key       string
	Validator func(value string) bool
}
type Parameter struct {
	Field
	Description string
}

type Argument struct {
	Field
	Value string
}

func NewParam(key, desc string, validator func(value string) bool) Parameter {
	return Parameter{
		Field:       Field{Key: key, Validator: validator},
		Description: desc,
	}
}

func init() {
	logger = log.New(os.Stderr, "[tgolf] ", log.LstdFlags|log.Lmsgprefix)
}
