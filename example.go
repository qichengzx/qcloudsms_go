package qcloudsms

import (
	"time"
)

var (
	appid  string = "yourappid"
	appkey string = "yourappkey"
	sign   string = "yoursign"
)

func ExampleNewClient() {
	opt := NewOptions(appid, appkey, sign)
	// 可以为 options 指定debug
	opt.Debug = true
	NewClient(opt)
}

func ExampleQcloudSMS_SendVoice() {
	opt := NewOptions(appid, appkey, sign)

	var client = NewClient(opt)
	//也可以在生成Client实例后设定 debug 模式
	client.SetDebug(true)

	var vr = VoiceReq{
		Promptfile: "您的验证码为：123。该验证码10分钟内有效。",
		Playtimes:  1,
		Prompttype: 2,
	}
	vr.Tel.Mobile = "86"
	vr.Tel.Mobile = "yourmobile"
	client.SendVoice(vr)
}

func ExampleQcloudSMS_NewSign() {
	opt := NewOptions(appid, appkey, sign)
	opt.Debug = true

	var client = NewClient(opt)

	var ns = SignReq{
		Remark:        "remark of sign",
		International: 0,
		Text:          "sign1",
	}

	client.NewSign(ns)
}

func ExampleQcloudSMS_ModTemplate() {
	opt := NewOptions(appid, appkey, sign)
	opt.Debug = true

	var client = NewClient(opt)

	var t = TemplateNew{
		TplID:  180101,
		Title:  "template title",
		Remark: "template remark",
		Text:   "here is {1} template",
		Type:   0,
	}

	client.ModTemplate(t)
}

func ExampleQcloudSMS_SendSMSSingle() {
	opt := NewOptions(appid, appkey, sign)
	opt.Debug = true

	var client = NewClient(opt)

	var sm = SMSSingleReq{
		Type: 0,
		Msg:  "短信内容",
		Tel:  SMSTel{Nationcode: "86", Mobile: "mobile"},
	}

	client.SendSMSSingle(sm)
}

func ExampleQcloudSMS_GetTemplateByPage() {
	opt := NewOptions(appid, appkey, sign)
	opt.Debug = true

	var client = NewClient(opt)

	client.GetTemplateByPage(0, 30)
}

func ExampleQcloudSMS_DelSign() {
	opt := NewOptions(appid, appkey, sign)
	opt.Debug = true

	var client = NewClient(opt)

	client.DelSign([]uint{171231, 171230})
}

//选择腾讯云上面配置的语音模板发送语音
func ExampleQcloudSMSVoiceTemplateSend() {
	opt := NewOptions(appid, appkey, sign)
	opt.Debug = true
	var client = NewClient(opt)
	tel := SMSTel{Nationcode: "86", Mobile: "18800000000"}
	var req = SMSVoiceTemplate{
		TplId:     123456,
		Playtimes: 2,
		Params:    []string{"11", "22", "33"},
		Tel:       tel,
		Time:      time.Now().Unix(),
		Ext:       "golang",
	}
	client.VoiceTemplateSend(req)
}
