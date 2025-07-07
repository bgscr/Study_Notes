# 连接数据库
### 安装
    go get -u gorm.io/gorm
    go get -u gorm.io/driver/sqlite
    go get -u gorm.io/driver/mysql
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

# CRUD操作
### 插入
    // 单条插入
    user := User{Name: "Alice", Age: 25}
    result := db.Create(&user)         // 同步插入
    fmt.Println(user.ID)               // 获取自增ID
    fmt.Println(result.Error)          // 错误检查

    // 批量插入
    users := []User{{Name: "Bob"}, {Name: "Charlie"}}
    db.CreateInBatches(users, 100)     // 每批100条

### 查询
    // 基础查询
    var user User
    db.First(&user)        // 按主键排序第一条
    db.Take(&user)         // 随机一条
    db.Last(&user)         // 按主键最后一条

    // 条件查询
    db.Where("age > ?", 20).Find(&users)                   // 条件查询
    db.Where(&User{Name: "Alice", Age: 25}).First(&user)   // 结构体条件（忽略零值）
    db.Where("name LIKE ?", "%Ali%").Find(&users)          // 模糊查询

    // 预加载关联（假设User有Orders关联）
    db.Preload("Orders").Find(&users)

    // 错误处理
    if err := db.First(&user, 999).Error; errors.Is(err, gorm.ErrRecordNotFound) {
        fmt.Println("记录不存在")
    }

### 更新
    // 更新单个字段
    db.Model(&User{}).Where("id = ?", 1).Update("name", "Alice Updated")

    // 更新多个字段（忽略零值）
    db.Model(&user).Updates(User{Name: "Alice", Age: 26})

    // 选择字段更新
    db.Model(&user).Select("Name").Updates(User{Name: "Bob", Age: 30}) // 仅更新Name
    db.Model(&user).Omit("Age").Updates(User{Name: "Bob", Age: 30})    // 排除Age

### 删除
    // 软删除（需模型包含DeletedAt字段）
    db.Delete(&user)       // UPDATE users SET deleted_at = NOW() WHERE id = ?

    // 物理删除
    db.Unscoped().Delete(&user)  // DELETE FROM users WHERE id = ?

### 事务
    tx := db.Begin()
    if err := tx.Create(&user).Error; err != nil {
        tx.Rollback()
    }
    tx.Commit()


### 调试SQL
    调试 SQL：通过 db.Debug() 查看生成的 SQL 语句。

# 关联
### BelongsTo
    type Dog struct {
        gorm.Model
        BeautyID uint   // 外键
        Beauty   Beauty // 所属模型
    }
    type Beauty struct {
        gorm.Model
    }

### HasOne
    type User struct {
        gorm.Model
        Profile Profile
    }
    type Profile struct {
        UserID uint // 外键
    }

### HasMany
    type User struct {
        gorm.Model
        Orders []Order
    }
    type Order struct {
        UserID uint // 外键
    }

### ManyToMany 通过中间表处理多对多关系。
    type User struct {
        Languages []Language `gorm:"many2many:user_languages;"`
    }
    type Language struct {
        Name string
    }

### Preload
    var user User
    db.Preload("Orders").First(&user, 1)      // 预加载所有订单
    db.Preload("Orders", "state = ?", "paid").First(&user, 1) // 条件筛选

### 动态关联Association

    // 添加关联
    db.Model(&user).Association("Languages").Append(&Language{Name: "Go"})

    // 替换关联
    db.Model(&user).Association("Languages").Replace([]Language{{Name: "Python"}})

    // 删除关联
    db.Model(&user).Association("Languages").Delete(&Language{Name: "Java"})

    // 清空关联
    db.Model(&user).Association("Languages").Clear()

### 多态关联Polymorphism
    // Comment 模型（多态关联的从表）
    type Comment struct {
        ID        uint   `gorm:"primaryKey"`
        Content   string
        TargetType string `gorm:"column:target_type"` // 通过 polymorphicType 自定义类型字段名
        TargetID   uint   `gorm:"column:target_id"`   // 通过 polymorphicId 自定义ID字段名
    }

    // Article 模型（主表）
    type Article struct {
        ID       uint   `gorm:"primaryKey"`
        Title    string
        Comments []Comment `gorm:"polymorphic:Owner; polymorphicType:TargetType; polymorphicId:TargetID; polymorphicValue:blog_post"`
    }

    // Video 模型（主表）
    type Video struct {
        ID       uint   `gorm:"primaryKey"`
        Title    string
        Comments []Comment `gorm:"polymorphic:Owner; polymorphicType:TargetType; polymorphicId:TargetID; polymorphicValue:video_clip"`
    }

    // 创建文章及其评论
    article := Article{
        Title: "GORM 多态关联指南",
        Comments: []Comment{
            {Content: "非常实用的教程！"},
            {Content: "期待更多案例！"},
        },
    }
    db.Create(&article)

    // 创建视频及其评论
    video := Video{
        Title: "Go 实战演示",
        Comments: []Comment{
            {Content: "讲解清晰！"},
        },
    }
    db.Create(&video)

    //实际数据
    ID	Content	TargetType	TargetID
    1	非常实用的教程！	blog_post	1
    2	期待更多案例！	blog_post	1
    3	讲解清晰！	video_clip	1
    