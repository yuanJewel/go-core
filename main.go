package main

import (
	"fmt"
	apiInterface "github.com/SmartLyu/go-core/api"
	"github.com/SmartLyu/go-core/asset"
	"github.com/SmartLyu/go-core/cmdb"
	_ "github.com/SmartLyu/go-core/docs"
	"github.com/SmartLyu/go-core/pkg/api"
	"github.com/SmartLyu/go-core/pkg/config"
	"github.com/SmartLyu/go-core/pkg/db"
	"github.com/alecthomas/kingpin/v2"
	"github.com/kataras/iris/v12"
	"github.com/prometheus/common/version"
	"log"
)

func init() {
	dirs := []string{"views", "docs"}
	for _, dir := range dirs {
		if err := asset.RestoreAssets("./", dir); err != nil {
			fmt.Printf("解压%s失败\n", dir)
		}
	}
}

// @title Swagger SmartLyu go-core API
// @version 1.0.2
// @description SmartLyu go-core API
// @contact.name SmartLyu go-core Support

// @contact.url https://github.com/SmartLyu/go-core/blob/main/README.md
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

	if err := config.LoadConfig(*configPath); err != nil {
		log.Fatal("Load Config Error...", err)
	}

	if *initDb {
		if err := db.SetupCmdb(&config.GlobalConfig.DataSourceDetail); err != nil {
			log.Fatal("Init Database Error...", err)
		}
		log.Println("Init Database Success...")
		return
	}

	if err := cmdb.InitCmdb(&config.GlobalConfig.DataSourceDetail); err != nil {
		log.Fatal("Init Database Error...", err)
	}

	app, _close := apiInterface.CreateApi(api.Object{}, config.GlobalConfig.Swagger)

	defer func() {
		_ = _close()
	}()

	if err := app.Run(iris.Addr(fmt.Sprintf(":%d", config.GlobalConfig.Server.Port))); err != nil {
		log.Fatal("Start Api Error...")
	}
}
