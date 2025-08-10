package main

import (
	"context"
	"errors"
	"fmt"
	"regexp"
    "encoding/json"
	"github.com/larksuite/oapi-sdk-go/v3"
	"github.com/larksuite/oapi-sdk-go/v3/core"
	"github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

// 定义飞书卡片结构（专业级样式）
type CardContent struct {
	Config struct {
		WideScreenMode bool `json:"wide_screen_mode"` // 宽屏模式
	} `json:"config"`
	Header struct {
		Title struct {
			Tag     string `json:"tag"`     // 标题类型
			Content string `json:"content"` // 标题内容
		} `json:"title"`
		Template string `json:"template"` // 标题配色模板
	} `json:"header"`
	Elements []struct {
		Tag  string `json:"tag"` // 元素类型
		Text struct {
			Content string `json:"content"` // 文本内容
			Lines   int    `json:"lines"`   // 显示行数
		} `json:"text"`
	} `json:"elements"`
}


// SDK 使用文档：开发前准备 - 服务端 API - 开发文档 - 飞书开放平台
// 复制该 Demo 后, 需要将 "YOUR_APP_ID", "YOUR_APP_SECRET" 替换为自己应用的 APP_ID, APP_SECRET.
// 以下示例代码默认根据文档示例值填充，如果存在代码问题，请在 API 调试台填上相关必要参数后再复制代码使用
func SendMsg (reqbody []byte) (string,error){
// 1. 验证并格式化输入内容
	if len(reqbody) == 0 {
		return "", errors.New("请求内容为空")
	}
	var prettyJSON interface{}
	if err := json.Unmarshal(reqbody, &prettyJSON); err != nil {
		return "", fmt.Errorf("内容格式错误: %v", err)
	}
	formattedJSON, _ := json.MarshalIndent(prettyJSON, "", "  ")

	// 2. 构造专业级卡片
	card := CardContent{}
	// 宽屏配置
	card.Config.WideScreenMode = true
	// 标题配置（使用警告配色模板）
	card.Header.Title.Tag = "plain_text"
	card.Header.Title.Content = "Hurry Up Important Alter Information"
	card.Header.Template = "yellow" // 红色标题模板（支持red/yellow/blue/grey）

    // 3. 添加内容区块（代码块+描述），重点修复转义和语法问题
    card.Elements = []struct {
        Tag  string `json:"tag"`
        Text struct {
            Content string `json:"content"`
            Lines   int    `json:"lines"`
        } `json:"text"`
    }{
        {
            Tag: "div",
            Text: struct {
                Content string `json:"content"`
                Lines   int    `json:"lines"`
            }{
                // 修复反斜杠转义，Go 中字符串里用 \\ 表示一个 \
                Content: fmt.Sprintf(
                   // "**告警详情**\n以下是系统检测到的告警数据\n```json\n%s\n```\n\n**处理建议**：\n1. 立即登录监控系统确认状态\n2. 检查相关主机资源使用情况\n3. 必要时触发应急预案",
                    "**告警详情**\n%s\n\n**处理建议**：\n1. 立即登录监控系统确认状态\n2. 检查相关主机资源使用情况\n3. 必要时触发应急预案",
                    formattedJSON,
                ),
                Lines: 30, // 限制最大显示行数
            },
        },
    }

	// 4. 转换为JSON字符串
	cardJSON, err := json.Marshal(card)
	if err != nil {
		return "", fmt.Errorf("构造卡片失败: %v", err)
	}

	// 5、创建 Client
	client := lark.NewClient(appconfig.AppID, appconfig.AppSecret)
	// 6、创建请求对象
	req := larkim.NewCreateMessageReqBuilder().
		ReceiveIdType(`chat_id`).
		Body(larkim.NewCreateMessageReqBodyBuilder().
			ReceiveId(appconfig.Chat_Id).
			//MsgType(`text`).
			//Content(fmt.Sprintf(`{"text":%q}`,string(reqbody))).
			MsgType("interactive").
			Content(string(cardJSON)).
			Uuid(``).
			Build()).
		Build()

	// 7、发起请求
	resp, err := client.Im.V1.Message.Create(context.Background(), req)

	// 8、处理错误
	if err != nil {
        LogPrint(LEVEL_FATAL,"resp failed,the error is %v\n",err)
		return "",errors.New("response failed")
	}

	// 9、服务端错误处理
	if !resp.Success() {
		LogPrint(LEVEL_FATAL,"logId: %s, error response: \n%s", resp.RequestId(), larkcore.Prettify(resp.CodeError))
		return "",errors.New("response failed")
	}

	// 10、业务处理
	LogPrint(LEVEL_INFO,"the response is %s\n",larkcore.Prettify(resp))
    
    re := regexp.MustCompile(`MessageId:\s*"([^"]+)"`)
    // 11、查找匹配的内容
    matches := re.FindStringSubmatch(larkcore.Prettify(resp))
    if len(matches) >= 2 {
        messageId := matches[1]
        LogPrint(LEVEL_INFO,"提取到的 MessageId 值为：%s\n", messageId)
        return messageId,nil
    } else {
        LogPrint(LEVEL_FATAL,"未找到 MessageId 对应的内容")
        return "",errors.New("not found the MessageId")
    }

}
