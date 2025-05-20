package adapter

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	mcp "github.com/metoro-io/mcp-golang"
	"github.com/tyrm/mcp-dbmem/internal/logic"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var serverTracer = otel.Tracer("internal/adapter.ServerAdapter")
var serverTracerAttrs []trace.SpanStartOption

type ServerAdapter struct {
	handlers map[string]RPCHandler
	logic    logic.Logic
}

func NewServerAdapter(logic logic.Logic) *ServerAdapter {
	return &ServerAdapter{
		handlers: make(map[string]RPCHandler),
		logic:    logic,
	}
}

// api

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

// websockets

type RPCHandler func(params json.RawMessage) (any, error)

type RPCMessage struct {
	ID     string          `json:"id"`
	Method string          `json:"method"`
	Params json.RawMessage `json:"params,omitempty"`
	Result json.RawMessage `json:"result,omitempty"`
	Error  *RPCError       `json:"error,omitempty"`
}

type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func (d *ServerAdapter) Register(method string, handler RPCHandler) {
	d.handlers[method] = handler
}
func (d *ServerAdapter) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Upgrade error: %v", err)
		return
	}
	defer conn.Close()

	for {
		var msg RPCMessage
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("Read error: %v", err)
			break
		}

		// Handle the RPC call
		response := d.handleRPC(msg)

		err = conn.WriteJSON(response)
		if err != nil {
			log.Printf("Write error: %v", err)
			break
		}
	}
}

func (d *ServerAdapter) handleRPC(req RPCMessage) RPCMessage {
	response := RPCMessage{ID: req.ID}

	handler, exists := d.handlers[req.Method]
	if !exists {
		response.Error = &RPCError{
			Code:    -32601,
			Message: "Method not found",
		}
		return response
	}

	result, err := handler(req.Params)
	if err != nil {
		response.Error = &RPCError{
			Code:    -32000,
			Message: err.Error(),
		}
	} else {
		resultBytes, _ := json.Marshal(result)
		response.Result = resultBytes
	}

	return response
}
