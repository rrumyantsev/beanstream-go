// beanstream project beanstream.go
package beanstream

import ()

type Gateway struct {
	Config Config
}

//PaymentsAPI returns a new beanstream.PaymentsAPI type struct with the config set.
func (v *Gateway) Payments() PaymentsAPI {
	api := PaymentsAPI{v.Config}

	return api
}

//ProfilesAPI returns a new beanstream.ProfilesAPI type struct with the config set.
func (v *Gateway) Profiles() ProfilesAPI {
	api := ProfilesAPI{v.Config}

	return api
}
