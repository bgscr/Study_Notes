package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	Auth struct {
		AccessSecret string
		AccessExpire int64
	}
	UserRPC   zrpc.RpcClientConf
	PostRPC   zrpc.RpcClientConf
	MySQLConf MySQLConf `json:",omitempty"`
}

type MySQLConf struct {
	Host     string   `json:",omitempty"`
	Port     int      `json:",omitempty"`
	User     string   `json:",omitempty"`
	Password string   `json:",omitempty"`
	Database string   `json:",omitempty"`
	Gorm     GormConf `json:",omitempty"`
}

type GormConf struct {
	TablePrefix   string `json:",omitempty"` // 表前缀
	SingularTable bool   `json:",omitempty"` // 是否单数表名
}
