// 签名

package qcloudsms

import (
	"encoding/json"
)

// 添加/修改签名的请求结构
type SignReq struct {
	Sig    string `json:"sig"`
	Time   int64  `json:"time"`
	Remark string `json:"remark"`
	// 是否为国际短信签名，1=国际，0=国内
	International int `json:"international"`
	// 签名内容，不带"【】"
	Text string `json:"text"`
	// 签名内容相关的证件截图base64格式字符串，非必传参数
	Pic string `json:"pic,omitempty"`
	// 要修改的签名ID
	SignID int `json:"sign_id,omitempty"`
}

// 查询/删除签名的请求结构
type SignDelGet struct {
	Sig    string `json:"sig"`
	Time   int64  `json:"time"`
	SignID []int  `json:"sign_id"`
}

// 添加/修改/删除 的返回结构
type SignResult struct {
	Result int    `json:"result"`
	Msg    string `json:"msg"`
	Data   struct {
		ID            int    `json:"id"`
		International int    `json:"international,omitempty"`
		Text          string `json:"text"`
		Status        int    `json:"status"`
	} `json:"data,omitempty"`
}

// 短信签名状态返回结构体
type SignStatusResult struct {
	Result int    `json:"result"`
	Msg    string `json:"msg"`
	Count  int    `json:"count"`
	Data   []struct {
		ID            int    `json:"id"`
		Text          string `json:"text"`
		International int    `json:"international,omitempty"`
		Status        int    `json:"status"`
		Reply         string `json:"reply"`
		ApplyTime     string `json:"apply_time"`
	} `json:"data"`
}

// 添加签名
//
// https://cloud.tencent.com/document/product/382/6038
func (c *QcloudSMS) NewSign(s SignReq) (SignResult, error) {
	c = c.NewSig("").NewUrl(ADDSIGN)

	s.Time = c.ReqTime
	s.Sig = c.Sig

	var res SignResult
	resp, err := c.NewRequest(s)
	if err != nil {
		return res, err
	}

	json.Unmarshal([]byte(resp), &res)

	return res, nil
}

// 修改签名
//
// https://cloud.tencent.com/document/product/382/8650
func (c *QcloudSMS) ModSign(s SignReq) (SignResult, error) {
	c = c.NewSig("").NewUrl(MODSIGN)

	s.Time = c.ReqTime
	s.Sig = c.Sig

	var res SignResult
	resp, err := c.NewRequest(s)
	if err != nil {
		return res, err
	}

	json.Unmarshal([]byte(resp), &res)

	return res, nil
}

// 短信签名状态查询
//
// https://cloud.tencent.com/document/product/382/6040
func (c *QcloudSMS) GetSign(signid []int) (SignStatusResult, error) {
	c = c.NewSig("").NewUrl(GETSIGN)

	var s = SignDelGet{
		SignID: signid,
		Time:   c.ReqTime,
		Sig:    c.Sig,
	}

	var res SignStatusResult
	resp, err := c.NewRequest(s)
	if err != nil {
		return res, err
	}

	json.Unmarshal([]byte(resp), &res)

	return res, nil
}

// 删除短信签名
//
// https://cloud.tencent.com/document/product/382/6039
func (c *QcloudSMS) DelSign(signid []int) (SignResult, error) {
	c = c.NewSig("").NewUrl(DELSIGN)

	var s = SignDelGet{
		Time:   c.ReqTime,
		Sig:    c.Sig,
		SignID: signid,
	}

	var res SignResult
	resp, err := c.NewRequest(s)
	if err != nil {
		return res, err
	}

	json.Unmarshal([]byte(resp), &res)

	return res, nil
}
