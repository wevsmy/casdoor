package sms

import (
	"fmt"

	sender "github.com/casdoor/go-sms-sender"
)

const (
	Aliyun       = "Aliyun SMS"
	TencentCloud = "Tencent Cloud SMS"
	VolcEngine   = "Volc Engine SMS"
	Huyi         = "Huyi SMS"
	HuaweiCloud  = "Huawei Cloud SMS"
	Twilio       = "Twilio SMS"
	SmsBao       = "SmsBao SMS"
	MockSms      = "Mock SMS"
	SUBMAIL      = "SUBMAIL SMS"
	Sinowel      = "Sinowel SMS"
)

type SmsClient interface {
	SendMessage(param map[string]string, targetPhoneNumber ...string) error
}

func NewSmsClient(provider string, accessId string, accessKey string, sign string, template string, other ...string) (SmsClient, error) {
	switch provider {
	case Aliyun:
		return sender.GetAliyunClient(accessId, accessKey, sign, template)
	case TencentCloud:
		return sender.GetTencentClient(accessId, accessKey, sign, template, other)
	case VolcEngine:
		return sender.GetVolcClient(accessId, accessKey, sign, template, other)
	case Huyi:
		return sender.GetHuyiClient(accessId, accessKey, template)
	case HuaweiCloud:
		return sender.GetHuaweiClient(accessId, accessKey, sign, template, other)
	case Twilio:
		return sender.GetTwilioClient(accessId, accessKey, template)
	case SmsBao:
		return sender.GetSmsbaoClient(accessId, accessKey, sign, template, other)
	case MockSms:
		return sender.NewMocker(accessId, accessKey, sign, template, other)
	case SUBMAIL:
		return sender.GetSubmailClient(accessId, accessKey, template)
	case Sinowel:
		return GetSmsSinowelClient(accessId, accessKey, sign, template, other)
	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}
}
