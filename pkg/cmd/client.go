package cmd

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	namespaces   = "namespaces"
	pods         = "pods"
	services     = "services"
	configs      = "config maps"
	pvcs         = "persistent volume claims"
	sas          = "service accounts"
	secrets      = "secrets"
	endpoints    = "endpoints"
	daemonsets   = "daemonsets"
	deploys      = "deployments"
	replicasets  = "replica sets"
	statefulsets = "stateful sets"
	jobs         = "jobs"
)

func (o *globalSettings) InitClient() {
	restConfig, err := o.configFlags.ToRESTConfig()
	if err != nil {
		panic(err)
	}
	restConfig.ContentType = "application/vnd.kubernetes.protobuf"
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

func (o *globalSettings) GetRessources() (map[string]int, error) {
	resources := make(map[string]int)
	opts := metav1.ListOptions{}
	ns, err := o.client.CoreV1().Namespaces().List(opts)
	if err == nil {
		resources[namespaces] = len(ns.Items)
	}
	po, err := o.client.CoreV1().Pods("").List(opts)
	if err == nil {
		resources[pods] = len(po.Items)
	}
	svc, err := o.client.CoreV1().Services("").List(opts)
	if err == nil {
		resources[services] = len(svc.Items)
	}
	cm, err := o.client.CoreV1().ConfigMaps("").List(opts)
	if err == nil {
		resources[configs] = len(cm.Items)
	}
	sec, err := o.client.CoreV1().Secrets("").List(opts)
	if err == nil {
		resources[secrets] = len(sec.Items)
	}
	sa, err := o.client.CoreV1().ServiceAccounts("").List(opts)
	if err == nil {
		resources[sas] = len(sa.Items)
	}
	end, err := o.client.CoreV1().Endpoints("").List(opts)
	if err == nil {
		resources[endpoints] = len(end.Items)
	}
	pvc, err := o.client.CoreV1().PersistentVolumeClaims("").List(opts)
	if err == nil {
		resources[pvcs] = len(pvc.Items)
	}
	ds, err := o.client.AppsV1().DaemonSets("").List(opts)
	if err == nil {
		resources[daemonsets] = len(ds.Items)
	}
	depl, err := o.client.AppsV1().Deployments("").List(opts)
	if err == nil {
		resources[deploys] = len(depl.Items)
	}
	rs, err := o.client.AppsV1().ReplicaSets("").List(opts)
	if err == nil {
		resources[replicasets] = len(rs.Items)
	}
	sts, err := o.client.AppsV1().StatefulSets("").List(opts)
	if err == nil {
		resources[statefulsets] = len(sts.Items)
	}
	j, err := o.client.BatchV1().Jobs("").List(opts)
	if err == nil {
		resources[jobs] = len(j.Items)
	}

	return resources, nil
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
