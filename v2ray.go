package nessielight

import (
	"bytes"
	"context"
	b64 "encoding/base64"
	"html/template"

	core "github.com/v2fly/v2ray-core/v4"
	"github.com/v2fly/v2ray-core/v4/app/proxyman"
	"github.com/v2fly/v2ray-core/v4/app/proxyman/command"
	statsService "github.com/v2fly/v2ray-core/v4/app/stats/command"
	"github.com/v2fly/v2ray-core/v4/common/net"
	"github.com/v2fly/v2ray-core/v4/common/protocol"
	"github.com/v2fly/v2ray-core/v4/common/serial"
	"github.com/v2fly/v2ray-core/v4/common/uuid"
	"github.com/v2fly/v2ray-core/v4/proxy/vmess"
	vmessInbound "github.com/v2fly/v2ray-core/v4/proxy/vmess/inbound"
	"github.com/v2fly/v2ray-core/v4/transport/internet"
	"github.com/v2fly/v2ray-core/v4/transport/internet/websocket"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// 调用 V2ray API 的客户端
// v2rayClient implements V2rayService
type v2rayClient struct {
	statClient statsService.StatsServiceClient
	handClient command.HandlerServiceClient
	// vmess settings
	inboundTag       string
	port, clientport int
	domain           string
	path             string
}

func (r *v2rayClient) AddVmessWsInbound(tag string, port uint16, wspath string) error {
	_, err := r.handClient.AddInbound(context.Background(), &command.AddInboundRequest{
		Inbound: &core.InboundHandlerConfig{
			Tag: tag,
			ReceiverSettings: serial.ToTypedMessage(&proxyman.ReceiverConfig{
				PortRange: net.SinglePortRange(net.Port(port)),
				Listen:    net.NewIPOrDomain(net.LocalHostIP), // 127.0.0.1
				StreamSettings: &internet.StreamConfig{
					ProtocolName: "websocket",
					TransportSettings: []*internet.TransportConfig{{
						ProtocolName: "websocket",
						Settings: serial.ToTypedMessage(&websocket.Config{
							Path: wspath,
						}),
					}},
				},
				SniffingSettings: &proxyman.SniffingConfig{
					Enabled:             true,
					DestinationOverride: []string{"http", "tls"},
				},
			}),
			ProxySettings: serial.ToTypedMessage(&vmessInbound.Config{
				User: []*protocol.User{},
			}),
		},
	})
	if err != nil {
		return err
	}
	logger.Printf("successfully add inbound %s, port=%d, path=%s", tag, port, wspath)
	return nil
}
func (r *v2rayClient) RemoveInbound(tag string) error {
	_, err := r.handClient.RemoveInbound(context.Background(), &command.RemoveInboundRequest{
		Tag: tag,
	})
	if err != nil {
		return err
	}
	logger.Printf("successfully remove inbound %s", tag)
	return nil
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
	logger.Printf("SetUser: email=%s id=%s", email, id)
	return nil
}

func (r *v2rayClient) AddUser(email string) (string, error) {
	userID := NewUUID()
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

type V2rayTrafficStat struct {
	Name  string
	Value int64
}

// reset is used to determine whether resetting traffic statistics
func (r *v2rayClient) QueryTraffic(pattern string, reset bool) ([]V2rayTrafficStat, error) {
	resp, err := r.statClient.QueryStats(context.Background(), &statsService.QueryStatsRequest{
		Pattern: pattern,
		Reset_:  reset,
	})
	if err != nil {
		return nil, err
	}

	stat := resp.GetStat()
	trafficStat := make([]V2rayTrafficStat, 0, len(stat))

	for _, v := range stat {
		if v != nil {
			trafficStat = append(trafficStat, V2rayTrafficStat{
				Name:  v.GetName(),
				Value: v.GetValue(),
			})
		}
	}

	return trafficStat, nil
}

// reset is used to determine whether resetting traffic statistics
func (r *v2rayClient) QueryUserTraffic(reset bool) ([]V2rayTrafficStat, error) {
	return r.QueryTraffic("user>>>"+r.inboundTag, reset)
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

// generate vmess link from vmessid
func (r *v2rayClient) VmessLink(vmessid string) string {
	if len(vmessid) < 6 {
		vmessid = "123456"
	}
	o := vConfig{
		Name:   r.domain + "_" + vmessid[:6],
		ID:     vmessid,
		Port:   r.clientport,
		Domain: r.domain,
		Path:   r.path,
	}
	var b2 bytes.Buffer
	VConfJson.Execute(&b2, o)
	str := b64.StdEncoding.EncodeToString(b2.Bytes())
	return "vmess://" + str
}

// generate vmess proxy description from vmessid
func (r *v2rayClient) VmessText(vmessid string) string {
	if len(vmessid) < 6 {
		vmessid = "123456"
	}
	o := vConfig{
		Name:   r.domain + "_" + vmessid[:6],
		ID:     vmessid,
		Port:   r.clientport,
		Domain: r.domain,
		Path:   r.path,
	}
	var b bytes.Buffer
	VConfText.Execute(&b, o)
	return b.String()
}

func (r *v2rayClient) NewProxy() Proxy {
	id := NewUUID()
	proxy := v2rayProxy{
		email: r.inboundTag + "-" + id,
		id:    id,
	}
	return &proxy
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

func NewUUID() string {
	return protocol.NewID(uuid.New()).String()
}

// implement Proxy
type v2rayProxy struct {
	email string
	id    string
}

func (r *v2rayProxy) ID() string {
	return r.id
}
func (r *v2rayProxy) Activate() error {
	V2rayServiceInstance.RemoveUser(r.email)
	return V2rayServiceInstance.SetUser(r.email, r.id)
}
func (r *v2rayProxy) Deactivate() error {
	return V2rayServiceInstance.RemoveUser(r.email)
}
func (r *v2rayProxy) Message() string {
	return "v2ray(vmess): <code>" + V2rayServiceInstance.VmessLink(r.id) + "</code>"
}

var _ Proxy = (*v2rayProxy)(nil)

// func (r *v2rayProxy) Value() (sqldriver.Value, error) {
// 	return r.id, nil
// }
// func (r *v2rayProxy) Scan(src interface{}) error {
// 	if id, ok := src.(string); ok {
// 		r.id = id
// 		return nil
// 	}
// 	return fmt.Errorf("invalid src type when scanning v2rayProxy")
// }
