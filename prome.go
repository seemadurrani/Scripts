package main

import (
	"context"
	"fmt"
	"os"

	v1 "k8s.io/api/core/v1"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	kubeconfig := os.Getenv("KUBECONFIG")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	namespace := "monitoring"

	// Create namespace
	_, err = clientset.CoreV1().Namespaces().Create(context.TODO(), &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespace,
		},
	}, metav1.CreateOptions{})
	if err != nil {
		fmt.Printf("Namespace creation error (may already exist): %v\n", err)
	}

	// Create a simple Prometheus Deployment (minimal config)
	deploy := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "prometheus",
			Namespace: namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": "prometheus"},
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"app": "prometheus"},
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:  "prometheus",
							Image: "prom/prometheus:latest",
							Ports: []v1.ContainerPort{
								{
									ContainerPort: 9090,
								},
							},
							Args: []string{"--config.file=/etc/prometheus/prometheus.yml"},
						},
					},
				},
			},
		},
	}

	_, err = clientset.AppsV1().Deployments(namespace).Create(context.TODO(), deploy, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Println("Prometheus deployment created.")

	// Create a simple service for Prometheus
	service := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "prometheus",
			Namespace: namespace,
		},
		Spec: v1.ServiceSpec{
			Selector: map[string]string{"app": "prometheus"},
			Ports: []v1.ServicePort{
				{
					Port:     9090,
					Protocol: v1.ProtocolTCP,
				},
			},
			Type: v1.ServiceTypeClusterIP,
		},
	}

	_, err = clientset.CoreV1().Services(namespace).Create(context.TODO(), service, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Println("Prometheus service created.")
}

func int32Ptr(i int32) *int32 { return &i }
