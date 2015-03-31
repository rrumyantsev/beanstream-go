// +build integration

package tests

import (
	beanstream "github.com/Beanstream-DRWP/beanstream-go"
	"github.com/Beanstream-DRWP/beanstream-go/paymentMethods"
)

func createGateway() beanstream.Gateway {
	config := beanstream.Config{
		"300200578",
		"4BaD82D9197b4cc4b70a221911eE9f70",
		"D97D3BE1EE964A6193D17A571D9FBC80",
		"4e6Ff318bee64EA391609de89aD4CF5d",
		"www",
		"api",
		"v1",
		"-8:00"}
	return beanstream.Gateway{config}
}

func createCardRequest() beanstream.PaymentRequest {
	request := beanstream.PaymentRequest{
		PaymentMethod: paymentMethods.CARD,
		OrderNumber:   beanstream.Util_randOrderId(6),
		Amount:        12.99,
		Card: beanstream.CreditCard{
			Name:        "John Doe",
			Number:      "5100000010001004",
			ExpiryMonth: "11",
			ExpiryYear:  "19",
			Cvd:         "123",
			Complete:    true}}
	return request
}
