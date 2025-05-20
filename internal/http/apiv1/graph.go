package apiv1

import "github.com/gin-gonic/gin"

func (a *API) graphGET(c *gin.Context) {
	ctx, span := tracer.Start(c, "graphGET", tracerAttrs...)
	defer span.End()
}
