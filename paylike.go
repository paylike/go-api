package paylike

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
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
	Name       string           `json:"name,omitempty"` // optional
	Currency   string           `json:"currency"`       // required, three letter ISO
	Test       bool             `json:"test,omitempty"` // optional, defaults to false
	Email      string           `json:"email"`          // required, contact email
	Website    string           `json:"website"`        // required, website with implementation
	Descriptor string           `json:"descriptor"`     // required, text on client bank statements
	Company    *MerchantCompany `json:"company"`        // required
	Bank       *MerchantBank    `json:"bank,omitempty"` // optional
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

// CreateMerchant ...
// https://github.com/paylike/api-docs#create-a-merchant
func (c Client) CreateMerchant(dto MerchantCreateDTO) (*Merchant, error) {
	b, err := json.Marshal(dto)
	if err != nil {
		return nil, err
	}
	return c.createMerchant(bytes.NewBuffer(b))
}

// getURL is to build the base API url along with the given dynamic route path
func (c Client) getURL(url string) string {
	return fmt.Sprintf("%s%s", c.baseAPI, url)
}

// createApp handles the API calls towards Paylike API
func (c Client) createApp(body io.Reader) (*App, error) {
	resp, err := c.client.Post(c.getURL("/apps"), "application/json", body)
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var marshalled map[string]*App
	if err := json.Unmarshal(b, &marshalled); err != nil {
		return nil, err
	}
	return marshalled["app"], nil
}

// getCurrentApp executes the request for fetching the current application
// along with marshalling the response
func (c Client) getCurrentApp() (*Identity, error) {
	req, err := http.NewRequest("GET", c.getURL("/me"), nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth("", c.Key)
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var marshalled map[string]*Identity
	if err := json.Unmarshal(b, &marshalled); err != nil {
		return nil, err
	}
	return marshalled["identity"], nil
}

// createMerchant handles the underlying logic of executing the API requests
// towards the merchant creation API
func (c Client) createMerchant(body io.Reader) (*Merchant, error) {
	req, err := http.NewRequest("POST", c.getURL("/merchants"), body)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth("", c.Key)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var marshalled map[string]*Merchant
	if err := json.Unmarshal(b, &marshalled); err != nil {
		return nil, err
	}
	return marshalled["merchant"], nil
}
