package test

import (
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
