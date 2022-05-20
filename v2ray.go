package nessielight

import (
	"bytes"
	"context"
	b64 "encoding/base64"
	"html/template"
	"log"
	"os"

	"github.com/v2fly/v2ray-core/v4/app/proxyman/command"
	statsService "github.com/v2fly/v2ray-core/v4/app/stats/command"
	"github.com/v2fly/v2ray-core/v4/common/protocol"
	"github.com/v2fly/v2ray-core/v4/common/serial"
	"github.com/v2fly/v2ray-core/v4/common/uuid"
	"github.com/v2fly/v2ray-core/v4/proxy/vmess"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var logger *log.Logger

// 调用 V2ray API 的客户端
// v2rayClient implements V2rayService
type v2rayClient struct {
	statClient statsService.StatsServiceClient
	handClient command.HandlerServiceClient
	// vmess settings
	inboundTag string
	port       int
	domain     string
	path       string
}

func (r *v2rayClient) SetUser(email, id string) error {
	_, err := r.handClient.AlterInbound(context.Background(), &command.AlterInboundRequest{
		Tag: r.inboundTag,
		Operation: serial.ToTypedMessage(&command.AddUserOperation{
			User: &protocol.User{
				Level: 0,
				Email: email,
				Account: serial.ToTypedMessage(&vmess.Account{
					Id:               id,
					AlterId:          0,
					SecuritySettings: &protocol.SecurityConfig{Type: protocol.SecurityType_AUTO},
				}),
			},
		}),
	})
	if err != nil {
		return err
	}
	log.Printf("SetUser: email=%s id=%s", email, id)
	return nil
}

func (r *v2rayClient) AddUser(email string) (string, error) {
	userID := protocol.NewID(uuid.New()).String()
	err := r.SetUser(email, userID)
	if err != nil {
		return "", err
	}
	return userID, nil
}

func (r *v2rayClient) RemoveUser(email string) error {
	_, err := r.handClient.AlterInbound(context.Background(), &command.AlterInboundRequest{
		Tag: r.inboundTag,
		Operation: serial.ToTypedMessage(&command.RemoveUserOperation{
			Email: email,
		}),
	})
	if err != nil {
		return err
	}
	logger.Printf("RemoveUser: email=%s", email)
	return nil
}

type UserTrafficStat struct {
	Name  string
	Value int64
}

func (r *v2rayClient) QueryUserTraffic(pattern string, reset bool) ([]UserTrafficStat, error) {
	resp, err := r.statClient.QueryStats(context.Background(), &statsService.QueryStatsRequest{
		Pattern: pattern,
		Reset_:  reset, // 查询完成后是否重置流量
	})
	if err != nil {
		return nil, err
	}
	// 获取返回值中的流量信息
	stat := resp.GetStat()
	trafficStat := make([]UserTrafficStat, 0, len(stat))

	for _, v := range stat {
		if v != nil {
			trafficStat = append(trafficStat, UserTrafficStat{
				Name:  v.GetName(),
				Value: v.GetValue(),
			})
		}
	}

	return trafficStat, nil
}

// 连接 v2ray API
func (r *v2rayClient) Start(listen string) error {
	defer logger.Printf("v2rayClient start on %s, inbound: %s", listen, r.inboundTag)
	conn, err := grpc.Dial(listen, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		return err
	}
	r.statClient = statsService.NewStatsServiceClient(conn)
	r.handClient = command.NewHandlerServiceClient(conn)

	if err != nil {
		return err
	}
	return nil
}

func (r *v2rayClient) GetVConfig(vmessid string) string {
	if len(vmessid) < 6 {
		vmessid = "123456"
	}
	o := vConfig{
		Name:   r.domain + "_" + vmessid[:6],
		ID:     vmessid,
		Port:   r.port,
		Domain: r.domain,
		Path:   r.path,
	}
	var b, b2 bytes.Buffer
	VConfText.Execute(&b, o)

	b.WriteString("\n<b>Vmess 订阅:</b>\n")
	VConfJson.Execute(&b2, o)
	str := b64.StdEncoding.EncodeToString(b2.Bytes())
	b.WriteString("<code>vmess://" + str + "</code>")

	return b.String()
}

var _ V2rayService = (*v2rayClient)(nil)

// vmess tls
type vConfig struct {
	Name   string
	ID     string
	Port   int
	Domain string
	Path   string
}

var VConfText = template.Must(template.New("conftext").Parse(`
<b>Proxy Settings:</b>
协议类型: vmess
地址: {{.Domain}}
伪装域名/SNI: {{.Domain}}
端口: <code>{{.Port}}</code>
用户ID: <code>{{.ID}}</code>
安全: tls
传输方式: ws
路径: <code>{{.Path}}</code>
`))

var VConfJson = template.Must(template.New("confjson").Parse(`
{
   "add":"{{.Domain}}",
   "aid":"0",
   "host":"{{.Domain}}",
   "id":"{{.ID}}",
   "net":"ws",
   "path":"{{.Path}}",
   "port":"{{.Port}}",
   "ps":"{{.Name}}",
   "scy":"auto",
   "sni":"{{.Domain}}",
   "tls":"tls",
   "type":"",
   "v":"2"
}
`))

func init() {
	logger = log.New(os.Stderr, "[v2ray] ", log.LstdFlags|log.Lmsgprefix)
}
