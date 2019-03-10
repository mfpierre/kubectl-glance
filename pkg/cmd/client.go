package cmd

import (
	"fmt"

	"k8s.io/client-go/kubernetes"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (o *globalSettings) InitClient() {
	restConfig, err := o.configFlags.ToRESTConfig()
	if err != nil {
		panic(err)
	}
	c := kubernetes.NewForConfigOrDie(restConfig)
	rawKubeConfig := o.configFlags.ToRawKubeConfigLoader()
	ns, _, _ := rawKubeConfig.Namespace()
	o.namespace = ns
	o.client = c
	o.restConfig = restConfig
}

// GeNodeForPod gets the node of a pod
func (o *globalSettings) GeNodeForPod(podName string) (string, error) {
	pod, err := o.client.CoreV1().Pods(o.namespace).Get(podName, metav1.GetOptions{})
	if err != nil {
		return "", fmt.Errorf("got an error while retrieving pod %s: %s", podName, err)
	}
	return pod.Spec.NodeName, nil
}

func (o *globalSettings) GetNamespaces() (int, error) {
	opts := metav1.ListOptions{}
	ns, err := o.client.CoreV1().Namespaces().List(opts)
	if err != nil {
		return 0, fmt.Errorf("got an error while getting namespaces: %s", err)
	}
	return len(ns.Items), nil
}

func (o *globalSettings) GetNodes() (int, int, error) {
	opts := metav1.ListOptions{}
	no, err := o.client.CoreV1().Nodes().List(opts)
	if err != nil {
		return 0, 0, fmt.Errorf("got an error while getting namespaces: %s", err)
	}
	unschedulable := 0
	for _, n := range no.Items {
		if n.Spec.Unschedulable {
			unschedulable++
		}
		// TODO capacity check
		// fmt.Println(n.Status.Allocatable)
		// fmt.Println(n.Status.Capacity)
	}
	return len(no.Items), unschedulable, nil
}
