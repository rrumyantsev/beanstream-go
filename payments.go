package beanstream

import (
	"beanstream/httpMethods"
	"fmt"
	"time"
)

const paymentUrl = "/payments"
const continueUrl = paymentUrl + "/%v/continue"
const completionUrl = paymentUrl + "/%v/completions"
const returnUrl = paymentUrl + "/%v/returns"
const voidUrl = paymentUrl + "/%v/void"
const getPaymentUrl = paymentUrl + "/%v"

type PaymentsAPI struct {
	Config Config
}

func (api PaymentsAPI) MakePayment(transaction interface{}) (*PaymentResponse, error) {
	url := api.Config.BaseUrl() + paymentUrl
	responseType := PaymentResponse{}
	res, err := ProcessBody(httpMethods.POST, url, api.Config.MerchantId, api.Config.PaymentsApiKey, transaction, &responseType)
	if err != nil {
		return nil, err
	}
	//fmt.Printf("MakePayment result: %T %v\n", res, res)
	pr := res.(*PaymentResponse)
	pr.CreatedTime = api.AsDate(pr.created)
	return pr, nil
}

func (api PaymentsAPI) CompletePayment(transId string, amount float32) (*PaymentResponse, error) {
	url := api.Config.BaseUrl() + completionUrl
	url = fmt.Sprintf(url, transId)
	responseType := PaymentResponse{}
	req := completionRequest{amount}
	res, err := ProcessBody(httpMethods.POST, url, api.Config.MerchantId, api.Config.PaymentsApiKey, req, &responseType)
	if err != nil {
		return nil, err
	}
	pr := res.(*PaymentResponse)
	pr.CreatedTime = api.AsDate(pr.created)
	return pr, nil
}

func (api PaymentsAPI) VoidPayment(transId string, amount float32) (*PaymentResponse, error) {
	url := api.Config.BaseUrl() + voidUrl
	url = fmt.Sprintf(url, transId)
	responseType := PaymentResponse{}
	req := voidRequest{amount}
	res, err := ProcessBody(httpMethods.POST, url, api.Config.MerchantId, api.Config.PaymentsApiKey, req, &responseType)
	if err != nil {
		return nil, err
	}
	pr := res.(*PaymentResponse)
	pr.CreatedTime = api.AsDate(pr.created)
	return pr, nil
}

func (api PaymentsAPI) ReturnPayment(transId string, amount float32) (*PaymentResponse, error) {
	url := api.Config.BaseUrl() + returnUrl
	url = fmt.Sprintf(url, transId)
	responseType := PaymentResponse{}
	req := returnRequest{amount}
	res, err := ProcessBody(httpMethods.POST, url, api.Config.MerchantId, api.Config.PaymentsApiKey, req, &responseType)
	if err != nil {
		return nil, err
	}
	pr := res.(*PaymentResponse)
	pr.CreatedTime = api.AsDate(pr.created)
	return pr, nil
}

type PaymentRequest struct {
	PaymentMethod string     `json:"payment_method"`
	OrderNumber   string     `json:"order_number"`
	Amount        float32    `json:"amount"`
	Card          CreditCard `json:"card"`
}

type PaymentRequestToken struct {
	PaymentMethod string  `json:"payment_method"`
	OrderNumber   string  `json:"order_number"`
	Amount        float32 `json:"amount"`
	Token         Token   `json:"token"`
}

type CreditCard struct {
	Name        string `json:"name"`
	Number      string `json:"number"`
	ExpiryMonth string `json:"expiry_month"`
	ExpiryYear  string `json:"expiry_year"`
	Cvd         string `json:"cvd"`
	Complete    bool   `json:"complete"`
}

type Token struct {
	Token    string `json:"code"`
	Name     string `json:"name"`
	Complete bool   `json:"complete"`
}

type CashPayment struct {
	PaymentMethod string  `json:"payment_method"`
	OrderNumber   string  `json:"order_number"`
	Amount        float32 `json:"amount"`
}

type ChequePayment struct {
	PaymentMethod string  `json:"payment_method"`
	OrderNumber   string  `json:"order_number"`
	Amount        float32 `json:"amount"`
}

type completionRequest struct {
	Amount float32 `json:"amount"`
}

type voidRequest struct {
	Amount float32 `json:"amount"`
}

type returnRequest struct {
	Amount float32 `json:"amount"`
}

// JSON:
//{
// 	"id":"10108462",
//	"approved":"1",
//	"message_id":"1",
//	"message":"Approved",
//	"auth_code":"TEST",
//	"created":"2015-03-13T08:59:24",
//	"order_number":"YFEJXU1426262363",
//	"type":"PA",
//	"payment_method":"CC",
//	"card":{
//		"card_type":"MC",
//		"last_four":"1004",
//		"cvd_match":0,
//		"address_match":0,
//		"postal_result":0},
//	"links":[
//		{"rel":"complete","href":"https://www.beanstream.com/api/v1/payments/10108462/completions","method":"POST"}
//	]
//}
//
type PaymentResponse struct {
	Approved int    `json:"approved,string"`
	AuthCode string `json:"auth_code"`
	Card     struct {
		AddressMatch int    `json:"address_match"`
		CardType     string `json:"card_type"`
		CvdMatch     int    `json:"cvd_match"`
		LastFour     string `json:"last_four"`
		PostalResult int    `json:"postal_result"`
	} `json:"card"`
	created     string `json:"created,omitempty"`
	CreatedTime time.Time
	ID          string `json:"id"`
	Links       []struct {
		Href   string `json:"href"`
		Method string `json:"method"`
		Rel    string `json:"rel"`
	} `json:"links"`
	Message       string `json:"message"`
	MessageID     int    `json:"message_id,string"`
	OrderNumber   string `json:"order_number"`
	PaymentMethod string `json:"payment_method"`
	Type          string `json:"type"`
}

func (api PaymentsAPI) AsDate(val string) time.Time {
	rfc3339Time := val + "Z" + api.Config.TimezoneOffset
	t, _ := time.Parse(time.RFC3339, rfc3339Time)
	return t
}
