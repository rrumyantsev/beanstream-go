// +build unit integration

package tests

import (
	beanstream "github.com/Beanstream-DRWP/beanstream-go"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestUnit_errors_BusinessRuleException(t *testing.T) {
	err := beanstream.BeanstreamApiException{302, 0, 0, "Test error message", "Test error", nil}
	assert.True(t, strings.Contains(err.Error(), "BusinessRuleException"), "Error is not a BusinessRuleException")

	err = beanstream.BeanstreamApiException{402, 0, 0, "Test error message", "Test error", nil}
	assert.True(t, strings.Contains(err.Error(), "BusinessRuleException"), "Error is not a BusinessRuleException")
}

func TestUnit_errors_UnexpectedException(t *testing.T) {
	err := beanstream.BeanstreamApiException{-1, 0, 0, "Test error message", "Test error", nil}
	assert.True(t, strings.Contains(err.Error(), "UnexpectedException"), "Error is not an UnexpectedException")
}

func TestUnit_errors_InvalidRequestException(t *testing.T) {
	err := beanstream.BeanstreamApiException{400, 0, 0, "Test error message", "Test error", nil}
	assert.True(t, strings.Contains(err.Error(), "InvalidRequestException"), "Error is not an InvalidRequestException")

	err = beanstream.BeanstreamApiException{405, 0, 0, "Test error message", "Test error", nil}
	assert.True(t, strings.Contains(err.Error(), "InvalidRequestException"), "Error is not an InvalidRequestException")

	err = beanstream.BeanstreamApiException{415, 0, 0, "Test error message", "Test error", nil}
	assert.True(t, strings.Contains(err.Error(), "InvalidRequestException"), "Error is not an InvalidRequestException")
}

func TestUnit_errors_UnauthorizedException(t *testing.T) {
	err := beanstream.BeanstreamApiException{401, 0, 0, "Test error message", "Test error", nil}
	assert.True(t, strings.Contains(err.Error(), "UnauthorizedException"), "Error is not an UnauthorizedException")
}

func TestUnit_errors_ForbiddenException(t *testing.T) {
	err := beanstream.BeanstreamApiException{403, 0, 0, "Test error message", "Test error", nil}
	assert.True(t, strings.Contains(err.Error(), "ForbiddenException"), "Error is not a ForbiddenException")
}

func TestUnit_errors_NotFoundException(t *testing.T) {
	err := beanstream.BeanstreamApiException{404, 0, 0, "Test error message", "Test error", nil}
	assert.True(t, strings.Contains(err.Error(), "NotFoundException"), "Error is not a NotFoundException")
}

func TestUnit_errors_InternalServerException(t *testing.T) {
	err := beanstream.BeanstreamApiException{123, 0, 0, "Test error message", "Test error", nil}
	assert.True(t, strings.Contains(err.Error(), "InternalServerException"), "Error is not an InternalServerException")
}
