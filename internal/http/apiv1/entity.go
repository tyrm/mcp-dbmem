package apiv1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tyrm/mcp-dbmem/pkg/api"
)

func (a *API) entityPOST(c *gin.Context) {
	_, span := tracer.Start(c, "entityPOST", tracerAttrs...)
	defer span.End()

	var entity api.Entity
	// Bind the JSON body to your struct
	if err := c.ShouldBindJSON(&entity); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNotImplemented, struct{}{})
}

func (a *API) entityDELETE(c *gin.Context) {
	_, span := tracer.Start(c, "entityDELETE", tracerAttrs...)
	defer span.End()

	c.JSON(http.StatusNotImplemented, struct{}{})
}
