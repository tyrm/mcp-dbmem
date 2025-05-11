package bun

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var tracer = otel.Tracer("internal/db/bun")
var tracerAttrs []trace.SpanStartOption
