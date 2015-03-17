// +build integration

package tests

import (
	"beanstream"
	"beanstream/paymentMethods"
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestIntegration_Profiles_CreateProfile(t *testing.T) {
	gateway := createGateway()
	request := beanstream.Profile{
		Card: beanstream.CreditCard{
			Name:        "John Doe",
			Number:      "5100000010001004",
			ExpiryMonth: "11",
			ExpiryYear:  "19",
			Cvd:         "123"},
		BillingAddress: beanstream.Address{
			"John Doe",
			"123 Fake St.",
			"suite 3",
			"Victoria",
			"BC",
			"CA",
			"V8T4M3",
			"12505550123",
			"test@example.com"}}

	res, err := gateway.Profiles().CreateProfile(request)
	assert.Nil(t, err)
	assert.NotNil(t, res)
	fmt.Println(res)
	assert.NotNil(t, res.Id)
	assert.NotEmpty(t, res.Id)

	// delete profile
	profile, err2 := gateway.Profiles().DeleteProfile(res.Id)
	assert.Nil(t, err2)
	assert.NotNil(t, profile)
}

func TestIntegration_Profiles_MakePayment(t *testing.T) {
	gateway := createGateway()
	request := beanstream.Profile{
		Card: beanstream.CreditCard{
			Name:        "John Doe",
			Number:      "5100000010001004",
			ExpiryMonth: "11",
			ExpiryYear:  "19",
			Cvd:         "123"},
		BillingAddress: beanstream.Address{
			"John Doe",
			"123 Fake St.",
			"suite 3",
			"Victoria",
			"BC",
			"CA",
			"V8T4M3",
			"12505550123",
			"test@example.com"}}

	res, err := gateway.Profiles().CreateProfile(request)
	assert.Nil(t, err)
	assert.NotNil(t, res)
	fmt.Println(res)
	assert.NotNil(t, res.Id)
	assert.NotEmpty(t, res.Id)

	payment := beanstream.PaymentRequest{
		PaymentMethod: paymentMethods.PROFILE,
		OrderNumber:   beanstream.Util_randOrderId(6),
		Amount:        12.99,
		Profile: beanstream.ProfilePayment{
			res.Id,
			1,
			true}}
	res2, err2 := gateway.Payments().MakePayment(payment)
	assert.Nil(t, err2)
	assert.NotNil(t, res2)
	assert.Equal(t, 1, res2.Approved)
	assert.Equal(t, "P", res2.Type)

	// delete profile
	profile, err3 := gateway.Profiles().DeleteProfile(res.Id)
	assert.Nil(t, err3)
	assert.NotNil(t, profile)
}

func TestIntegration_Profiles_CreateProfileFromToken(t *testing.T) {
	// step 1: get the token
	token, t_err := beanstream.LegatoTokenizeCard(
		"5100000010001004",
		"11",
		"19",
		"123")
	assert.Nil(t, t_err)
	assert.NotNil(t, token)
	assert.NotEmpty(t, token)

	// step 2: create the profile
	gateway := createGateway()
	request := beanstream.Profile{
		Token: beanstream.Token{
			Name:  "John Doe",
			Token: token},
		BillingAddress: beanstream.Address{
			"John Doe",
			"123 Fake St.",
			"suite 3",
			"Victoria",
			"BC",
			"CA",
			"V8T4M3",
			"12505550123",
			"test@example.com"}}

	res, err := gateway.Profiles().CreateProfile(request)
	assert.Nil(t, err)
	assert.NotNil(t, res)
	fmt.Println(res)
	assert.NotNil(t, res.Id)
	assert.NotEmpty(t, res.Id)

	// step 3: Make a payment
	payment := beanstream.PaymentRequest{
		PaymentMethod: paymentMethods.PROFILE,
		OrderNumber:   beanstream.Util_randOrderId(6),
		Amount:        14.99,
		Profile: beanstream.ProfilePayment{
			res.Id,
			1,
			true},
		Comment: "Payment with a token profile"}
	res2, err2 := gateway.Payments().MakePayment(payment)
	assert.Nil(t, err2)
	assert.NotNil(t, res2)
	assert.Equal(t, 1, res2.Approved)
	assert.Equal(t, "P", res2.Type)

	// step 4: Make another payment
	payment.OrderNumber = beanstream.Util_randOrderId(6)
	payment.Amount = 1.89
	payment.Comment = "A 2nd payment with the same token profile"
	res3, err3 := gateway.Payments().MakePayment(payment)
	assert.Nil(t, err3)
	assert.NotNil(t, res3)
	assert.Equal(t, 1, res3.Approved)
	assert.Equal(t, "P", res3.Type)

	// clean up: delete profile
	profile, err4 := gateway.Profiles().DeleteProfile(res.Id)
	assert.Nil(t, err4)
	assert.NotNil(t, profile)
}

func TestIntegration_Profiles_DeleteProfile(t *testing.T) {
	gateway := createGateway()
	request := beanstream.Profile{
		Card: beanstream.CreditCard{
			Name:        "John Doe",
			Number:      "5100000010001004",
			ExpiryMonth: "11",
			ExpiryYear:  "19",
			Cvd:         "123"},
		BillingAddress: beanstream.Address{
			"John Doe",
			"999 Fake St.",
			"suite 3",
			"Victoria",
			"BC",
			"CA",
			"V8T4M3",
			"12505550123",
			"test@example.com"}}
	res, err := gateway.Profiles().CreateProfile(request)
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.NotEmpty(t, res.Id)

	// delete profile
	profile, err2 := gateway.Profiles().DeleteProfile(res.Id)
	assert.Nil(t, err2)
	assert.NotNil(t, profile)

	// get profile
	profile2, err3 := gateway.Profiles().GetProfile(res.Id)
	assert.Nil(t, profile2)
	assert.NotNil(t, err3)
	assert.True(t, strings.Contains(err3.Error(), "NotFoundException"))
}

func TestIntegration_Profiles_GetProfile(t *testing.T) {

	gateway := createGateway()
	request := beanstream.Profile{
		Card: beanstream.CreditCard{
			Name:        "John Doe",
			Number:      "5100000010001004",
			ExpiryMonth: "11",
			ExpiryYear:  "19",
			Cvd:         "123"},
		BillingAddress: beanstream.Address{
			"John Doe",
			"123 Fake St.",
			"suite 3",
			"Victoria",
			"BC",
			"CA",
			"V8T4M3",
			"12505550123",
			"test@example.com"}}
	res, err := gateway.Profiles().CreateProfile(request)
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.NotEmpty(t, res.Id)

	profile, err2 := gateway.Profiles().GetProfile(res.Id)
	assert.Nil(t, err2)
	assert.NotNil(t, profile)
}

func TestIntegration_Profiles_UpdateProfile(t *testing.T) {
	gateway := createGateway()
	request := beanstream.Profile{
		Card: beanstream.CreditCard{
			Name:        "John Doe",
			Number:      "5100000010001004",
			ExpiryMonth: "11",
			ExpiryYear:  "19",
			Cvd:         "123"},
		BillingAddress: beanstream.Address{
			"John Doe",
			"123 Fake St.",
			"suite 3",
			"Victoria",
			"BC",
			"CA",
			"V8T4M3",
			"12505550123",
			"test@example.com"}}
	res, err := gateway.Profiles().CreateProfile(request)
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.NotEmpty(t, res.Id)

	// get profile
	profile, err2 := gateway.Profiles().GetProfile(res.Id)
	assert.Nil(t, err2)
	assert.NotNil(t, profile)

	// update profile
	street := "456 Dingle Bingle Road"
	profile.BillingAddress.AddressLine1 = street
	res2, err2 := gateway.Profiles().UpdateProfile(profile)
	assert.Nil(t, err2)
	assert.NotNil(t, res2)

	// get profile again
	profile2, err3 := gateway.Profiles().GetProfile(res2.Id)
	assert.Nil(t, err3)
	assert.NotNil(t, profile2)
	assert.Equal(t, street, profile2.BillingAddress.AddressLine1)

	// delete profile
	profile3, err4 := gateway.Profiles().DeleteProfile(profile2.Id)
	assert.Nil(t, err4)
	assert.NotNil(t, profile3)
}

func TestIntegration_Profiles_GetCards(t *testing.T) {
	gateway := createGateway()
	request := beanstream.Profile{
		Card: beanstream.CreditCard{
			Name:        "John Doe",
			Number:      "5100000010001004",
			ExpiryMonth: "11",
			ExpiryYear:  "19",
			Cvd:         "123"},
		BillingAddress: beanstream.Address{
			"John Doe",
			"999 Fake St.",
			"suite 3",
			"Victoria",
			"BC",
			"CA",
			"V8T4M3",
			"12505550123",
			"test@example.com"}}
	res, err := gateway.Profiles().CreateProfile(request)
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.NotEmpty(t, res.Id)

	// get profile
	profile, err2 := gateway.Profiles().GetProfile(res.Id)
	assert.Nil(t, err2)
	assert.NotNil(t, profile)

	// get cards on the profile
	cards, err3 := gateway.Profiles().GetCards(profile.Id)
	assert.Nil(t, err3)
	assert.NotNil(t, cards)
	assert.Equal(t, 1, len(cards))
	assert.Equal(t, "510000XXXXXX1004", cards[0].Number)
	assert.Equal(t, 1, cards[0].Id)
	assert.Equal(t, "DEF", cards[0].Function)

	// delete profile
	res2, err4 := gateway.Profiles().DeleteProfile(profile.Id)
	assert.Nil(t, err4)
	assert.NotNil(t, res2)
}

func TestIntegration_Profiles_AddCard(t *testing.T) {
	gateway := createGateway()
	request := beanstream.Profile{
		Card: beanstream.CreditCard{
			Name:        "John Doe",
			Number:      "5100000010001004",
			ExpiryMonth: "11",
			ExpiryYear:  "19",
			Cvd:         "123"},
		BillingAddress: beanstream.Address{
			"John Doe",
			"999 Fake St.",
			"suite 3",
			"Victoria",
			"BC",
			"CA",
			"V8T4M3",
			"12505550123",
			"test@example.com"}}
	res, err := gateway.Profiles().CreateProfile(request)
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.NotEmpty(t, res.Id)

	// get profile
	profile, err2 := gateway.Profiles().GetProfile(res.Id)
	assert.Nil(t, err2)
	assert.NotNil(t, profile)

	// add a 2nd card
	card2 := beanstream.CreditCard{
		Name:        "Jane Doe",
		Number:      "4030000010001234",
		ExpiryMonth: "03",
		ExpiryYear:  "18",
		Cvd:         "123"}
	res2, err3 := gateway.Profiles().AddCard(profile.Id, card2)
	assert.Nil(t, err3)
	assert.NotNil(t, res2)

	// get cards
	cards, err4 := profile.GetCards(gateway.Profiles())
	assert.Nil(t, err4)
	assert.NotNil(t, cards)
	assert.Equal(t, 2, len(cards))
	assert.Equal(t, "510000XXXXXX1004", cards[0].Number)
	assert.Equal(t, 1, cards[0].Id)
	assert.Equal(t, "DEF", cards[0].Function)
	assert.Equal(t, "403000XXXXXX1234", cards[1].Number)
	assert.Equal(t, 2, cards[1].Id)
	assert.Equal(t, "SEC", cards[1].Function)

	// delete profile
	res3, err5 := gateway.Profiles().DeleteProfile(profile.Id)
	assert.Nil(t, err5)
	assert.NotNil(t, res3)
}

func TestIntegration_Profiles_DeleteCard(t *testing.T) {
	gateway := createGateway()
	request := beanstream.Profile{
		Card: beanstream.CreditCard{
			Name:        "John Doe",
			Number:      "5100000010001004",
			ExpiryMonth: "11",
			ExpiryYear:  "19",
			Cvd:         "123"},
		BillingAddress: beanstream.Address{
			"John Doe",
			"999 Fake St.",
			"suite 3",
			"Victoria",
			"BC",
			"CA",
			"V8T4M3",
			"12505550123",
			"test@example.com"}}
	res, err := gateway.Profiles().CreateProfile(request)
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.NotEmpty(t, res.Id)

	// get profile
	profile, err2 := gateway.Profiles().GetProfile(res.Id)
	assert.Nil(t, err2)
	assert.NotNil(t, profile)

	// add a 2nd card
	card2 := beanstream.CreditCard{
		Name:        "Jane Doe",
		Number:      "4030000010001234",
		ExpiryMonth: "03",
		ExpiryYear:  "18",
		Cvd:         "123"}
	res2, err3 := gateway.Profiles().AddCard(profile.Id, card2)
	assert.Nil(t, err3)
	assert.NotNil(t, res2)

	// get cards
	cards, err4 := profile.GetCards(gateway.Profiles())
	assert.Nil(t, err4)
	assert.NotNil(t, cards)
	assert.Equal(t, 2, len(cards))
	assert.Equal(t, "510000XXXXXX1004", cards[0].Number)
	assert.Equal(t, 1, cards[0].Id)
	assert.Equal(t, "DEF", cards[0].Function)
	assert.Equal(t, "403000XXXXXX1234", cards[1].Number)
	assert.Equal(t, 2, cards[1].Id)
	assert.Equal(t, "SEC", cards[1].Function)

	// delete card
	res3, err5 := profile.DeleteCard(gateway.Profiles(), cards[1].Id)
	assert.Nil(t, err5)
	assert.NotNil(t, res3)

	// get cards again
	cards2, err6 := profile.GetCards(gateway.Profiles())
	assert.Nil(t, err6)
	assert.NotNil(t, cards2)
	assert.Equal(t, 1, len(cards2))
	assert.Equal(t, "510000XXXXXX1004", cards2[0].Number)
	assert.Equal(t, 1, cards2[0].Id)

	// delete profile
	res4, err7 := gateway.Profiles().DeleteProfile(profile.Id)
	assert.Nil(t, err7)
	assert.NotNil(t, res4)
}

func TestIntegration_Profiles_UpdateCard(t *testing.T) {
	gateway := createGateway()
	request := beanstream.Profile{
		Card: beanstream.CreditCard{
			Name:        "John Doe",
			Number:      "5100000010001004",
			ExpiryMonth: "11",
			ExpiryYear:  "19",
			Cvd:         "123"},
		BillingAddress: beanstream.Address{
			"John Doe",
			"999 Fake St.",
			"suite 3",
			"Victoria",
			"BC",
			"CA",
			"V8T4M3",
			"12505550123",
			"test@example.com"}}
	res, err := gateway.Profiles().CreateProfile(request)
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.NotEmpty(t, res.Id)

	// get profile
	profile, err2 := gateway.Profiles().GetProfile(res.Id)
	assert.Nil(t, err2)
	assert.NotNil(t, profile)

	// get card
	card, err3 := profile.GetCard(gateway.Profiles(), 1) // the first card is always #1
	assert.Nil(t, err3)
	assert.NotNil(t, card)

	// update card
	card.ExpiryMonth = "04"
	res2, err4 := profile.UpdateCard(gateway.Profiles(), *card)
	assert.Nil(t, err4)
	assert.NotNil(t, res2)

	// get card again
	card2, err5 := profile.GetCard(gateway.Profiles(), card.Id)
	assert.Nil(t, err5)
	assert.NotNil(t, card2)
	assert.Equal(t, "04", card.ExpiryMonth)

	// delete profile
	res3, err6 := gateway.Profiles().DeleteProfile(profile.Id)
	assert.Nil(t, err6)
	assert.NotNil(t, res3)
}
