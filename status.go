package qcloudsms

import (
	"encoding/json"
)

// StatusReq 查询回执数据请求结构
type StatusReq struct {
	Sig       string `json:"sig"`
	Time      int64  `json:"time"`
	BeginDate uint32 `json:"begin_date"`
	EndDate   uint32 `json:"end_date"`
}

// StatusResult 回执数据结构
type StatusResult struct {
	Result uint   `json:"result"`
	Errmsg string `json:"errmsg"`
	Data   struct {
		// 短信提交成功量
		Success uint `json:"success"`
		// 短信回执量
		Status uint `json:"status"`
		// 短信回执成功量
		StatusSuccess uint `json:"status_success"`
		// 短信回执失败量
		StatusFail uint `json:"status_fail"`
		// 运营商内部错误
		StatusFail0 uint `json:"status_fail_0"`
		// 号码无效或空号
		StatusFail1 uint `json:"status_fail_1"`
		// 停机、关机等
		StatusFail2 uint `json:"status_fail_2"`
		// 黑名单
		StatusFail3 uint `json:"status_fail_3"`
		// 运营商频率限制
		StatusFail4 uint `json:"status_fail_4"`
	} `json:"data"`
}

// SendStatusReq 发送数据统计请求结构
type SendStatusReq struct {
	Sig       string `json:"sig"`
	Time      int64  `json:"time"`
	BeginDate uint32 `json:"begin_date"`
	EndDate   uint32 `json:"end_date"`
}

// SendStatusResult 发送数据统计返回结构
type SendStatusResult struct {
	Result uint   `json:"result"`
	Errmsg string `json:"errmsg"`
	Data   struct {
		Request    uint `json:"request"`
		Success    uint `json:"success"`
		BillNumber uint `json:"bill_number"`
	} `json:"data"`
}

// GetStatus 回执数据统计
//
// https://cloud.tencent.com/document/product/382/7756
func (c *QcloudSMS) GetStatus(begin, end uint32) (StatusResult, error) {
	c = c.NewSig("").NewURL(PULLCBSTATUS)

	var cbs = StatusReq{
		Sig:       c.Sig,
		Time:      c.ReqTime,
		BeginDate: begin,
		EndDate:   end,
	}

	var res StatusResult
	resp, err := c.NewRequest(cbs)
	if err != nil {
		return res, err
	}

	json.Unmarshal([]byte(resp), &res)

	return res, nil
}

// GetSendStatus 发送数据统计
//
// https://cloud.tencent.com/document/product/382/7755
func (c *QcloudSMS) GetSendStatus(begin, end uint32) (SendStatusResult, error) {
	c = c.NewSig("").NewURL(PULLSENDSTATUS)
	var cs = SendStatusReq{
		Sig:       c.Sig,
		Time:      c.ReqTime,
		BeginDate: begin,
		EndDate:   end,
	}

	var res SendStatusResult
	resp, err := c.NewRequest(cs)
	if err != nil {
		return res, err
	}

	json.Unmarshal([]byte(resp), &res)

	return res, nil
}
