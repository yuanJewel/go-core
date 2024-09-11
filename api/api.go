package api

import (
	"github.com/SmartLyu/go-core/logger"
	"github.com/iris-contrib/middleware/cors"
	"github.com/iris-contrib/swagger/v12"
	"github.com/iris-contrib/swagger/v12/swaggerFiles"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/core/router"
	"github.com/kataras/iris/v12/middleware/jwt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/version"
	"log"
)

type UserClaims struct {
	Username string `json:"username"`
}

var (
	// Reverser 全局调用自身路由的对象
	Reverser *router.RoutePathReverser
	// Auth ApiAuth 存储Api认证
	Auth string
	// AuthenticateUrl 存储Authenticate组件地址
	AuthenticateUrl string
)

type Service interface {
	// GetAuth 传入Api使用的Auth信息
	GetAuth() string
	// GetAuthenticateUrl 传入Authenticate组件地址
	GetAuthenticateUrl() string
	// Party 自定义路由逻辑
	Party(iris.Party)
	// Health 定义健康检查接口需要检查的内容
	Health() func() map[string]error
	// Dot 接口请求访问的需要打点逻辑，可以定义:
	Dot(...interface{})

	// AuthenticateApi 自定义登录认证的接口
	AuthenticateApi(iris.Context)
	// Authenticate 请求Authenticate组件的方法
	Authenticate(ctx iris.Context) error
}

func CreateApi(service Service, isSwagger bool) (*iris.Application, func() error) {
	app := iris.New()
	app.Get("/", func(ctx iris.Context) {
		_ = ctx.JSON(map[string]string{
			"version":      version.Version,
			"go version":   version.GoVersion,
			"writer email": version.BuildUser,
			"build time":   version.BuildDate,
			"git branch":   version.Branch,
			"git revision": version.Revision,
			"swagger path": "/swagger/index.html",
		})
	}).Name = "/"
	// 调用本地接口对象
	Reverser = router.NewRoutePathReverser(app)

	// 配置访问日志的对象，如果需要打点逻辑
	r, _close := logger.NewRequestLogger(service.Dot)
	app.Use(r)

	// 支持跨域访问
	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // allows everything, use that to change the hosts.
		AllowedMethods:   []string{"HEAD", "OPTIONS", "GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})
	app.UseRouter(crs)
	app.AllowMethods(iris.MethodOptions)

	// 处理静态资源文件
	app.RegisterView(iris.HTML("./views", ".html"))

	// swagger
	if isSwagger {
		swaggerConfig := &swagger.Config{
			URL:         "swagger/doc.json", //The url pointing to API definition
			DeepLinking: true,
		}
		swaggerUI := swagger.CustomWrapHandler(swaggerConfig, swaggerFiles.Handler)
		app.Get("/swagger", swaggerUI)
		app.Get("/swagger/{any:path}", swaggerUI)
	}

	// monitor
	app.Get("/health", healthCheckHandle(service.Health())).Name = "get_health_check"
	app.Get("/metrics", iris.FromStd(promhttp.Handler())).Name = "get_metrics"

	// 权限管理
	app.Post("/authenticate", service.AuthenticateApi).Name = "post_authenticate"

	// 增加token 验证
	Auth = service.GetAuth()
	AuthenticateUrl = service.GetAuthenticateUrl()
	verifier := jwt.NewVerifier(jwt.HS256, Auth)
	verifyMiddleware := verifier.Verify(func() interface{} {
		return new(UserClaims)
	})
	service.Party(app.Party("/api/v1", verifyMiddleware, irisAuthenticate(service.Authenticate)))

	return app, _close
}

type Object struct{}

func (Object) GetAuth() string { return "" }

func (Object) GetAuthenticateUrl() string { return "" }

func (Object) Party(iris.Party) {}

func (Object) Health() func() map[string]error {
	return func() map[string]error {
		return map[string]error{}
	}
}

func (Object) Dot(i ...interface{}) {
	var (
	//traceid = i[0].(string)
	//username = i[1].(string)
	//ctx      = i[2].(iris.Context)
	//
	//uri = ctx.Path()
	//method = ctx.Method()
	)
	log.Println(i)
}
