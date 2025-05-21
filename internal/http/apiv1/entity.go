package apiv1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tyrm/mcp-dbmem/pkg/api"
)

func (a *API) entityPOST(c *gin.Context) {
	ctx, span := tracer.Start(c, "entityPOST", tracerAttrs...)
	defer span.End()

	var entity api.CreateEntitiesRequest
	// Bind the JSON body to your struct
	if err := c.ShouldBindJSON(&entity); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate the request body
	if err := a.validate.StructCtx(ctx, entity); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// create entity models
	entities := make([]api.Entity, len(entity.Entities))
	for i, e := range entity.Entities {
		entities[i] = api.Entity{
			Name:         e.Name,
			Type:         e.Type,
			Observations: e.Observations,
		}
	}

	c.JSON(http.StatusNotImplemented, struct{}{})
}

func (a *API) entityDELETE(c *gin.Context) {
	_, span := tracer.Start(c, "entityDELETE", tracerAttrs...)
	defer span.End()

	c.JSON(http.StatusNotImplemented, struct{}{})
}
