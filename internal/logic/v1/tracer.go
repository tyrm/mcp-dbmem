package v1

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var tracer = otel.Tracer("internal/logic/v1")
var tracerAttrs []trace.SpanStartOption
