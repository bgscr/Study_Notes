/cmd	存放 可执行程序入口，每个子目录对应一个 main 包	app1/main.go, app2/main.go

/internal	存储 私有代码（仅当前项目可访问），Go 语言强制禁止外部导入	内部工具类、业务逻辑实现代码

/pkg	存放 公开库代码（允许外部项目导入）	通用工具包、第三方库适配层

/api	定义 API 协议（如 Protobuf、OpenAPI/Swagger 文件）	.proto, swagger.yaml

/web	Web 应用专属资源（模板、静态文件等）	static/, templates/

/configs	配置文件（YAML/TOML/环境变量等）	config.yaml, .env

/test	全局测试代码和测试数据	集成测试、E2E 测试

/docs	项目文档	design.md, api-spec.md

/scripts	构建、部署、代码生成等脚本	build.sh, deploy.go
