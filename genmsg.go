package main

import (
  "encoding/json"
  "errors"
)

/*the requst json is send by the altermanage
{
  "receiver": "web\\.hook\\.prometheusalert",
  "status": "firing",
  "alerts": [
    {
      "status": "firing",
      "labels": {
        "alertname": "HostMemoryUsage",
        "group": "超算openstack",
        "hostname": "mn4",
        "instance": "20.20.0.4:9101",
        "job": "超算",
        "level": "3",
        "nodename": "mn4",
        "severity": "critical"
      },
      "annotations": {
        "description": "超算openstack主机: 【20.20.0.4:9101】 内存使用率超过95% (当前使用率：70.31211974246727%)",
        "title": "主机内存使用率超过95"
      },
      "startsAt": "2025-08-01T11:47:18.831Z",
      "endsAt": "0001-01-01T00:00:00Z",
      "generatorURL": "http://hz-center-promethues:9090/graph?g0.expr=%28%281+-+%28%28node_memory_Buffers_bytes+%2B+node_memory_Cached_bytes+%2B+node_memory_MemFree_bytes%29+%2F+node_memory_MemTotal_bytes%29%29+%2A+100+%3E+70%29+%2A+on+%28instance%29+group_left+%28nodename%29+max+by+%28instance%2C+nodename%29+%28node_uname_info%29\u0026g0.tab=1",
      "fingerprint": "2072b1e4351e7281"
    }
  ],
  "groupLabels": {
    "instance": "20.20.0.4:9101"
  },
}
*/

// the requst json is send by the altermanage
type AlertManageMessage struct {
  Receiver          string            `json:"receiver"`
  Status            string            `json:"status"`
  Alerts            []AlertItem       `json:"alerts"`
  GroupLabels       GroupLab `json:"groupLabels"`
  CommonLabels      map[string]string `json:"commonLabels"`
  CommonAnnotations map[string]string `json:"commonAnnotations"`
  ExternalURL       string            `json:"externalURL"`
  Version           string            `json:"version"`
  GroupKey          string            `json:"groupKey"`
  TruncatedAlerts   int               `json:"truncatedAlerts"`
}

// define the alter details struct
type AlertItem struct {
  Status       string            `json:"status"`
  Labels       Labe              `json:"labels"`
  Annotations  Anno              `json:"annotations"`
  StartsAt     string            `json:"startsAt"`
  EndsAt       string            `json:"endsAt"`
  GeneratorURL string            `json:"generatorURL"`
  Fingerprint  string            `json:"fingerprint"`
}

//define the anno struct
type Anno struct {
    Description string `json:"description"`
    Title string `json:"title"`
}

//define the alter info label
type Labe struct {
    Status   string `json:"status"`   // 告警状态（如 firing）
	Title    string `json:"title"`    // 告警标题
	Severity string `json:"severity"` // 告警级别（如 critical）
	Job      string `json:"job"`      // 任务名称
	Group    string `json:"group"`    // 所属分组
	Hostname string `json:"hostname"` // 主机名
	Instance string `json:"instance"` // 实例地址（IP:端口）
}

//define the alter group label
type GroupLab struct {
    Instance string `json:'instance'`
}

/*the request body will be send by the interface of the SendMsg as feishuapp request_body
{
  "status": "firing",
  "title": "主机内存使用率超过95",
  "severity": "critical",
  "job": "超算",
  "group": "超算openstack",
  "hostname": "mn4",
  "instance": "20.20.0.4:9101",
}
*/

// used for SendMsg as a request body
type SendMsgReqBody struct {
  Status   string `json:"status"`   // 告警状态（如 firing）
  Title    string `json:"title"`    // 告警标题
  Severity string `json:"severity"` // 告警级别（如 critical）
  Startsat string `json:"startsAt"` // 告警开始时间
  Job      string `json:"job"`      // 任务名称
  Group    string `json:"group"`    // 所属分组
  Hostname string `json:"hostname"` // 主机名
  Instance string `json:"instance"` // 实例地址（IP:端口）
}

var AlterMnageInfo AlertManageMessage
var SendMsgReqInfo SendMsgReqBody

// 解析示例
func GenerateRequestbody(requestbody []byte) ([]byte,error) {

  //parse the altermanage request body to the structure of the AlterMnageInfo
  if err := json.Unmarshal(requestbody, &AlterMnageInfo); err != nil {
    LogPrint(LEVEL_FATAL,"parse the altermanage request body failed: %v\n", err)
    return []byte(""),errors.New("get request body failed from altermanage request")
  }

  //make a the structure of the SendMsgReqBody,used to test with the first alterinfo
  SendMsgReqInfo.Status = AlterMnageInfo.Alerts[0].Status
  SendMsgReqInfo.Title = AlterMnageInfo.Alerts[0].Annotations.Title
  SendMsgReqInfo.Severity = AlterMnageInfo.Alerts[0].Labels.Severity
  SendMsgReqInfo.Job = AlterMnageInfo.Alerts[0].Labels.Job
  SendMsgReqInfo.Group = AlterMnageInfo.Alerts[0].Labels.Group
  SendMsgReqInfo.Hostname = AlterMnageInfo.Alerts[0].Labels.Hostname
  SendMsgReqInfo.Instance = AlterMnageInfo.GroupLabels.Instance

  //transform the structure of the SendMsgReqBody to the json byte
  requestbody,err := json.Marshal(SendMsgReqInfo)

  if err != nil {
    return []byte(""),errors.New("transform the request body for the request-body of the SendMsg failed")
  }

  return requestbody,nil
}
