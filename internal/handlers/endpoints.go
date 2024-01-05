package handlers

import (
	"github.com/gin-gonic/gin"
)

const (
	HealthzPath     = "/healthz"
	DeploymentsPath = "/deployments"
	PodsPath        = "/pods"
)

func AddUnsecuredEndpoints(
	r *gin.RouterGroup,
	httpHandler *Handler,
) {
	r.GET(HealthzPath, httpHandler.Healthz())
	r.HEAD(HealthzPath, httpHandler.Healthz())
	r.GET(DeploymentsPath, httpHandler.Deployments())
	r.GET(PodsPath, httpHandler.Pods())
}

func AddSecuredEndpoints(
	r *gin.RouterGroup,
	httpHandler *Handler,
) {
}
