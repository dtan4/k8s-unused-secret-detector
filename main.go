package main

import (
	"fmt"

	flag "github.com/spf13/pflag"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	var context, namespace string

	flag.StringVar(&context, "context", "", "kubeconfig context")
	flag.StringVarP(&namespace, "namespace", "n", "default", "namespace")

	flag.Parse()

	fmt.Printf("context:   %s\n", context)
	fmt.Printf("namespace: %s\n", namespace)

	kubeconfig := clientcmd.RecommendedHomeFile

	clientConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfig},
		&clientcmd.ConfigOverrides{CurrentContext: context},
	)

	config, err := clientConfig.ClientConfig()
	if err != nil {
		panic(err)
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	secrets, err := client.CoreV1().Secrets(namespace).List(metav1.ListOptions{})
	if errors.IsNotFound(err) {
		panic("secrets not found")
	}

	if err != nil {
		panic(err)
	}

	for _, secret := range secrets.Items {
		fmt.Printf("%s\t%s\t%s\n", secret.Namespace, secret.Name, secret.Type)
	}
}
