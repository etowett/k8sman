package handlers

import (
	"fmt"
	"log/slog"
	"net/http"

	"k8sman/internal/schema"

	"github.com/gin-gonic/gin"
)

func (h Handler) Deployments() func(c *gin.Context) {
	return func(c *gin.Context) {
		var form schema.RequestForm
		err := c.ShouldBind(&form)
		if err != nil {
			errMessage := fmt.Sprintf("Error with request: %v", err)
			slog.Error(errMessage)
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": errMessage,
			})
			return
		}
		deployments, err := h.k8sProvider.GetDeployments(form.Namespace)
		if err != nil {
			message := fmt.Sprintf("failed to get deployments: %v", err)
			slog.Error(message)
			c.JSON(http.StatusInternalServerError, gin.H{
				"success":     false,
				"message":     message,
				"deployments": deployments,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success":     true,
			"message":     "success",
			"deployments": deployments,
		})
	}
}

func (h Handler) Pods() func(c *gin.Context) {
	return func(c *gin.Context) {
		var form schema.RequestForm
		err := c.ShouldBind(&form)
		if err != nil {
			errMessage := fmt.Sprintf("Error with request: %v", err)
			slog.Error(errMessage)
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": errMessage,
			})
			return
		}

		pods, err := h.k8sProvider.GetPods(form.Namespace)
		if err != nil {
			message := fmt.Sprintf("failed to get pods: %v", err)
			slog.Error(message)
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": message,
				"pods":    pods,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "success",
			"pods":    pods,
		})
	}
}
