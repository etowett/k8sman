package providers

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type (
	K8SProvider interface {
		GetPods(string) ([]map[string]interface{}, error)
		GetDeployments(string) ([]map[string]interface{}, error)
	}

	AppK8SProvider struct {
		client kubernetes.Interface
	}
)

func NewK8SProvider(env string) *AppK8SProvider {
	var (
		kubeConfig *rest.Config
		err        error
	)

	if env == "local" {
		userHomeDir, err := os.UserHomeDir()
		if err != nil {
			log.Fatalf("could not user home dir: %v", err)
		}
		kubeConfigPath := filepath.Join(userHomeDir, ".kube", "config")
		slog.Debug("Using kubeconfig", "file", kubeConfigPath)

		kubeConfig, err = clientcmd.BuildConfigFromFlags("", kubeConfigPath)
		if err != nil {
			log.Fatalf("could not get kubernetes config from home: %v", err)
		}
	} else {
		// creates the in-cluster config
		kubeConfig, err = rest.InClusterConfig()
		if err != nil {
			log.Fatalf("could not get kubernetes config from rest: %v", err)
		}
	}

	// creates the clientset
	clientset, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		log.Fatalf("could not kubernetes new config: %v", err)
	}

	return &AppK8SProvider{
		client: clientset,
	}
}

func (p *AppK8SProvider) GetDeployments(namespace string) ([]map[string]interface{}, error) {
	deployments, err := p.client.AppsV1().Deployments(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("error getting deployments: %v", err)
	}

	if err != nil {
		return nil, fmt.Errorf("could not list deployments: %v", err)
	}

	theDeployments := make([]map[string]interface{}, 0)
	for _, deploy := range deployments.Items {
		deployContainers := make([]map[string]interface{}, 0)
		for _, container := range deploy.Spec.Template.Spec.Containers {
			deployContainers = append(deployContainers, map[string]interface{}{
				"name":      container.Name,
				"image":     container.Image,
				"resources": container.Resources,
			})
		}

		theDeployments = append(theDeployments, map[string]interface{}{
			"namespace":  deploy.Namespace,
			"name":       deploy.Name,
			"containers": deployContainers,
		})
	}

	return theDeployments, nil
}

func (p *AppK8SProvider) GetPods(namespace string) ([]map[string]interface{}, error) {
	pods, err := p.client.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed getting pods: %v", err)
	}

	thePods := make([]map[string]interface{}, 0)
	for _, pod := range pods.Items {
		fmt.Printf("NS: %v, Phase: %v, Name: %v, 0-Image: %v\n", pod.Namespace, pod.Status.Phase, pod.Name, pod.Spec.Containers[0].Image)

		theContainers := make([]map[string]interface{}, 0)
		for _, container := range pod.Spec.Containers {
			theContainers = append(theContainers, map[string]interface{}{
				"name":      container.Name,
				"image":     container.Image,
				"resources": container.Resources,
			})
		}

		thePods = append(thePods, map[string]interface{}{
			"namespace":  pod.Namespace,
			"phase":      pod.Status.Phase,
			"name":       pod.Name,
			"containers": theContainers,
		})
	}

	return thePods, nil
}
