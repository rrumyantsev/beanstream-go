package beanstream

/*
This feature is not yet implemented!
*/

import (
	//"github.com/Beanstream/beanstream-go/httpMethods"
	"time"
)

const batchUrl = "/batchpayments"

/*
 */
type BatchAPI struct {
	Config Config
}

/*
func (api BatchAPI) MakeBatchPayment(batchCriteria BatchCriteria, batchFile string) (*PaymentResponse, error) {
	batchCriteria.batchType = 0 // 0 == standard batch, 1 == direct deposit
	batchCriteria.processDate = batchCriteria.ProcessDate.Format("20060102")

	url := api.Config.BaseUrl() + batchUrl
	responseType := PaymentResponse{}
	res, err := ProcessMultiPart(httpMethods.POST, url, api.Config.MerchantId, api.Config.BatchApiKey, &responseType, batchCriteria, batchFile)
	if err != nil {
		return nil, err
	}
	//fmt.Printf("MakePayment result: %T %v\n", res, res)
	pr := res.(*PaymentResponse)
	pr.CreatedTime = AsDate(pr.created, api.Config)
	return pr, nil
}*/

//func (api BatchAPI) MakeDirectDeposit(transaction interface{}) (*PaymentResponse, error) {
//	batchCriteria.batchType = 1 // 0 == standard batch, 1 == direct deposit

type BatchCriteria struct {
	ProcessDate time.Time
	processDate string `json:"process_date,omitempty"`
	batchType   int    `json:"batch_type,omitempty"`
	ProcessNow  bool   `json:"process_now,omitempty"`
}
