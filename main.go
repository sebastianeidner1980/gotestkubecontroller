package main

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/seidner/app-controller/pkg/config"
)

type Client struct {
	dynamicClient dynamic.Interface
}

var (
	Logger = logrus.New()
	cfg    = config.NewConfiguration()
)

func init() {
	Logger.Out = os.Stdout

	loglevel, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		loglevel = logrus.InfoLevel
	}
	Logger.SetLevel(loglevel)
	Logger.Formatter = &logrus.JSONFormatter{
		PrettyPrint:     true,
		TimestampFormat: "2006-01-02 15:04:05",
	}
	Logger.Debug("Logging successful launched")
}

func main() {

	kubeconfig := os.Getenv("HOME") + "/.kube/yigit_kubeconfig"
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
		1*time.Minute,
		"memcached-operator-system",
		nil,
	)
	var handler cache.ResourceEventHandlerFuncs
	handler.AddFunc = func(obj interface{}) {
		Logger.Debug("add event")
	}
	handler.UpdateFunc = func(old, new interface{}) {
		oldObj := old.(*unstructured.Unstructured)
		newObj := new.(*unstructured.Unstructured)
		oldReplica, _, _ := unstructured.NestedInt64(oldObj.Object, "spec", "replicas")
		newReplica, _, _ := unstructured.NestedInt64(newObj.Object, "spec", "replicas")
		if oldReplica != newReplica {
			Logger.Debug("Replicas is changed: %v\n", newReplica)
		}

	}
	handler.DeleteFunc = func(obj interface{}) {
		u := obj.(*unstructured.Unstructured)

		Logger.Debug(u.GetName())
	}

	gvr, _ := schema.ParseResourceArg("deployments.v1.apps")
	i := informer.ForResource(*gvr)
	i.Informer().AddEventHandler(handler)
	i.Informer().Run(ch)

}
