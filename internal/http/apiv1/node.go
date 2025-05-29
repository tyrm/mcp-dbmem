package apiv1

import "github.com/gin-gonic/gin"

func (a *API) nodeGET(c *gin.Context) {
	_, span := tracer.Start(c, "nodeGET", tracerAttrs...)
	defer span.End()
}

func (a *API) nodeSearchGET(c *gin.Context) {
	_, span := tracer.Start(c, "nodeSearchGET", tracerAttrs...)
	defer span.End()
}
