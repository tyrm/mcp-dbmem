package apiv1

import "github.com/gin-gonic/gin"

func (a *API) relationPOST(c *gin.Context) {
	_, span := tracer.Start(c, "relationPOST", tracerAttrs...)
	defer span.End()
}

func (a *API) relationDELETE(c *gin.Context) {
	_, span := tracer.Start(c, "relationDELETE", tracerAttrs...)
	defer span.End()
}
