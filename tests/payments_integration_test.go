// +build integration

package tests

import (
	"beanstream"
	"beanstream/paymentMethods"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
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

func TestIntegration_Payments_MakePayment(t *testing.T) {
	gateway := createGateway()
	request := beanstream.PaymentRequest{
		paymentMethods.CARD,
		beanstream.Util_randOrderId(6),
		12.99,
		beanstream.CreditCard{
			"John Doe",
			"5100000010001004",
			"11",
			"19",
			"123",
			true}}
	res, err := gateway.Payments().MakePayment(request) //returns a pointer to PaymentResponse
	assert.Nil(t, err, "Unexpected error occurred.", err)
	assert.NotNil(t, res, "Result was nil")
	assert.Equal(t, 1, res.Approved)
	assert.Equal(t, "P", res.Type)
}

func TestIntegration_Payments_PreAuthComplete(t *testing.T) {
	gateway := createGateway()
	request := beanstream.PaymentRequest{
		paymentMethods.CARD,
		beanstream.Util_randOrderId(6),
		50.00,
		beanstream.CreditCard{
			"John Doe",
			"5100000010001004",
			"11",
			"19",
			"123",
			false}} // pre-auth (complete=false)
	res, err := gateway.Payments().MakePayment(request)
	assert.Nil(t, err, "Unexpected error occurred.", err)
	assert.NotNil(t, res, "Result was nil")
	assert.Equal(t, 1, res.Approved)
	assert.Equal(t, "PA", res.Type)

	res2, err2 := gateway.Payments().CompletePayment(res.ID, 24.67)
	assert.Nil(t, err2, "Unexpected error occurred.", err2)
	assert.NotNil(t, res2, "Result was nil")
	assert.Equal(t, 1, res2.Approved)
	assert.Equal(t, "PAC", res2.Type)
}

func TestIntegration_Payments_Void(t *testing.T) {
	gateway := createGateway()
	request := beanstream.PaymentRequest{
		paymentMethods.CARD,
		beanstream.Util_randOrderId(6),
		22.00,
		beanstream.CreditCard{
			"John Doe",
			"5100000010001004",
			"11",
			"19",
			"123",
			true}}
	res, err := gateway.Payments().MakePayment(request)
	assert.Nil(t, err, "Unexpected error occurred.", err)
	assert.NotNil(t, res, "Result was nil")
	assert.Equal(t, 1, res.Approved)
	assert.Equal(t, "P", res.Type)

	res2, err2 := gateway.Payments().VoidPayment(res.ID, 22.00)
	assert.Nil(t, err2, "Unexpected error occurred.", err2)
	assert.NotNil(t, res2, "Result was nil")
	assert.Equal(t, 1, res2.Approved)
	assert.Equal(t, "VP", res2.Type)
}

func TestIntegration_Payments_Return(t *testing.T) {
	gateway := createGateway()
	request := beanstream.PaymentRequest{
		paymentMethods.CARD,
		beanstream.Util_randOrderId(6),
		100.00,
		beanstream.CreditCard{
			"John Doe",
			"5100000010001004",
			"11",
			"19",
			"123",
			true}}
	res, err := gateway.Payments().MakePayment(request)
	assert.Nil(t, err, "Unexpected error occurred.", err)
	assert.NotNil(t, res, "Result was nil")
	assert.Equal(t, 1, res.Approved)
	assert.Equal(t, "P", res.Type)

	res2, err2 := gateway.Payments().ReturnPayment(res.ID, 100.00)
	assert.Nil(t, err2, "Unexpected error occurred.", err2)
	assert.NotNil(t, res2, "Result was nil")
	assert.Equal(t, 1, res2.Approved)
	assert.Equal(t, "R", res2.Type)
}

func TestIntegration_Payments_ReturnError(t *testing.T) {
	gateway := createGateway()
	request := beanstream.PaymentRequest{
		paymentMethods.CARD,
		beanstream.Util_randOrderId(6),
		100.00,
		beanstream.CreditCard{
			"John Doe",
			"5100000010001004",
			"11",
			"19",
			"123",
			true}}
	res, err := gateway.Payments().MakePayment(request)
	assert.Nil(t, err, "Unexpected error occurred.", err)
	assert.NotNil(t, res, "Result was nil")
	assert.Equal(t, 1, res.Approved)
	assert.Equal(t, "P", res.Type)

	// return more than we charged so we get an error
	res2, err2 := gateway.Payments().ReturnPayment(res.ID, 105.00)
	assert.Nil(t, res2, "Did not expect a proper result", res2)
	assert.NotNil(t, err2, "Error was nil and shouldn't have been")
	bicError := err2.(*beanstream.BeanstreamApiException)
	assert.True(t, strings.Contains(bicError.Error(), "InvalidRequestException"), "Error is not an InvalidRequestException")
	assert.Equal(t, 194, bicError.Code)
	assert.Equal(t, 2, bicError.Category)
}
func TestIntegration_Payments_Token(t *testing.T) {
	// step 1: get the token
	token, err := beanstream.LegatoTokenizeCard(
		"5100000010001004",
		"11",
		"19",
		"123")
	assert.Nil(t, err, "Unexpected error occurred.", err)
	assert.NotNil(t, token, "No token returned")
	assert.NotEmpty(t, token, "Legato token was empty")

	// step 2: make the purchase
	gateway := createGateway()
	request := beanstream.PaymentRequestToken{
		paymentMethods.TOKEN,
		beanstream.Util_randOrderId(6),
		15.99,
		beanstream.Token{
			token,
			"John Doe",
			true}}
	res, err2 := gateway.Payments().MakePayment(request)
	assert.Nil(t, err2, "Unexpected error occurred.", err2)
	assert.NotNil(t, res, "Result was nil")
	assert.Equal(t, 1, res.Approved)
	assert.Equal(t, "P", res.Type)

}

func TestIntegration_Payments_TokenPreAuth(t *testing.T) {
	// step 1: get the token
	token, err := beanstream.LegatoTokenizeCard(
		"5100000010001004",
		"11",
		"19",
		"123")
	assert.Nil(t, err, "Unexpected error occurred.", err)
	assert.NotNil(t, token, "No token returned")
	assert.NotEmpty(t, token, "Legato token was empty")

	// step 2: make the purchase
	gateway := createGateway()
	request := beanstream.PaymentRequestToken{
		paymentMethods.TOKEN,
		beanstream.Util_randOrderId(6),
		30.00,
		beanstream.Token{
			token,
			"John Doe",
			false}} // pre-auth (complete=false)
	res, err2 := gateway.Payments().MakePayment(request)
	assert.Nil(t, err2, "Unexpected error occurred.", err2)
	assert.NotNil(t, res, "Result was nil")
	assert.Equal(t, 1, res.Approved)
	assert.Equal(t, "PA", res.Type)

	res2, err2 := gateway.Payments().CompletePayment(res.ID, 12.01)
	assert.Nil(t, err2, "Unexpected error occurred.", err2)
	assert.NotNil(t, res2, "Result was nil")
	assert.Equal(t, 1, res2.Approved)
	assert.Equal(t, "PAC", res2.Type)
}

func TestIntegration_Payments_Cash(t *testing.T) {
	gateway := createGateway()
	request := beanstream.CashPayment{
		paymentMethods.CASH,
		beanstream.Util_randOrderId(6),
		10}
	res, err := gateway.Payments().MakePayment(request) //returns a pointer to PaymentResponse
	assert.Nil(t, err, "Unexpected error occurred.", err)
	assert.NotNil(t, res, "Result was nil")
	assert.Equal(t, 1, res.Approved)
	assert.Equal(t, "P", res.Type)
}

func TestIntegration_Payments_Cheque(t *testing.T) {
	gateway := createGateway()
	request := beanstream.ChequePayment{
		paymentMethods.CHEQUE,
		beanstream.Util_randOrderId(6),
		14.35}
	res, err := gateway.Payments().MakePayment(request) //returns a pointer to PaymentResponse
	assert.Nil(t, err, "Unexpected error occurred.", err)
	assert.NotNil(t, res, "Result was nil")
	assert.Equal(t, 1, res.Approved)
	assert.Equal(t, "P", res.Type)
}
