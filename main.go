package main

import (
	"log"
	"os"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	cache "k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	// log.Printf("informer app started")

	kubeconfig := os.Getenv("KUBECONFIG")

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatalf("Error building kubeconfig: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)

	if err != nil {
		log.Fatalf("Error building clientset: %v", err)
	}

	factory := informers.NewSharedInformerFactory(clientset, 0)

	informer := factory.Core().V1().Pods().Informer()

	stopper := make(chan struct{})

	defer close(stopper)
	defer runtime.HandleCrash()

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{})

	go informer.Run(stopper)

	if !cache.WaitForCacheSync(stopper, informer.HasSynced) {
		log.Fatalf("Error waiting for caches to sync")
		return
	}

	<-stopper

}

func onUpdate(old, new interface{}) {

	oldPods := old.(*v1.Pod)
	newPods := new.(*v1.Pod)

	log.Printf("Pod Updated: %v --> %v", oldPods.ResourceVersion, newPods.ResourceVersion)
}
