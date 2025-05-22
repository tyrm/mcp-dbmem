package apiv1

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var tracer = otel.Tracer("internal/http/apiv1")
var tracerAttrs = []trace.SpanStartOption{}
