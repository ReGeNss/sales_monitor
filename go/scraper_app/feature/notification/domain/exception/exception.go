package exception

type IDomainError interface {
	Error() string
}

type domainError struct {
	message string
}

func (e domainError) Error() string {
	return e.message
}

func NewDomainError(message string) IDomainError {
	return domainError{message: message}
}
