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

/*
ProfilesAPI lets you store customer information so it can be
re-used. When you create a payment profile you receive a customer code,
also known as a multi-use token. You can then use this token to make
payments. Profiles have standard CRUD operations available to them
as well as the ability to add more credit cards to the profile.
*/
type ProfilesAPI struct {
	Config Config
}

// CreateProfile Creates a new profile.
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

// GetProfile Retrieves a profile using the profile ID. This ID is returned when you create
// a profile.
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
	pr.ModifiedDate = AsDate(pr.modified, api.Config)
	pr.Id = profileId
	return pr, nil
}

// UpdateProfile Updates a profile
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

// DeleteProfile Deletes a profile
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

// GetCards gets all cards on a profile
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

// GetCard Gets a single card from a profile. Cards are indexed starting with id 1 (not zero)
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

// AddCard Add a card to a profile
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

// UpdateCard Deletes a card from a profile
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

// UpdateCard Updates a card stored on a profile. This will NOT update the card number. To update
// a card number you must remove the old card and add the new one.
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

/*
Profile stores the information needed to make purchases and provide a means
for saving this information when a customer returns to your store.

A profile can be created with a Credit Card or a single-use Legato token (thus
making the single-use token multi-use).

Profiles can have more than one card stored on them. The amount of cards has a limit
that is configurable in the Beanstream backoffice control panel of your account, located
under the Profiles section.
*/
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

// GetCards Retrieves all cards from a profile
func (p *Profile) GetCards(pAPI ProfilesAPI) ([]CreditCard, error) {
	return pAPI.GetCards(p.Id)
}

// GetCard Get a single card from a profile. Cards are indexed starting with id 1 (not zero)
func (p *Profile) GetCard(pAPI ProfilesAPI, cardId int) (*CreditCard, error) {
	return pAPI.GetCard(p.Id, cardId)
}

// AddCard Add a card to a profile
func (p *Profile) AddCard(pAPI ProfilesAPI, card CreditCard) (*ProfileResponse, error) {
	return pAPI.AddCard(p.Id, card)
}

// UpdateCard Updates a card stored on a profile. This will NOT update the card number. To update
// a card number you must remove the old card and add the new one.
func (p *Profile) UpdateCard(pAPI ProfilesAPI, card CreditCard) (*ProfileResponse, error) {
	return pAPI.UpdateCard(p.Id, card)
}

// DeleteCard Deletes a card from a profile
func (p *Profile) DeleteCard(pAPI ProfilesAPI, cardId int) (*ProfileResponse, error) {
	return pAPI.DeleteCard(p.Id, cardId)
}

// used internally when adding a new card to a profile
type cardWrapper struct {
	Card CreditCard `json:"card"`
}

// ProfileResponse is the response from profile CRUD operations.
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
