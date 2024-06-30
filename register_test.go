package needle_test

import (
	"testing"

	"github.com/goplexhq/needle"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNeedle_Register_InvalidType(t *testing.T) {
	t.Cleanup(needle.Reset)

	regErr := needle.Register[int](needle.Transient)
	assert.ErrorIs(t, regErr, needle.ErrInvalidServiceType)
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

func TestNeedle_RegisterToRegistry(t *testing.T) {
	t.Cleanup(needle.Reset)

	type testStruct struct{}

	registry := needle.NewRegistry()

	err := needle.RegisterToRegistry[testStruct](registry, needle.Singleton)
	require.NoError(t, err)

	services := registry.RegisteredServices()
	assert.Len(t, services, 1)
	assert.Contains(t, services, "github.com/goplexhq/needle_test.testStruct")
}

func TestNeedle_RegisterInstanceToRegistry(t *testing.T) {
	t.Cleanup(needle.Reset)

	type testStruct struct{ name string }

	registry := needle.NewRegistry()

	instance := &testStruct{name: "myStruct"}
	regErr := needle.RegisterInstanceToRegistry(registry, instance)
	require.NoError(t, regErr)

	val, resErr := needle.ResolveFromRegistry[testStruct](registry)
	require.NoError(t, resErr)
	assert.NotNil(t, val)
	assert.Equal(t, "myStruct", val.name)
}

func TestNeedle_RegisterScoped(t *testing.T) {
	t.Cleanup(needle.Reset)

	type testStruct struct{}

	opt := needle.WithScope("my custom scope")
	err := needle.Register[testStruct](needle.Scoped, opt)

	require.NoError(t, err)

	services := needle.RegisteredServices()
	assert.Len(t, services, 1)
	assert.Contains(t, services, "github.com/goplexhq/needle_test.testStruct")
}

func TestNeedle_RegisterScoped_EmptyScope(t *testing.T) {
	t.Cleanup(needle.Reset)

	type testStruct struct{}

	err := needle.Register[testStruct](needle.Scoped)
	require.ErrorIs(t, err, needle.ErrEmptyScope)

	assert.Empty(t, needle.RegisteredServices())
}

func TestNeedle_RegisterThreadLocal(t *testing.T) {
	t.Cleanup(needle.Reset)

	type testStruct struct{}

	require.NoError(t, needle.Register[testStruct](needle.ThreadLocal))

	services := needle.RegisteredServices()
	assert.Len(t, services, 1)
	assert.Contains(t, services, "github.com/goplexhq/needle_test.testStruct")
}

func TestNeedle_RegisterThreadLocal_CustomThreadID(t *testing.T) {
	t.Cleanup(needle.Reset)

	type testStruct struct{}

	opt := needle.WithThreadID("my custom thread id")
	err := needle.Register[testStruct](needle.ThreadLocal, opt)

	require.NoError(t, err)

	services := needle.RegisteredServices()
	assert.Len(t, services, 1)
	assert.Contains(t, services, "github.com/goplexhq/needle_test.testStruct")
}
