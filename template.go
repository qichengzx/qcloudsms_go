package qcloudsms

import (
	"encoding/json"
)

// TemplateGetReq 查询模板状态请求结构
type TemplateGetReq struct {
	Sig     string `json:"sig"`
	Time    int64  `json:"time"`
	TplID   []uint `json:"tpl_id,omitempty"`
	TplPage struct {
		Offset uint `json:"offset"`
		Max    uint `json:"max"`
	} `json:"tpl_page,omitempty"`
}

// TemplateGetResult 查询模板状态返回结构
type TemplateGetResult struct {
	Result uint       `json:"result"`
	Msg    string     `json:"msg"`
	Total  uint       `json:"total"`
	Count  uint       `json:"count"`
	Data   []Template `json:"data"`
}

// Template 模板结构
type Template struct {
	ID            uint   `json:"id"`
	Text          string `json:"text"`
	Status        uint   `json:"status"`
	Reply         string `json:"reply"`
	Type          uint   `json:"type"`
	International uint   `json:"international"`
	ApplyTime     string `json:"apply_time"`
}

// TemplateNew 新增，修改模板的请求结构
// 其中，新增时TplID为空，修改时传入指定要修改的模板ID
type TemplateNew struct {
	Sig           string `json:"sig"`
	Time          int64  `json:"time"`
	Title         string `json:"title"`
	Remark        string `json:"remark"`
	International uint   `json:"international"`
	Text          string `json:"text"`
	TplID         uint   `json:"tpl_id,omitempty"`
	Type          uint   `json:"type"`
}

// TemplateDelReq 删除模板请求结构
type TemplateDelReq struct {
	Sig   string `json:"sig"`
	Time  int64  `json:"time"`
	TplID []uint `json:"tpl_id"`
}

// TemplateResult 新增，修改，删除模板返回结构
// 以上几种操作共用此结构，唯一区别在于返回值中的data可能不同
type TemplateResult struct {
	Result uint     `json:"result"`
	Msg    string   `json:"msg"`
	Data   Template `json:"data,omitempty"`
}

// GetTemplateByID 根据模板ID查询模板状态
// 将需要查询的模板ID以数组方式传入
//
// https://cloud.tencent.com/document/product/382/5819
func (c *QcloudSMS) GetTemplateByID(id []uint) (TemplateGetResult, error) {
	c = c.NewSig("").NewURL(GETTEMPLATE)

	var t = TemplateGetReq{
		Sig:   c.Sig,
		Time:  c.ReqTime,
		TplID: id,
	}

	var res TemplateGetResult
	resp, err := c.NewRequest(t)
	if err != nil {
		return res, err
	}

	json.Unmarshal([]byte(resp), &res)

	return res, nil
}

// GetTemplateByPage 用于批量获取模板数据
// 参数为偏移量，拉取条数
func (c *QcloudSMS) GetTemplateByPage(offset, max uint) (TemplateGetResult, error) {
	c = c.NewSig("").NewURL(GETTEMPLATE)

	var t = TemplateGetReq{
		Sig:  c.Sig,
		Time: c.ReqTime,
	}

	t.TplPage.Offset = offset
	t.TplPage.Max = max

	var res TemplateGetResult
	resp, err := c.NewRequest(t)
	if err != nil {
		return res, err
	}

	json.Unmarshal([]byte(resp), &res)

	return res, nil
}

// NewTemplate 新建模板
// 参数是一个 TemplateNew 结构
//
// https://cloud.tencent.com/document/product/382/5817
func (c *QcloudSMS) NewTemplate(t TemplateNew) (TemplateResult, error) {
	c = c.NewSig("").NewURL(ADDTEMPLATE)

	t.Time = c.ReqTime
	t.Sig = c.Sig

	var res TemplateResult
	resp, err := c.NewRequest(t)
	if err != nil {
		return res, err
	}

	json.Unmarshal([]byte(resp), &res)

	return res, nil
}

// ModTemplate 修改模板
// 参数是一个 TemplateNew 结构
//
// https://cloud.tencent.com/document/product/382/8649
func (c *QcloudSMS) ModTemplate(t TemplateNew) (TemplateResult, error) {
	c = c.NewSig("").NewURL(MODTEMPLATE)

	t.Time = c.ReqTime
	t.Sig = c.Sig

	var res TemplateResult
	resp, err := c.NewRequest(t)
	if err != nil {
		return res, err
	}

	json.Unmarshal([]byte(resp), &res)

	return res, nil
}

// DelTemplate 删除模板
// 参数是一个uint数组，表示模板ID
//
// https://cloud.tencent.com/document/product/382/5818
func (c *QcloudSMS) DelTemplate(id []uint) (TemplateResult, error) {
	c = c.NewSig("").NewURL(DELTEMPLATE)

	var t = TemplateDelReq{
		Time:  c.ReqTime,
		Sig:   c.Sig,
		TplID: id,
	}

	var res TemplateResult
	resp, err := c.NewRequest(t)
	if err != nil {
		return res, err
	}

	json.Unmarshal([]byte(resp), &res)

	return res, nil
}
