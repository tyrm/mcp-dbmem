package v1

import (
	"context"
	"encoding/json"

	mcp "github.com/metoro-io/mcp-golang"
	"github.com/tyrm/mcp-dbmem/internal/db"
	"github.com/tyrm/mcp-dbmem/internal/logic"
)

// Logic implements the program logic
type Logic struct {
	DB db.DB
}

var _ logic.Logic = (*Logic)(nil)

func toolJSONResponse(ctx context.Context, response any) (*mcp.ToolResponse, error) {
	ctx, span := tracer.Start(ctx, "toolJSONResponse", tracerAttrs...)
	defer span.End()

	// convert response to json string
	jsonResponse, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return mcp.NewToolResponse(
		mcp.NewTextContent(string(jsonResponse)),
	), nil
}
