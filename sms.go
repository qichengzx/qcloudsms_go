package qcloudsms

import (
	"encoding/json"
	"errors"
	"strings"
)

/*
SMSSingleReq 发送单条短信请求结构

将单发短信和模板短信统一到了一个结构体里，构造请求时为对应请求的必传字段赋值即可

单发 https://cloud.tencent.com/document/product/382/5808

模板 https://cloud.tencent.com/document/product/382/5976
*/
type SMSSingleReq struct {
	Tel    SMSTel   `json:"tel"`
	Type   int      `json:"type,omitempty"`
	Sign   string   `json:"sign,omitempty"`
	TplID  int      `json:"tpl_id,omitempty"`
	Params []string `json:"params"`
	Msg    string   `json:"msg,omitempty"`
	Sig    string   `json:"sig"`
	Time   int64    `json:"time"`
	Extend string   `json:"extend"`
	Ext    string   `json:"ext"`
}

// SMSTel 国家码，手机号
type SMSTel struct {
	Nationcode string `json:"nationcode"`
	Mobile     string `json:"mobile"`
}

// SMSResult 发送短信返回结构
type SMSResult struct {
	Result uint   `json:"result"`
	Errmsg string `json:"errmsg"`
	Ext    string `json:"ext"`
	Sid    string `json:"sid,omitempty"`
	Fee    uint   `json:"fee,omitempty"`
}

// SendSMSSingle 发送单条短信
func (c *QcloudSMS) SendSMSSingle(ss SMSSingleReq) (bool, error) {
	c = c.NewSig(ss.Tel.Mobile).NewURL(SENDSMS)

	ss.Time = c.ReqTime
	ss.Sig = c.Sig

	resp, err := c.NewRequest(ss)
	if err != nil {
		return false, err
	}

	var res SMSResult
	json.Unmarshal([]byte(resp), &res)

	if res.Result == SUCCESS {
		return true, nil
	}

	return false, errors.New(res.Errmsg)
}

/*
SMSMultiReq 群发短信请求结构

将普通短信和模板短信统一到了一个结构体里，构造请求时为对应请求的必传字段赋值即可

普通 https://cloud.tencent.com/document/product/382/5806

模板 https://cloud.tencent.com/document/product/382/5977
*/
type SMSMultiReq struct {
	Tel    []SMSTel `json:"tel"`
	Type   uint     `json:"type,omitempty"`
	Sign   string   `json:"sign,omitempty"`
	TplID  uint     `json:"tpl_id,omitempty"`
	Params []string `json:"params"`
	Msg    string   `json:"msg,omitempty"`
	Sig    string   `json:"sig"`
	Time   int64    `json:"time"`
	Extend string   `json:"extend"`
	Ext    string   `json:"ext"`
}

// SMSMultiResult 群发短信返回结构
type SMSMultiResult struct {
	Result uint   `json:"result"`
	Errmsg string `json:"errmsg"`
	Ext    string `json:"ext"`
	Detail []struct {
		Result     uint   `json:"result"`
		Errmsg     string `json:"errmsg"`
		Mobile     string `json:"mobile"`
		Nationcode string `json:"nationcode"`
		Sid        string `json:"sid,omitempty"`
		Fee        uint   `json:"fee,omitempty"`
	} `json:"detail"`
}

// SendSMSMulti 群发短信
func (c *QcloudSMS) SendSMSMulti(sms SMSMultiReq) (bool, error) {
	var sigMobile []string

	if len(sms.Tel) > MULTISMSMAX {
		return false, ErrMultiCount
	}

	for _, m := range sms.Tel {
		sigMobile = append(sigMobile, m.Mobile)
	}

	mobileStr := strings.Join(sigMobile, ",")
	c = c.NewSig(mobileStr).NewURL(MULTISMS)

	sms.Time = c.ReqTime
	sms.Sig = c.Sig

	resp, err := c.NewRequest(sms)
	if err != nil {
		return false, err
	}

	var res SMSMultiResult
	json.Unmarshal([]byte(resp), &res)

	if res.Result == SUCCESS {
		return true, nil
	}

	return false, errors.New(res.Errmsg)
}

// StatusMobileReq 拉取单个手机短信状态请求结构
type StatusMobileReq struct {
	Sig  string `json:"sig"`
	Time int64  `json:"time"`
	// 0 1分别代表 短信下发状态，短信回复
	Type int `json:"type"`
	// 最大条数 最多100
	Max        int    `json:"max"`
	BeginTime  int64  `json:"begin_time"`
	EndTime    int64  `json:"end_time"`
	Nationcode string `json:"nationcode"`
	Mobile     string `json:"mobile"`
}

// StatusMobileResult 拉取单个手机短信状态的返回结构
type StatusMobileResult struct {
	Result int               `json:"result"`
	Errmsg string            `json:"errmsg"`
	Count  int               `json:"count"`
	Data   []SMSStatusResult `json:"data"`
}

// StatusReplyResult 短信回复结构
type StatusReplyResult struct {
	Result int              `json:"result"`
	Errmsg string           `json:"errmsg"`
	Count  int              `json:"count"`
	Data   []SMSReplyResult `json:"data"`
}

// SMSStatusResult 单个手机短信状态结构
type SMSStatusResult struct {
	UserReceiveTime string `json:"user_receive_time"`
	Nationcode      string `json:"nationcode"`
	Mobile          string `json:"mobile"`
	ReportStatus    string `json:"report_status"`
	Errmsg          string `json:"errmsg"`
	Description     string `json:"description"`
	Sid             string `json:"sid"`
}

// SMSReplyResult 短信回复列表结构
type SMSReplyResult struct {
	Nationcode string `json:"nationcode"`
	Mobile     string `json:"mobile"`
	Text       string `json:"text"`
	Sign       string `json:"sign"`
	Time       int64  `json:"time"`
	Extend     string `json:"extend"`
}

// GetStatusForMobile 拉取单个手机短信状态（下发状态）
//
// https://cloud.tencent.com/document/product/382/5811
func (c *QcloudSMS) GetStatusForMobile(smr StatusMobileReq) (StatusMobileResult, error) {
	c = c.NewSig("").NewURL(MOBILESTATUS)

	smr.Time = c.ReqTime
	smr.Sig = c.Sig

	var res StatusMobileResult
	resp, err := c.NewRequest(smr)
	if err != nil {
		return res, err
	}

	json.Unmarshal([]byte(resp), &res)

	return res, nil
}

// GetReplyForMobile 拉取单个手机短信状态（短信回复）
//
// https://cloud.tencent.com/document/product/382/5811
func (c *QcloudSMS) GetReplyForMobile(smr StatusMobileReq) (StatusReplyResult, error) {
	c = c.NewSig("").NewURL(MOBILESTATUS)

	smr.Time = c.ReqTime
	smr.Sig = c.Sig

	var res StatusReplyResult
	resp, err := c.NewRequest(smr)
	if err != nil {
		return res, err
	}

	json.Unmarshal([]byte(resp), &res)

	return res, nil
}

// PullStatusReq 拉取短信状态请求结构
type PullStatusReq struct {
	Sig  string `json:"sig"`
	Time int64  `json:"time"`
	// 0 1分别代表 短信下发状态，短信回复
	Type int `json:"type"`
	// 最大条数 最多100
	Max int `json:"max"`
}

// PullStatusResult 拉取下发状态返回数据结构
type PullStatusResult struct {
	Result int    `json:"result"`
	Errmsg string `json:"errmsg"`
	// result为0时有效，返回的信息条数
	Count int               `json:"count"`
	Data  []SMSStatusResult `json:"data"`
}

// PullReplyResult 拉取短信回复返回数据结构
type PullReplyResult struct {
	Result int    `json:"result"`
	Errmsg string `json:"errmsg"`
	// result为0时有效，返回的信息条数
	Count int              `json:"count"`
	Data  []SMSReplyResult `json:"data"`
}

// GetStatusMQ 拉取短信状态（下发状态,短信回复）
// 已拉取过的数据将不会再返回
//
// https://cloud.tencent.com/document/product/382/5810
func (c *QcloudSMS) GetStatusMQ(psr PullStatusReq) (StatusMobileResult, error) {
	c = c.NewSig("").NewURL(PULLSTATUS)

	psr.Time = c.ReqTime
	psr.Sig = c.Sig

	resp, err := c.NewRequest(psr)
	if err != nil {
		return StatusMobileResult{}, err
	}

	var res StatusMobileResult
	json.Unmarshal([]byte(resp), &res)

	return res, nil
}
