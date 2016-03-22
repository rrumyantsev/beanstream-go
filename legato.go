package beanstream

import (
	"github.com/Beanstream/beanstream-go/httpMethods"
)

const url = "https://www.beanstream.com/scripts/tokenization/tokens"

type legatoCardRequest struct {
	Number       string `json:"number"`
	Expiry_month string `json:"expiry_month"`
	Expiry_year  string `json:"expiry_year"`
	Cvd          string `json:"cvd"`
}

type legatoTokenResponse struct {
	Token   string `json:"token"`
	Code    int    `json:"code"`
	Version int    `json:"version"`
	Message string `json:"message"`
}

// Turn a credit card into a single-use token.
// This should not be used from a production environment. The point
// of using a token is to not have the credit card info go to your server,
// thus increasing the scope of your PCI compliance. The token should be
// collected on the client-side app.
func LegatoTokenizeCard(cardNumber string, expMo string, expYr string, cvd string) (string, error) {
	req := legatoCardRequest{cardNumber, expMo, expYr, cvd}
	responseType := legatoTokenResponse{}
	res, err := ProcessBody(httpMethods.POST, url, "", "", req, &responseType)
	if err != nil {
		return "", err
	}
	token := res.(*legatoTokenResponse)
	return token.Token, nil
}
