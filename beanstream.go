// beanstream project beanstream.go
package beanstream

import (
//"time"
)

type Gateway struct {
	Config Config
}

//PaymentsAPI returns a new beanstream.PaymentsAPI type struct the Config.
func (v *Gateway) Payments() PaymentsAPI {
	api := PaymentsAPI{v.Config}

	return api
}

/*type bicTime struct {
	time.Time
	f string
}

func (t *bicTime) UnmarshalJSON(b []byte) (err error) {

}*/
