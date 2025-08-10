package main

import (
	"encoding/json"
	"errors"
	"os"
)

const (
    ConfigPath = "/etc/altervoice/altervoice.json"
)

// define config struct
type ConfigJson struct {
	AppID     string `json:"app_id"`
	AppSecret string `json:"app_secret"`
    Chat_Id string `json:"chat_id"`
    UserList []string `json:"user_list"`
}

var appconfig = ConfigJson{}

func GetConfig(config string) (error) {
    data,err := os.ReadFile(config)

    if err != nil {
        LogPrint(LEVEL_ERROR,"the json config path is nil")        
        return errors.New("config is nil")
    }
    
    err = json.Unmarshal(data,&appconfig)
    if err != nil {
        LogPrint(LEVEL_ERROR,"format the appid and appsecret failed")
        return errors.New("format the appid and appsecret failed")
    }

    return nil
}
