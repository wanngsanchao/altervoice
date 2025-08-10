package main

import (
	"context"
    "errors"
	"github.com/larksuite/oapi-sdk-go/v3"
	"github.com/larksuite/oapi-sdk-go/v3/core"
	"github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

// SDK 使用文档：开发前准备 - 服务端 API - 开发文档 - 飞书开放平台
// 复制该 Demo 后, 需要将 "YOUR_APP_ID", "YOUR_APP_SECRET" 替换为自己应用的 APP_ID, APP_SECRET.
// 以下示例代码默认根据文档示例值填充，如果存在代码问题，请在 API 调试台填上相关必要参数后再复制代码使用
func SendVoice(msgid string) error {
	// 创建 Client
	client := lark.NewClient(appconfig.AppID, appconfig.AppSecret)
	// 创建请求对象
	req := larkim.NewUrgentPhoneMessageReqBuilder().
		MessageId(msgid).
		UserIdType(`user_id`).
		UrgentReceivers(larkim.NewUrgentReceiversBuilder().
			UserIdList(appconfig.UserList).
			Build()).
		Build()

	// 发起请求
	resp, err := client.Im.V1.Message.UrgentPhone(context.Background(), req)

	// 处理错误
	if err != nil {
        LogPrint(LEVEL_INFO,"the err is %v",err)
		return errors.New("resp failed")
	}

	// 服务端错误处理
	if !resp.Success() {
		LogPrint(LEVEL_INFO,"logId: %s, error response: \n%s", resp.RequestId(), larkcore.Prettify(resp.CodeError))
		return errors.New("resp failed")
	}

	// 业务处理
	LogPrint(LEVEL_INFO,"%s\n",larkcore.Prettify(resp))
    return nil
}
