package internal_test

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/goplexhq/needle/internal"
)

func TestServiceName(t *testing.T) {
	type testStruct struct{}

	name := internal.ServiceName(reflect.TypeOf(testStruct{}))
	assert.Equal(t, "github.com/goplexhq/needle/internal_test.testStruct", name)
	assert.Equal(t, "", internal.ServiceName(reflect.TypeOf(struct{}{})))
	assert.Equal(t, "", internal.ServiceName(reflect.TypeOf(nil)))
}

func TestIsStructType(t *testing.T) {
	assert.True(t, internal.IsStructType(reflect.TypeOf(struct{}{})))
	assert.False(t, internal.IsStructType(reflect.TypeOf([]string{})))
	assert.False(t, internal.IsStructType(reflect.TypeOf(map[string]any{})))
	assert.False(t, internal.IsStructType(reflect.TypeOf(0)))
	assert.False(t, internal.IsStructType(reflect.TypeOf(nil)))
}
