// +build unit

package tests

import (
	"beanstream"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Transaction_GenerateAuthCode(t *testing.T) {
	auth := beanstream.GenerateAuthCode("300200578", "4BaD82D9197b4cc4b70a221911eE9f70")
	assert.EqualValues(t, auth, "MzAwMjAwNTc4OjRCYUQ4MkQ5MTk3YjRjYzRiNzBhMjIxOTExZUU5Zjcw")
}
