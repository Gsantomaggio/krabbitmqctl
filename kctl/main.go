package kctl

import (
	"bytes"
	"context"
	"errors"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes/scheme"
	typev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/remotecommand"
	kubernetes "krabbitmqctl/kubernets"
	"strings"
)

func Run(ctx context.Context, config *Config) (string, string, error) {
	clientConfig := kubernetes.NewClientConfig(config.KubeConfig, config.ContextName)
	clientset, err := kubernetes.NewClientSet(clientConfig)
	if err != nil {
		return "", "", err
	}

	services := clientset.CoreV1().Services(config.NameSpace)
	servicesList, _ := services.List(metav1.ListOptions{})

	ss := corev1.Service{}
	for _, service := range servicesList.Items {
		if service.Name == strings.ToLower(config.ServiceName) {
			ss = service
			break
		}
	}

	pods, errSVC := getPodsForSvc(&ss, config.NameSpace, clientset.CoreV1())
	if errSVC != nil {
		return "", "", errSVC
	}

	podToQuery, errPod := getPodToQuery(config, pods)
	if errPod != nil {
		return "", "", errPod
	}

	req := clientset.CoreV1().RESTClient().Post().Resource("pods").Name(podToQuery).
		Namespace(config.NameSpace).SubResource("exec")

	buf := &bytes.Buffer{}
	errBuf := &bytes.Buffer{}
//	log.Printf("Going to query %s", podToQuery)
	ctlCommand := []string{"rabbitmqctl"}
	for _, parameters := range config.CtlCommand {
		ctlCommand = append(ctlCommand, parameters)
	}

	ctlCommand = append(ctlCommand, "-p")
	ctlCommand = append(ctlCommand, config.VirtualHost)

	option := &v1.PodExecOptions{
		Command: ctlCommand,
		Stdin:   false,
		Stdout:  true,
		Stderr:  true,
		TTY:     true,
	}
	req.VersionedParams(
		option,
		scheme.ParameterCodec,
	)
	restConfig, errConf := clientConfig.ClientConfig()
	if errConf != nil {
		return "", "", errConf
	}

	exec, err := remotecommand.NewSPDYExecutor(restConfig, "POST", req.URL())
	if err != nil {
		return "", "", err
	}
	err = exec.Stream(remotecommand.StreamOptions{
		Stdout: buf,
		Stderr: errBuf,
	})

	if err != nil {
		return "", "", err
	}

	return buf.String(), errBuf.String(), nil

}

func getPodToQuery(config *Config, pods *v1.PodList) (string, error) {
	if config.PodName != "" {
		return strings.ToLower(config.PodName), nil
	} else {
		for _, pod := range pods.Items {
			return pod.Name, nil
		}
	}
	return "", errors.New("no pod selected")
}

func getPodsForSvc(svc *corev1.Service, namespace string, k8sClient typev1.CoreV1Interface) (*corev1.PodList, error) {
	set := labels.Set(svc.Spec.Selector)
	listOptions := metav1.ListOptions{LabelSelector: set.AsSelector().String()}
	pods, err := k8sClient.Pods(namespace).List(listOptions)
	return pods, err
}
