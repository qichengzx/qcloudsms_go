// 腾讯云短信平台
// https://cloud.tencent.com/document/product/382

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

type QcloudSMS struct {
	Random  string
	Sig     string
	Url     string
	ReqTime int64
	Options Options
	Logger  *log.Logger
}

type Options struct {
	// 腾讯云短信appid
	APPID string
	// 腾讯云短信appkey
	APPKEY string
	// 表示短信签名
	SIGN string

	RandomLen int
	UserAgent string

	Http struct {
		Timeout time.Duration
	}

	Debug bool
}

const (
	SDKName    = "qcloudsms-go-sdk"
	SDKVersion = "0.3.1"

	// API
	SVR string = "https://yun.tim.qq.com/v5/"

	// 腾讯云短信业务主URL
	TLSSMSSVR string = "tlssmssvr/"

	// 腾讯云语音URL
	VOICESVR string = "tlsvoicesvr/"

	// 短信业务URL附加内容
	TLSSMSSVRAfter string = "?sdkappid=%s&random=%s"

	// 发送短信
	SENDSMS string = "sendsms"

	// 群发
	MULTISMS string = "sendmultisms2"

	// 语音验证码
	SENDVOICE string = "sendvoice"

	// 语音通知
	PROMPTVOICE string = "sendvoiceprompt"

	// 添加模板
	ADDTEMPLATE string = "add_template"

	// 查询模板状态
	GETTEMPLATE string = "get_template"

	// 查询模板
	DELTEMPLATE string = "del_template"

	// 修改模板
	MODTEMPLATE string = "mod_template"

	// 添加签名
	ADDSIGN string = "add_sign"

	// 查询签名状态
	GETSIGN string = "get_sign"

	// 查询签名状态
	MODSIGN string = "mod_sign"

	// 查询签名状态
	DELSIGN string = "del_sign"

	// 拉取短信状态
	PULLSTATUS string = "pullstatus"

	// 拉取单个手机短信状态（下发状态，短信回复等）
	MOBILESTATUS string = "pullstatus4mobile"

	// 发送数据统计
	PULLSENDSTATUS string = "pullsendstatus"

	// 回执数据统计
	PULLCBSTATUS string = "pullcallbackstatus"

	// 请求成功的状态码
	SUCCESS int = 0

	// 短信类型，0=普通短信，1=营销短信
	MSGTYPE   int = 0
	MSGTYPEAD int = 1

	// 群发短信单批次最大手机号数量
	MULTISMSMAX int = 200

	// 语音类型，为2表示语音通知
	PROMPTVOICETYPE int = 2
)

var (
	ErrMultiCount = errors.New("单次提交不超过200个手机号")
	ErrRequest = errors.New("请求失败")
)

func NewOptions() *Options {
	opt := &Options{
		APPID:  "",
		APPKEY: "",
		SIGN:   "",

		RandomLen: 6,
		UserAgent: SDKName + "/" + SDKVersion,

		Debug: false,
	}

	opt.Http.Timeout = 10 * time.Second

	return opt
}

func NewClient(o *Options) *QcloudSMS {
	c := &QcloudSMS{}
	c.Options = *o

	c.NewRandom(c.Options.RandomLen)
	c.ReqTime = time.Now().Unix()

	c.Logger = log.New(os.Stderr, "["+SDKName+"]", log.LstdFlags)
	return c
}

func (c *QcloudSMS) SetAPPID(appid string) *QcloudSMS {
	c.Options.APPID = appid
	return c
}

func (c *QcloudSMS) SetAPPKEY(appkey string) *QcloudSMS {
	c.Options.APPKEY = appkey
	return c
}

func (c *QcloudSMS) SetSIGN(sign string) *QcloudSMS {
	c.Options.SIGN = sign
	return c
}

func (c *QcloudSMS) SetLogger(logger *log.Logger) *QcloudSMS {
	c.Logger = logger
	return c
}

func (c *QcloudSMS) SetDebug(debug bool) *QcloudSMS {
	if debug {
		c.Options.Debug = debug
	}

	return c
}

// 为 client 生成新的随机数
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

// 为请求构造sig
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

// 为请求构造URL
func (c *QcloudSMS) NewUrl(api string) *QcloudSMS {
	url := ""
	if api == SENDVOICE || api == PROMPTVOICE {
		url = VOICESVR
	} else {
		url = TLSSMSSVR
	}

	c.Url = SVR + url + api + fmt.Sprintf(TLSSMSSVRAfter, c.Options.APPID, c.Random)

	return c
}

// 发送请求
func (c *QcloudSMS) NewRequest(params interface{}) ([]byte, error) {
	j, _ := json.Marshal(params)

	req, err := http.NewRequest("POST", c.Url, bytes.NewBuffer([]byte(j)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", c.Options.UserAgent)

	httpClient := &http.Client{
		Timeout: c.Options.Http.Timeout,
	}
	resp, err := httpClient.Do(req)
	defer resp.Body.Close()

	if err != nil {
		return []byte{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return []byte{}, ErrRequest
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}

	if c.Options.Debug {
		c.Logger.Printf("Request Url : %s, Request Params : %s, Request Res : %s\n", c.Url, string(j), string(body))
	}

	return body, err
}
