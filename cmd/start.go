package cmd

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"varconf/core/dao"
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
	IP   string `json:"ip"`
	Port int    `json:"port"`
}

type ServiceInfo struct {
	Cron string `json:"cron"`
}

type ConfigInfo struct {
	ServerInfo   ServerInfo   `json:"server"`
	DatabaseInfo DatabaseInfo `json:"database"`
	ServiceInfo  ServiceInfo  `json:"service"`
}

func Start(configPath, initFile string, recreate bool) error {
	configInfo := loadConfigFile(configPath)
	if configInfo == nil {
		return errors.New("can't read config")
	}

	dbConnect := initDatabase(configInfo.DatabaseInfo, recreate)
	if dbConnect == nil {
		return errors.New("database connect error")
	}

	routeMux := initRouter(configInfo.ServerInfo)
	if routeMux == nil {
		return errors.New("router init error")
	}

	if initFile != "" {
		fmt.Println("Apply init data...")
		applyInitData(initFile, dbConnect)
	}

	initMVC(routeMux, dbConnect, configInfo.ServiceInfo)

	return routeMux.Run()
}

func loadConfigFile(configPath string) *ConfigInfo {
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

func initDatabase(database DatabaseInfo, recreateTable bool) *sql.DB {
	db, err := sql.Open(database.Driver, database.DataSource)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic("Ping DB err: " + err.Error())
	}

	if recreateTable {
		fmt.Println("Clear tables...")
		// drop database if exist
		for tableName, sql := range dao.DropTableSql {
			_, err := db.Exec(sql)
			if err != nil {
				panic(fmt.Sprintf("Drop table %s err: %s", tableName, err.Error()))
			}
		}
	}

	// create table if not exist
	fmt.Println("Create tables...")
	for tableName, sql := range dao.CreateTableSql {
		_, err := db.Exec(sql)
		if err != nil {
			panic(fmt.Sprintf("Create table %s err: %s", tableName, err.Error()))
		}
	}

	// Insert default user
	_, err = db.Exec(dao.InsertDefaultUserSql)
	if err != nil {
		panic("Insert default user err:" + err.Error())
	}

	return db
}

func initRouter(serverInfo ServerInfo) *router.Router {
	routeMux := router.NewRouter()
	routeMux.SetAddress(serverInfo.IP, serverInfo.Port)
	routeMux.Get("/", func(w http.ResponseWriter, r *http.Request, c *router.Context) {
		http.Redirect(w, r, "/static/html/index.html", http.StatusFound)
	})
	routeMux.Static("/static(.*)", "./varconf-ui", "index.html")

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
