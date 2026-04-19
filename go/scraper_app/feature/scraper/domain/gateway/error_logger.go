package gateway

type ErrorContext struct {
	Context string
	URL     string
	Index   int
}

type ErrorLogger interface {
	LogError(err error, ctx ErrorContext) string
}
