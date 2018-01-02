package qcloudsms

import (
	"time"
)

func ExampleNewClient() {
	opt := NewOptions()
	opt.APPID = "yourappid"
	opt.APPKEY = "yourappkey"
	opt.Debug = true
	opt.HTTP.Timeout = 10 * time.Second
	opt.SIGN = "yoursign"

	NewClient(opt)
}

func ExampleQcloudSMS_SendVoice() {
	opt := NewOptions()
	opt.APPID = "yourappid"
	opt.APPKEY = "yourappkey"
	opt.Debug = true
	opt.HTTP.Timeout = 10 * time.Second
	opt.SIGN = "yoursign"

	var client = NewClient(opt)

	var vr = VoiceReq{
		Promptfile: "您的验证码为：123。该验证码10分钟内有效。",
		Playtimes:  1,
		Prompttype: 2,
	}

	vr.Tel.Nationcode = "86"
	vr.Tel.Mobile = "yourmobile"

	client.SendVoice(vr)
}

func ExampleQcloudSMS_NewSign() {
	opt := NewOptions()
	opt.APPID = "yourappid"
	opt.APPKEY = "yourappkey"
	opt.Debug = true
	opt.HTTP.Timeout = 10 * time.Second
	opt.SIGN = "yoursign"

	var client = NewClient(opt)

	var ns = SignReq{
		Remark:        "remark of sign",
		International: 0,
		Text:          "sign1",
	}

	client.NewSign(ns)
}

func ExampleQcloudSMS_ModTemplate() {
	opt := NewOptions()
	opt.APPID = "yourappid"
	opt.APPKEY = "yourappkey"
	opt.Debug = true
	opt.HTTP.Timeout = 10 * time.Second
	opt.SIGN = "yoursign"

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
	opt := NewOptions()
	opt.APPID = "yourappid"
	opt.APPKEY = "yourappkey"
	opt.Debug = true
	opt.HTTP.Timeout = 10 * time.Second
	opt.SIGN = "yoursign"

	var client = NewClient(opt)

	var sm = SMSSingleReq{
		Type: 0,
		Msg:  "短信内容",
		Tel:  SMSTel{Nationcode: "86", Mobile: "mobile"},
	}

	client.SendSMSSingle(sm)
}

func ExampleQcloudSMS_GetTemplateByPage() {
	opt := NewOptions()
	opt.APPID = "yourappid"
	opt.APPKEY = "yourappkey"
	opt.Debug = true
	opt.HTTP.Timeout = 10 * time.Second
	opt.SIGN = "yoursign"

	var client = NewClient(opt)

	client.GetTemplateByPage(0, 30)
}

func ExampleQcloudSMS_DelSign() {
	opt := NewOptions()
	opt.APPID = "yourappid"
	opt.APPKEY = "yourappkey"
	opt.Debug = true
	opt.HTTP.Timeout = 10 * time.Second
	opt.SIGN = "yoursign"

	var client = NewClient(opt)

	client.DelSign([]uint{171231, 171230})
}
