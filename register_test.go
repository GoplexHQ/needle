package needle_test

import (
	"testing"

	"github.com/goplexhq/needle"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNeedle_RegisterInvalidType(t *testing.T) {
	t.Cleanup(needle.Reset)

	regErr := needle.Register[int](needle.Transient)
	assert.ErrorIs(t, regErr, needle.ErrInvalidType)
}

func TestNeedle_RegisterDuplicate(t *testing.T) {
	t.Cleanup(needle.Reset)

	type testStruct struct{}

	require.NoError(t, needle.Register[testStruct](needle.Transient))
	assert.ErrorIs(t, needle.Register[testStruct](needle.Transient), needle.ErrRegistered)
}

func TestNeedle_Register(t *testing.T) {
	t.Cleanup(needle.Reset)

	type testStruct struct{}

	require.NoError(t, needle.Register[testStruct](needle.Transient))

	services := needle.RegisteredServices()
	assert.Len(t, services, 1)
	assert.Contains(t, services, "github.com/goplexhq/needle_test.testStruct")
}

func TestNeedle_RegisterInstance(t *testing.T) {
	t.Cleanup(needle.Reset)

	type testStruct struct{ name string }

	regErr := needle.RegisterInstance(&testStruct{name: "myStruct"})
	require.NoError(t, regErr)

	val, resErr := needle.Resolve[testStruct]()
	require.NoError(t, resErr)
	assert.NotNil(t, val)
	assert.Equal(t, "myStruct", val.name)
}

func TestNeedle_RegisterToStore(t *testing.T) {
	t.Cleanup(needle.Reset)

	type testStruct struct{}

	store := needle.NewStore()

	err := needle.RegisterToStore[testStruct](store, needle.Singleton)
	require.NoError(t, err)

	services := store.RegisteredServices()
	assert.Len(t, services, 1)
	assert.Contains(t, services, "github.com/goplexhq/needle_test.testStruct")
}

func TestNeedle_RegisterInstanceToStore(t *testing.T) {
	t.Cleanup(needle.Reset)

	type testStruct struct{ name string }

	store := needle.NewStore()

	instance := &testStruct{name: "myStruct"}
	regErr := needle.RegisterInstanceToStore(store, instance)
	require.NoError(t, regErr)

	val, resErr := needle.ResolveFromStore[testStruct](store)
	require.NoError(t, resErr)
	assert.NotNil(t, val)
	assert.Equal(t, "myStruct", val.name)
}
