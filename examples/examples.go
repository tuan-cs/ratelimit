package main

import (
	"context"
	"fmt"
	"time"

	"github.com/ratelimit/limiter"
	"github.com/ratelimit/runner"
	"github.com/ratelimit/service"
)

func Example_basic() {
	limitTotalSMS := limiter.Config{
		IdentifierExtrator: func(ctx context.Context) (string, error) {
			return "total", nil
		},

		ErrorHandler: func(_ context.Context, err error) error {
			fmt.Printf("error: %v\n", err)
			return err
		},

		DenyHandler: func(_ context.Context, identifier string, err error) error {
			fmt.Printf("deny: identifier: %v, err: %v\n", identifier, err)
			return err
		},

		Store: limiter.NewTotalMemoryStoreWithConfig(limiter.TotalMemoryStoreConfig{Rate: 5, Burst: 5, ExpiresIn: 1}),
	}

	cmd := limiter.New(limitTotalSMS)

	config := service.SMSConfig{
		Endpoint:    "http://localhost:10000/sms",
		AccessToken: "access_token",
	}
	smsService := service.NewSMSService(config)
	smsReq := service.SMSRequest{
		RequestId:   "a47b9c69-0c10-4211-a175-b22331810e1",
		PhoneNumber: "0966666666",
		Template:    "Ma xac thuc OTP la: 932781",
	}

	for i := 0; i < 7; i++ {
		var result service.SMSResponse
		err := cmd.Run(context.TODO(), func(ctx context.Context) (err error) {
			result, err = smsService.SendSMS(ctx, smsReq)
			if err != nil {
				fmt.Printf("limiter error: %v\n", err)
				return
			}

			return
		})
		if err != nil {
			fmt.Printf("error in limiter: %v\n", err)
		} else {
			fmt.Printf("send sms result: %v\n", result)
		}
	}
}

func Example_chain() {
	limitTotalSMS := limiter.Config{
		IdentifierExtrator: func(ctx context.Context) (string, error) {
			return "total", nil
		},

		ErrorHandler: func(_ context.Context, err error) error {
			fmt.Printf("error: %v\n", err)
			return err
		},

		DenyHandler: func(_ context.Context, identifier string, err error) error {
			fmt.Printf("deny: identifier: %v, err: %v\n", identifier, err)
			return err
		},

		Store: limiter.NewTotalMemoryStoreWithConfig(limiter.TotalMemoryStoreConfig{Rate: 10, Burst: 10, ExpiresIn: 1 * time.Second}),
	}

	limitPhoneSMS := limiter.Config{
		IdentifierExtrator: func(ctx context.Context) (string, error) {
			request := ctx.Value("Request")
			smsRequest := request.(service.SMSRequest)
			return smsRequest.PhoneNumber, nil
		},

		ErrorHandler: func(_ context.Context, err error) error {
			fmt.Printf("error: %v\n", err)
			return err
		},

		DenyHandler: func(_ context.Context, identifier string, err error) error {
			fmt.Printf("deny: phone number %v, error: %v\n", identifier, err)
			return err
		},

		Store: limiter.NewIndividualMemoryStoreWithConfig(limiter.IndividualMemoryStoreConfig{Rate: 2, Burst: 2, ExpiresIn: 1 * time.Second}),
	}

	cmd := runner.RunnerChain(
		limiter.NewMiddleware(limitTotalSMS),
		limiter.NewMiddleware(limitPhoneSMS),
	)

	config := service.SMSConfig{
		Endpoint:    "http://localhost:10000/sms",
		AccessToken: "access_token",
	}
	smsService := service.NewSMSService(config)
	smsReq := service.SMSRequest{
		RequestId:   "a47b9c69-0c10-4211-a175-b22331810e1",
		PhoneNumber: "0966666666",
		Template:    "Ma xac thuc OTP la: 932781",
	}

	for i := 0; i < 7; i++ {
		var result service.SMSResponse
		ctx := context.TODO()
		ctx = context.WithValue(ctx, "Request", smsReq)
		err := cmd.Run(ctx, func(ctx context.Context) (err error) {
			result, err = smsService.SendSMS(ctx, smsReq)
			if err != nil {
				fmt.Printf("send sms error: %v\n", err)
				return
			}
			return
		})

		if err != nil {
			fmt.Printf("error in limiter: %v\n", err)
		} else {
			fmt.Printf("send sms result: %v\n", result)
		}
	}
}

func main() {
	Example_basic()
	Example_chain()
}
