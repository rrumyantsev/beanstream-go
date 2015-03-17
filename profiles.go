package beanstream

import (
	"beanstream/httpMethods"
	"fmt"
	"time"
)

const profilesBaseUrl = "/profiles"
const profileUrl = profilesBaseUrl + "/%v"

type ProfilesAPI struct {
	Config Config
}

func (api ProfilesAPI) CreateProfile(profile Profile) (*ProfileResponse, error) {
	url := api.Config.BaseUrl() + profilesBaseUrl
	responseType := ProfileResponse{}
	res, err := ProcessBody(httpMethods.POST, url, api.Config.MerchantId, api.Config.ProfilesApiKey, profile, &responseType)
	if err != nil {
		return nil, err
	}
	//fmt.Printf("CreateProfile result: %T %v\n", res, res)
	pr := res.(*ProfileResponse)
	return pr, nil
}

func (api ProfilesAPI) GetProfile(profileId string) (*Profile, error) {
	url := api.Config.BaseUrl() + profileUrl
	url = fmt.Sprintf(url, profileId)

	responseType := Profile{}
	res, err := Process(httpMethods.GET, url, api.Config.MerchantId, api.Config.ProfilesApiKey, &responseType)
	if err != nil {
		return nil, err
	}
	//fmt.Printf("GetProfile result: %T %v\n", res, res)
	pr := res.(*Profile)
	pr.ModifiedDate = api.AsDate(pr.modified)
	pr.Id = profileId
	return pr, nil
}

func (api ProfilesAPI) UpdateProfile(profile *Profile) (*ProfileResponse, error) {
	url := api.Config.BaseUrl() + profileUrl
	url = fmt.Sprintf(url, profile.Id)
	profile.Card = CreditCard{} // do not update cards here. To modify the cards use UpdateCard
	profile.Token = Token{}     // can only create a profile with a token, cannot update the token

	responseType := ProfileResponse{}
	res, err := ProcessBody(httpMethods.PUT, url, api.Config.MerchantId, api.Config.ProfilesApiKey, profile, &responseType)
	if err != nil {
		return nil, err
	}
	//fmt.Printf("UpdateProfile result: %T %v\n", res, res)
	pr := res.(*ProfileResponse)
	return pr, nil
}

func (api ProfilesAPI) DeleteProfile(profileId string) (*ProfileResponse, error) {
	url := api.Config.BaseUrl() + profileUrl
	url = fmt.Sprintf(url, profileId)

	responseType := ProfileResponse{}
	res, err := Process(httpMethods.DELETE, url, api.Config.MerchantId, api.Config.ProfilesApiKey, &responseType)
	if err != nil {
		return nil, err
	}
	//fmt.Printf("GetProfile result: %T %v\n", res, res)
	pr := res.(*ProfileResponse)
	return pr, nil
}

type Profile struct {
	Id              string
	Card            CreditCard   `json:"card,omitempty"`
	Token           Token        `json:"token,omitempty"`
	BillingAddress  Address      `json:"billing,omitempty"`
	Custom          CustomFields `json:"custom,omitempty"`
	Language        string       `json:"language,omitempty"`
	Comment         string       `json:"comment,omitempty"`
	modified        string       `json:"modified_date,omitempty"`
	LastTransaction string       `json:"last_transaction,omitempty"`
	Status          string       `json:"status,omitempty"`
	ModifiedDate    time.Time
}

type ProfileResponse struct {
	Id      string `json:"customer_code,omitempty"`
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

func (p ProfileResponse) String() string {
	return fmt.Sprintf("Profile Id: %v", p.Id)
}

func (api ProfilesAPI) AsDate(val string) time.Time {
	rfc3339Time := val + "Z" + api.Config.TimezoneOffset
	t, _ := time.Parse(time.RFC3339, rfc3339Time)
	return t
}
