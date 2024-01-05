package handlers

import (
	"k8sman/internal/config"
	"k8sman/internal/providers"

	"github.com/gin-gonic/gin"
)

type (
	SMSHandler interface {
		Healthz() func(c *gin.Context)
	}

	Handler struct {
		cfg         *config.ServiceConfig
		k8sProvider *providers.AppK8SProvider
	}
)

func NewHandler(
	cfg *config.ServiceConfig,
	k8sProvider *providers.AppK8SProvider,
) *Handler {
	return &Handler{
		cfg:         cfg,
		k8sProvider: k8sProvider,
	}
}
