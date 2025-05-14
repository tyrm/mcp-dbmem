package adapter

import (
	"context"

	mcp "github.com/metoro-io/mcp-golang"
	"github.com/tyrm/mcp-dbmem/internal/logic"
)

type DirectAdapter struct {
	logic logic.Logic
}

func NewDirectAdapter(logic logic.Logic) *DirectAdapter {
	return &DirectAdapter{
		logic: logic,
	}
}

func (d *DirectAdapter) CreateEntities(ctx context.Context, args CreateEntitiesArgs) (*mcp.ToolResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (d *DirectAdapter) DeleteEntities(ctx context.Context, args DeleteEntitiesArgs) (*mcp.ToolResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (d *DirectAdapter) ReadGraph(ctx context.Context, args ReadGraphArgs) (*mcp.ToolResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (d *DirectAdapter) OpenNodes(ctx context.Context, args OpenNodesArgs) (*mcp.ToolResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (d *DirectAdapter) SearchNodes(ctx context.Context, args SearchNodesArgs) (*mcp.ToolResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (d *DirectAdapter) AddObservations(ctx context.Context, args AddObservationsArgs) (*mcp.ToolResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (d *DirectAdapter) DeleteObservations(ctx context.Context, args DeleteObservationsArgs) (*mcp.ToolResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (d *DirectAdapter) CreateRelations(ctx context.Context, args CreateRelationsArgs) (*mcp.ToolResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (d *DirectAdapter) DeleteRelations(ctx context.Context, args DeleteRelationsArgs) (*mcp.ToolResponse, error) {
	//TODO implement me
	panic("implement me")
}

var _ Adapter = (*DirectAdapter)(nil)
