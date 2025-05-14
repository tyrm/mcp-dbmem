package util

import (
	"context"
	"encoding/json"

	mcp "github.com/metoro-io/mcp-golang"
)

func ToolJSONResponse(ctx context.Context, response any) (*mcp.ToolResponse, error) {
	_, span := tracer.Start(ctx, "ToolJSONResponse", tracerAttrs...)
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
