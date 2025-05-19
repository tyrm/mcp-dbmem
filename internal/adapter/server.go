package adapter

import (
	"context"

	mcp "github.com/metoro-io/mcp-golang"
	"github.com/tyrm/mcp-dbmem/internal/logic"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var serverTracer = otel.Tracer("internal/adapter.ServerAdapter")
var serverTracerAttrs = []trace.SpanStartOption{
	trace.WithAttributes(
		attribute.String("service.name", "mcp-dbmem"),
		attribute.String("component", "server"),
		attribute.String("span.kind", "server"),
	),
}

type ServerAdapter struct {
	logic logic.Logic
}

func NewServerAdapter(logic logic.Logic) *ServerAdapter {
	return &ServerAdapter{
		logic: logic,
	}
}

func (d *ServerAdapter) Apply(server *mcp.Server) error {
	return apply(d, server)
}

func (d *ServerAdapter) CreateEntities(ctx context.Context, args CreateEntitiesArgs) (*mcp.ToolResponse, error) {
	ctx, span := serverTracer.Start(ctx, "CreateEntities", serverTracerAttrs...)
	defer span.End()

	return nil, nil
}

func (d *ServerAdapter) DeleteEntities(ctx context.Context, args DeleteEntitiesArgs) (*mcp.ToolResponse, error) {
	ctx, span := serverTracer.Start(ctx, "DeleteEntities", serverTracerAttrs...)
	defer span.End()

	return nil, nil
}

func (d *ServerAdapter) ReadGraph(ctx context.Context, args ReadGraphArgs) (*mcp.ToolResponse, error) {
	ctx, span := serverTracer.Start(ctx, "ReadGraph", serverTracerAttrs...)
	defer span.End()

	return nil, nil
}

func (d *ServerAdapter) OpenNodes(ctx context.Context, args OpenNodesArgs) (*mcp.ToolResponse, error) {
	ctx, span := serverTracer.Start(ctx, "OpenNodes", serverTracerAttrs...)
	defer span.End()

	return nil, nil
}

func (d *ServerAdapter) SearchNodes(ctx context.Context, args SearchNodesArgs) (*mcp.ToolResponse, error) {
	ctx, span := serverTracer.Start(ctx, "SearchNodes", serverTracerAttrs...)
	defer span.End()

	return nil, nil
}

func (d *ServerAdapter) AddObservations(ctx context.Context, args AddObservationsArgs) (*mcp.ToolResponse, error) {
	ctx, span := serverTracer.Start(ctx, "AddObservations", serverTracerAttrs...)
	defer span.End()

	return nil, nil
}

func (d *ServerAdapter) DeleteObservations(ctx context.Context, args DeleteObservationsArgs) (*mcp.ToolResponse, error) {
	ctx, span := serverTracer.Start(ctx, "DeleteObservations", serverTracerAttrs...)
	defer span.End()

	return nil, nil
}

func (d *ServerAdapter) CreateRelations(ctx context.Context, args CreateRelationsArgs) (*mcp.ToolResponse, error) {
	ctx, span := serverTracer.Start(ctx, "CreateRelations", serverTracerAttrs...)
	defer span.End()

	return nil, nil
}

func (d *ServerAdapter) DeleteRelations(ctx context.Context, args DeleteRelationsArgs) (*mcp.ToolResponse, error) {
	ctx, span := serverTracer.Start(ctx, "DeleteRelations", serverTracerAttrs...)
	defer span.End()

	return nil, nil
}
