package main 

import (
	namespace "httpServer/pkg/namespace"
	node "httpServer/pkg/node"
	deployment "httpServer/pkg/deployment"
	devicemodel "httpServer/pkg/device_model"
	device "httpServer/pkg/device"
    "k8s.io/klog"
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

	ClientCrd, err := NewCRDClient(cfg)
    if err != nil {
        klog.Fatalf("Error building CRDClient error: %s", err.Error())
    }

    router := gin.Default()

    nsRouter := namespace.Namespace{ClientSet: kubeClient}
    nodeRouter := node.Node{ClientSet: kubeClient}
	deploymentRouter := deployment.Deployment{ClientSet: kubeClient}
	devicemodelRouter := devicemodel.DeviceModel{Client: ClientCrd}
	deviceRouter := device.Device{Client: ClientCrd}

    router.POST("/v1/datacenter/id/cluster/id/namespace", nsRouter.AddNamespace)
    router.DELETE("/v1/datacenter/id/cluster/id/namespace/delete", nsRouter.DeleteNamespace)
    router.GET("/v1/datacenter/id/cluster/id/namespace/list", nsRouter.ListNamespace)
    router.GET("/v1/datacenter/id/cluster/id/namespace/get", nsRouter.GetNamespace)
    router.PUT("/v1/datacenter/id/cluster/id/namespace/update", nsRouter.UpdateNamespace)

	router.POST("vl/datacenter/id/cluster/id/node/create", nodeRouter.AddNode)
	router.DELETE("vl/datacenter/id/cluster/id/node/delete", nodeRouter.DeleteNode)
	router.GET("vl/datacenter/id/cluster/id/node/list", nodeRouter.ListNode)
	router.GET("vl/datacenter/id/cluster/id/node/get", nodeRouter.GetNode)
	router.PUT("vl/datacenter/id/cluster/id/node/update", nodeRouter.UpdateNode)

	router.POST("vl/datacenter/id/cluster/id/deployment/create", deploymentRouter.AddDeployment)
	router.DELETE("vl/datacenter/id/cluster/id/deployment/delete", deploymentRouter.DeleteDeployment)
	router.GET("vl/datacenter/id/cluster/id/deployment/list", deploymentRouter.ListDeployment)
	router.GET("vl/datacenter/id/cluster/id/deployment/get", deploymentRouter.GetDeployment)
	router.PUT("vl/datacenter/id/cluster/id/deployment/update", deploymentRouter.UpdateDeployment)

	router.POST("vl/datacenter/id/cluster/id/devicemodel/create", devicemodelRouter.AddDeviceModel)
	router.GET("vl/datacenter/id/cluster/id/devicemodel/get", devicemodelRouter.GetDeviceModel)
	router.GET("vl/datacenter/id/cluster/id/devicemodel/list", devicemodelRouter.ListDeviceModel)
	router.DELETE("vl/datacenter/id/cluster/id/devicemodel/delete", devicemodelRouter.DeleteDeviceModel)
	router.PUT("vl/datacenter/id/cluster/id/devicemodel/update", devicemodelRouter.UpdateDeviceModel)

	router.POST("vl/datacenter/id/cluster/id/deviceinstance/create", deviceRouter.AddDevice)
	router.GET("vl/datacenter/id/cluster/id/deviceinstance/get", deviceRouter.GetDevice)
	router.GET("vl/datacenter/id/cluster/id/deviceinstance/list", deviceRouter.ListDevice)
	router.DELETE("vl/datacenter/id/cluster/id/deviceinstance/delete", deviceRouter.DeleteDevice)
	router.PUT("vl/datacenter/id/cluster/id/deviceinstance/update", deviceRouter.UpdateDevice)
    router.Run(":8000")
}

func init() {
    flag.StringVar(&kubeconfig, "kubeconfig", "", "Path to a kubeconfig. Only required if out-of-cluster.")
    flag.StringVar(&masterURL, "master", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluste        r.")
}


