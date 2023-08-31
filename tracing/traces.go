package tracing

type ErrorTrace struct {
	Error error
}

func (t ErrorTrace) TraceInterface() {}

type WarningTrace struct {
	Warning interface{}
}

func (t WarningTrace) TraceInterface() {}
