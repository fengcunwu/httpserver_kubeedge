package main 

import (
	//namespace "httpServer/pkg/namespace"
	node "httpServer/pkg/node"
	"k8s.io/client-go/tools/clientcmd"
    "k8s.io/klog"
	"flag"
	"github.com/gin-gonic/gin"
	"k8s.io/client-go/kubernetes"
)

var (
    masterURL  string
    kubeconfig string
)


func main(){
    klog.InitFlags(nil)
    flag.Parse()
    //stopCh := signals.SetupSignalHandler()

    cfg, err := clientcmd.BuildConfigFromFlags(masterURL, kubeconfig)
    if err != nil {
        klog.Fatalf("Error building kubeconfig: %s", err.Error())
    }

    kubeClient, err := kubernetes.NewForConfig(cfg)
    if err != nil {
        klog.Fatalf("Error building kubernetes clientset: %s", err.Error())
    }

    router := gin.Default()
    //nsRouter := namespace.Namespace{ClientSet: kubeClient}
	nodeRouter := node.Node{ClientSet: kubeClient}

    //router.POST("/v1/datacenter/id/cluster/id/namespace", nsRouter.AddNamespace)
    //router.DELETE("/v1/datacenter/id/cluster/id/namespace/delete", nsRouter.DeleteNamespace)
    //router.GET("/v1/datacenter/id/cluster/id/namespace/list", nsRouter.ListNamespace)
    //router.GET("/v1/datacenter/id/cluster/id/namespace/get", nsRouter.GetNamespace)
    //router.PUT("/v1/datacenter/id/cluster/id/namespace/update", nsRouter.UpdateNamespace)

	router.POST("vl/datacenter/id/cluster/id/node/create", nodeRouter.AddNode)
	router.DELETE("vl/datacenter/id/cluster/id/node/delete", nodeRouter.DeleteNode)
	router.GET("vl/datacenter/id/cluster/id/node/list", nodeRouter.ListNode)
	router.GET("vl/datacenter/id/cluster/id/node/get", nodeRouter.GetNode)
	router.PUT("vl/datacenter/id/cluster/id/node/update", nodeRouter.UpdateNode)
    router.Run(":8000")
}

func init() {
    flag.StringVar(&kubeconfig, "kubeconfig", "", "Path to a kubeconfig. Only required if out-of-cluster.")
    flag.StringVar(&masterURL, "master", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluste        r.")
}


