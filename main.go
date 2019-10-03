package main

import (
	"fmt"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	kubeconfig := clientcmd.RecommendedHomeFile
	// context := ""
	namespace := "default"

	clientConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfig},
		// &clientcmd.ConfigOverrides{CurrentContext: context},
		&clientcmd.ConfigOverrides{},
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
