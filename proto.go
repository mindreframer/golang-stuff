// proto gives Go operations like Map, Reduce, Filter, De/Multiplex, etc.
// without sacrificing idiomatic harmony or speed.
//
// The `Proto` type is a stand-in approximation for dynamic typing. Due to
// Go's powerful casting and type inference idioms, we can approximate the
// flexibility of dynamic typing even though Go is a statically typed language.
// Doing so sacrifices some of the benefits of static typing AND some of the
// benefits of dynamic typing, but this sacrifice is fundamentally required by
// Go until such time as a true 'Generic' type is implemented.
//
// In order to use a Proto-typed variable (from here on out, simply a 'Proto'),
// you will generally have to cast it to a type that you will know to use based
// on the semantics of your program.
//
// This package (specifically, the other files in this package) provide
// operations on Proto variables as well as some that make Proto variables out
// of 'traditionally typed' variables. Many of the operations will require the
// use of higher-order functions which you will need to provide, and those
// functions commonly will need you to manually "unbox" (cast-from-Proto) the
// variable to perform useful operations.
//
// Examples of the use of this package can be found in the "*_test.go" files,
// which contain testing code. A good example of a higher-order function which
// will commonly need manual-unboxing is the `Filter` function, found in
// "filter.go". `Filter` takes as its first argument a filter-function which
// will almost certainly require you to un-box the Proto channel values that it
// receives to perform the filtering action.
//
// Finally, a word on the entire point of this package: while it is named after
// the Proto type that pervades it and guides its syntax, the true nature of
// the `proto` package lies in cascading channels, rather than in dynamic
// typing. In fact this package might be more appropriately named after
// channels. Maybe `canal` would have been a better name. I wanted to bring
// the syntax and familiar patterns of functional programming idioms to the
// power and scalability of Go's goroutines and channels, and found that the
// syntax made this task very simple.
//
// You may find, as I did, that the majority of the code in this package is very
// 'obvious'. At first I was concerned by this - much of the code is very
// trivial - but now I feel pleased by the re-usability and natural
// 'correctness' of `proto`. Look at this package not as some monumental
// time-saving framework, but rather as a light scaffold for a useful and
// idiomatic style of programming within the existing constructs of Go.
//
// Ultimately, though, you're going to be typing the word Proto an awful lot,
// and thus the type became the eponym.
package proto

// The Proto type. (Get it?)
type Proto interface{}
