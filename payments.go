package beanstream

import (
	"fmt"
	"github.com/Beanstream/beanstream-go/httpMethods"
	"time"
)

const paymentUrl = "/payments"
const continueUrl = paymentUrl + "/%v/continue"
const completionUrl = paymentUrl + "/%v/completions"
const returnUrl = paymentUrl + "/%v/returns"
const voidUrl = paymentUrl + "/%v/void"
const getPaymentUrl = paymentUrl + "/%v"

/*
Through the payments API you can create payments, get payments,
return payments, void payments, as well as pre-authorize and complete
payments.
*/
type PaymentsAPI struct {
	Config Config
}

/*
Create a payment. Either a Credit Card, Profile, Cash, or Cheque payment request. Cash and Cheque payments
are just for your own record keeping.
You must supply it a PaymentRequest that is defined in this package
*/
func (api PaymentsAPI) MakePayment(transaction PaymentRequest) (*PaymentResponse, error) {
	url := api.Config.BaseUrl() + paymentUrl
	responseType := PaymentResponse{}
	res, err := ProcessBody(httpMethods.POST, url, api.Config.MerchantId, api.Config.PaymentsApiKey, transaction, &responseType)
	if err != nil {
		return nil, err
	}
	//fmt.Printf("MakePayment result: %T %v\n", res, res)
	pr := res.(*PaymentResponse)
	pr.CreatedTime = AsDate(pr.created, api.Config)
	return pr, nil
}

// Complete a pre-authorized payment for some or all of the pre-authorized amount.
func (api PaymentsAPI) CompletePayment(transId string, request PaymentRequest) (*PaymentResponse, error) {
	url := api.Config.BaseUrl() + completionUrl
	url = fmt.Sprintf(url, transId)
	responseType := PaymentResponse{}
	res, err := ProcessBody(httpMethods.POST, url, api.Config.MerchantId, api.Config.PaymentsApiKey, request, &responseType)
	if err != nil {
		return nil, err
	}
	pr := res.(*PaymentResponse)
	pr.CreatedTime = AsDate(pr.created, api.Config)
	return pr, nil
}

// VoidPayment cancels a payment for all of the original amount.
// In order to void a payment you must not wait too long.
// The amount must equal the original amount.
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
	pr.CreatedTime = AsDate(pr.created, api.Config)
	return pr, nil
}

// ReturnPayment returns the money to the customer for all or some of the original amount.
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
	pr.CreatedTime = AsDate(pr.created, api.Config)
	return pr, nil
}

// GetTransaction retrieves a transaction and all adjustments that were performed on it.
func (api PaymentsAPI) GetTransaction(transId string) (*Transaction, error) {
	url := api.Config.BaseUrl() + getPaymentUrl
	url = fmt.Sprintf(url, transId)

	responseType := Transaction{}
	res, err := Process(httpMethods.GET, url, api.Config.MerchantId, api.Config.PaymentsApiKey, &responseType)
	if err != nil {
		return nil, err
	}
	pr := res.(*Transaction)
	pr.CreatedTime = AsDate(pr.created, api.Config)
	if pr.Adjustments != nil {
		for _, adj := range pr.Adjustments {
			adj.CreatedTime = AsDate(adj.created, api.Config)
		}
	}
	return pr, nil
}

// PaymentRequest is the main struct for making a payment. The mandatory fields are:
//	PaymentMethod
//	Ordernumber
//	Amount
// For specific types of payments you will need: Card, Token, or Profile
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

// CreditCard info for making a payment.
// You can pre-authorize a purchase by setting Complete to false.
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

// Token is a single-use Legato token for making a payment.
// You can pre-authorize a purchase by setting Complete to false.
type Token struct {
	Token    string `json:"code"`
	Name     string `json:"name"`
	Complete bool   `json:"complete"`
}

// ProfilePayment allows you to make a payment with a Payment Profile. You need the Profile ID
// as well as the Card id. If there is just one card on a profile
// use card ID = 1
// You can pre-authorize a purchase by setting Complete to false.
type ProfilePayment struct {
	ProfileId string `json:"customer_code"`
	CardId    int    `json:"card_id"`
	Complete  bool   `json:"complete"`
}

type completionRequest struct {
	Amount float32      `json:"amount"`
	Custom CustomFields `json:"custom"`
}

type voidRequest struct {
	Amount float32 `json:"amount"`
}

type returnRequest struct {
	Amount float32 `json:"amount"`
}

// Address is either a billing or shipping address
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

// CustomFields that can be added to a transaction on creation.
type CustomFields struct {
	Ref1 string `json:"ref1,omitempty"`
	Ref2 string `json:"ref2,omitempty"`
	Ref3 string `json:"ref3,omitempty"`
	Ref4 string `json:"ref4,omitempty"`
	Ref5 string `json:"ref5,omitempty"`
}

/*
Transaction is the struct you receive when you call GetTransaction().
To check if a transaction is approved you can call the method IsApproved()
*/
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
	Adjustments      []Adjustment `json:"adjusted_by,omitempty"`
	Links            []Link       `json:"links,omitempty"`
}

// IsApproved will test if a Payment was approved
func (t *Transaction) IsApproved() bool {
	if t.Approved == 1 {
		return true
	}
	return false
}

/*
Adjustment to a payment, often a return or void.
*/
type Adjustment struct {
	Id          int     `json:"id,omitempty"`
	Type        string  `json:"type,omitempty"`
	Approval    int     `json:"approval,omitempty"`
	Message     string  `json:"message,omitempty"`
	Amount      float32 `json:"amount,omitempty"`
	created     string  `json:"created,omitempty"`
	CreatedTime time.Time
	Url         string `json:"url,omitempty"`
}

/*
Link for http access for returns and voids. Not useful with
the SDK be here none the less.
*/
type Link struct {
	Rel    string `json:"rel,omitempty"`
	Href   string `json:"href,omitempty"`
	Method string `json:"method,omitempty"`
}

// PaymentResponse is the response from a successful transaction. Some fields might be empty.
// To check if a transaction is approved you can call the method IsApproved().
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

// Example json for PaymentResponse
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

// IsApproved will test if a Payment was approved
func (t *PaymentResponse) IsApproved() bool {
	if t.Approved == 1 {
		return true
	}
	return false
}
