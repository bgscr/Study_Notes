
# Gin框架

### 路由

    r := gin.Default()
    r.GET("/", func(c *gin.Context) {
        c.String(200, "Hello, Gin!")
    })
    r.Run(":8080") // 默认监听8080端口

    //路由分组
    v1 := r.Group("/v1")
    {
        v1.GET("/posts", listPosts)
        v1.POST("/posts", createPost)
    }

    //url路由方式 /welcome?firstname=John
    c.Query("firstname") 
    c.DefaultQuery("firstname", "Guest")
    // 获取路径参数
    r.GET("/user/:id", func(c *gin.Context) {
        id := c.Param("id")  
        c.JSON(200, gin.H{"id": id})
    })
    //表单数据
    type LoginForm struct {
        User     string `form:"user" binding:"required"`
        Password string `form:"password" binding:"required"`
    }
    func login(c *gin.Context) {
        var form LoginForm
        if err := c.ShouldBind(&form); err != nil { // 自动绑定表单数据
            c.JSON(400, gin.H{"error": err.Error()})
            return
        }
        // 处理表单数据
    }

    //渲染模板
    r.LoadHTMLGlob("templates/**/*")
    r.GET("/index", func(c *gin.Context) {
        c.HTML(200, "index.html", gin.H{"title": "Home"})
    })
    /*
    html中可以定义路径{{ define "templates/index.tmpl" }}    
    内容使用{{ .title }}填充
    */
### 参数验证

    ShouldBind: 自动识别参数，并绑定对应字段到结构体中
    ShouldBindJSON	tag:json
    ShouldBindXML		tag:xml
    ShouldBindQuery	tag:form
    ShouldBindYAML	tag:yaml
    ShouldBindTOML	tag:toml
    ShouldBindHeader	tag:header	无法自动识别
    ShouldBindUri	tag:uri	无法自动识别

    type User struct {
        Name  string `json:"name" binding:"required,min=3"`
        Email string `json:"email" binding:"required,email"`
    }
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

#### ShouldBindBodyWith 多次绑定对象
    objA := formA{}
    objB := formB{}
    // 读取 c.Request.Body 并将结果存入上下文。
    if errA := c.ShouldBindBodyWith(&objA, binding.JSON); errA == nil {
        c.String(http.StatusOK, `the body should be formA`)
    // 这时, 复用存储在上下文中的 body。
    } else if errB := c.ShouldBindBodyWith(&objB, binding.JSON); errB == nil {
        c.String(http.StatusOK, `the body should be formB JSON`)
    // 可以接受其他格式
    } else if errB2 := c.ShouldBindBodyWith(&objB, binding.XML); errB2 == nil {
        c.String(http.StatusOK, `the body should be formB XML`)
    } else {
        ...
    }

### MustBindWith 绑定错误直接返回400


### 中间件

    // 全局中间件，日志中间件示例
    r.Use(gin.Logger())

    //分组中间件
    authGroup := r.Group("/admin")
    authGroup.Use(AuthMiddleware())  // 仅/admin路由组使用


    //自定义中间件
    func AuthMiddleware() gin.HandlerFunc {
        return func(c *gin.Context) {
            token := c.GetHeader("Authorization")
            if token != "valid-token" {
                c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
            }
            c.Next()  // 继续后续处理
        }
    }

### HTTPS
    //本地认证证书
    router.RunTLS(":8080", "./testdata/server.pem", "./testdata/server.key")

    //Let Encrypt
    package main

    import (
    "log"

    "github.com/gin-gonic/autotls"
    "github.com/gin-gonic/gin"
    "golang.org/x/crypto/acme/autocert"
    )

    func main() {
    router := gin.Default()

    // Ping handler
    router.GET("/ping", func(c *gin.Context) {
        c.String(200, "pong")
    })

    m := autocert.Manager{
        Prompt:     autocert.AcceptTOS,
        HostPolicy: autocert.HostWhitelist("example1.com", "example2.com"),
        Cache:      autocert.DirCache("/var/www/.cache"),
    }

    log.Fatal(autotls.RunWithManager(r, &m))
    }