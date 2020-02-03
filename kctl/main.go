package kctl

import (
	"bytes"
	"context"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes/scheme"
	typev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/remotecommand"
	kubernetes "krabbitmqctl/kubernets"
	"log"
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

	pods, _ := getPodsForSvc(&ss, config.NameSpace, clientset.CoreV1())

	podToQuery := ""
	for _, pod := range pods.Items {
		podToQuery = pod.Name
	}

	req := clientset.CoreV1().RESTClient().Post().Resource("pods").Name(podToQuery).
		Namespace(config.NameSpace).SubResource("exec")

	buf := &bytes.Buffer{}
	errBuf := &bytes.Buffer{}

	ctlCommand := []string{"rabbitmqctl"}
	for _, parameters := range config.CtlCommand {
		ctlCommand = append(ctlCommand, parameters)

	}

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
	log.Printf(" executing %v command on POD: %s, namespace: %s", ctlCommand, podToQuery, config.NameSpace)
	restConfig, _ := clientConfig.ClientConfig()
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

func getPodsForSvc(svc *corev1.Service, namespace string, k8sClient typev1.CoreV1Interface) (*corev1.PodList, error) {
	set := labels.Set(svc.Spec.Selector)
	listOptions := metav1.ListOptions{LabelSelector: set.AsSelector().String()}
	pods, err := k8sClient.Pods(namespace).List(listOptions)
	return pods, err
}
