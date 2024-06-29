package needle_test

import (
	"testing"

	"github.com/goplexhq/needle"
	"github.com/stretchr/testify/assert"
)

func TestNeedle_InjectStructFields(t *testing.T) {
	t.Cleanup(needle.Reset)

	type Dep struct{ name string }
	type TestStruct struct{ Dep *Dep }

	assert.NoError(t, needle.RegisterInstance(&Dep{name: "myDep"}))

	var testStruct TestStruct
	assert.NoError(t, needle.InjectStructFields(&testStruct))
	assert.NotNil(t, testStruct.Dep)
	assert.Equal(t, "myDep", testStruct.Dep.name)
}

func TestNeedle_InjectStructFieldsFromStore(t *testing.T) {
	t.Cleanup(needle.Reset)

	store := needle.NewStore()

	type Dep struct{ name string }
	type TestStruct struct{ Dep *Dep }

	assert.NoError(t, needle.RegisterInstanceToStore(store, &Dep{name: "myDep"}))

	var testStruct TestStruct
	assert.NoError(t, needle.InjectStructFieldsFromStore(store, &testStruct))
	assert.NotNil(t, testStruct.Dep)
	assert.Equal(t, "myDep", testStruct.Dep.name)
}

func TestNeedle_InjectStructFieldsInvalidType(t *testing.T) {
	t.Cleanup(needle.Reset)

	var invalidType int

	injErr := needle.InjectStructFields(&invalidType)
	assert.ErrorIs(t, injErr, needle.ErrInvalidType)
}

func TestNeedle_InjectStructFieldsNotRegistered(t *testing.T) {
	t.Cleanup(needle.Reset)

	type Dep struct{}
	type TestStruct struct{ Dep *Dep }

	var testStruct TestStruct
	assert.NoError(t, needle.InjectStructFields(&testStruct))
	assert.Nil(t, testStruct.Dep)
}
