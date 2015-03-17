package beanstream

import (
	"beanstream/httpMethods"
	//"beanstream/paymentMethods"
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

func (api PaymentsAPI) GetTransaction(transId string) (*Transaction, error) {
	url := api.Config.BaseUrl() + getPaymentUrl
	url = fmt.Sprintf(url, transId)

	responseType := Transaction{}
	res, err := Process(httpMethods.GET, url, api.Config.MerchantId, api.Config.PaymentsApiKey, &responseType)
	if err != nil {
		return nil, err
	}
	pr := res.(*Transaction)
	pr.CreatedTime = api.AsDate(pr.created)
	return pr, nil
}

type PaymentRequest struct {
	PaymentMethod   string         `json:"payment_method"`
	OrderNumber     string         `json:"order_number,omitempty"`
	Amount          float32        `json:"amount"`
	Card            CreditCard     `json:"card,omitempty"`
	Token           Token          `json:"token,omitempty"`
	Profile         ProfilePayment `json:"payment_profile,omitempty"`
	BillingAddress  Address        `json:"billing,omitempty"`
	ShippingAddress Address        `json:"shipping,omitempty"`
	Comment         string         `json:"comments,omitempty"`
	Language        string         `json:"language,omitempty"`
	CustomerIp      string         `json:"customer_ip,omitempty"`
	TermUrl         string         `json:"term_url,omitempty"`
	Custom          CustomFields   `json:"custom,omitempty"`
}

type CreditCard struct {
	Name        string `json:"name"`
	Number      string `json:"number"`
	ExpiryMonth string `json:"expiry_month"`
	ExpiryYear  string `json:"expiry_year"`
	Cvd         string `json:"cvd"`
	Complete    bool   `json:"complete"`
	Function    string `json:"function,omitempty"`
	Type        string `json:"card_type,omitempty"`
	Id          int    `json:"card_id,string,omitempty"`
	AvsResult   string `json:"avs_result,omitempty"`
	CvdResult   string `json:"cvd_result,omitempty"`
}

type Token struct {
	Token    string `json:"code"`
	Name     string `json:"name"`
	Complete bool   `json:"complete"`
}

type ProfilePayment struct {
	ProfileId string `json:"customer_code"`
	CardId    int    `json:"card_id"`
	Complete  bool   `json:"complete"`
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

type Address struct {
	Name         string `json:"name,omitempty"`
	AddressLine1 string `json:"address_line1,omitempty"`
	AddressLine2 string `json:"address_line2,omitempty"`
	City         string `json:"city,omitempty"`
	Province     string `json:"province,omitempty"`
	Country      string `json:"country,omitempty"`
	PostalCode   string `json:"postal_code,omitempty"`
	PhoneNumber  string `json:"phone_number,omitempty"`
	EmailAddress string `json:"email_address,omitempty"`
}

type CustomFields struct {
	Ref1 string `json:"ref1,omitempty"`
	Ref2 string `json:"ref2,omitempty"`
	Ref3 string `json:"ref3,omitempty"`
	Ref4 string `json:"ref4,omitempty"`
	Ref5 string `json:"ref5,omitempty"`
}

type Transaction struct {
	Id               int    `json:"id,omitempty"`
	Approved         int    `json:"approved,omitempty"`
	MessageId        int    `json:"message_id,omitempty"`
	Message          string `json:"message,omitempty"`
	AuthCode         string `json:"auth_code,omitempty"`
	created          string `json:"created,omitempty"`
	CreatedTime      time.Time
	OrderNumber      string       `json:"order_number,omitempty"`
	Amount           float32      `json:"amount,omitempty"`
	Type             string       `json:"type,omitempty"`
	Comment          string       `json:"comments,omitempty"`
	BatchNumber      string       `json:"batch_number,omitempty"`
	TotalRefunds     float32      `json:"total_refunds,omitempty"`
	TotalCompletions float32      `json:"total_completions,omitempty"`
	PaymentMethod    string       `json:"payment_method,omitempty"`
	Card             CreditCard   `json:"card,omitempty"`
	BillingAddress   Address      `json:"billing,omitempty"`
	ShippingAddress  Address      `json:"shipping,omitempty"`
	Custom           CustomFields `json:"custom,omitempty"`
	Adjustments      []Adjustment `json:"adjustments,omitempty"`
	Links            []Link       `json:"links,omitempty"`
}

func (t *Transaction) IsApproved() bool {
	if t.Amount == 1 {
		return true
	}
	return false
}

type Adjustment struct {
	Id       string  `json:"id,omitempty"`
	Type     string  `json:"type,omitempty"`
	Approval string  `json:"approval,omitempty"`
	Message  string  `json:"message,omitempty"`
	Amount   float32 `json:"amount,omitempty"`
	Url      string  `json:"url,omitempty"`
}

type Link struct {
	Rel    string `json:"rel,omitempty"`
	Href   string `json:"href,omitempty"`
	Method string `json:"method,omitempty"`
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
