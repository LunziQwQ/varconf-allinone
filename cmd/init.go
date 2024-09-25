package cmd

import (
	"database/sql"
	"fmt"
	"os"
	"varconf/core/dao"
	"varconf/core/service"

	"gopkg.in/yaml.v3"
)

type initConfigItem struct {
	Key   string `yaml:"key"`
	Value string `yaml:"value"`
}

type initAppConfig struct {
	AppId    int64             `yaml:"appId"`
	AppName  string            `yaml:"appName"`
	AppDesc  string            `yaml:"appDesc"`
	AppToken string            `yaml:"appToken"`
	Configs  []*initConfigItem `yaml:"configs"`
}

type initConfigData struct {
	Apps []*initAppConfig `yaml:"apps"`
}

func applyInitData(initFile string, dbConnect *sql.DB) error {
	if initFile == "" {
		return nil
	}

	// load yaml content
	yamlData, err := os.ReadFile(initFile)
	if err != nil {
		panic(fmt.Sprintf("Read init file err: %s", err.Error()))
	}

	var initData initConfigData
	err = yaml.Unmarshal(yamlData, &initData)
	if err != nil {
		panic(fmt.Sprintf("Unmarshal init file err: %s", err.Error()))
	}

	// apply init data
	appService := service.NewAppService(dbConnect)
	configService := service.NewConfigService(dbConnect)
	for _, appConfig := range initData.Apps {
		// check app exist and create
		app := appService.QueryApp(appConfig.AppId)
		if app != nil {
			fmt.Printf("App %d already exist, skip init.", appConfig.AppId)
			continue
		}
		appService.CreateApp(&dao.AppData{
			AppId:        appConfig.AppId,
			Name:         appConfig.AppName,
			Desc:         appConfig.AppDesc,
			ApiKey:       appConfig.AppToken,
			Code:         appConfig.AppName,
			ReleaseIndex: 0,
		})

		// create config
		for _, configItem := range appConfig.Configs {
			configService.CreateConfig(&dao.ConfigData{
				AppId:    appConfig.AppId,
				ConfigId: 1,
				Key:      configItem.Key,
				Value:    configItem.Value,
				Operate:  dao.OPERATE_NEW,
			})

			configService.ReleaseConfig(appConfig.AppId, "admin")
		}
	}
	return nil
}
