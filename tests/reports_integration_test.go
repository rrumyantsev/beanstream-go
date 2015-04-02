// +build integration

package tests

import (
	"fmt"
	beanstream "github.com/Beanstream-DRWP/beanstream-go"
	"github.com/Beanstream-DRWP/beanstream-go/fields"
	"github.com/Beanstream-DRWP/beanstream-go/operators"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
	"time"
)

func TestIntegration_Reports_TimeQuery(t *testing.T) {
	gateway := createGateway()
	request := createCardRequest()

	// make a test payment
	trans, err := gateway.Payments().MakePayment(request)
	assert.Nil(t, err)
	assert.NotNil(t, trans)

	startTime := time.Now().Add(-1 * time.Minute)
	endTime := time.Now().Add(1 * time.Minute)

	res, err2 := gateway.Reports().Query(startTime, endTime, 1, 1000)
	assert.NotNil(t, res)
	assert.Nil(t, err2)
	fmt.Printf("Report: %v", res)
	assert.True(t, len(res) > 0)
	found := false
	for _, r := range res {
		if strconv.Itoa(r.TransactionId) == trans.ID {
			found = true
		}
	}
	assert.True(t, found, "Could not find our transaction in the report!")
}

func TestIntegration_Reports_QueryCriteria(t *testing.T) {
	gateway := createGateway()

	// make a first test payment
	request1 := createCardRequest()
	request1.Amount = 100.00
	trans, err := gateway.Payments().MakePayment(request1)
	assert.Nil(t, err)
	assert.NotNil(t, trans)

	// make a second test payment
	request2 := createCardRequest()
	request2.Amount = 200.00
	trans2, err2 := gateway.Payments().MakePayment(request2)
	assert.Nil(t, err2)
	assert.NotNil(t, trans2)

	startTime := time.Now().Add(-1 * time.Minute)
	endTime := time.Now().Add(1 * time.Minute)

	criteria1 := beanstream.Criteria{
		fields.Amount,
		operators.GreaterThanEqual,
		"200.00"}
	res, err3 := gateway.Reports().Query(startTime, endTime, 1, 10, criteria1)
	assert.NotNil(t, res)
	assert.Nil(t, err3)
	fmt.Printf("Report: %v", res)
	assert.True(t, len(res) > 0)
	found := false
	for _, r := range res {
		if strconv.Itoa(r.TransactionId) == trans2.ID {
			found = true
		}
	}
	assert.True(t, found, "Could not find our transaction in the report!")
}

func TestIntegration_Reports_QueryCriteriaStringEquals(t *testing.T) {
	gateway := createGateway()

	// make a first test payment
	request1 := createCardRequest()
	request1.Amount = 100.00
	orderNum := request1.OrderNumber
	trans, err := gateway.Payments().MakePayment(request1)
	assert.Nil(t, err)
	assert.NotNil(t, trans)

	// make a second test payment
	request2 := createCardRequest()
	request2.Amount = 200.00
	trans2, err2 := gateway.Payments().MakePayment(request2)
	assert.Nil(t, err2)
	assert.NotNil(t, trans2)

	startTime := time.Now().Add(-4 * time.Hour)
	endTime := time.Now().Add(2 * time.Hour)

	fmt.Printf("Searching for order number:\n", orderNum)

	criteria1 := beanstream.Criteria{
		fields.OrderNumber,
		operators.Equals,
		orderNum}
	res, err3 := gateway.Reports().Query(startTime, endTime, 1, 2, criteria1)
	assert.NotNil(t, res)
	assert.Nil(t, err3)
	//fmt.Printf("Report: %v", res)
	assert.True(t, len(res) > 0)
	found := false
	for _, r := range res {
		if strconv.Itoa(r.TransactionId) == trans.ID {
			found = true
		}
	}
	assert.True(t, found, "Could not find our transaction in the report!")
}
