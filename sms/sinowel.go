package sms

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func _md5(str string) string {
	data := []byte(str)
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has)
	return md5str
}

type SmsSinowelClient struct {
	username string
	password string
	sign     string
	template string
}

type SinowelResult struct {
	Detail []Detail `json:"detail"`
	Status string   `json:"status"`
}
type Detail struct {
	GatewayID string `json:"gatewayId"`
	Phone     string `json:"phone"`
	Status    string `json:"status"`
}

func GetSmsSinowelClient(username string, password string, sign string, template string, other []string) (*SmsSinowelClient, error) {
	return &SmsSinowelClient{
		username: username,
		password: password,
		sign:     sign,
		template: template,
	}, nil
}

func (c *SmsSinowelClient) SendMessage(param map[string]string, targetPhoneNumber ...string) error {
	code, ok := param["code"]
	if !ok {
		return fmt.Errorf("missing parameter: msg code")
	}

	if len(targetPhoneNumber) < 1 {
		return fmt.Errorf("missing parameter: targetPhoneNumber")
	}

	smsContent := "【" + c.sign + "】" + fmt.Sprintf(c.template, code)

	for _, mobile := range targetPhoneNumber {
		if strings.HasPrefix(mobile, "+86") {
			mobile = mobile[3:]
		} else if strings.HasPrefix(mobile, "+") {
			return fmt.Errorf("unsupported country code")
		}
		timestamp := time.Now().Unix() * 1000
		sign_md5 := _md5(fmt.Sprintf("%s%s%v", c.username, _md5(c.password), timestamp))
		url := "http://sms.sinowel.com/smsservice/json/SendSMS"
		str := fmt.Sprintf("{\"timestamp\": \"%v\",\"userId\": \"%s\",\"sign\": \"%s\",\"messageList\": [{\"phone\": \"%s\",\"message\": \"%s\"}]}", timestamp, c.username, sign_md5, mobile, smsContent)
		payload := strings.NewReader(str)
		client := &http.Client{}
		req, _ := http.NewRequest("POST", url, payload)
		req.Header.Add("content-type", "application/json")
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		result, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		// {"detail":[{"gatewayId":"4990152764117886716","phone":"15736933211","status":"R01"}],"status":"R01"}
		var sinowelSuccessResult SinowelResult
		err = json.Unmarshal(result, &sinowelSuccessResult)
		if err != nil {
			return fmt.Errorf("json unmarshal error")
		}
		if sinowelSuccessResult.Status != "R01" {
			return fmt.Errorf("sinowel sms error")
		}
	}
	return nil
}
