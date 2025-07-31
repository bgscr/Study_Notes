package svc

import (
	"fmt"
	"testGoZero/internal/config"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type ServiceContext struct {
	Config config.Config
	GormDB *gorm.DB // 添加 GORM 连接
}

func NewServiceContext(c config.Config) *ServiceContext {
	db, _ := gorm.Open(mysql.Open(getDSN(&c)), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   c.MySQLConf.Gorm.TablePrefix,
			SingularTable: c.MySQLConf.Gorm.SingularTable,
		},
	})
	// 获取通用数据库对象 sql.DB ，然后使用其提供的功能
	sqlDB, err := db.DB()
	// SetMaxIdleConns 用于设置连接池中空闲连接的最大数量。
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err != nil {
		logx.Must(err)
	}

	return &ServiceContext{
		Config: c,
		GormDB: db,
	}
}

// 构建 DSN 字符串
func getDSN(c *config.Config) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.MySQLConf.User,
		c.MySQLConf.Password,
		c.MySQLConf.Host,
		c.MySQLConf.Port,
		c.MySQLConf.Database,
	)
}
