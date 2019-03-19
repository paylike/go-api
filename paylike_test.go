package paylike

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const TestKey = "4c8453c3-8285-4984-ab72-216e324372e6"
const TestEmail = "john@example.com"
const TestSite = "https://example.com"

func TestCreateApp(t *testing.T) {
	client := NewClient(TestKey)
	app, err := client.CreateApp()
	assert.Nil(t, err)
	assert.NotEmpty(t, app)
	assert.Empty(t, app.Name)
}

func TestCreateAppWithName(t *testing.T) {
	client := NewClient(TestKey)
	app, err := client.CreateAppWithName("Macilaci")
	assert.Nil(t, err)
	assert.NotEmpty(t, app)
	assert.Equal(t, "Macilaci", app.Name)
}

func TestGetApp(t *testing.T) {
	client := NewClient(TestKey)
	app, err := client.CreateApp()
	assert.Nil(t, err)
	assert.NotEmpty(t, app)

	identity, err := client.SetKey(app.Key).GetCurrentApp()
	assert.Nil(t, err)
	assert.NotEmpty(t, identity)
}

func TestCreateMerchant(t *testing.T) {
	client := NewClient(TestKey)
	app, err := client.CreateApp()
	assert.Nil(t, err)
	assert.NotEmpty(t, app)

	dto := MerchantCreateDTO{
		Test:       true,
		Currency:   "HUF",
		Email:      TestEmail,
		Website:    TestSite,
		Descriptor: "1234567897891234",
		Company: &MerchantCompany{
			Country: "HU",
		},
	}
	merchant, err := client.SetKey(app.Key).CreateMerchant(dto)
	assert.Nil(t, err)
	assert.NotEmpty(t, merchant)
}

func TestFetchMerchants(t *testing.T) {
	client := NewClient(TestKey)
	app, err := client.CreateApp()
	assert.Nil(t, err)
	assert.NotEmpty(t, app)

	dto := MerchantCreateDTO{
		Test:       true,
		Currency:   "HUF",
		Email:      TestEmail,
		Website:    TestSite,
		Descriptor: "1234567897891234",
		Company: &MerchantCompany{
			Country: "HU",
		},
	}
	merchant, err := client.SetKey(app.Key).CreateMerchant(dto)
	assert.Nil(t, err)
	assert.NotEmpty(t, merchant)

	merchants, err := client.FetchMerchants(app.ID, 5)
	assert.Nil(t, err)
	assert.NotEmpty(t, merchants)
	assert.Equal(t, merchant, merchants[0])
}

func TestGetMerchant(t *testing.T) {
	client := NewClient(TestKey)
	app, err := client.CreateApp()
	assert.Nil(t, err)
	assert.NotEmpty(t, app)

	dto := MerchantCreateDTO{
		Test:       true,
		Currency:   "HUF",
		Email:      TestEmail,
		Website:    TestSite,
		Descriptor: "1234567897891234",
		Company: &MerchantCompany{
			Country: "HU",
		},
	}
	merchant, err := client.SetKey(app.Key).CreateMerchant(dto)
	assert.Nil(t, err)
	assert.NotEmpty(t, merchant)

	fetchedMerchant, err := client.GetMerchant(merchant.ID)
	assert.Nil(t, err)
	assert.NotEmpty(t, fetchedMerchant)
	assert.Equal(t, fetchedMerchant, merchant)
}

func TestUpdateMerchant(t *testing.T) {
	client := NewClient(TestKey)
	app, err := client.CreateApp()
	assert.Nil(t, err)
	assert.NotEmpty(t, app)

	dto := MerchantCreateDTO{
		Name:       "NotTest",
		Test:       true,
		Currency:   "HUF",
		Email:      TestEmail,
		Website:    TestSite,
		Descriptor: "1234567897891234",
		Company: &MerchantCompany{
			Country: "HU",
		},
	}
	merchant, err := client.SetKey(app.Key).CreateMerchant(dto)
	assert.Nil(t, err)
	assert.NotEmpty(t, merchant)

	updateDTO := MerchantUpdateDTO{
		Name:       "Test",
		Descriptor: "NotNumbers",
		Email:      fmt.Sprintf("not_%s", dto.Email),
	}
	err = client.UpdateMerchant(merchant.ID, updateDTO)
	assert.Nil(t, err)
	updatedMerchant, err := client.GetMerchant(merchant.ID)
	assert.Nil(t, err)
	assert.NotEmpty(t, updatedMerchant)
	assert.Equal(t, updatedMerchant.Email, updateDTO.Email)
	assert.Equal(t, updatedMerchant.Name, updateDTO.Name)
	assert.Equal(t, updatedMerchant.Descriptor, updateDTO.Descriptor)
}

func TestInviteUserToMerchant(t *testing.T) {
	client := NewClient(TestKey)
	app, err := client.CreateApp()
	assert.Nil(t, err)
	assert.NotEmpty(t, app)

	dto := MerchantCreateDTO{
		Name:       "NotTest",
		Test:       true,
		Currency:   "HUF",
		Email:      TestEmail,
		Website:    TestSite,
		Descriptor: "1234567897891234",
		Company: &MerchantCompany{
			Country: "HU",
		},
	}
	merchant, err := client.SetKey(app.Key).CreateMerchant(dto)
	assert.Nil(t, err)
	assert.NotEmpty(t, merchant)

	data, err := client.InviteUserToMerchant(merchant.ID, "one@example.com")
	assert.Nil(t, err)
	assert.False(t, data.IsMember)
}

func TestFetchUsersToMerchant(t *testing.T) {
	client := NewClient(TestKey)
	app, err := client.CreateApp()
	assert.Nil(t, err)
	assert.NotEmpty(t, app)

	dto := MerchantCreateDTO{
		Name:       "NotTest",
		Test:       true,
		Currency:   "HUF",
		Email:      TestEmail,
		Website:    TestSite,
		Descriptor: "1234567897891234",
		Company: &MerchantCompany{
			Country: "HU",
		},
	}
	merchant, err := client.SetKey(app.Key).CreateMerchant(dto)
	assert.Nil(t, err)
	assert.NotEmpty(t, merchant)

	data, err := client.InviteUserToMerchant(merchant.ID, "one@example.com")
	assert.Nil(t, err)
	assert.False(t, data.IsMember)

	users, err := client.FetchUsersToMerchant(merchant.ID, 3)
	assert.Nil(t, err)
	assert.NotEmpty(t, users)
	assert.Equal(t, "one@example.com", users[0].Email)
}

func TestRevokeUserFromMerchant(t *testing.T) {
	client := NewClient(TestKey)
	app, err := client.CreateApp()
	assert.Nil(t, err)
	assert.NotEmpty(t, app)

	dto := MerchantCreateDTO{
		Name:       "NotTest",
		Test:       true,
		Currency:   "HUF",
		Email:      TestEmail,
		Website:    TestSite,
		Descriptor: "1234567897891234",
		Company: &MerchantCompany{
			Country: "HU",
		},
	}
	merchant, err := client.SetKey(app.Key).CreateMerchant(dto)
	assert.Nil(t, err)
	assert.NotEmpty(t, merchant)

	data, err := client.InviteUserToMerchant(merchant.ID, "one@example.com")
	assert.Nil(t, err)
	assert.False(t, data.IsMember)

	users, err := client.FetchUsersToMerchant(merchant.ID, 3)
	assert.Nil(t, err)
	assert.NotEmpty(t, users)
	assert.Equal(t, "one@example.com", users[0].Email)

	err = client.RevokeUserFromMerchant(merchant.ID, users[0].ID)
	assert.Nil(t, err)

	users, err = client.FetchUsersToMerchant(merchant.ID, 3)
	assert.Nil(t, err)
	assert.Empty(t, users)
}

func TestAddAppToMerchant(t *testing.T) {
	client := NewClient(TestKey)
	app, err := client.CreateApp()
	assert.Nil(t, err)
	assert.NotEmpty(t, app)

	dto := MerchantCreateDTO{
		Name:       "NotTest",
		Test:       true,
		Currency:   "HUF",
		Email:      TestEmail,
		Website:    TestSite,
		Descriptor: "1234567897891234",
		Company: &MerchantCompany{
			Country: "HU",
		},
	}
	merchant, err := client.SetKey(app.Key).CreateMerchant(dto)
	assert.Nil(t, err)
	assert.NotEmpty(t, merchant)

	err = client.AddAppToMerchant(merchant.ID, app.ID)
	assert.Nil(t, err)
}

func TestFetchAppsToMerchant(t *testing.T) {
	client := NewClient(TestKey)
	app, err := client.CreateApp()
	assert.Nil(t, err)
	assert.NotEmpty(t, app)

	dto := MerchantCreateDTO{
		Name:       "NotTest",
		Test:       true,
		Currency:   "HUF",
		Email:      TestEmail,
		Website:    TestSite,
		Descriptor: "1234567897891234",
		Company: &MerchantCompany{
			Country: "HU",
		},
	}
	merchant, err := client.SetKey(app.Key).CreateMerchant(dto)
	assert.Nil(t, err)
	assert.NotEmpty(t, merchant)

	err = client.AddAppToMerchant(merchant.ID, app.ID)
	assert.Nil(t, err)

	apps, err := client.FetchAppsToMerchant(merchant.ID, 2)
	assert.Nil(t, err)
	assert.NotEmpty(t, apps)
	assert.Equal(t, app, apps[0])
}

func TestRevokeAppFromMerchant(t *testing.T) {
	client := NewClient(TestKey)
	app, err := client.CreateApp()
	assert.Nil(t, err)
	assert.NotEmpty(t, app)

	dto := MerchantCreateDTO{
		Name:       "NotTest",
		Test:       true,
		Currency:   "HUF",
		Email:      TestEmail,
		Website:    TestSite,
		Descriptor: "1234567897891234",
		Company: &MerchantCompany{
			Country: "HU",
		},
	}
	merchant, err := client.SetKey(app.Key).CreateMerchant(dto)
	assert.Nil(t, err)
	assert.NotEmpty(t, merchant)

	err = client.AddAppToMerchant(merchant.ID, app.ID)
	assert.Nil(t, err)

	err = client.RevokeAppFromMerchant(merchant.ID, app.ID)
	assert.Nil(t, err)

	apps, err := client.FetchAppsToMerchant(merchant.ID, 2)
	assert.Nil(t, err)
	assert.Empty(t, apps)
}

func TestFetchLinesToMerchant(t *testing.T) {
	client := NewClient(TestKey).SetKey("4ff7de37-dddf-4e51-8cc9-48b61a102923")
	lines, err := client.FetchLinesToMerchant("55006bdfe0308c4cbfdbd0e1", 1)
	assert.Nil(t, err)
	assert.NotEmpty(t, lines)
	assert.Len(t, lines, 1)
}

func TestCreateTransaction(t *testing.T) {
	client := NewClient(TestKey).SetKey("4ff7de37-dddf-4e51-8cc9-48b61a102923")
	dto := TransactionDTO{
		TransactionID: "560fd96b7973ff3d2362a78c",
		Currency:      "EUR",
		Amount:        200,
		Custom:        map[string]string{"source": "test"},
	}
	data, err := client.CreateTransaction("55006bdfe0308c4cbfdbd0e1", dto)
	assert.Nil(t, err)
	assert.NotEmpty(t, data)
}

func TestListTransactions(t *testing.T) {
	client := NewClient(TestKey).SetKey("4ff7de37-dddf-4e51-8cc9-48b61a102923")
	transactions, err := client.ListTransactions("55006bdfe0308c4cbfdbd0e1", 20)
	assert.Nil(t, err)
	assert.NotEmpty(t, transactions)
	assert.Len(t, transactions, 20)
}

func TestCaptureTransactions(t *testing.T) {
	client := NewClient(TestKey).SetKey("4ff7de37-dddf-4e51-8cc9-48b61a102923")
	transactionDTO := TransactionDTO{
		TransactionID: "560fd96b7973ff3d2362a78c",
		Currency:      "EUR",
		Amount:        200,
		Custom:        map[string]string{"source": "test"},
	}
	data, err := client.CreateTransaction("55006bdfe0308c4cbfdbd0e1", transactionDTO)
	assert.Nil(t, err)
	assert.NotEmpty(t, data)

	captureDTO := TransactionCaptureDTO{
		Amount:     2,
		Currency:   "EUR",
		Descriptor: "Testing",
	}
	transaction, err := client.CaptureTransaction(data.ID, captureDTO)
	assert.Nil(t, err)
	assert.NotEmpty(t, transaction)
	assert.Len(t, transaction.Trail, 1)
	assert.Equal(t, transaction.Trail[0].Amount, captureDTO.Amount)
	assert.Equal(t, transaction.Trail[0].Descriptor, captureDTO.Descriptor)
}
