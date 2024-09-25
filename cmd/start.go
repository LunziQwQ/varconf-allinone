package cmd

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"

	"varconf/core/moudle/router"
	"varconf/core/service"
	"varconf/core/web/controller"
	"varconf/core/web/interceptor"
	"varconf/core/web/resolver"

	_ "github.com/go-sql-driver/mysql"
)

type DatabaseInfo struct {
	Driver     string `json:"driver"`
	DataSource string `json:"dataSource"`
}

type ServerInfo struct {
	IP     string `json:"ip"`
	Port   int    `json:"port"`
	Static string `json:"static"`
}

type ServiceInfo struct {
	Cron string `json:"cron"`
}

type ConfigInfo struct {
	ServerInfo   ServerInfo   `json:"server"`
	DatabaseInfo DatabaseInfo `json:"database"`
	ServiceInfo  ServiceInfo  `json:"service"`
}

func Start(configPath, initFile string) error {
	configInfo := initConfig(configPath)
	if configInfo == nil {
		return errors.New("can't read config")
	}

	dbConnect := initDatabase(configInfo.DatabaseInfo)
	if dbConnect == nil {
		return errors.New("database connect error")
	}

	routeMux := initRouter(configInfo.ServerInfo)
	if routeMux == nil {
		return errors.New("router init error")
	}

	initMVC(routeMux, dbConnect, configInfo.ServiceInfo)

	if initFile != "" {
		loadInitFile(initFile)
	}

	return routeMux.Run()
}

func initConfig(configPath string) *ConfigInfo {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil
	}

	configInfo := ConfigInfo{}
	err = json.Unmarshal(data, &configInfo)
	if err != nil {
		return nil
	}

	return &configInfo
}

func initDatabase(database DatabaseInfo) *sql.DB {
	db, err := sql.Open(database.Driver, database.DataSource)
	if err != nil {
		panic(err)
	}
	db.Ping()
	return db
}

func initRouter(serverInfo ServerInfo) *router.Router {
	routeMux := router.NewRouter()
	routeMux.SetAddress(serverInfo.IP, serverInfo.Port)
	routeMux.Get("/", func(w http.ResponseWriter, r *http.Request, c *router.Context) {
		http.Redirect(w, r, "/static/html/index.html", http.StatusFound)
	})
	routeMux.Static("/static(.*)", serverInfo.Static, "index.html")

	return routeMux
}

func initMVC(routeMux *router.Router, dbConnect *sql.DB, serviceInfo ServiceInfo) {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	routeMux.SetLogger(logger)

	homeService := service.NewHomeService(dbConnect)
	authService := service.NewAuthService(dbConnect)
	userService := service.NewUserService(dbConnect)
	appService := service.NewAppService(dbConnect)
	configService := service.NewConfigService(dbConnect)

	interceptor.InitApiAuthInterceptor(routeMux, authService)
	interceptor.InitUserAuthInterceptor(routeMux, authService)
	resolver.InitErrorRecover(routeMux)

	controller.InitHomeController(routeMux, homeService)
	controller.InitApiController(routeMux, authService, configService)
	controller.InitUserController(routeMux, authService, userService)
	controller.InitAppController(routeMux, appService, configService)
	controller.InitConfigController(routeMux, configService)

	configService.CronRelease(serviceInfo.Cron)
}

func loadInitFile(initFile string) {
	return
}
