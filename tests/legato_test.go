// +build integration

package tests

import (
	"fmt"
	beanstream "github.com/Beanstream-DRWP/beanstream-go"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIntegration_Legato_TokenizeCard(t *testing.T) {
	res, err := beanstream.LegatoTokenizeCard(
		"5100000010001004",
		"11",
		"19",
		"123")
	assert.Nil(t, err, "Legato returned an error")
	assert.NotNil(t, res, "No token returned")
	assert.NotEmpty(t, res, "Legato token was empty: "+res)
	fmt.Println("Legato token for Card: ", res)
}
