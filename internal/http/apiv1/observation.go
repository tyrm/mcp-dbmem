package apiv1

import "github.com/gin-gonic/gin"

func (a *API) observationPOST(c *gin.Context) {
	ctx, span := tracer.Start(c, "observationPOST", tracerAttrs...)
	defer span.End()
}

func (a *API) observationDELETE(c *gin.Context) {
	ctx, span := tracer.Start(c, "observationDELETE", tracerAttrs...)
	defer span.End()
}
