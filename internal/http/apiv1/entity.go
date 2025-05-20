package apiv1

import "github.com/gin-gonic/gin"

func (a *API) entityPOST(c *gin.Context) {
	ctx, span := tracer.Start(c, "entityPOST", tracerAttrs...)
	defer span.End()
}

func (a *API) entityDELETE(c *gin.Context) {
	ctx, span := tracer.Start(c, "entityDELETE", tracerAttrs...)
	defer span.End()
}
