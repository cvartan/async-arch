package strhelper

type StringWriter struct {
	source *string
}

func CreateStringWriter(str *string) *StringWriter {
	return &StringWriter{
		source: str,
	}
}

func (sw *StringWriter) Write(p []byte) (int, error) {
	*sw.source = string(p)
	return len(p), nil
}

func (sw *StringWriter) String() string {
	return *sw.source
}
