package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type SMSRequest struct {
	RequestId   string `json:"request_id"`
	PhoneNumber string `json:"phone_number"`
	Template    string `json:"template"`
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Data struct {
	Id string `json:"id"`
}

type SMSResponse struct {
	Error Error `json:"error"`
	Data  Data  `json:"data"`
}

func main() {
	e := echo.New()

	e.POST("/sms", func(c echo.Context) error {
		smsReq := new(SMSRequest)
		if err := c.Bind(smsReq); err != nil {
			fmt.Printf("convert post data error: %v\n", err)
			return c.JSON(http.StatusBadRequest,
				SMSResponse{
					Error: Error{
						Code:    10,
						Message: "Invalid",
					},
					Data: Data{
						Id: "",
					},
				})
		}

		fmt.Printf("Received data: %v\n", smsReq)
		return c.JSON(http.StatusBadRequest,
			SMSResponse{
				Error: Error{
					Code:    0,
					Message: "",
				},
				Data: Data{
					Id: "1111",
				},
			})
	})

	e.Logger.Fatal(e.Start(":10000"))
}
