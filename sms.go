// 短信API

package qcloudsms

import (
	"encoding/json"
	"errors"
	"strings"
)

/*
短信请求结构体

将单发短信和模板短信统一到了一个结构体里，构造请求时为对应请求的必传字段赋值即可

单发 https://cloud.tencent.com/document/product/382/5808

模板 https://cloud.tencent.com/document/product/382/5976
*/
type SMSSingleReq struct {
	Tel    SMSTel   `json:"tel"`
	Type   int      `json:"type,omitempty"`
	Sign   string   `json:"sign,omitempty"`
	TplID  int      `json:"tpl_id,omitempty"`
	Params []string `json:"params,omitempty"`
	Msg    string   `json:"msg,omitempty"`
	Sig    string   `json:"sig"`
	Time   int64    `json:"time"`
	Extend string   `json:"extend"`
	Ext    string   `json:"ext"`
}

// 国家码，手机号
type SMSTel struct {
	Nationcode string `json:"nationcode"`
	Mobile     string `json:"mobile"`
}

// 发送短信返回结构
type SMSResult struct {
	Result int    `json:"result"`
	Errmsg string `json:"errmsg"`
	Ext    string `json:"ext"`
	Sid    string `json:"sid,omitempty"`
	Fee    int    `json:"fee,omitempty"`
}

// 发送短信
func (c *QcloudSMS) SendSMSSingle(ss SMSSingleReq) (bool, error) {
	c = c.NewSig(ss.Tel.Mobile).NewUrl(SENDSMS)

	ss.Time = c.ReqTime
	ss.Sig = c.Sig

	resp, err := c.NewRequest(ss)
	if err != nil {
		return false, err
	}

	var res SMSResult
	json.Unmarshal([]byte(resp), &res)

	if res.Result == SUCCESS {
		return true, errors.New("发送成功")
	}

	return false, errors.New(res.Errmsg)
}

/*
群发短信请求结构体

将普通短信和模板短信统一到了一个结构体里，构造请求时为对应请求的必传字段赋值即可

普通 https://cloud.tencent.com/document/product/382/5806

模板 https://cloud.tencent.com/document/product/382/5977
*/
type SMSMultiReq struct {
	Tel    []SMSTel `json:"tel"`
	Type   int      `json:"type,omitempty"`
	Sign   string   `json:"sign,omitempty"`
	TplID  int      `json:"tpl_id,omitempty"`
	Params []string `json:"params,omitempty"`
	Msg    string   `json:"msg,omitempty"`
	Sig    string   `json:"sig"`
	Time   int64    `json:"time"`
	Extend string   `json:"extend"`
	Ext    string   `json:"ext"`
}

// 群发短信返回结构
type SMSMultiResult struct {
	Result int    `json:"result"`
	Errmsg string `json:"errmsg"`
	Ext    string `json:"ext"`
	Detail []struct {
		Result     int    `json:"result"`
		Errmsg     string `json:"errmsg"`
		Mobile     string `json:"mobile"`
		Nationcode string `json:"nationcode"`
		Sid        string `json:"sid,omitempty"`
		Fee        int    `json:"fee,omitempty"`
	} `json:"detail"`
}

// 群发短信
func (c *QcloudSMS) SendSMSMulti(sms SMSMultiReq) (bool, error) {
	var sigMobile []string

	if len(sms.Tel) > MULTISMSMAX {
		return false, ErrMultiCount
	}

	for _, m := range sms.Tel {
		sigMobile = append(sigMobile, m.Mobile)
	}

	mobileStr := strings.Join(sigMobile, ",")
	c = c.NewSig(mobileStr).NewUrl(MULTISMS)

	sms.Time = c.ReqTime
	sms.Sig = c.Sig

	resp, err := c.NewRequest(sms)
	if err != nil {
		return false, err
	}

	var res SMSMultiResult
	json.Unmarshal([]byte(resp), &res)

	if res.Result == SUCCESS {
		return true, errors.New("发送成功")
	}

	return false, errors.New(res.Errmsg)
}

// 拉取单个手机短信状态请求结构
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

// 拉取单个手机短信状态的返回结构
type StatusMobileResult struct {
	Result int               `json:"result"`
	Errmsg string            `json:"errmsg"`
	Count  int               `json:"count"`
	Data   []SMSStatusResult `json:"data"`
}

// 短信回复结构
type StatusReplyResult struct {
	Result int              `json:"result"`
	Errmsg string           `json:"errmsg"`
	Count  int              `json:"count"`
	Data   []SMSReplyResult `json:"data"`
}

// 单个手机短信状态结构
type SMSStatusResult struct {
	UserReceiveTime string `json:"user_receive_time"`
	Nationcode      string `json:"nationcode"`
	Mobile          string `json:"mobile"`
	ReportStatus    string `json:"report_status"`
	Errmsg          string `json:"errmsg"`
	Description     string `json:"description"`
	Sid             string `json:"sid"`
}

// 短信回复列表结构
type SMSReplyResult struct {
	Nationcode string `json:"nationcode"`
	Mobile     string `json:"mobile"`
	Text       string `json:"text"`
	Sign       string `json:"sign"`
	Time       int64  `json:"time"`
	Extend     string `json:"extend"`
}

/*
 拉取单个手机短信状态（下发状态）

 https://cloud.tencent.com/document/product/382/5811
*/
func (c *QcloudSMS) GetStatusForMobile(smr StatusMobileReq) (StatusMobileResult, error) {
	c = c.NewSig("").NewUrl(MOBILESTATUS)

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

/*
 拉取单个手机短信状态（短信回复）

 https://cloud.tencent.com/document/product/382/5811
*/
func (c *QcloudSMS) GetReplyForMobile(smr StatusMobileReq) (StatusReplyResult, error) {
	c = c.NewSig("").NewUrl(MOBILESTATUS)

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

// 拉取短信状态请求结构体
type PullStatusReq struct {
	Sig  string `json:"sig"`
	Time int64  `json:"time"`
	// 0 1分别代表 短信下发状态，短信回复
	Type int `json:"type"`
	// 最大条数 最多100
	Max int `json:"max"`
}

type PullStatusResult struct {
	Result int    `json:"result"`
	Errmsg string `json:"errmsg"`
	// result为0时有效，返回的信息条数
	Count int               `json:"count"`
	Data  []SMSStatusResult `json:"data"`
}

type PullReplyResult struct {
	Result int    `json:"result"`
	Errmsg string `json:"errmsg"`
	// result为0时有效，返回的信息条数
	Count int              `json:"count"`
	Data  []SMSReplyResult `json:"data"`
}

// 拉取短信状态（下发状态,短信回复）
// 已拉取过的数据将不会再返回
//
// https://cloud.tencent.com/document/product/382/5810
func (c *QcloudSMS) GetStatusMQ(psr PullStatusReq) (StatusMobileResult, error) {
	c = c.NewSig("").NewUrl(PULLSTATUS)

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
