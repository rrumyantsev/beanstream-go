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

// Create a config object with the default details.
// You must still supply it with your merchant ID and API keys.
func DefaultConfig() Config {
	cfg := Config{
		UrlPrefix:      "www",
		UrlApi:         "api",
		UrlApiVersion:  "v1",
		TimezoneOffset: "0:00"}
	return cfg
}

func (v Config) BaseUrl() string {
	// https://www.beanstream.com/api/v1
	return bicUrlProtocol + v.UrlPrefix + "." + bicUrl + "/" + v.UrlApi + "/" + v.UrlApiVersion
}
