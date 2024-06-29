package needle_test

import (
	"testing"

	"github.com/goplexhq/needle"
	"github.com/stretchr/testify/assert"
)

func TestNeedle_ResolveSingleton(t *testing.T) {
	t.Cleanup(needle.Reset)

	type testStruct struct{ name string }
	regErr := needle.Register[testStruct](needle.Singleton)
	assert.NoError(t, regErr)

	{
		val, resErr := needle.Resolve[testStruct]()
		assert.NoError(t, resErr)
		assert.NotNil(t, val)
		val.name = "new name" // modify the singleton instance
	}

	{
		val, resErr := needle.Resolve[testStruct]()
		assert.NoError(t, resErr)
		assert.NotNil(t, val)
		assert.Equal(t, "new name", val.name) // verify the modification exists
	}
}

func TestNeedle_ResolveTransient(t *testing.T) {
	t.Cleanup(needle.Reset)

	type testStruct struct{ name string }

	regErr := needle.Register[testStruct](needle.Transient)
	assert.NoError(t, regErr)

	{
		val, resErr := needle.Resolve[testStruct]()
		assert.NoError(t, resErr)
		assert.NotNil(t, val)
		val.name = "new name" // modify the transient instance
	}

	{
		val, resErr := needle.Resolve[testStruct]()
		assert.NoError(t, resErr)
		assert.NotNil(t, val)
		assert.Equal(t, "", val.name) // verify new instance
	}
}

func TestNeedle_ResolveFromStore(t *testing.T) {
	t.Cleanup(needle.Reset)

	type testStruct struct{ name string }
	store := needle.NewStore()

	err := needle.RegisterInstanceToStore(store, &testStruct{name: "myStruct"})
	assert.NoError(t, err)

	val, resErr := needle.ResolveFromStore[testStruct](store)
	assert.NoError(t, resErr)
	assert.NotNil(t, val)
	assert.Equal(t, "myStruct", val.name)
}

func TestNeedle_ResolveNotRegistered(t *testing.T) {
	t.Cleanup(needle.Reset)

	type testStruct struct{}
	val, err := needle.Resolve[testStruct]()
	assert.ErrorIs(t, err, needle.ErrNotRegistered)
	assert.Nil(t, val)
}
