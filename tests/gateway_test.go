// +build unit integration

package tests

import (
	beanstream "github.com/Beanstream-DRWP/beanstream-go"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUnit_Gateway_Config_BaseUrl(t *testing.T) {
	config := beanstream.Config{
		"",
		"",
		"",
		"",
		"www",
		"api",
		"v1",
		"-8:00"}
	assert.EqualValues(t, config.BaseUrl(), "https://www.beanstream.com/api/v1")
}
