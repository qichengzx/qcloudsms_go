// 统计
// 发送数据统计 & 回执数据统计

package qcloudsms

import (
	"encoding/json"
)

// 查询回执数据请求结构
type StatusReq struct {
	Sig       string `json:"sig"`
	Time      int64  `json:"time"`
	BeginDate uint32 `json:"begin_date"`
	EndDate   uint32 `json:"end_date"`
}

// 回执数据结构
type StatusResult struct {
	Result int    `json:"result"`
	Errmsg string `json:"errmsg"`
	Data   struct {
		// 短信提交成功量
		Success int `json:"success"`
		// 短信回执量
		Status int `json:"status"`
		// 短信回执成功量
		StatusSuccess int `json:"status_success"`
		// 短信回执失败量
		StatusFail int `json:"status_fail"`
		// 运营商内部错误
		StatusFail0 int `json:"status_fail_0"`
		// 号码无效或空号
		StatusFail1 int `json:"status_fail_1"`
		// 停机、关机等
		StatusFail2 int `json:"status_fail_2"`
		// 黑名单
		StatusFail3 int `json:"status_fail_3"`
		// 运营商频率限制
		StatusFail4 int `json:"status_fail_4"`
	} `json:"data"`
}

// 发送数据统计请求结构体
type SendStatusReq struct {
	Sig       string `json:"sig"`
	Time      int64  `json:"time"`
	BeginDate uint32 `json:"begin_date"`
	EndDate   uint32 `json:"end_date"`
}

// 数据统计返回结构
type SendStatusResult struct {
	Result int    `json:"result"`
	Errmsg string `json:"errmsg"`
	Data   struct {
		Request    int `json:"request"`
		Success    int `json:"success"`
		BillNumber int `json:"bill_number"`
	} `json:"data"`
}

// 查询回执数据
//
// https://cloud.tencent.com/document/product/382/7756
func (c *QcloudSMS) GetStatus(begin, end uint32) (StatusResult, error) {
	c = c.NewSig("").NewUrl(PULLCBSTATUS)

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

// 拉取统计数据
//
// https://cloud.tencent.com/document/product/382/7755
func (c *QcloudSMS) GetSendStatus(begin, end uint32) (SendStatusResult, error) {
	c = c.NewSig("").NewUrl(PULLSENDSTATUS)
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
