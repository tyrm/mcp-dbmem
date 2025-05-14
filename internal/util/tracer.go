package util

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var tracer = otel.Tracer("internal/util")
var tracerAttrs []trace.SpanStartOption
