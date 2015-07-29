package beanstream

import (
	//"fmt"
	"github.com/Beanstream-DRWP/beanstream-go/httpMethods"
	"strconv"
	"time"
)

const reportsBaseUrl = "/reports"

/*
The Reports API lets you search for transactions based on date ranges
and search criteria.
*/
type ReportsAPI struct {
	Config Config
}

/*
Search/Query for transactions.
Transactions must be bounded by a date range. You must also supply a startRow and an endRow
to page the results as there is a limit of 1000 results returned per query.
Finally you can supply zero or more search Criteria. These Criteria are ANDed together.

Criteria have 3 parameters: field, operator, and value. For details on these refer to the
Criteria struct's documentation.

For paging just one row, use the values: 1, 2.
Paging index starts inclusively at the first number and non-inclusively at the 2nd number:
[start,end).
The lowest paging index number is 1.
*/
func (api ReportsAPI) Query(startTime time.Time, endTime time.Time, startRow int, endRow int, criteria ...Criteria) ([]TransactionRecord, error) {
	url := api.Config.BaseUrl() + reportsBaseUrl

	q := query{
		"Search",
		startTime.Format("2006-01-02T15:04:05"),
		endTime.Format("2006-01-02T15:04:05"),
		strconv.Itoa(startRow),
		strconv.Itoa(endRow),
		criteria}
	responseType := RecordsResult{}
	res, err := ProcessBody(httpMethods.POST, url, api.Config.MerchantId, api.Config.ReportingApiKey, &q, &responseType)
	if err != nil {
		return nil, err
	}

	pr := res.(*RecordsResult)
	for _, r := range pr.Records {
		r.DateTime = AsDate(r.dateTime, api.Config)
	}
	return pr.Records, nil
}

type query struct {
	Name      string     `json:"name"`
	StartDate string     `json:"start_date"`
	EndDate   string     `json:"end_date"`
	StartRow  string     `json:"start_row"`
	EndRow    string     `json:"end_row"`
	Criteria  []Criteria `json:"criteria"`
}

/*
Criteria let you narrow down your search results. Each criteria is ANDed together.

The Field is the field on a transaction record that you are testing against.
The Operator is one of Equals, less than, Greater than, etc...
The Value is what is being compared to.

For example if you want to search for amounts less than $100 you would
set the Field as: fields.Amount
the operator as: operators.LessThan
and the value as: "100"
*/
type Criteria struct {
	Field    int    `json:"field,omitempty"`
	Operator string `json:"operator,omitempty"`
	Value    string `json:"value,omitempty"`
}

// The query result that contains an array of Transaction records
type RecordsResult struct {
	Records []TransactionRecord `json:"records,omitempty"`
}

// The transaction in a query RecordsResult
type TransactionRecord struct {
	RowId            int    `json:"row_id,omitempty"`
	TransactionId    int    `json:"trn_id,omitempty"`
	dateTime         string `json:"trn_date_time,omitempty"`
	DateTime         time.Time
	Type             string  `json:"trn_type,omitempty"`
	OrderNumber      string  `json:"trn_order_number,omitempty"`
	PaymentMethod    string  `json:"trn_payment_method,omitempty"`
	Comments         string  `json:"trn_comments,omitempty"`
	MaskedCard       string  `json:"trn_masked_card,omitempty"`
	Amount           float32 `json:"trn_amount,string,omitempty"`
	Returns          float32 `json:"trn_returns,string,omitempty"`
	Completions      float32 `json:"trn_completions,string,omitempty"`
	Voided           int     `json:"trn_voided,omitempty"`
	Response         int     `json:"trn_response,omitempty"`
	CardType         string  `json:"trn_card_type,omitempty"`
	BatchNumber      int     `json:"trn_batch_no,omitempty"`
	AvsResult        string  `json:"trn_avs_result,omitempty"`
	CvdResult        int     `json:"trn_cvd_result,omitempty"`
	CardExpiry       string  `json:"trn_card_expiry,omitempty"`
	MessageId        int     `json:"message_id,omitempty"`
	MessageText      string  `json:"message_text,omitempty"`
	CardOwner        string  `json:"trn_card_owner,omitempty"`
	IpAddress        string  `json:"trn_ip,omitempty"`
	ApprovalCode     string  `json:"trn_approval_code,omitempty"`
	Reference        int     `json:"trn_reference,omitempty"`
	BillingName      string  `json:"b_name,omitempty"`
	BillingEmail     string  `json:"b_email,omitempty"`
	BillingPhone     string  `json:"b_phone,omitempty"`
	BillingAddress1  string  `json:"b_address1,omitempty"`
	BillingAddress2  string  `json:"b_address2,omitempty"`
	BillingCity      string  `json:"b_city,omitempty"`
	BillingProvince  string  `json:"b_province,omitempty"`
	BillingPostal    string  `json:"b_postal,omitempty"`
	BillingCountry   string  `json:"b_country,omitempty"`
	ShippingName     string  `json:"s_name,omitempty"`
	ShippingEmail    string  `json:"s_email,omitempty"`
	ShippingPhone    string  `json:"s_phone,omitempty"`
	ShippingAddress1 string  `json:"s_address1,omitempty"`
	ShippingAddress2 string  `json:"s_address2,omitempty"`
	ShippingCity     string  `json:"s_city,omitempty"`
	ShippingProvince string  `json:"s_province,omitempty"`
	ShippingPostal   string  `json:"s_postal,omitempty"`
	ShippingCountry  string  `json:"s_country,omitempty"`
	Ref1             string  `json:"ref1,omitempty"`
	Ref2             string  `json:"ref2,omitempty"`
	Ref3             string  `json:"ref3,omitempty"`
	Ref4             string  `json:"ref4,omitempty"`
	ProductName      string  `json:"product_name,omitempty"`
	ProductId        string  `json:"product_id,omitempty"`
	CustomerCode     string  `json:"customer_code,omitempty"`
}
