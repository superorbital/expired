package main

import (
	"fmt"
	"context"
	"path/filepath"

	//"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	coreV1 "k8s.io/api/core/v1"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/util/homedir"
	"k8s.io/apimachinery/pkg/watch"
	k8srestclient "k8s.io/client-go/rest"
	k8sclientcmd "k8s.io/client-go/tools/clientcmd"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

func loadConfig() (*k8srestclient.Config, error) {
	// try to load an in cluster config, fall back to ~/.kube/config file.
	config, err := k8srestclient.InClusterConfig()
	if err != nil {

		kubeconfig := filepath.Join(homedir.HomeDir(), ".kube", "config")
		config, err := k8sclientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, err
		}
		return config, nil
	}

	return config, nil
}


func watchSecrets(c *kubernetes.Clientset) <-chan watch.Event {
	api := c.CoreV1().Secrets("")
	secrets, err := api.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	resourceVersion := secrets.ListMeta.ResourceVersion

	watcher, err := api.Watch(context.TODO(), metav1.ListOptions{ResourceVersion: resourceVersion})

	if err != nil {
		panic(err.Error())
	}

	return watcher.ResultChan()
}

func main() {
	fmt.Println("initializing")
	config, err := loadConfig()
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("creating client")
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("creating channel")
	channel := watchSecrets(clientset)


	fmt.Println("watching for secrets")
	for {

		event := <-channel
		secret, ok := event.Object.(*coreV1.Secret)
		if !ok {
			panic("Could not cast to Secret")
		}

		fmt.Printf("%v\n", secret.ObjectMeta.Name)
		fmt.Printf("%+v\n", secret)
	}

	/*pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("Found %d pods.\n", len(pods.Items))
	*/
}
