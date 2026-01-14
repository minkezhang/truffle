package form

type Key struct {
	// Label is a human-readable display text
	Label string

	// Key is a machine label for the key.
	Key string
}

type Value[T any] struct {
	Key   Key
	Value T
}
