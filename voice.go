package qcloudsms

import (
	"encoding/json"
	"errors"
)

// VoiceReq 语音接口请求结构
type VoiceReq struct {
	// 手机号码结构
	Tel struct {
		Nationcode string `json:"nationcode"`
		Mobile     string `json:"mobile"`
	} `json:"tel"`

	// 语音类型，为2表示语音通知
	Prompttype uint `json:"prompttype,omitempty"`

	// 语音内容，语音类型为通知时有效
	Promptfile string `json:"promptfile,omitempty"`

	// 语音内容，语音类型为验证码通知时有效
	Msg string `json:"msg,omitempty"`

	// 播放次数
	Playtimes uint   `json:"playtimes"`
	Sig       string `json:"sig"`
	Time      int64  `json:"time"`
	Ext       string `json:"ext"`
}

// VoiceResult 语音发送后返回结果的结构
type VoiceResult struct {
	Result uint   `json:"result"`
	Errmsg string `json:"errmsg"`
	Ext    string `json:"ext"`
	Callid string `json:"callid"`
}

// SendVoice 执行发送语音的逻辑
//
// 此接口整合了语音验证码和语音通知，使用时根据相应的参数构造请求体即可。
func (c *QcloudSMS) SendVoice(v VoiceReq) (bool, error) {
	var api string
	// 根据Prompttype类型验证是验证码还是普通通知，构造不同的请求URL
	if v.Prompttype == PROMPTVOICETYPE {
		api = PROMPTVOICE
	} else {
		api = SENDVOICE
	}

	c = c.NewSig(v.Tel.Mobile).NewURL(api)

	v.Sig = c.Sig
	v.Time = c.ReqTime

	resp, err := c.NewRequest(v)
	if err != nil {
		return false, err
	}

	var res VoiceResult
	json.Unmarshal([]byte(resp), &res)

	if res.Result == SUCCESS {
		return true, errors.New("发送成功")
	}

	return false, errors.New(res.Errmsg)
}
