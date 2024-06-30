package needle_test

import (
	"testing"

	"github.com/goplexhq/needle"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNeedle_InjectStructFields(t *testing.T) {
	t.Cleanup(needle.Reset)

	type Dep struct{ name string }

	type TestStruct struct {
		Dep *Dep `needle:"inject"`
	}

	require.NoError(t, needle.RegisterSingletonInstance(&Dep{name: "myDep"}))

	var testStruct TestStruct

	require.NoError(t, needle.InjectStructFields(&testStruct))
	assert.NotNil(t, testStruct.Dep)
	assert.Equal(t, "myDep", testStruct.Dep.name)
}

func TestNeedle_InjectStructFieldsFromRegistry(t *testing.T) {
	t.Cleanup(needle.Reset)

	registry := needle.NewRegistry()

	type Dep struct{ name string }

	type TestStruct struct {
		Dep *Dep `needle:"inject"`
	}

	require.NoError(t, needle.RegisterSingletonInstanceToRegistry(registry, &Dep{name: "myDep"}))

	var testStruct TestStruct

	require.NoError(t, needle.InjectStructFieldsFromRegistry(registry, &testStruct))
	assert.NotNil(t, testStruct.Dep)
	assert.Equal(t, "myDep", testStruct.Dep.name)
}

func TestNeedle_InjectStructFieldsInvalidType(t *testing.T) {
	t.Cleanup(needle.Reset)

	var invalidType int

	require.ErrorIs(t, needle.InjectStructFields(&invalidType), needle.ErrInvalidDestType)

	type Dep struct{}

	type TestStruct struct {
		Dep Dep `needle:"inject"`
	}

	require.ErrorIs(t, needle.InjectStructFields(&TestStruct{}), needle.ErrFieldPtr) //nolint:exhaustruct
}

func TestNeedle_InjectStructFieldsNotRegistered(t *testing.T) {
	t.Cleanup(needle.Reset)

	type Dep struct{}

	type TestStruct struct{ Dep *Dep }

	var testStruct TestStruct

	require.NoError(t, needle.InjectStructFields(&testStruct))
	assert.Nil(t, testStruct.Dep)
}
