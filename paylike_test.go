package paylike

import (
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
