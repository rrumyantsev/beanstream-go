// +build integration-ignore

package tests

import (
	//beanstream "github.com/Beanstream-DRWP/beanstream-go"
	//	"beanstream"
	//"github.com/Beanstream-DRWP/beanstream-go/paymentMethods"
	//	"github.com/stretchr/testify/assert"
	//"strings"
	"testing"
	//	"time"
)

func TestIntegration_Batch_MakePayment(t *testing.T) {
	gateway := createGateway()
	gateway.Config.MerchantId = "334350000"
	gateway.Config.UrlPrefix = "qa1"
	gateway.Config.PaymentsApiKey = "beanstream"
	gateway.Config.ProfilesApiKey = "beanstream"
	gateway.Config.ReportingApiKey = "beanstream"
	gateway.Config.BatchApiKey = "beanstream"

	/*request := beanstream.BatchCriteria{
		ProcessDate: time.Now(),
		ProcessNow:  true}
	res, err := gateway.Batch().MakeBatchPayment(request, "./batch.csv")
	assert.Nil(t, err, "Unexpected error occurred.", err)
	assert.NotNil(t, res, "Result was nil")
	*/
}
