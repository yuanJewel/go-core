package main

import (
	"fmt"
	"github.com/alecthomas/kingpin/v2"
	"github.com/kataras/iris/v12"
	"github.com/prometheus/common/version"
	apiInterface "github.com/yuanJewel/go-core/api"
	"github.com/yuanJewel/go-core/asset"
	"github.com/yuanJewel/go-core/db/service"
	_ "github.com/yuanJewel/go-core/docs"
	"github.com/yuanJewel/go-core/logger"
	"github.com/yuanJewel/go-core/pkg/api"
	"github.com/yuanJewel/go-core/pkg/config"
	"github.com/yuanJewel/go-core/pkg/db"
	"log"
	"time"
)

func init() {
	dirs := []string{"views", "docs"}
	for _, dir := range dirs {
		if err := asset.RestoreAssets("./", dir); err != nil {
			fmt.Printf("解压%s失败\n", dir)
		}
	}
}

// @title Swagger yuanJewel go-core API
// @version 1.3.5
// @description yuanJewel go-core API
// @contact.name yuanJewel go-core Support

// @contact.url https://github.com/yuanJewel/go-core/blob/main/README.md
// @contact.email luyu151111@gmail.com
// @securityDefinitions.apikey  ApiKeyAuth
// @in header
// @name Authorization
// @host
// @BasePath /
func main() {
	var (
		configPath = kingpin.Flag("config", "go-core config file,default application.yml").
				Default("application.yml").Short('c').String()
		initDb = kingpin.Flag("init", "是否初始化数据库").Short('i').Bool()
	)
	kingpin.Version(version.Print("go-core"))
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	logger.PrintLogStatus()
	if err := config.LoadConfig(*configPath); err != nil {
		log.Fatal("Load Config Error...", err)
	}

	if err := service.InitDb(&config.GlobalConfig.DataSourceDetail); err != nil {
		log.Fatal("Init Database Error...", err)
	}

	if *initDb {
		if err := db.SetupCmdb(); err != nil {
			log.Fatal("Init Database Error...", err)
		}
		log.Println("Init Database Success...")
		return
	}

	app, _close := apiInterface.CreateApi(api.Object{}, config.GlobalConfig.Swagger)
	app.Configure(iris.WithConfiguration(iris.Configuration{
		Timeout:           time.Duration(config.GlobalConfig.HttpTimeout) * time.Second,
		LogLevel:          config.GlobalConfig.LogLevel,
		DisableStartupLog: config.GlobalConfig.DisableStartupLog,
	}))
	defer func() {
		_ = _close()
	}()

	logger.Log.Info("服务已运行...")
	if err := app.Run(iris.Addr(fmt.Sprintf(":%d", config.GlobalConfig.Server.Port))); err != nil {
		log.Fatal("Start Api Error...")
	}
}
