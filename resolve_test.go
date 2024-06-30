package needle_test

import (
	"fmt"
	"sync"
	"testing"

	"github.com/goplexhq/needle"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNeedle_ResolveSingleton(t *testing.T) {
	t.Cleanup(needle.Reset)

	type testStruct struct{ name string }

	regErr := needle.Register[testStruct](needle.Singleton)
	require.NoError(t, regErr)

	{
		val, resErr := needle.Resolve[testStruct]()
		require.NoError(t, resErr)
		assert.NotNil(t, val)
		val.name = "new name" // modify the singleton instance
	}

	{
		val, resErr := needle.Resolve[testStruct]()
		require.NoError(t, resErr)
		assert.NotNil(t, val)
		assert.Equal(t, "new name", val.name) // verify the modification exists
	}
}

func TestNeedle_ResolveTransient(t *testing.T) {
	t.Cleanup(needle.Reset)

	type testStruct struct{ name string }

	regErr := needle.Register[testStruct](needle.Transient)
	require.NoError(t, regErr)

	{
		val, resErr := needle.Resolve[testStruct]()
		require.NoError(t, resErr)
		assert.NotNil(t, val)
		val.name = "new name" // modify the transient instance
	}

	{
		val, resErr := needle.Resolve[testStruct]()
		require.NoError(t, resErr)
		assert.NotNil(t, val)
		assert.Equal(t, "", val.name) // verify new instance
	}
}

func TestNeedle_ResolveScoped(t *testing.T) {
	t.Cleanup(needle.Reset)

	type testStruct struct{ name string }

	optA := needle.WithScope("Scope A")
	require.NoError(t, needle.Register[testStruct](needle.Scoped, optA))

	optB := needle.WithScope("Scope B")
	require.NoError(t, needle.Register[testStruct](needle.Scoped, optB))

	{
		val, resErr := needle.Resolve[testStruct](optA)
		require.NoError(t, resErr)
		assert.NotNil(t, val)
		val.name = "Scope A" // modify scope A instance
	}

	{
		val, resErr := needle.Resolve[testStruct](optA)
		require.NoError(t, resErr)
		assert.NotNil(t, val)
		assert.Equal(t, "Scope A", val.name) // verify scope A instance modification
	}

	{
		val, resErr := needle.Resolve[testStruct](optB)
		require.NoError(t, resErr)
		assert.NotNil(t, val)
		assert.Equal(t, "", val.name) // verify non-modified scope B instance
	}

	{
		val, resErr := needle.Resolve[testStruct](optB)
		require.NoError(t, resErr)
		assert.NotNil(t, val)
		val.name = "Scope B" // modify scope B instance
	}

	{
		val, resErr := needle.Resolve[testStruct](optA)
		require.NoError(t, resErr)
		assert.NotNil(t, val)
		assert.Equal(t, "Scope A", val.name) // verify scope A instance modification

		val, resErr = needle.Resolve[testStruct](optB)
		require.NoError(t, resErr)
		assert.NotNil(t, val)
		assert.Equal(t, "Scope B", val.name) // verify scope B instance modification
	}
}

func TestNeedle_ResolveScoped_EmptyScope(t *testing.T) {
	t.Cleanup(needle.Reset)

	type testStruct struct{}

	regErr := needle.Register[testStruct](needle.Scoped)
	require.ErrorIs(t, regErr, needle.ErrEmptyScope)
}

type testNeedleResolveThreadLocalTestStruct struct{ name string }

func testNeedleResolveThreadLocalModifyHelper(expectedName string) error {
	val, err := needle.Resolve[testNeedleResolveThreadLocalTestStruct]()
	if err != nil {
		return err
	}

	val.name = expectedName

	return nil
}

func testNeedleResolveThreadLocalVerifyHelper(expectedName string) error {
	val, err := needle.Resolve[testNeedleResolveThreadLocalTestStruct]()
	if err != nil {
		return err
	}

	if val.name != expectedName {
		return fmt.Errorf("expected %q, got %q", expectedName, val.name) //nolint:err113
	}

	return nil
}

func TestNeedle_ResolveThreadLocal(t *testing.T) {
	t.Cleanup(needle.Reset)

	require.NoError(t, needle.Register[testNeedleResolveThreadLocalTestStruct](needle.ThreadLocal))
	require.NoError(t, testNeedleResolveThreadLocalModifyHelper("main goroutine"))
	require.NoError(t, testNeedleResolveThreadLocalVerifyHelper("main goroutine"))

	var waitGroup sync.WaitGroup

	errChan := make(chan error, 1)

	waitGroup.Add(1)

	go func() {
		defer waitGroup.Done()

		if err := needle.Register[testNeedleResolveThreadLocalTestStruct](needle.ThreadLocal); err != nil {
			errChan <- err

			return
		}

		if err := testNeedleResolveThreadLocalVerifyHelper(""); err != nil { // initial state.
			errChan <- err

			return
		}

		if err := testNeedleResolveThreadLocalModifyHelper("sub goroutine"); err != nil {
			errChan <- err

			return
		}

		if err := testNeedleResolveThreadLocalVerifyHelper("sub goroutine"); err != nil {
			errChan <- err

			return
		}

		close(errChan)
	}()

	waitGroup.Wait()

	require.NoError(t, <-errChan)
	require.NoError(t, testNeedleResolveThreadLocalVerifyHelper("main goroutine"))
}

func TestNeedle_ResolveFromRegistry(t *testing.T) {
	t.Cleanup(needle.Reset)

	type testStruct struct{ name string }

	registry := needle.NewRegistry()

	err := needle.RegisterSingletonInstanceToRegistry(registry, &testStruct{name: "myStruct"})
	require.NoError(t, err)

	val, resErr := needle.ResolveFromRegistry[testStruct](registry)
	require.NoError(t, resErr)
	assert.NotNil(t, val)
	assert.Equal(t, "myStruct", val.name)
}

func TestNeedle_ResolveNotRegistered(t *testing.T) {
	t.Cleanup(needle.Reset)

	type testStruct struct{}

	val, err := needle.Resolve[testStruct]()
	require.ErrorIs(t, err, needle.ErrNotRegistered)
	assert.Nil(t, val)
}
