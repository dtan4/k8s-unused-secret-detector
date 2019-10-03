package main

import (
	"context"
	"fmt"
	"log"

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

func detectUnusedSecrets(pods []*v1.Pod, secrets []*v1.Secret) ([]*v1.Secret, error) {
	usedSecretNames := map[string]bool{}

	for _, pod := range pods {
		for _, container := range pod.Spec.Containers {
			for _, envFrom := range container.EnvFrom {
				if envFrom.SecretRef != nil {
					usedSecretNames[envFrom.SecretRef.Name] = true
				}
			}

			for _, env := range container.Env {
				if env.ValueFrom != nil && env.ValueFrom.SecretKeyRef != nil {
					usedSecretNames[env.ValueFrom.SecretKeyRef.Name] = true
				}
			}
		}

		for _, volume := range pod.Spec.Volumes {
			if volume.Secret != nil {
				usedSecretNames[volume.Secret.SecretName] = true
			}
		}
	}

	unused := []*v1.Secret{}

	for _, secret := range secrets {
		if secret.Type != v1.SecretTypeOpaque {
			continue
		}

		if !usedSecretNames[secret.Name] {
			unused = append(unused, secret)
		}
	}

	return unused, nil
}

func main() {
	var kubecontext, namespace string

	flag.StringVar(&kubecontext, "context", "", "kubeconfig context")
	flag.StringVarP(&namespace, "namespace", "n", "default", "namespace")

	flag.Parse()

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

	log.Printf("Retrieving Pods in %s...", namespace)

	pods, err := k8sClient.ListPods(namespace)
	if err != nil {
		panic(err)
	}

	log.Printf("Retrieving Secrets in %s...", namespace)

	secrets, err := k8sClient.ListSecrets(namespace)
	if err != nil {
		panic(err)
	}

	log.Printf("Detecting unused Secrets in %s...", namespace)

	unused, err := detectUnusedSecrets(pods, secrets)
	if err != nil {
		panic(err)
	}

	for _, secret := range unused {
		fmt.Println(secret.Name)
	}
}
