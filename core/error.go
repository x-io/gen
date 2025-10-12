package core

// Error Error
type Error interface {
	error

	HttpStatus() int
}
