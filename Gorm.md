# 连接数据库
### 方法一
    import (
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    )

    func main() {
    // 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
    dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    }
### 方法二
    db, err := gorm.Open(mysql.New(mysql.Config{
        DSN: "gorm:gorm@tcp(127.0.0.1:3306)/gorm?charset=utf8&parseTime=True&loc=Local", // DSN data source name
        DefaultStringSize: 256, // string 类型字段的默认长度
        DisableDatetimePrecision: true, // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
        DontSupportRenameIndex: true, // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
        DontSupportRenameColumn: true, // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
        SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
        }), &gorm.Config{})
### 连接池
    // 获取通用数据库对象 sql.DB ，然后使用其提供的功能
    sqlDB, err := db.DB()

    // SetMaxIdleConns 用于设置连接池中空闲连接的最大数量。
    sqlDB.SetMaxIdleConns(10)

    // SetMaxOpenConns 设置打开数据库连接的最大数量。
    sqlDB.SetMaxOpenConns(100)

    // SetConnMaxLifetime 设置了连接可复用的最大时间。
    sqlDB.SetConnMaxLifetime(time.Hour)
# 模型定义
    type User struct {
        ID      uint   `gorm:"primaryKey"` // 主键
        Name    string `gorm:"size:50"`    // 列名默认 snake_case，如 name
        Email   *string                   // 指针类型支持NULL
        Age     int    `gorm:"default:18"`// 默认值
        MemberNumber sql.NullString // Uses sql.NullString to handle nullable strings
        ActivatedAt  sql.NullTime   // Uses sql.NullTime for nullable time fields

        Author //嵌套Author的字段
        Author  Author `gorm:"embedded;embeddedPrefix:author_"` //嵌套并设置字段前缀
    }

    type Author struct {
        Name  string
        Email string
   }

    *string 和 *time.Time 类型的指针表示可空字段

    来自 database/sql 包的 sql.NullString 和 sql.NullTime 用于具有更多控制的可空字段。例如**user.MemberNumber.Valid =true**

    首字母小写会被忽略字段

[更多标签应用](https://gorm.io/zh_CN/docs/models.html#%E5%AD%97%E6%AE%B5%E6%A0%87%E7%AD%BE)

