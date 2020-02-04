# Beanstream-GO (Legacy) [![GoDoc](http://img.shields.io/badge/godoc-reference-blue.svg)](http://godoc.org/github.com/Beanstream/beanstream-go) [![Build Status](https://travis-ci.org/Beanstream/beanstream-go.svg?branch=master)](https://travis-ci.org/Beanstream/beanstream-go)
Go lang SDK for processing payments through Beanstream

The Go Lang SDK for Beanstream lets you take payments, save payment profiles, and run reports on your transactions. It's easy to get started, just follow the steps below.

The master version of this SDK requires GoLang 1.6+. There is a v1.4+ available on this [branch](https://github.com/Beanstream/beanstream-go/tree/golang-v1.4).

# Get Started

### Step 1) Import the Code
Import the project directly from Github:
```go
import beanstream "github.com/Beanstream/beanstream-go"
```
Run Go Get (or let your IDE do it)
```
go get
```

### Step 2) Create a Payment

```go
  import (
    beanstream "github.com/Beanstream/beanstream-go"
    "github.com/Beanstream/beanstream-go/paymentMethods"
  )
  ...
  config := beanstream.DefaultConfig()
	config.MerchantId = "YOUR_MERCHANT_ID"
	config.PaymentsApiKey = "YOUR_PAYMENTS_API_KEY"
	config.ProfilesApiKey = "YOUR_PROFILES_API_KEY"
	config.ReportingApiKey = "YOUR_REPORTS_API_KEY"
	
	gateway := beanstream.Gateway{config}
	request := beanstream.PaymentRequest{
		PaymentMethod: paymentMethods.CARD,
		OrderNumber:   beanstream.Util_randOrderId(6),
		Amount:        12.99,
		Card: beanstream.CreditCard{
			Name:        "John Doe",
			Number:      "5100000010001004",
			ExpiryMonth: "11",
			ExpiryYear:  "19",
			Cvd:         "123",
			Complete:    true}}
	res, err := gateway.Payments().MakePayment(request)
```

# Want to help improve the SDK?
Whether it is a bug fix or a new feature improvement feel free to fork the project and send us pull requests. We are always excited to work with the community to improve the project. We even have code bounties as a reward for great contributions. So never hesitate to send improvements our way!

## Developer setup
If you want to help improve the project, follow the steps below to get your dev environment set up.

### 1) Checkout the source code from Github
```
git clone https://github.com/Beanstream/beanstream-go
```

### 2) Install Testify for unit testing
Run:
```
go get github.com/stretchr/testify
```
Or if you use LiteIDE you can just add the Testify import (defined below in the 
next step) and hit the G (Get) Button on the toolbar


### 3) Write Unit test cases
Import 'testing' and 'testify/assert' in your test file:
```
	import (
	  "testing"
	  "github.com/stretchr/testify/assert"
	)
```
If you are using LiteIDE you can add:
```
	// +build unit integration	
```
to the top of your unit test files. Make SURE to add a blank line after that!!
Then add 
```
-tags=unit integration
```
to the TESTARGS settings parameter.

The unit tests will now run during the unit and integration phases. Ie. all the time.
	
### 4) Write Integration tests
If you are using LiteIDE then you will want to configure it to have a separate
integration test button.

a) open [liteIde_install_dir]/share/liteide.litebuild/gosec.xml

b) add this custom tag:
```
  <custom id="IntegrationTestArgs" name="INTEGRATIONTESTARGS" value="-v -tags=integration"/>
```
c) add this action tag:
```
  <action id="TestIntegration" menu="Test" img="blue/test.png" key = "Ctrl+Shift+T" cmd="$(GO)" args="test $(INTEGRATIONTESTARGS) ./..." save="all" output="true" codec="utf-8" regex="$(ERRREGEX)" takeall="true" navigate="true"/>
```
d) save the file, close and re-open LiteIDE. You should see a TestIntegration button under T

