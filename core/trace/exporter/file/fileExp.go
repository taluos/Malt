package file

import (
	"fmt"
	"os"

	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
)

func NewFileExporter(endpoint string) (*stdouttrace.Exporter, error) {
	f, err := os.OpenFile(endpoint, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("file exporter endpoint error: %s", err.Error())
	}
	return stdouttrace.New(stdouttrace.WithWriter(f))
}
