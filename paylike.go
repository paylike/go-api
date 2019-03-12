package paylike

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/davecgh/go-spew/spew"
)

// Client describes all information regarding the API
type Client struct {
	Key     string
	client  *http.Client
	baseAPI string
}

// App describes information about the application
type App struct {
	ID   string
	Name string
	Key  string
}

// Identity describes information about the current application that has
// been created
type Identity struct {
	ID      string
	Name    string
	Created string
}

// MerchantCreateDTO describes options for creating a merchant
type MerchantCreateDTO struct {
	Name       string           `json:"name,omitempty"` // optional, name of merchant
	Currency   string           `json:"currency"`       // required, three letter ISO
	Test       bool             `json:"test,omitempty"` // optional, defaults to false
	Email      string           `json:"email"`          // required, contact email
	Website    string           `json:"website"`        // required, website with implementation
	Descriptor string           `json:"descriptor"`     // required, text on client bank statements
	Company    *MerchantCompany `json:"company"`        // required, company information
	Bank       *MerchantBank    `json:"bank,omitempty"` // optional, bank information
}

// MerchantUpdateDTO describes options to update a given merchant
// If you cannot find your desired option here, create a new merchant instead
type MerchantUpdateDTO struct {
	Name       string `json:"name,omitempty"`       // optional, name of merchant
	Email      string `json:"email,omitempty"`      // optional, contact email
	Descriptor string `json:"descriptor,omitempty"` // optional, text on client bank statements
}

// Amount ...
type Amount struct {
	Currency string
	Amount   float64
}

// MerchantTransfer ...
type MerchantTransfer struct {
	ToCard Pricing
}

// Pricing ...
type Pricing struct {
	Rate    float64
	Flat    Amount
	Dispute Amount
}

// MerchantPricing ...
type MerchantPricing struct {
	Pricing
	Transfer MerchantTransfer
}

// MerchantTDS ...
type MerchantTDS struct {
	Mode string
}

// Merchant ...
type Merchant struct {
	ID         string
	Name       string
	Company    MerchantCompany
	Claim      MerchantClaim
	Pricing    MerchantPricing
	Currency   string
	Email      string
	TDS        MerchantTDS
	Key        string
	Bank       MerchantBank
	Created    string
	Test       bool
	Descriptor string
	Website    string
	Balance    float64
}

// MerchantClaim ...
type MerchantClaim struct {
	CanChargeCard     bool
	CanSaveCard       bool
	CanTransferToCard bool
	CanCapture        bool
	CanRefund         bool
	CanVoid           bool
}

// MerchantCompany ...
type MerchantCompany struct {
	Country string `json:"country"`          // required, ISO 3166 code (e.g. DK)
	Number  string `json:"number,omitempty"` // optional, registration number ("CVR" in Denmark)
}

// MerchantBank ...
type MerchantBank struct {
	Iban string `json:"iban,omitempty"` // optional, (format: XX00000000, XX is country code, length varies)
}

// NewClient creates a new client
func NewClient(key string) *Client {
	return &Client{key, &http.Client{}, "https://api.paylike.io"}
}

// SetKey provides an elegent way to deal with
// setting the key and calling other methods after that
func (c *Client) SetKey(key string) *Client {
	c.Key = key
	return c
}

// CreateApp creates a new application
// https://github.com/paylike/api-docs#create-an-app
func (c Client) CreateApp() (*App, error) {
	return c.createApp(nil)
}

// CreateAppWithName creates a new application with the given name
// https://github.com/paylike/api-docs#create-an-app
func (c Client) CreateAppWithName(name string) (*App, error) {
	return c.createApp(
		bytes.NewBuffer([]byte(fmt.Sprintf(`{"name":"%s"}`, name))),
	)
}

// GetCurrentApp is to fetch information about the current application
// https://api.paylike.io/me
func (c Client) GetCurrentApp() (*Identity, error) {
	return c.getCurrentApp()
}

// CreateMerchant creates a new merchant under a given app
// https://github.com/paylike/api-docs#create-a-merchant
func (c Client) CreateMerchant(dto MerchantCreateDTO) (*Merchant, error) {
	b, err := json.Marshal(dto)
	if err != nil {
		return nil, err
	}
	return c.createMerchant(bytes.NewBuffer(b))
}

// GetMerchant gets a merchant based on it's ID
// https://github.com/paylike/api-docs#fetch-a-merchant
func (c Client) GetMerchant(id string) (*Merchant, error) {
	return c.getMerchant(id)
}

// FetchMerchants fetches all merchants for given app ID
// https://github.com/paylike/api-docs#fetch-all-merchants
func (c Client) FetchMerchants(appID string, limit int) ([]*Merchant, error) {
	return c.fetchMerchants(appID, limit)
}

// UpdateMerchant updates a merchant with given parameters
// https://github.com/paylike/api-docs#update-a-merchant
func (c Client) UpdateMerchant(id string, dto MerchantUpdateDTO) error {
	b, err := json.Marshal(dto)
	if err != nil {
		return err
	}
	return c.updateMerchant(id, bytes.NewBuffer(b))
}

// getURL is to build the base API url along with the given dynamic route path
func (c Client) getURL(url string) string {
	return fmt.Sprintf("%s%s", c.baseAPI, url)
}

// createApp handles the underlying logic of executing the API requests
// towards the app creation API
func (c Client) createApp(body io.Reader) (*App, error) {
	req, err := http.NewRequest("POST", c.getURL("/apps"), body)
	if err != nil {
		return nil, err
	}
	var marshalled map[string]*App
	err = c.executeRequestAndMarshal(req, &marshalled)
	return marshalled["app"], err
}

// getCurrentApp handles the underlying logic of executing the API requests
// towards the app API to get the currently used app
func (c Client) getCurrentApp() (*Identity, error) {
	req, err := http.NewRequest("GET", c.getURL("/me"), nil)
	if err != nil {
		return nil, err
	}
	var marshalled map[string]*Identity
	err = c.executeRequestAndMarshal(req, &marshalled)
	return marshalled["identity"], err
}

// createMerchant handles the underlying logic of executing the API requests
// towards the merchant creation API
func (c Client) createMerchant(body io.Reader) (*Merchant, error) {
	req, err := http.NewRequest("POST", c.getURL("/merchants"), body)
	if err != nil {
		return nil, err
	}
	var marshalled map[string]*Merchant
	err = c.executeRequestAndMarshal(req, &marshalled)
	return marshalled["merchant"], err
}

// fetchMerchants handles the underlying logic of executing the API requests
// towards the merchant fetching API
func (c Client) fetchMerchants(appID string, limit int) ([]*Merchant, error) {
	path := fmt.Sprintf("/identities/%s/merchants?limit=%d", appID, limit)
	req, err := http.NewRequest("GET", c.getURL(path), nil)
	if err != nil {
		return nil, err
	}
	var marshalled []*Merchant
	return marshalled, c.executeRequestAndMarshal(req, &marshalled)
}

// getMerchant handles the underlying logic of executing the API requests
// towards the merchant API and gets a merchant based on it's ID
func (c Client) getMerchant(id string) (*Merchant, error) {
	path := fmt.Sprintf("/merchants/%s", id)
	req, err := http.NewRequest("GET", c.getURL(path), nil)
	if err != nil {
		return nil, err
	}
	var marshalled map[string]*Merchant
	err = c.executeRequestAndMarshal(req, &marshalled)
	return marshalled["merchant"], err
}

// updateMerchant handles the underlying logic of executing the API requests
// towards the merchant API and updates a given merchant
func (c Client) updateMerchant(id string, body io.Reader) error {
	path := fmt.Sprintf("/merchants/%s", id)
	req, err := http.NewRequest("PUT", c.getURL(path), body)
	if err != nil {
		return nil
	}
	return c.executeRequestAndMarshal(req, nil)
}

// executeRequestAndMarshal sets the correct headers, then executes the request and tries to marshal
// the response from the body into the given interface{} value
func (c Client) executeRequestAndMarshal(req *http.Request, value interface{}) error {
	req.SetBasicAuth("", c.Key)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	// Skip marshaling if the interface value is nil to avoid
	// unnecessary errors
	if value == nil {
		return nil
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	// c.exploreAPI(b)
	return json.Unmarshal(b, &value)
}

// temporary function
func (c Client) exploreAPI(b []byte) {
	var t interface{}
	if err := json.Unmarshal(b, &t); err != nil {
		log.Fatal(err)
	}
	spew.Dump(t)
}
