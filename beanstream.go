package beanstream

import ()

/*
Gateway is the entry point for making payments. It stores the configuration
for your merchant account, such as merchant ID and passcode, in a Config struct.

The Gateway will give you access to the 3 APIs: Payments, Profiles, and Reporting.
Each time you call one of those APIs you get a new API object. It is recommended
that you always call these methods if you are going to process payments in a
multi-threaded environment using go routines. Do not share them across threads if
possible.
*/
type Gateway struct {
	Config Config
}

//Payments returns a new beanstream.PaymentsAPI type struct with the config set.
func (v *Gateway) Payments() PaymentsAPI {
	api := PaymentsAPI{v.Config}

	return api
}

//Profiles returns a new beanstream.ProfilesAPI type struct with the config set.
func (v *Gateway) Profiles() ProfilesAPI {
	api := ProfilesAPI{v.Config}

	return api
}

//Reports returns a new beanstream.ReportsAPI type struct with the config set.
func (v *Gateway) Reports() ReportsAPI {
	api := ReportsAPI{v.Config}

	return api
}
