package core

//Error Error
type Error interface {
	error

	Code() int
}

//HTTPError HttpError
type HTTPError interface {
	error

	Status() int
}
