package main

import (
	"os"
	"time"
	"log"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/cache"
)

type Client struct {
	dynamicClient dynamic.Interface
}
var (
    InfoLogger    *log.Logger
)

func init() {
    file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
    if err != nil {
        log.Fatal(err)
    }

    InfoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func main(){
	kubeconfig := os.Getenv("HOME") + "/.crc/machines/crc/kubeconfig"
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	clientset, err := dynamic.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	ch := make(chan struct{})
	defer close(ch)
	informer := dynamicinformer.NewFilteredDynamicSharedInformerFactory(
		clientset,
		1 * time.Minute,
		"memcached-operator-system",
		nil,
	)
	var handler cache.ResourceEventHandlerFuncs
	handler.AddFunc = func(obj interface{}) {
		InfoLogger.Println("add event")
	}
	handler.UpdateFunc = func(old, new interface{}) {
		oldObj := old.(*unstructured.Unstructured)
		//newObj := new.(*unstructured.Unstructured)

        /*data, _, err := unstructured.NestedMap(oldObj.Object, "spec", "template","spec","containers")
		if err != nil {
			InfoLogger.Println(err.Error())
			
		}*/
		InfoLogger.Println("update event")
	}
	handler.DeleteFunc = func(obj interface{}) {
		u := obj.(*unstructured.Unstructured)

		InfoLogger.Println(u.GetName())
	}

	gvr, _ := schema.ParseResourceArg("deployments.v1.apps")
	i := informer.ForResource(*gvr)
	i.Informer().AddEventHandler(handler)
	i.Informer().Run(ch)

}