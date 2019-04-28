package cmd

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const namespaces = "namespaces"
const pods = "pods"
const services = "services"
const configs = "config maps"
const pvcs = "persistent volume claims"
const sas = "service accounts"
const secrets = "secrets"
const endpoints = "endpoints"
const daemonsets = "daemonsets"
const deploys = "deployments"
const replicasets = "replica sets"
const statefulsets = "stateful sets"
const jobs = "jobs"

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
	namespacedResources[daemonsets] = 0
	namespacedResources[deploys] = 0
	namespacedResources[replicasets] = 0
	namespacedResources[statefulsets] = 0
	namespacedResources[jobs] = 0

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
		ds, err := o.client.AppsV1().DaemonSets(n.Name).List(opts)
		if err == nil {
			namespacedResources[daemonsets] = len(ds.Items)
		}
		depl, err := o.client.AppsV1().Deployments(n.Name).List(opts)
		if err == nil {
			namespacedResources[deploys] = len(depl.Items)
		}
		rs, err := o.client.AppsV1().ReplicaSets(n.Name).List(opts)
		if err == nil {
			namespacedResources[replicasets] = len(rs.Items)
		}
		sts, err := o.client.AppsV1().StatefulSets(n.Name).List(opts)
		if err == nil {
			namespacedResources[statefulsets] = len(sts.Items)
		}
		j, err := o.client.BatchV1().Jobs(n.Name).List(opts)
		if err == nil {
			namespacedResources[jobs] = len(j.Items)
		}
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

func (o *globalSettings) GetNodes() (int, int, string, string, error) {
	opts := metav1.ListOptions{}
	no, err := o.client.CoreV1().Nodes().List(opts)
	if err != nil {
		return 0, 0, "", "", fmt.Errorf("got an error while getting namespaces: %s", err)
	}
	unschedulable := 0
	cpuAllocatable, _ := resource.ParseQuantity("0")
	memAllocatable, _ := resource.ParseQuantity("0")
	for _, n := range no.Items {
		if n.Spec.Unschedulable {
			unschedulable++
		}
		cpuAllocatable.Add(n.Status.Allocatable[corev1.ResourceName("cpu")])
		memAllocatable.Add(n.Status.Capacity[corev1.ResourceName("memory")])
	}
	return len(no.Items), unschedulable, cpuAllocatable.String(), memAllocatable.String(), nil
}
