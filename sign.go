package qcloudsms

import (
	"encoding/json"
)

// SignReq 添加/修改签名的请求结构
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
	SignID uint `json:"sign_id,omitempty"`
}

// SignDelGet 查询/删除签名的请求结构
type SignDelGet struct {
	Sig    string `json:"sig"`
	Time   int64  `json:"time"`
	SignID []uint `json:"sign_id"`
}

// SignResult 添加/修改/删除 的返回结构
type SignResult struct {
	Result uint   `json:"result"`
	Msg    string `json:"msg"`
	Data   struct {
		ID            uint   `json:"id"`
		International uint   `json:"international,omitempty"`
		Text          string `json:"text"`
		Status        uint   `json:"status"`
	} `json:"data,omitempty"`
}

// SignStatusResult 短信签名状态返回结构体
type SignStatusResult struct {
	Result uint   `json:"result"`
	Msg    string `json:"msg"`
	Count  uint   `json:"count"`
	Data   []struct {
		ID            uint   `json:"id"`
		Text          string `json:"text"`
		International uint   `json:"international,omitempty"`
		Status        uint   `json:"status"`
		Reply         string `json:"reply"`
		ApplyTime     string `json:"apply_time"`
	} `json:"data"`
}

// NewSign 添加签名
//
// https://cloud.tencent.com/document/product/382/6038
func (c *QcloudSMS) NewSign(s SignReq) (SignResult, error) {
	c = c.NewSig("").NewURL(ADDSIGN)

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

// ModSign 修改签名
// 参数是一个 SignReq 结构
//
// https://cloud.tencent.com/document/product/382/8650
func (c *QcloudSMS) ModSign(s SignReq) (SignResult, error) {
	c = c.NewSig("").NewURL(MODSIGN)

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

// GetSign 短信签名状态查询
// 参数是要删除签名的ID数组
//
// https://cloud.tencent.com/document/product/382/6040
func (c *QcloudSMS) GetSign(signid []uint) (SignStatusResult, error) {
	c = c.NewSig("").NewURL(GETSIGN)

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

// DelSign 删除短信签名
// 参数是要删除签名的ID数组
//
// https://cloud.tencent.com/document/product/382/6039
func (c *QcloudSMS) DelSign(signid []uint) (SignResult, error) {
	c = c.NewSig("").NewURL(DELSIGN)

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
