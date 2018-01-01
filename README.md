腾讯云短信 Go SDK
===

## Overview

> 此 SDK 为非官方版本，命名和结构上与官方版本有一些区别。

> 海外短信和国内短信使用同一接口，只需替换相应的国家码与手机号码，每次请求群发接口手机号码需全部为国内或者海外手机号码。

> 语音通知目前支持语音验证码以及语音通知功能。

## Features

##### 短信
- [x] 单发短信
- [x] 指定模板单发短信
- [x] 群发短信
- [x] 群发模板短信
- [ ] 短信下发状态通知
- [ ] 短信回复
- [x] 拉取短信状态
- [x] 拉取单个手机短信状态

##### 语音
- [x] 发送语音验证码
- [x] 发送语音通知
- [ ] 语音验证码状态通知
- [ ] 语音通知状态通知
- [ ] 语音通知按键通知
- [ ] 语音送达失败原因推送

##### 模板
- [x] 添加模板
- [x] 修改模板
- [x] 删除模板
- [x] 模板状态查询

##### 签名
- [x] 添加签名
- [x] 修改签名
- [x] 删除签名
- [x] 短信签名状态查询

##### 统计
- [x] 发送数据统计
- [x] 回执数据统计

# Getting Start

## 准备

- [ ] 申请APPID以及APPKey

> 在开始使用之前，需要先获取APPID和APPkey，如尚未申请，请到https://console.qcloud.com/sms/smslist 中添加应用，应用添加成功后您将获得APPID以及APPKey，注意APPID是以14xxxxx开头。

- [ ] 申请签名

> 下发短信必须携带签名，在相应服务模块 *短信内容配置*  中进行申请。

- [ ] 申请模板

> 下发短信内容必须经过审核，在相应服务 *短信内容配置* 中进行申请。

完成以上三项便可开始代码开发。

## 安装

```
go get github.com/qichengzx/qcloudsms_go
```

## 用法

```Go

import "github.com/qichengzx/qcloudsms_go"

opt := NewOptions()
opt.APPID = "yourappid"
opt.APPKEY = "yourappkey"
opt.Debug = true
opt.Http.Timeout = 10 * time.Second
opt.SIGN = "yoursign"

var client = NewClient(opt)

```
