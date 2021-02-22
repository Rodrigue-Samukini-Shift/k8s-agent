package main

import (
	"bytes"
	"context"
	"net/http"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/json"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type TelemetryData struct {
	ClusterProvider string `json:"clusterProvider"`
	CustomerID string `json:"customerId"`
	ClusterName string `json:"clusterName"`
	ClusterRegion string `json:"clusterRegion"`
	NodeList *v1.NodeList `json:"nodeList"`
	PodList *v1.PodList `json:"podList"`
}

func sendTelemetry(clusterName string, telemetry *TelemetryData) error {
	b, err := json.Marshal(telemetry)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		"POST",
		os.Getenv("TELEMETRY_API_URL"),
		bytes.NewBuffer(b),
	)

	if err != nil {
		return err
	}

	req.Header.Set("X-API-Key", os.Getenv("TELEMETRY_API_KEY"))
	req.Header.Set("Content-Type", "application/json")


	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	_, err = client.Do(req)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	log := logrus.New()
	log.Info("starting the agent")
	ctx := context.Background()
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	if err != nil {
		panic(err)
	}

	const interval = 10 * time.Second
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
		case <-ctx.Done():
		}

		nodes, err := clientset.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
		if err != nil {
			log.Errorf("failed: %v", err)
			panic(err)
		}

		pods, err := clientset.CoreV1().Pods("").List(ctx, metav1.ListOptions{})
		if err != nil {
			log.Errorf("failed: %v", err)
		}


		log.Infof("nodes[%d], pods[%d] in the cluster", len(nodes.Items), len(pods.Items))

		node1 := nodes.Items[0]
		clusterName := node1.Labels["alpha.eksctl.io/cluster-name"]
		clusterRegion := node1.Labels["topology.kubernetes.io/region"]

		err = sendTelemetry(clusterName, &TelemetryData{
			ClusterProvider: "EKS",
			ClusterName: clusterName,
			ClusterRegion: clusterRegion,
			CustomerID: os.Getenv("TELEMETRY_CUSTOMER_ID"),
			NodeList: nodes,
			PodList:  pods,
		})

		if err != nil {
			log.Errorf("failed to send data: %v", err)
		}
	}
}