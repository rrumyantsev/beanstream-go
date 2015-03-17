package beanstream

import (
	"beanstream/httpMethods"
	"fmt"
	"strings"
	"time"
)

const profilesBaseUrl = "/profiles"
const profileUrl = profilesBaseUrl + "/%v"
const cardsBaseUrl = profileUrl + "/cards"
const cardUrl = cardsBaseUrl + "/%v"

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

func (api ProfilesAPI) GetCards(profileId string) ([]CreditCard, error) {
	url := api.Config.BaseUrl() + cardsBaseUrl
	url = fmt.Sprintf(url, profileId)

	responseType := profileCardsResponse{}
	res, err := Process(httpMethods.GET, url, api.Config.MerchantId, api.Config.ProfilesApiKey, &responseType)
	if err != nil {
		return nil, err
	}
	//fmt.Printf("GetCards result: %T %v\n", res, res)
	pr := res.(*profileCardsResponse)
	return pr.Cards, nil
}

func (api ProfilesAPI) GetCard(profileId string, cardId int) (*CreditCard, error) {
	url := api.Config.BaseUrl() + cardsBaseUrl
	url = fmt.Sprintf(url, profileId)

	responseType := profileCardsResponse{}
	res, err := Process(httpMethods.GET, url, api.Config.MerchantId, api.Config.ProfilesApiKey, &responseType)
	if err != nil {
		return nil, err
	}
	//fmt.Printf("GetCard result: %T %v\n", res, res)
	pr := res.(*profileCardsResponse)
	if pr.Cards == nil || cardId < 1 || cardId > len(pr.Cards) {
		return nil, &BeanstreamApiException{400, 0, 0, "cardId not in the range of available cards!", "", nil}
	}

	return &pr.Cards[cardId-1], nil
}

func (api ProfilesAPI) AddCard(profileId string, card CreditCard) (*ProfileResponse, error) {
	url := api.Config.BaseUrl() + cardsBaseUrl
	url = fmt.Sprintf(url, profileId)

	wrapper := cardWrapper{card}
	responseType := ProfileResponse{}
	res, err := ProcessBody(httpMethods.POST, url, api.Config.MerchantId, api.Config.ProfilesApiKey, &wrapper, &responseType)
	if err != nil {
		return nil, err
	}
	fmt.Printf("AddCard result: %T %v\n", res, res)
	pr := res.(*ProfileResponse)

	return pr, nil
}

func (api ProfilesAPI) DeleteCard(profileId string, cardId int) (*ProfileResponse, error) {
	url := api.Config.BaseUrl() + cardUrl
	url = fmt.Sprintf(url, profileId, cardId)

	responseType := ProfileResponse{}
	res, err := Process(httpMethods.DELETE, url, api.Config.MerchantId, api.Config.ProfilesApiKey, &responseType)
	if err != nil {
		return nil, err
	}
	//fmt.Printf("GetCard result: %T %v\n", res, res)
	pr := res.(*ProfileResponse)
	return pr, nil
}

func (api ProfilesAPI) UpdateCard(profileId string, card CreditCard) (*ProfileResponse, error) {
	url := api.Config.BaseUrl() + cardUrl
	url = fmt.Sprintf(url, profileId, card.Id)

	if strings.Contains(card.Number, "X") || strings.Contains(card.Number, "x") {
		card.Number = "" // do not save masked card numbers
	}

	wrapper := cardWrapper{card}
	responseType := ProfileResponse{}
	res, err := ProcessBody(httpMethods.PUT, url, api.Config.MerchantId, api.Config.ProfilesApiKey, &wrapper, &responseType)
	if err != nil {
		return nil, err
	}
	//fmt.Printf("GetCard result: %T %v\n", res, res)
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

func (p *Profile) GetCards(pAPI ProfilesAPI) ([]CreditCard, error) {
	return pAPI.GetCards(p.Id)
}

func (p *Profile) GetCard(pAPI ProfilesAPI, cardId int) (*CreditCard, error) {
	return pAPI.GetCard(p.Id, cardId)
}

func (p *Profile) AddCard(pAPI ProfilesAPI, card CreditCard) (*ProfileResponse, error) {
	return pAPI.AddCard(p.Id, card)
}

func (p *Profile) UpdateCard(pAPI ProfilesAPI, card CreditCard) (*ProfileResponse, error) {
	return pAPI.UpdateCard(p.Id, card)
}

func (p *Profile) DeleteCard(pAPI ProfilesAPI, cardId int) (*ProfileResponse, error) {
	return pAPI.DeleteCard(p.Id, cardId)
}

// used internally when adding a new card to a profile
type cardWrapper struct {
	Card CreditCard `json:"card"`
}

type ProfileResponse struct {
	Id      string `json:"customer_code,omitempty"`
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

func (p ProfileResponse) String() string {
	return fmt.Sprintf("Profile Id: %v", p.Id)
}

type profileCardsResponse struct {
	Code         int          `json:"code,omitempty"`
	Message      string       `json:"message,omitempty"`
	CustomerCode string       `json:"customer_code,omitempty"`
	Cards        []CreditCard `json:"card,omitempty"`
}

func (api ProfilesAPI) AsDate(val string) time.Time {
	rfc3339Time := val + "Z" + api.Config.TimezoneOffset
	t, _ := time.Parse(time.RFC3339, rfc3339Time)
	return t
}
