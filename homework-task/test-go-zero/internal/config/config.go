package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	Salt string
	Auth struct {
		AccessSecret string
		AccessExpire int64
	}

	RpcServer zrpc.RpcServerConf `json:",optional"`
	MySQLConf MySQLConf
}

type MySQLConf struct {
	Host     string `json:",optional"`
	Port     int    `json:",optional"`
	User     string `json:",optional"`
	Password string `json:",optional"`
	Database string `json:",optional"`
	Gorm     GormConf
}

type GormConf struct {
	TablePrefix   string `json:",optional"` // 表前缀
	SingularTable bool   `json:",optional"` // 是否单数表名
}
