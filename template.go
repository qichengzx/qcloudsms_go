// 模板

package qcloudsms

import (
	"encoding/json"
)

// 查询模板状态请求结构
type TemplateGetReq struct {
	Sig     string `json:"sig"`
	Time    int64  `json:"time"`
	TplID   []uint `json:"tpl_id,omitempty"`
	TplPage struct {
		Offset uint `json:"offset"`
		Max    uint `json:"max"`
	} `json:"tpl_page,omitempty"`
}

// 查询模板状态返回结构
type TemplateGetResult struct {
	Result uint       `json:"result"`
	Msg    string     `json:"msg"`
	Total  uint       `json:"total"`
	Count  uint       `json:"count"`
	Data   []Template `json:"data"`
}

// 模板结构
type Template struct {
	Id            uint   `json:"id"`
	Text          string `json:"text"`
	Status        uint   `json:"status"`
	Reply         string `json:"reply"`
	Type          uint   `json:"type"`
	International uint   `json:"international"`
	ApplyTime     string `json:"apply_time"`
}

// 新增，修改模板的请求结构
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

// 删除模板请求结构
type TemplateDelReq struct {
	Sig   string `json:"sig"`
	Time  int64  `json:"time"`
	TplID []uint `json:"tpl_id"`
}

// 新增，修改，删除模板返回结构
// 以上几种操作共用此结构，唯一区别在于返回值中的data可能不同
type TemplateResult struct {
	Result uint     `json:"result"`
	Msg    string   `json:"msg"`
	Data   Template `json:"data,omitempty"`
}

// 查询模板状态
// 将需要查询的模板ID以数组方式传入
//
// https://cloud.tencent.com/document/product/382/5819
func (c *QcloudSMS) GetTemplateByID(id []uint) (TemplateGetResult, error) {
	c = c.NewSig("").NewUrl(GETTEMPLATE)

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

func (c *QcloudSMS) GetTemplateByPage(offset, max uint) (TemplateGetResult, error) {
	c = c.NewSig("").NewUrl(GETTEMPLATE)

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

// 新建模板
//
// https://cloud.tencent.com/document/product/382/5817
func (c *QcloudSMS) NewTemplate(t TemplateNew) (TemplateResult, error) {
	c = c.NewSig("").NewUrl(ADDTEMPLATE)

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

// 修改模板
//
// https://cloud.tencent.com/document/product/382/8649
func (c *QcloudSMS) ModTemplate(t TemplateNew) (TemplateResult, error) {
	c = c.NewSig("").NewUrl(MODTEMPLATE)

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

// 删除模板
// TplID是一个uint数组
//
// https://cloud.tencent.com/document/product/382/5818
func (c *QcloudSMS) DelTemplate(id []uint) (TemplateResult, error) {
	c = c.NewSig("").NewUrl(DELTEMPLATE)

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
