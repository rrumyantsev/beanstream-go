package beanstream

const bicUrlProtocol = "https://"
const bicUrl = "beanstream.com"

type Config struct {
	MerchantId      string
	PaymentsApiKey  string
	ProfilesApiKey  string
	ReportingApiKey string
	UrlPrefix       string //"www"
	UrlApi          string //"api"
	UrlApiVersion   string //"v1"
	TimezoneOffset  string //eg -8:00
}

func (v Config) BaseUrl() string {
	// https://www.beanstream.com/api/v1
	return bicUrlProtocol + v.UrlPrefix + "." + bicUrl + "/" + v.UrlApi + "/" + v.UrlApiVersion
}
