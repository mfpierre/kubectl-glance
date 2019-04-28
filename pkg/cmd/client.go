package cmd

import (
	"fmt"

	"k8s.io/client-go/kubernetes"

	// corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const namespaces = "namespaces"
const pods = "pods"
const services = "services"
const configs = "config maps"
const pvcs = "persistent volume claims"
const sas = "service accounts"
const secrets = "secrets"
const endpoints = "endpoints"

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

func (o *globalSettings) GetNamespacedRessources() (map[string]int, error) {
	namespacedResources := make(map[string]int)
	opts := metav1.ListOptions{}
	ns, err := o.client.CoreV1().Namespaces().List(opts)
	if err != nil {
		return namespacedResources, fmt.Errorf("got an error while getting namespaces: %s", err)
	}
	namespacedResources[namespaces] = len(ns.Items)
	namespacedResources[pods] = 0
	namespacedResources[services] = 0
	namespacedResources[configs] = 0
	namespacedResources[pvcs] = 0
	namespacedResources[secrets] = 0
	namespacedResources[sas] = 0
	namespacedResources[endpoints] = 0
	for _, n := range ns.Items {
		opts := metav1.ListOptions{}
		po, err := o.client.CoreV1().Pods(n.Name).List(opts)
		if err == nil {
			namespacedResources[pods] += len(po.Items)
		}
		svc, err := o.client.CoreV1().Services(n.Name).List(opts)
		if err == nil {
			namespacedResources[services] += len(svc.Items)
		}
		cm, err := o.client.CoreV1().ConfigMaps(n.Name).List(opts)
		if err == nil {
			namespacedResources[configs] += len(cm.Items)
		}
		sec, err := o.client.CoreV1().Secrets(n.Name).List(opts)
		if err == nil {
			namespacedResources[secrets] += len(sec.Items)
		}
		sa, err := o.client.CoreV1().ServiceAccounts(n.Name).List(opts)
		if err == nil {
			namespacedResources[sas] += len(sa.Items)
		}
		end, err := o.client.CoreV1().Endpoints(n.Name).List(opts)
		if err == nil {
			namespacedResources[endpoints] += len(end.Items)
		}
		pvc, err := o.client.CoreV1().PersistentVolumeClaims(n.Name).List(opts)
		if err == nil {
			namespacedResources[pvcs] += len(pvc.Items)
		}
		// secr, sa, endpoint
	}
	return namespacedResources, nil
}

func (o *globalSettings) GetPersistentVolumes() (int, error) {
	opts := metav1.ListOptions{}
	pv, err := o.client.CoreV1().PersistentVolumes().List(opts)
	if err != nil {
		return 0, fmt.Errorf("got an error while getting pv: %s", err)
	}
	return len(pv.Items), nil
}

func (o *globalSettings) GetNodes() (int, int, error) {
	opts := metav1.ListOptions{}
	no, err := o.client.CoreV1().Nodes().List(opts)
	if err != nil {
		return 0, 0, fmt.Errorf("got an error while getting namespaces: %s", err)
	}
	unschedulable := 0
	// cpuAllocatable := int64(0)
	// cpuCapacity := int64(0)
	for _, n := range no.Items {
		if n.Spec.Unschedulable {
			unschedulable++
		}
		// TODO capacity check
		// q := n.Status.Allocatable[corev1.ResourceName("cpu")]
		// v, _ := q.AsInt64()
		// fmt.Println(v)
		// cpuAllocatable += v
		// q = n.Status.Capacity[corev1.ResourceName("cpu")]
		// cpuCapacity += q.Value()
		// fmt.Println(n.Status.Allocatable)
		// fmt.Println(n.Status.Capacity)
	}
	// fmt.Println(cpuAllocatable)
	// fmt.Println(cpuCapacity)
	return len(no.Items), unschedulable, nil
}
