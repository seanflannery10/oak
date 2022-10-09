package middleware

import "github.com/justinas/alice"

func Chain(constructors ...alice.Constructor) alice.Chain {
	return alice.New(constructors...)
}
