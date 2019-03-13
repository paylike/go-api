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

// InviteUserToMerchantResponse describes the response when a user
// is being invited to a given merchant
type InviteUserToMerchantResponse struct {
	IsMember bool
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

// User describes a user in the system
type User struct {
	ID    string `json:"id"`
	Email string `json:"email"`
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

// Line desccribes a given item in the history of the merchant balance
type Line struct {
	ID            string `json:"id"`
	Created       string `json:"created"`
	MerchantID    string `json:"merchantId"`
	Balance       int    `json:"balance"`
	Fee           int    `json:"fee"`
	TransactionID string `json:"transactionId"`
	Amount        Amount `json:"amount"`
	Refund        bool   `json:"refund"`
	Test          bool   `json:"test"`
}

// TransactionDTO describes options in terms of the transaction
// creation API
type TransactionDTO struct {
	TransactionID string      `json:"transactionId"`        // required
	Descriptor    string      `json:"descriptor,omitempty"` // optional, will fallback to merchant descriptor
	Currency      string      `json:"currency"`             // required, three letter ISO
	Amount        int         `json:"amount"`               // required, amount in minor units
	Custom        interface{} `json:"custom,omitempty"`     // optional, any custom data

}

// TransactionID describes the ID for a given unique transaction used for referencing
type TransactionID struct {
	ID string `json:"id"`
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

// InviteUserToMerchant invites given user to use the given merchant account
// https://github.com/paylike/api-docs#invite-user-to-a-merchant
func (c Client) InviteUserToMerchant(merchantID string, email string) (*InviteUserToMerchantResponse, error) {
	return c.inviteUserToMerchant(merchantID, email)
}

// FetchUsersToMerchant fetches users for a given merchant
// https://github.com/paylike/api-docs#fetch-all-users-on-a-merchant
func (c Client) FetchUsersToMerchant(merchantID string, limit int) ([]*User, error) {
	return c.fetchUsersToMerchant(merchantID, limit)
}

// RevokeUserFromMerchant revokes a given user from a given merchant
// https://github.com/paylike/api-docs#revoke-user-from-a-merchant
func (c Client) RevokeUserFromMerchant(merchantID string, userID string) error {
	return c.revokeUserFromMerchant(merchantID, userID)
}

// AddAppToMerchant revokes a given user from a given merchant
// https://github.com/paylike/api-docs#add-app-to-a-merchant
func (c Client) AddAppToMerchant(merchantID string, appID string) error {
	return c.addAppToMerchant(merchantID, appID)
}

// FetchAppsToMerchant fetches apps for a given merchant
// https://github.com/paylike/api-docs#fetch-all-apps-on-a-merchant
func (c Client) FetchAppsToMerchant(merchantID string, limit int) ([]*App, error) {
	return c.fetchAppsToMerchant(merchantID, limit)
}

// RevokeAppFromMerchant revokes a given app from a given merchant
// https://github.com/paylike/api-docs#revoke-app-from-a-merchant
func (c Client) RevokeAppFromMerchant(merchantID string, appID string) error {
	return c.revokeAppFromMerchant(merchantID, appID)
}

// FetchLinesToMerchant fetches the history that makes up a given merchant's balance
// https://github.com/paylike/api-docs#merchants-lines
func (c Client) FetchLinesToMerchant(merchantID string, limit int) ([]*Line, error) {
	return c.fetchLinesToMerchant(merchantID, limit)
}

// CreateTransaction creates a new transaction based on previous transaction informations
// https://github.com/paylike/api-docs#using-a-previous-transaction
func (c Client) CreateTransaction(merchantID string, dto TransactionDTO) (*TransactionID, error) {
	b, err := json.Marshal(dto)
	if err != nil {
		return nil, err
	}
	return c.createTransaction(merchantID, bytes.NewBuffer(b))
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

// inviteUserToMerchant handles the underlying logic of executing the API requests
// towards the merchant API and invites a given user in the system
// to use the given merchant
func (c Client) inviteUserToMerchant(id string, email string) (*InviteUserToMerchantResponse, error) {
	data := []byte(fmt.Sprintf(`{"email":"%s"}`, email))
	path := fmt.Sprintf("/merchants/%s/users", id)
	req, err := http.NewRequest("POST", c.getURL(path), bytes.NewBuffer(data))
	if err != nil {
		return nil, nil
	}
	var marshalled InviteUserToMerchantResponse
	err = c.executeRequestAndMarshal(req, &marshalled)
	return &marshalled, err
}

// fetchUsersToMerchant handles the underlying logic of executing the API requests
// towards the merchant API and lists all users that are related for the given merchant
func (c Client) fetchUsersToMerchant(id string, limit int) ([]*User, error) {
	path := fmt.Sprintf("/merchants/%s/users?limit=%d", id, limit)
	req, err := http.NewRequest("GET", c.getURL(path), nil)
	if err != nil {
		return nil, nil
	}
	var marshalled []*User
	return marshalled, c.executeRequestAndMarshal(req, &marshalled)
}

// revokeUserFromMerchant handles the underlying logic of executing the API requests
// towards the merchant API and revokes a given user from a given merchant
func (c Client) revokeUserFromMerchant(merchantID string, userID string) error {
	path := fmt.Sprintf("/merchants/%s/users/%s", merchantID, userID)
	req, err := http.NewRequest("DELETE", c.getURL(path), nil)
	if err != nil {
		return err
	}
	return c.executeRequestAndMarshal(req, nil)
}

// addAppToMerchant handles the underlying logic of executing the API requests
// towards the merchant API and adds the given app to the given merchant
func (c Client) addAppToMerchant(merchantID string, appID string) error {
	data := []byte(fmt.Sprintf(`{"appId":"%s"}`, appID))
	path := fmt.Sprintf("/merchants/%s/apps", merchantID)
	req, err := http.NewRequest("POST", c.getURL(path), bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	return c.executeRequestAndMarshal(req, nil)
}

// fetchAppsToMerchant handles the underlying logic of executing the API requests
// towards the merchant API and lists all apps related to the merchant
func (c Client) fetchAppsToMerchant(merchantID string, limit int) ([]*App, error) {
	path := fmt.Sprintf("/merchants/%s/apps?limit=%d", merchantID, limit)
	req, err := http.NewRequest("GET", c.getURL(path), nil)
	if err != nil {
		return nil, err
	}
	var marshalled []*App
	return marshalled, c.executeRequestAndMarshal(req, &marshalled)
}

// revokeAppFromMerchant handles the underlying logic of executing the API requests
// towards the merchant API and revokes a given app from a given merchant
func (c Client) revokeAppFromMerchant(merchantID string, appID string) error {
	path := fmt.Sprintf("/merchants/%s/apps/%s", merchantID, appID)
	req, err := http.NewRequest("DELETE", c.getURL(path), nil)
	if err != nil {
		return err
	}
	return c.executeRequestAndMarshal(req, nil)
}

// fetchLinesToMerchant handles the underlying logic of executing the API requests
// towards the merchant API and fetches all lines related to a merchant's history
func (c Client) fetchLinesToMerchant(merchantID string, limit int) ([]*Line, error) {
	path := fmt.Sprintf("/merchants/%s/lines?limit=%d", merchantID, limit)
	req, err := http.NewRequest("GET", c.getURL(path), nil)
	if err != nil {
		return nil, err
	}
	var marshalled []*Line
	return marshalled, c.executeRequestAndMarshal(req, &marshalled)
}

// createTransaction handles the underlying logic of executing the API requests
// towards the merchant API and creates a new transaction
func (c Client) createTransaction(merchantID string, body io.Reader) (*TransactionID, error) {
	path := fmt.Sprintf("/merchants/%s/transactions", merchantID)
	req, err := http.NewRequest("POST", c.getURL(path), body)
	if err != nil {
		return nil, err
	}
	var marshalled map[string]*TransactionID
	return marshalled["transaction"], c.executeRequestAndMarshal(req, &marshalled)
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
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	c.exploreResponse(resp, b)
	if len(b) == 0 {
		return nil
	}
	return json.Unmarshal(b, &value)
}

// Temporary function
func (c Client) exploreResponse(resp *http.Response, b []byte) {
	fmt.Println(resp.Status)
	fmt.Println(string(b))
}
