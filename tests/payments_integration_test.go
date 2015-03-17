// +build integration

package tests

import (
	"beanstream"
	"beanstream/paymentMethods"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestIntegration_Payments_MakePayment(t *testing.T) {
	gateway := createGateway()
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
	res, err := gateway.Payments().MakePayment(request) //returns a pointer to PaymentResponse
	assert.Nil(t, err, "Unexpected error occurred.", err)
	assert.NotNil(t, res, "Result was nil")
	assert.Equal(t, 1, res.Approved)
	assert.Equal(t, "P", res.Type)
}

func TestIntegration_Payments_MakePaymentFullDetails(t *testing.T) {
	gateway := createGateway()
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
			Complete:    true},
		BillingAddress: beanstream.Address{
			"John Doe",
			"123 Fake St.",
			"suite 3",
			"Victoria",
			"BC",
			"CA",
			"V8T4M3",
			"12505550123",
			"test@example.com"},
		ShippingAddress: beanstream.Address{
			"John Doe",
			"456 Jinglepot Rd.",
			"",
			"Nanaimo",
			"BC",
			"CA",
			"V9T1R9",
			"12505550123",
			"test@example.com"},
		Comment:    "a comment",
		Language:   "ENG",
		CustomerIp: "127.0.0.1",
		Custom: beanstream.CustomFields{
			"ref1 something",
			"ref2 something",
			"ref3 something",
			"ref4 something",
			"ref5 something"}}
	res, err := gateway.Payments().MakePayment(request) //returns a pointer to PaymentResponse
	assert.Nil(t, err, "Unexpected error occurred.", err)
	assert.NotNil(t, res, "Result was nil")
	assert.Equal(t, 1, res.Approved)
	assert.Equal(t, "P", res.Type)
}

func TestIntegration_Payments_PreAuthComplete(t *testing.T) {
	gateway := createGateway()
	request := createCardRequest()
	request.Card.Complete = false // pre-auth (complete=false)

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
	request := createCardRequest()
	res, err := gateway.Payments().MakePayment(request)
	assert.Nil(t, err, "Unexpected error occurred.", err)
	assert.NotNil(t, res, "Result was nil")
	assert.Equal(t, 1, res.Approved)
	assert.Equal(t, "P", res.Type)

	res2, err2 := gateway.Payments().VoidPayment(res.ID, 79.99)
	assert.Nil(t, err2, "Unexpected error occurred.", err2)
	assert.NotNil(t, res2, "Result was nil")
	assert.Equal(t, 1, res2.Approved)
	assert.Equal(t, "VP", res2.Type)
}

func TestIntegration_Payments_Return(t *testing.T) {
	gateway := createGateway()
	request := createCardRequest()
	res, err := gateway.Payments().MakePayment(request)
	assert.Nil(t, err, "Unexpected error occurred.", err)
	assert.NotNil(t, res, "Result was nil")
	assert.Equal(t, 1, res.Approved)
	assert.Equal(t, "P", res.Type)

	res2, err2 := gateway.Payments().ReturnPayment(res.ID, 55.00)
	assert.Nil(t, err2, "Unexpected error occurred.", err2)
	assert.NotNil(t, res2, "Result was nil")
	assert.Equal(t, 1, res2.Approved)
	assert.Equal(t, "R", res2.Type)
}

func TestIntegration_Payments_ReturnError(t *testing.T) {
	gateway := createGateway()
	request := createCardRequest()
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
	request := beanstream.PaymentRequest{
		PaymentMethod: paymentMethods.TOKEN,
		OrderNumber:   beanstream.Util_randOrderId(6),
		Amount:        15.99,
		Token: beanstream.Token{
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
	request := beanstream.PaymentRequest{
		PaymentMethod: paymentMethods.TOKEN,
		OrderNumber:   beanstream.Util_randOrderId(6),
		Amount:        50.00,
		Token: beanstream.Token{
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
	request := beanstream.PaymentRequest{
		PaymentMethod: paymentMethods.CASH,
		OrderNumber:   beanstream.Util_randOrderId(6),
		Amount:        12.00}
	res, err := gateway.Payments().MakePayment(request) //returns a pointer to PaymentResponse
	assert.Nil(t, err, "Unexpected error occurred.", err)
	assert.NotNil(t, res, "Result was nil")
	assert.Equal(t, 1, res.Approved)
	assert.Equal(t, "P", res.Type)
}

func TestIntegration_Payments_Cheque(t *testing.T) {
	gateway := createGateway()
	request := beanstream.PaymentRequest{
		PaymentMethod: paymentMethods.CHEQUE,
		OrderNumber:   beanstream.Util_randOrderId(6),
		Amount:        15.01}
	res, err := gateway.Payments().MakePayment(request) //returns a pointer to PaymentResponse
	assert.Nil(t, err, "Unexpected error occurred.", err)
	assert.NotNil(t, res, "Result was nil")
	assert.Equal(t, 1, res.Approved)
	assert.Equal(t, "P", res.Type)
}

func TestIntegration_Payments_GetTransaction(t *testing.T) {
	gateway := createGateway()
	request := createCardRequest()
	res, err := gateway.Payments().MakePayment(request)
	transId := res.ID

	trans, err := gateway.Payments().GetTransaction(transId)
	assert.Nil(t, err)
	assert.NotNil(t, trans)

}

func TestIntegration_Payments_GetTransactionWithAdjustments(t *testing.T) {
	gateway := createGateway()
	request := createCardRequest()
	res, err := gateway.Payments().MakePayment(request)
	transId := res.ID

	gateway.Payments().ReturnPayment(transId, 1.00) // a small return so we can see it in the transaction
	gateway.Payments().VoidPayment(transId, 5.00)

	trans, err := gateway.Payments().GetTransaction(transId)
	assert.Nil(t, err)
	assert.NotNil(t, trans)
	assert.NotNil(t, trans.Adjustments)
	assert.Equal(t, 2, len(trans.Adjustments))
}
