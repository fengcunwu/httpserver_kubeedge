package main

import (
	namespace "httpServer/pkg/namespace"
	node "httpServer/pkg/node"
	deployment "httpServer/pkg/deployment"
	devicemodel "httpServer/pkg/device_model"
	device "httpServer/pkg/device"
    "k8s.io/klog"
	"github.com/gin-contrib/cors"
	"flag"
	"github.com/gin-gonic/gin"
	//"k8s.io/client-go/kubernetes"
	//"github.com/kubeedge/beehive/pkg/common/log"
	//"github.com/kubeedge/kubeedge/cloud/pkg/apis/devices/v1alpha1"	
	v1alpha1 "httpServer/pkg/apis/devices/v1alpha1"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

var (
    masterURL  string
    kubeconfig string
)

func NewCRDClient(cfg *rest.Config) (*rest.RESTClient, error) {
	scheme := runtime.NewScheme()
	schemeBuilder := runtime.NewSchemeBuilder(addDeviceCrds)

	err := schemeBuilder.AddToScheme(scheme)
	if err != nil {
		return nil, err
	}

	config := *cfg
	config.APIPath = "/apis"
	config.GroupVersion = &v1alpha1.SchemeGroupVersion
	config.ContentType = runtime.ContentTypeJSON
	config.NegotiatedSerializer = serializer.DirectCodecFactory{CodecFactory: serializer.NewCodecFactory(scheme)}

	client, err := rest.RESTClientFor(&config)
	if err != nil {
		//log.LOGGER.Errorf("Failed to create REST Client due to error %v", err)
		return nil, err
	}

	return client, nil
}

func addDeviceCrds(scheme *runtime.Scheme) error {
	// Add Device
	scheme.AddKnownTypes(v1alpha1.SchemeGroupVersion, &v1alpha1.Device{}, &v1alpha1.DeviceList{})
	v1.AddToGroupVersion(scheme, v1alpha1.SchemeGroupVersion)
	// Add DeviceModel
	scheme.AddKnownTypes(v1alpha1.SchemeGroupVersion, &v1alpha1.DeviceModel{}, &v1alpha1.DeviceModelList{})
	v1.AddToGroupVersion(scheme, v1alpha1.SchemeGroupVersion)
	return nil
}

func main(){
    klog.InitFlags(nil)
    flag.Parse()
    //stopCh := signals.SetupSignalHandler()

    cfg, err := clientcmd.BuildConfigFromFlags(masterURL, kubeconfig)
    if err != nil {
        klog.Fatalf("Error building kubeconfig: %s", err.Error())
    }

    kubeClient, err := kubernetes.NewForConfig(cfg)
	//rest.RESTClientFor(&config)
    if err != nil {
        klog.Fatalf("Error building kubernetes clientset: %s", err.Error())
    }

	crdClient, err := NewCRDClient(cfg)
    if err != nil {
        klog.Fatalf("Error building CRDClient error: %s", err.Error())
    }

    router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "PATCH", "DELETE", "POST", "GET"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
	}))

    nsRouter := namespace.Namespace{ClientSet: kubeClient}
    nodeRouter := node.Node{ClientSet: kubeClient}
	deploymentRouter := deployment.Deployment{ClientSet: kubeClient}
	devicemodelRouter := devicemodel.DeviceModel{Client: crdClient}
	deviceRouter := device.Device{Client: crdClient}

    router.POST("v1/datacenter/id/cluster/id/namespace/create", nsRouter.AddNamespace)
    router.DELETE("v1/datacenter/id/cluster/id/namespace/delete", nsRouter.DeleteNamespace)
    router.GET("v1/datacenter/id/cluster/id/namespace/list", nsRouter.ListNamespace)
    router.GET("v1/datacenter/id/cluster/id/namespace/get", nsRouter.GetNamespace)
    router.PUT("v1/datacenter/id/cluster/id/namespace/update", nsRouter.UpdateNamespace)

	router.POST("v1/datacenter/id/cluster/id/node", nodeRouter.AddNode)
	router.DELETE("v1/datacenter/id/cluster/id/node/:name", nodeRouter.DeleteNode)
	router.GET("v1/datacenter/id/cluster/id/nodes", nodeRouter.ListNode)
	router.GET("v1/datacenter/id/cluster/id/node/:name", nodeRouter.GetNode)
	router.PUT("v1/datacenter/id/cluster/id/node", nodeRouter.UpdateNode)

	router.POST("v1/cluster/id/deployment", deploymentRouter.AddDeployment)
	router.DELETE("v1/cluster/id/namespace/:namespace/deployment/:name", deploymentRouter.DeleteDeployment)
	router.GET("v1/cluster/id/namespace/:namespace/deployments", deploymentRouter.ListDeployment)
	router.GET("v1/cluster/id/namespace/:namespace/deployment/:name", deploymentRouter.GetDeployment)
	router.PUT("v1/cluster/id/deployment", deploymentRouter.UpdateDeployment)

	router.POST("v1/cluster/id/devicemodel", devicemodelRouter.AddDeviceModel)
	router.GET("v1/cluster/id/namespace/:namespace/devicemodel/:name", devicemodelRouter.GetDeviceModel)
	router.GET("v1/cluster/id/namespace/:namespace/devicemodels", devicemodelRouter.ListDeviceModel)
	router.DELETE("v1/cluster/id/namespace/:namespace/devicemodel/:name", devicemodelRouter.DeleteDeviceModel)
	router.PUT("v1/cluster/id/devicemodel", devicemodelRouter.UpdateDeviceModel)

	router.POST("v1/cluster/id/device", deviceRouter.AddDevice)
	router.GET("v1/cluster/id/namespace/:namespace/device/:name", deviceRouter.GetDevice)
	router.GET("v1/cluster/id/namespace/:namespace/device", deviceRouter.ListDevice)
	router.DELETE("v1/cluster/id/namespace/:namespace/device/:name", deviceRouter.DeleteDevice)
	router.PUT("v1/cluster/id/device", deviceRouter.UpdateDevice)

    router.Run(":8000")
}

func init() {
    flag.StringVar(&kubeconfig, "kubeconfig", "", "Path to a kubeconfig. Only required if out-of-cluster.")
    flag.StringVar(&masterURL, "master", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluste        r.")
}


