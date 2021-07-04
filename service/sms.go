package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/avast/retry-go"
)

type SMSService interface {
	SendSMS(ctx context.Context, smsRequest SMSRequest) (SMSResponse, error)
}

type SMSConfig struct {
	AccessToken string `json:"access_token"`
	Endpoint    string `json:"endpoint"`
}

type SMSRequest struct {
	RequestId   string `json:"request_id"`
	PhoneNumber string `json:"phone_number"`
	Template    string `json:"template"`
}

type SMSResponse struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`

	Data struct {
		Id string `json:"id"`
	}
}

type SMSServiceImpl struct {
	config SMSConfig
}

func NewSMSService(config SMSConfig) *SMSServiceImpl {
	return &SMSServiceImpl{config}
}

func (s *SMSServiceImpl) SendSMS(ctx context.Context, smsReq SMSRequest) (smsResp SMSResponse, err error) {
	body, err := json.Marshal(smsReq)
	if err != nil {
		fmt.Printf("convert sms request error: %v\n", err)
		return
	}
	fmt.Printf("send sms request: %v\n", smsReq)

	req, err := http.NewRequestWithContext(ctx, "POST", s.config.Endpoint, bytes.NewBuffer(body))
	req.Header.Add("Authorization", s.config.AccessToken)
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		fmt.Printf("create request error: %v\n", err)
		return
	}

	client := &http.Client{}

	err = retry.Do(
		func() error {
			resp, err := client.Do(req)
			if err != nil {
				fmt.Printf("call sms service error: %v\n", err)
				return err
			}
			defer resp.Body.Close()

			if err = json.NewDecoder(resp.Body).Decode(&smsResp); err != nil {
				fmt.Printf("convert sms response error: %v\n", err)
				return err
			}

			return nil
		},
		retry.Attempts(3),
	)
	if err != nil {
		return
	}

	return smsResp, nil
}
