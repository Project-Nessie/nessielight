package nessielight

import (
	"fmt"

	"github.com/Project-Nessie/nessielight/utils"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type simpleUser struct {
	gorm.Model
	Registerid int
	Nam        string
	// proxy id
	V2rayProxyID pq.Int32Array `gorm:"type:text[]"`
	Traff        TrafficValue  `gorm:"embedded"`
}

func (r *simpleUser) TelegramID() int {
	return r.Registerid
}

func (r *simpleUser) Name() string {
	return r.Nam
}

func (r *simpleUser) Proxy() []Proxy {
	return utils.Map(r.V2rayProxyID, func(id int32) Proxy {
		var v2rayproxy v2rayProxy
		DataBase.First(&v2rayproxy, id)
		return &v2rayproxy
	})
}

func (r *simpleUser) SetProxy(proxy []Proxy) error {
	r.V2rayProxyID = make([]int32, 0, len(r.V2rayProxyID))
	for _, v := range proxy {
		if v2rayproxy, ok := v.(*v2rayProxy); ok {
			r.V2rayProxyID = append(r.V2rayProxyID, int32(v2rayproxy.ID))
		} else {
			panic(fmt.Errorf("unknown proxy type"))
		}
	}
	return nil
}

func (r *simpleUser) SetName(name string) error {
	r.Nam = name
	return nil
}

func (r *simpleUser) Traffic() TrafficValue {
	return r.Traff
}
func (r *simpleUser) SetTraffic(val TrafficValue) error {
	r.Traff = val
	return nil
}

var _ User = (*simpleUser)(nil)
