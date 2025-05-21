package apiv1

import (
	ihttp "github.com/tyrm/mcp-dbmem/internal/http"
	"github.com/tyrm/mcp-dbmem/internal/http/path"
	"github.com/tyrm/mcp-dbmem/internal/logic"
)

type API struct {
	logic logic.Logic
}

func New(logic logic.Logic) *API {
	return &API{
		logic: logic,
	}
}

func (a *API) Name() string {
	return "apiv1"
}

// Route attaches routes to the web server.
func (a *API) Route(s *ihttp.Server) {
	api := s.Group(path.V1)

	api.POST(path.Entities, a.entityPOST)
	api.DELETE(path.Entities, a.entityDELETE)
	api.POST(path.Graph, a.graphGET)
	api.GET(path.Nodes, a.nodeGET)
	api.GET(path.NodesSearch, a.nodeSearchGET)
	api.POST(path.Observations, a.observationPOST)
	api.DELETE(path.Observations, a.observationDELETE)
	api.POST(path.Relations, a.relationPOST)
	api.DELETE(path.Relations, a.relationDELETE)
}
