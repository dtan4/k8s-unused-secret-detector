package main

import (
	"context"
	"fmt"

	flag "github.com/spf13/pflag"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/pager"
)

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

	p := pager.New(pager.SimplePageFunc(func(opts metav1.ListOptions) (runtime.Object, error) {
		list, err := client.CoreV1().Secrets(namespace).List(opts)
		if err != nil {
			return nil, err
		}

		return list, err
	}))
	p.PageSize = 500

	ctx := context.Background()

	secrets := []*v1.Secret{}

	err = p.EachListItem(ctx, metav1.ListOptions{}, func(obj runtime.Object) error {
		secret, ok := obj.(*v1.Secret)
		if !ok {
			return err
		}

		secrets = append(secrets, secret)

		return nil
	})
	if err != nil {
		panic(err)
	}

	for _, secret := range secrets {
		fmt.Printf("%s\t%s\t%s\n", secret.Namespace, secret.Name, secret.Type)
	}
}
