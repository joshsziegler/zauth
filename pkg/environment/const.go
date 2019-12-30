// Package environment simply defines three constants for development, testing,
// and production. This allows other modules to agree on the same constants
// without importing unused code.

package environment

const (
	Dev = 1 << iota
	Test
	Prod
)
