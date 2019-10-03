package main

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	flag "github.com/spf13/pflag"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/pager"
)

type k8sClient struct {
	client kubernetes.Interface
}

func (c *k8sClient) ListPods(namespace string) ([]*v1.Pod, error) {
	p := pager.New(pager.SimplePageFunc(func(opts metav1.ListOptions) (runtime.Object, error) {
		list, err := c.client.CoreV1().Pods(namespace).List(opts)
		if err != nil {
			return nil, errors.Wrap(err, "cannot retrieve pods")
		}

		return list, nil
	}))
	p.PageSize = 500

	ctx := context.Background()

	pods := []*v1.Pod{}

	err := p.EachListItem(ctx, metav1.ListOptions{}, func(obj runtime.Object) error {
		pod, ok := obj.(*v1.Pod)
		if !ok {
			return errors.Errorf("this is not a pod: %#v", obj)
		}

		pods = append(pods, pod)

		return nil
	})
	if err != nil {
		return []*v1.Pod{}, errors.Wrap(err, "cannot iterate secrets")
	}

	return pods, nil
}

func (c *k8sClient) ListSecrets(namespace string) ([]*v1.Secret, error) {
	p := pager.New(pager.SimplePageFunc(func(opts metav1.ListOptions) (runtime.Object, error) {
		list, err := c.client.CoreV1().Secrets(namespace).List(opts)
		if err != nil {
			return nil, errors.Wrap(err, "cannot retrieve secrets")
		}

		return list, nil
	}))
	p.PageSize = 500

	ctx := context.Background()

	secrets := []*v1.Secret{}

	err := p.EachListItem(ctx, metav1.ListOptions{}, func(obj runtime.Object) error {
		secret, ok := obj.(*v1.Secret)
		if !ok {
			return errors.Errorf("this is not a secret: %#v", obj)
		}

		secrets = append(secrets, secret)

		return nil
	})
	if err != nil {
		return []*v1.Secret{}, errors.Wrap(err, "cannot iterate secrets")
	}

	return secrets, nil
}

func main() {
	var kubecontext, namespace string

	flag.StringVar(&kubecontext, "context", "", "kubeconfig context")
	flag.StringVarP(&namespace, "namespace", "n", "default", "namespace")

	flag.Parse()

	fmt.Printf("context:   %s\n", kubecontext)
	fmt.Printf("namespace: %s\n", namespace)

	kubeconfig := clientcmd.RecommendedHomeFile

	clientConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfig},
		&clientcmd.ConfigOverrides{CurrentContext: kubecontext},
	)

	config, err := clientConfig.ClientConfig()
	if err != nil {
		panic(err)
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	k8sClient := &k8sClient{
		client: client,
	}

	secrets, err := k8sClient.ListSecrets(namespace)
	if err != nil {
		panic(err)
	}

	for _, secret := range secrets {
		fmt.Printf("%s\t%s\t%s\n", secret.Namespace, secret.Name, secret.Type)
	}

	pods, err := k8sClient.ListPods(namespace)
	if err != nil {
		panic(err)
	}

	for _, pod := range pods {
		fmt.Printf("%s\t%s\n", pod.Namespace, pod.Name)
	}
}
