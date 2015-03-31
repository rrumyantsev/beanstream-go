/*
Package beanstream supplies the 3 APIs for processing payments:
	PaymentsAPI
	ProfilesAPI
	ReportingAPI

Each API has its own Passcode and processes against your Merchant ID.

To start using an API you must create a Gateway and supply it the configuration
it needs to run:
	gateway := beanstream.Gateway{beanstream.Config{
		"300200578", // merchant ID
		"4BaD82D9197b4cc4b70a221911eE9f70", // Payments Passcode
		"D97D3BE1EE964A6193D17A571D9FBC80", // Profiles Passcode
		"4e6Ff318bee64EA391609de89aD4CF5d", // Reporting Passcode
		"www", // url prefix
		"api", // url suffix
		"v1", // api version
		"-8:00"}} // timezone offset

The above values use a Beanstream Test account.

To Create a new payment (credit card, cash, cheque...) use the Payments API and supply
it with a PaymentRrequest:
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
	res, err := gateway.Payments().MakePayment(request)

For more details visit the documentation for each particular API.
*/
package beanstream
