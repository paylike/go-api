# Paylike client (Go)

Writing your own client? Checkout the raw [HTTP service](https://github.com/paylike/api-docs).

**Make sure to [subscribe to our mailling list](http://eepurl.com/bCGmg1) for
deprecation notices, API changes and new features**

[Godoc for the project](https://godoc.org/github.com/paylike/go-api)

## Getting an API key

An API key can be obtained by creating a merchant and adding an app through
our [dashboard](https://app.paylike.io). If your app's target audience is
third parties, please reach out and we will make your app's API key hidden.

## Install

```shell
dep ensure -add github.com/paylike/node-api
```

```golang
import paylike "github.com/paylike/go-api"

client := paylike.NewClient(os.Getenv("PAYLIKE_APP_KEY"))
```

## Methods

```golang
// change key for authentication
client.SetKey("key")

// this command is also chainable
app, err := client.SetKey("key").FetchApp()

// create an app (requires no authentication)
createdApp, err := client.CreateApp()

// create an app with a dedicated name
createdAppWithName, err := client.CreateAppWithName("myApplication")

// fetch current app (based on key)
app, err := client.FetchApp()

// list app's merchants with limit
merchants, err := client.FetchMerchants("appID", 10)

// create merchant
merchant, err := client.CreateMerchant(MerchantCreateDTO{
    Test:       true,
    Currency:   "HUF",
    Email:      TestEmail,
    Website:    TestSite,
    Descriptor: "1234567897891234",
    Company: &MerchantCompany{
        Country: "HU",
    },
})

// update merchant
err := client.UpdateMerchant(MerchantUpdateDTO{
    Name:       "Test",
    Descriptor: "Test",
    Email:      "test@test.com",
})

// get merchant
fetchedMerchant, err := client.GetMerchant(merchant.ID)

// add users
data, err := client.InviteUserToMerchant(merchant.ID, "test@test.com")

// revoke users
err := client.RevokeUserFromMerchant(merchant.ID, users[0].ID)

// fetch users with limit
users, err := client.FetchUsersToMerchant(merchant.ID, 3)

// add app
err := client.AddAppToMerchant(merchant.ID, app.ID)

// revoke app
err := client.RevokeAppFromMerchant(merchant.ID, app.ID)

// fetch apps with limit
apps, err := client.FetchAppsToMerchant(merchant.ID, 2)

// fetch lines with limit
lines, err := client.FetchLinesToMerchant(merchant.ID, 1)

// create transaction
data, err := client.CreateTransaction(TestMerchant, TransactionDTO{
    TransactionID: "560fd96b7973ff3d2362a78c",
    Currency:      "EUR",
    Amount:        200,
    Custom:        map[string]interface{}{"source": "test"},
})

// fetch transactions with limit
transactions, err := client.ListTransactions(merchant.ID, 20)

// transaction capture
dto := TransactionTrailDTO{
    Amount:     2,
    Currency:   "EUR",
    Descriptor: "Testing",
}
transaction, err := client.CaptureTransaction(transaction.ID, dto)

// transaction refund
dto := TransactionTrailDTO{
    Amount:     1,
    Descriptor: "Testing Refund",
}
transaction, err := client.RefundTransaction(data.ID, dto)

// transaction void
dto := TransactionTrailDTO{
    Amount: 1,
}
transaction, err := client.VoidTransaction(data.ID, dto)

// transaction find
transaction, err := client.FindTransaction(data.ID)

// card create
dto := CardDTO{
    TransactionID: "560fd96b7973ff3d2362a78c",
}
data, err := client.CreateCard(TestMerchant, dto)

// card find
card, err := client.FetchCard(data.ID)
```

A webshop would typically need only `CaptureTransaction`, `RefundTransaction` and `VoidTransaction`. Some might
as well use `ListTransactions` and for recurring subscriptions
`CreateTransaction`.
