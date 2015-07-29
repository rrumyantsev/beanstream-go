// +build integration

package beanstream

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIntegration_Legato_TokenizeCard(t *testing.T) {
	res, err := LegatoTokenizeCard(
		"5100000010001004",
		"11",
		"19",
		"123")
	assert.Nil(t, err, "Legato returned an error")
	assert.NotNil(t, res, "No token returned")
	assert.NotEmpty(t, res, "Legato token was empty: "+res)
	fmt.Println("Legato token for Card: ", res)
}
