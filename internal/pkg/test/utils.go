package test

import (
	"fmt"

	"github.com/stretchr/testify/require"
)

type TestingT interface {
	require.TestingT

	Cleanup(f func())
	Log(args ...any)
	Logf(format string, args ...any)
	Name() string
	Failed() bool
}

type MockTestingT struct{}

func (m *MockTestingT) Logf(format string, args ...any) {
	fmt.Printf(format+"\n", args...)
}

func (m *MockTestingT) Cleanup(fn func()) {
	// Just run the function immediately for simplicity
	fn()
}

func (m *MockTestingT) Fatalf(format string, args ...any) {
	panic(fmt.Sprintf(format, args...))
}

func (m *MockTestingT) Name() string {
	return "MockTest"
}
