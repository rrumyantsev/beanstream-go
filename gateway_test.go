// +build unit integration

package beanstream

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUnit_Gateway_Config_BaseUrl(t *testing.T) {
	config := Config{
		"",
		"",
		"",
		"",
		"www",
		"api",
		"v1",
		"-8:00"}
	assert.EqualValues(t, "https://www.beanstream.com/api/v1", config.BaseUrl())
}
