package main

import (
    "net/http"
    "encoding/json"
    "io/ioutil"
)

func AlterVoice(w http.ResponseWriter,r * http.Request) {
    //get the request body from the alter-manage
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		LogPrint(LEVEL_FATAL,"read the request body failed: %v\n", err)
		http.Error(w, "read the request body failed", http.StatusInternalServerError)
		return
	}

	//print the alter info from alter-manage with string
	LogPrint(LEVEL_INFO,"the request body fomr the alter-manage  is: %s\n", string(body))
    
    //genarate the alter info send to freshu-app
    reqcontext,err := GenerateRequestbody(body)
    if err != nil {
		LogPrint(LEVEL_FATAL,"GenerateRequestbody failed: %v\n", err)
        reqcontext = []byte("there is alterinfo from alter-manage,but the info is parsing failed")
    }
    LogPrint(LEVEL_INFO,"here is the request nody:%s\n",string(reqcontext))

    //get the feishu config that includes app-id,app-secret,chat-id,user-list from the config file altervoice.json
    if err := GetConfig(ConfigPath); err != nil {
		LogPrint(LEVEL_FATAL,"read feishu app config failed: %v\n", err)
		http.Error(w, "GetConfig failed", http.StatusInternalServerError)
    }
    LogPrint(LEVEL_INFO,"read feishu app config success,the config is %v\n",appconfig)

    //sending msg with feishu api is to get the message-id
    messageid,err := SendMsg(reqcontext);
    if err != nil {
		LogPrint(LEVEL_FATAL,"SendMsg failed: %v\n", err)
		http.Error(w, "GetConfig failed", http.StatusInternalServerError)
    }
    LogPrint(LEVEL_INFO,"SendMsg success,and the messageid is: %s\n", messageid)

    //send voice with feishu api with the above message-id
    if err := SendVoice(messageid); err != nil {
		LogPrint(LEVEL_FATAL,"SendVoice failed: %v\n", err)
		http.Error(w, "SendVoice failed", http.StatusInternalServerError)
    }

    LogPrint(LEVEL_INFO,"SendVoice success\n")

    w.WriteHeader(200)
    w.Header().Set("Content-Type","application/json")
    data := map[string]string{"send_voice_status":"success"}

    jsondata,err := json.Marshal(data)
    if err != nil {
        http.Error(w,err.Error(),http.StatusInternalServerError)
        return
    }

    w.Write(jsondata)
} 
