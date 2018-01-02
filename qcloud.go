// Package qcloudsms 是针对 腾讯云短信平台 开发的一套 Go 语言 SDK
//
// 产品文档：https://cloud.tencent.com/document/product/382
package qcloudsms

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

// QcloudClient 用来构造请求，设置各项参数，和执行请求的接口
type QcloudClient interface {
	NewRandom(l int) *QcloudSMS
	NewSig(m string) *QcloudSMS
	NewUrl(api string) *QcloudSMS
	NewRequest(params interface{}) ([]byte, error)

	SetAPPID(appid string) *QcloudSMS
	SetAPPKEY(appkey string) *QcloudSMS
	SetSIGN(sign string) *QcloudSMS
	SetLogger(logger *log.Logger) *QcloudSMS
}

// QcloudSMS 是请求的结构
// 一次请求具体功能由 QcloudClient 接口实现
type QcloudSMS struct {
	Random  string
	Sig     string
	URL     string
	ReqTime int64
	Options Options
	Logger  *log.Logger
}

// Options 用来构造请求的参数结构
type Options struct {
	// 腾讯云短信appid
	APPID string
	// 腾讯云短信appkey
	APPKEY string
	// 表示短信签名
	SIGN string

	RandomLen int
	UserAgent string

	HTTP struct {
		Timeout time.Duration
	}

	Debug bool
}

const (
	//SDKName SDK名称，当前主要用于 log 中
	SDKName = "qcloudsms-go-sdk"
	// SDKVersion 版本
	SDKVersion = "0.3.2"

	// SVR 是腾讯云短信各请求结构的基本 URL
	SVR string = "https://yun.tim.qq.com/v5/"

	// TLSSMSSVR 腾讯云短信业务主URL
	TLSSMSSVR string = "tlssmssvr/"

	// VOICESVR 腾讯云语音URL
	VOICESVR string = "tlsvoicesvr/"

	// TLSSMSSVRAfter 短信业务URL附加内容
	TLSSMSSVRAfter string = "?sdkappid=%s&random=%s"

	// SENDSMS 发送短信
	SENDSMS string = "sendsms"

	// MULTISMS 群发短信
	MULTISMS string = "sendmultisms2"

	// SENDVOICE 发送语音验证码
	SENDVOICE string = "sendvoice"

	// PROMPTVOICE 发送语音通知
	PROMPTVOICE string = "sendvoiceprompt"

	// ADDTEMPLATE 添加模板
	ADDTEMPLATE string = "add_template"

	// GETTEMPLATE 查询模板状态
	GETTEMPLATE string = "get_template"

	// DELTEMPLATE 查询模板
	DELTEMPLATE string = "del_template"

	// MODTEMPLATE 修改模板
	MODTEMPLATE string = "mod_template"

	// ADDSIGN 添加签名
	ADDSIGN string = "add_sign"

	// GETSIGN 查询签名状态
	GETSIGN string = "get_sign"

	// MODSIGN 查询签名状态
	MODSIGN string = "mod_sign"

	// DELSIGN 查询签名状态
	DELSIGN string = "del_sign"

	// PULLSTATUS 拉取短信状态
	PULLSTATUS string = "pullstatus"

	// MOBILESTATUS 拉取单个手机短信状态（下发状态，短信回复等）
	MOBILESTATUS string = "pullstatus4mobile"

	// PULLSENDSTATUS 发送数据统计
	PULLSENDSTATUS string = "pullsendstatus"

	// PULLCBSTATUS 回执数据统计
	PULLCBSTATUS string = "pullcallbackstatus"

	// SUCCESS 请求成功的状态码
	SUCCESS uint = 0

	// MSGTYPE 普通短信类型
	MSGTYPE uint = 0
	// MSGTYPEAD 营销短信类型
	MSGTYPEAD uint = 1

	// MULTISMSMAX 群发短信单批次最大手机号数量
	MULTISMSMAX int = 200

	// PROMPTVOICETYPE 语音类型，为2表示语音通知
	PROMPTVOICETYPE uint = 2
)

var (
	//ErrMultiCount 群发号码数量错误
	ErrMultiCount = errors.New("单次提交不超过200个手机号")
	//ErrRequest 请求失败
	ErrRequest = errors.New("请求失败")
)

// NewOptions 返回一个新的 *Options
func NewOptions() *Options {
	opt := &Options{
		APPID:  "",
		APPKEY: "",
		SIGN:   "",

		RandomLen: 6,
		UserAgent: SDKName + "/" + SDKVersion,

		Debug: false,
	}

	opt.HTTP.Timeout = 10 * time.Second

	return opt
}

// NewClient 生成一个新的 client 实例
func NewClient(o *Options) *QcloudSMS {
	c := &QcloudSMS{}
	c.Options = *o

	c.NewRandom(c.Options.RandomLen)
	c.ReqTime = time.Now().Unix()

	c.Logger = log.New(os.Stderr, "["+SDKName+"]", log.LstdFlags)
	return c
}

// SetAPPID 为实例设置 APPID
func (c *QcloudSMS) SetAPPID(appid string) *QcloudSMS {
	c.Options.APPID = appid
	return c
}

// SetAPPKEY 为实例设置 APPKEY
func (c *QcloudSMS) SetAPPKEY(appkey string) *QcloudSMS {
	c.Options.APPKEY = appkey
	return c
}

// SetSIGN 为实例设置 SIGN
func (c *QcloudSMS) SetSIGN(sign string) *QcloudSMS {
	c.Options.SIGN = sign
	return c
}

// SetLogger 为实例设置 logger
func (c *QcloudSMS) SetLogger(logger *log.Logger) *QcloudSMS {
	c.Logger = logger
	return c
}

// SetDebug 为实例设置调试模式
func (c *QcloudSMS) SetDebug(debug bool) *QcloudSMS {
	if debug {
		c.Options.Debug = debug
	}

	return c
}

// NewRandom 为实例生成新的随机数
func (c *QcloudSMS) NewRandom(l int) *QcloudSMS {
	str := "0123456789"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	c.Random = string(result)

	return c
}

// NewSig 为实例生成 sig
func (c *QcloudSMS) NewSig(m string) *QcloudSMS {
	var t = strconv.FormatInt(c.ReqTime, 10)
	var sigContent = "appkey=" + c.Options.APPKEY + "&random=" + c.Random + "&time=" + t

	if len(m) > 0 {
		sigContent += "&mobile=" + m
	}
	h := sha256.New()
	h.Write([]byte(sigContent))

	c.Sig = fmt.Sprintf("%x", h.Sum(nil))

	return c
}

// NewURL 为实例设置 URL
func (c *QcloudSMS) NewURL(api string) *QcloudSMS {
	url := ""
	if api == SENDVOICE || api == PROMPTVOICE {
		url = VOICESVR
	} else {
		url = TLSSMSSVR
	}

	c.URL = SVR + url + api + fmt.Sprintf(TLSSMSSVRAfter, c.Options.APPID, c.Random)

	return c
}

// NewRequest 执行实例发送请求
func (c *QcloudSMS) NewRequest(params interface{}) ([]byte, error) {
	j, _ := json.Marshal(params)

	req, err := http.NewRequest("POST", c.URL, bytes.NewBuffer([]byte(j)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", c.Options.UserAgent)

	httpClient := &http.Client{
		Timeout: c.Options.HTTP.Timeout,
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []byte{}, ErrRequest
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}

	if c.Options.Debug {
		c.Logger.Printf("Request Url : %s, Request Params : %s, Request Res : %s\n", c.URL, string(j), string(body))
	}

	return body, err
}
