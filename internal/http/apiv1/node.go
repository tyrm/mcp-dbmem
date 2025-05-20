package apiv1

import "github.com/gin-gonic/gin"

func (a *API) nodeGET(c *gin.Context) {
	ctx, span := tracer.Start(c, "nodeGET", tracerAttrs...)
	defer span.End()
}

func (a *API) nodeSearchGET(c *gin.Context) {
	ctx, span := tracer.Start(c, "nodeSearchGET", tracerAttrs...)
	defer span.End()
}
