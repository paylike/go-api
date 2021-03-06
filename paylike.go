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

// PricingAmount describes the currency and the amount
type PricingAmount struct {
	Currency string
	Amount   float64
}

// MerchantTransfer describes a transfer to a given card
type MerchantTransfer struct {
	ToCard Pricing
}

// Pricing describes the exact amounts for a given item
type Pricing struct {
	Rate    float64
	Flat    PricingAmount
	Dispute PricingAmount
}

// MerchantPricing describes a pricing included in the merchant
type MerchantPricing struct {
	Pricing
	Transfer MerchantTransfer
}

// MerchantTDS either "attempt" or "full" based on 3-D secure
type MerchantTDS struct {
	Mode string
}

// Merchant describes information about a given merchant
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

// MerchantClaim describes claims for a given merchant
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

// MerchantCompany describes the company of a given merchant
type MerchantCompany struct {
	Country string `json:"country"`          // required, ISO 3166 code (e.g. DK)
	Number  string `json:"number,omitempty"` // optional, registration number ("CVR" in Denmark)
}

// MerchantBank describes a bank for a given merchant
type MerchantBank struct {
	Iban string `json:"iban,omitempty"` // optional, (format: XX00000000, XX is country code, length varies)
}

// Line desccribes a given item in the history of the merchant balance
type Line struct {
	ID            string        `json:"id"`
	Created       string        `json:"created"`
	MerchantID    string        `json:"merchantId"`
	Balance       int           `json:"balance"`
	Fee           int           `json:"fee"`
	TransactionID string        `json:"transactionId"`
	Amount        PricingAmount `json:"amount"`
	Refund        bool          `json:"refund"`
	Test          bool          `json:"test"`
}

// TransactionDTO describes options in terms of the transaction
// creation API
type TransactionDTO struct {
	CardID        string                 `json:"cardId,omitempty"`        // required if no TransactionID is present
	TransactionID string                 `json:"transactionId,omitempty"` // required if no CardID is present
	Descriptor    string                 `json:"descriptor,omitempty"`    // optional, will fallback to merchant descriptor
	Currency      string                 `json:"currency"`                // required, three letter ISO
	Amount        int                    `json:"amount"`                  // required, amount in minor units
	Custom        map[string]interface{} `json:"custom,omitempty"`        // optional, any custom data

}

// TransactionID describes the ID for a given unique transaction used for referencing
type TransactionID struct {
	ID string `json:"id"`
}

// TransactionTrailDTO describes information about the the capturing / refunding / voiding amount
type TransactionTrailDTO struct {
	Amount     int    `json:"amount"`               // required, amount in minor units (100 = DKK 1,00)
	Currency   string `json:"currency,omitempty"`   // optional, expected currency (for additional verification)
	Descriptor string `json:"descriptor,omitempty"` // optional, text on client bank statement
}

// CardCode describes if a given code is present to the card or not
type CardCode struct {
	Present bool `json:"present"`
}

// TransactionCard describes card information that can be found in transactions
type TransactionCard struct {
	Bin    string   `json:"bin"`
	Last4  string   `json:"last4"`
	Expiry string   `json:"expiry"`
	Scheme string   `json:"scheme"`
	Code   CardCode `json:"code"`
}

// Transaction describes information about a given transaction
type Transaction struct {
	TransactionID
	Test           bool                   `json:"test"`
	MerchantID     string                 `json:"merchantId"`
	Created        string                 `json:"created"`
	Amount         int                    `json:"amount"`
	RefundedAmount int                    `json:"refundedAmount"`
	CapturedAmount int                    `json:"capturedAmount"`
	VoidedAmount   int                    `json:"voidedAmount"`
	PendingAmount  int                    `json:"pendingAmount"`
	DisputedAmount int                    `json:"disputedAmount"`
	Card           TransactionCard        `json:"card"`
	TDS            string                 `json:"tds"`
	Currency       string                 `json:"currency"`
	Custom         map[string]interface{} `json:"custom"`
	Recurring      bool                   `json:"recurring"`
	Successful     bool                   `json:"successful"`
	Error          bool                   `json:"error"`
	Descriptor     string                 `json:"descriptor"`
	Trail          []*TransactionTrail    `json:"trail"`
}

// TransactionTrailFee describes fee included in the given trail
type TransactionTrailFee struct {
	Flat int `json:"flat"`
	Rate int `json:"rate"`
}

// TransactionTrail describes a given trail element in the transactions
type TransactionTrail struct {
	Fee        TransactionTrailFee `json:"fee"`
	Amount     int                 `json:"amount"`
	Balance    int                 `json:"balance"`
	Created    string              `json:"created"`
	Capture    bool                `json:"captrue"`
	Descriptor string              `json:"descriptor"`
	LineID     string              `json:"lineId"`
	Dispute    TrailDispute        `json:"dispute"`
}

// TrailDispute describes a given dispute in a given trail
type TrailDispute struct {
	ID   string `json:"id"`
	Won  bool   `json:"won,omitempty"`
	Lost bool   `json:"lost,omitempty"`
}

// Card describes the full information about a given card
type Card struct {
	TransactionCard
	CardID
	MerchantID string `json:"merchantId"`
	Created    string `json:"created"`
}

// CardDTO describes required information to create a new card
type CardDTO struct {
	TransactionID string `json:"transactionId"`
	Notes         string `json:"notes"`
}

// CardID describes a given card's ID
type CardID struct {
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

// FetchApp is to fetch information about the current application
// https://api.paylike.io/me
func (c Client) FetchApp() (*Identity, error) {
	return c.fetchApp()
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

// ListTransactions lists all transactions available under the given merchantID
// https://github.com/paylike/api-docs#fetch-all-transactions
func (c Client) ListTransactions(merchantID string, limit int) ([]*Transaction, error) {
	return c.listTransactions(merchantID, limit)
}

// CaptureTransaction captures a new amount for the given transaction
// https://github.com/paylike/api-docs#capture-a-transaction
func (c Client) CaptureTransaction(transactionID string, dto TransactionTrailDTO) (*Transaction, error) {
	b, err := json.Marshal(dto)
	if err != nil {
		return nil, err
	}
	return c.captureTransaction(transactionID, bytes.NewBuffer(b))
}

// RefundTransaction refunds a given amount for the given transaction
// https://github.com/paylike/api-docs#refund-a-transaction
func (c Client) RefundTransaction(transactionID string, dto TransactionTrailDTO) (*Transaction, error) {
	b, err := json.Marshal(dto)
	if err != nil {
		return nil, err
	}
	return c.refundTransaction(transactionID, bytes.NewBuffer(b))
}

// VoidTransaction cancels a given amount completely or partially
// https://github.com/paylike/api-docs#void-a-transaction
func (c Client) VoidTransaction(transactionID string, dto TransactionTrailDTO) (*Transaction, error) {
	b, err := json.Marshal(dto)
	if err != nil {
		return nil, err
	}
	return c.voidTransaction(transactionID, bytes.NewBuffer(b))
}

// FindTransaction finds the given transaction by ID
// https://github.com/paylike/api-docs#fetch-a-transaction
func (c Client) FindTransaction(transactionID string) (*Transaction, error) {
	return c.findTransaction(transactionID)
}

// FetchCard finds the given card by ID
// https://github.com/paylike/api-docs#fetch-a-card
func (c Client) FetchCard(cardID string) (*Card, error) {
	return c.fetchCard(cardID)
}

// CreateCard saves a new record for a given card
// https://github.com/paylike/api-docs#save-a-card
func (c Client) CreateCard(merchantID string, dto CardDTO) (*CardID, error) {
	b, err := json.Marshal(dto)
	if err != nil {
		return nil, err
	}
	return c.createCard(merchantID, bytes.NewBuffer(b))
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

// fetchApp handles the underlying logic of executing the API requests
// towards the app API to get the currently used app
func (c Client) fetchApp() (*Identity, error) {
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

// listTransactions handles the underlying logic of executing the API requests
// towards the merchant API and lists all related transactions
func (c Client) listTransactions(merchantID string, limit int) ([]*Transaction, error) {
	path := fmt.Sprintf("/merchants/%s/transactions?limit=%d", merchantID, limit)
	req, err := http.NewRequest("GET", c.getURL(path), nil)
	if err != nil {
		return nil, err
	}
	var marshalled []*Transaction
	return marshalled, c.executeRequestAndMarshal(req, &marshalled)
}

// captureTransaction handles the underlying logic of executing the API requests
// towards the merchant API and captures a new amount for a given transaction
func (c Client) captureTransaction(transactionID string, body io.Reader) (*Transaction, error) {
	path := fmt.Sprintf("/transactions/%s/captures", transactionID)
	req, err := http.NewRequest("POST", c.getURL(path), body)
	if err != nil {
		return nil, err
	}
	var marshalled map[string]*Transaction
	return marshalled["transaction"], c.executeRequestAndMarshal(req, &marshalled)
}

// refundTransaction handles the underlying logic of executing the API requests
// towards the merchant API and refunds a given amount for a given transaction
func (c Client) refundTransaction(transactionID string, body io.Reader) (*Transaction, error) {
	path := fmt.Sprintf("/transactions/%s/refunds", transactionID)
	req, err := http.NewRequest("POST", c.getURL(path), body)
	if err != nil {
		return nil, err
	}
	var marshalled map[string]*Transaction
	return marshalled["transaction"], c.executeRequestAndMarshal(req, &marshalled)
}

// voidTransaction handles the underlying logic of executing the API requests
// towards the merchant API and cancels a given amount payment partially or completely
func (c Client) voidTransaction(transactionID string, body io.Reader) (*Transaction, error) {
	path := fmt.Sprintf("/transactions/%s/voids", transactionID)
	req, err := http.NewRequest("POST", c.getURL(path), body)
	if err != nil {
		return nil, err
	}
	var marshalled map[string]*Transaction
	return marshalled["transaction"], c.executeRequestAndMarshal(req, &marshalled)
}

// findTransaction handles the underlying logic of executing the API requests
// towards the merchant API and tries to search for a given transaction
func (c Client) findTransaction(transactionID string) (*Transaction, error) {
	path := fmt.Sprintf("/transactions/%s", transactionID)
	req, err := http.NewRequest("GET", c.getURL(path), nil)
	if err != nil {
		return nil, err
	}
	var marshalled map[string]*Transaction
	return marshalled["transaction"], c.executeRequestAndMarshal(req, &marshalled)
}

// fetchCard handles the underlying logic of executing the API requests
// towards the cards API and tries to find a given card by ID
func (c Client) fetchCard(cardID string) (*Card, error) {
	path := fmt.Sprintf("/cards/%s", cardID)
	req, err := http.NewRequest("GET", c.getURL(path), nil)
	if err != nil {
		return nil, err
	}
	var marshalled map[string]*Card
	return marshalled["card"], c.executeRequestAndMarshal(req, &marshalled)
}

// createCard handles the underlying logic of executing the API requests
// towards the cards API and tries to find a given card by ID
func (c Client) createCard(merchantID string, body io.Reader) (*CardID, error) {
	path := fmt.Sprintf("/merchants/%s/cards", merchantID)
	req, err := http.NewRequest("POST", c.getURL(path), body)
	if err != nil {
		return nil, err
	}
	var marshalled map[string]*CardID
	return marshalled["card"], c.executeRequestAndMarshal(req, &marshalled)
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
	if len(b) == 0 {
		return nil
	}
	return json.Unmarshal(b, &value)
}
