package device

import (
    "github.com/gin-gonic/gin"
    "fmt"
	"github.com/json-iterator/go"
    "io/ioutil"
    //v1 "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    //"k8s.io/client-go/kubernetes"
    rest "k8s.io/client-go/rest"
    v1alpha1 "httpServer/pkg/apis/devices/v1alpha1"
)

type Device struct{
	Client *rest.RESTClient
}

func (di *Device) AddDevice(ctx *gin.Context){
    body, _ := ioutil.ReadAll(ctx.Request.Body)

	device := &v1alpha1.Device{}
    var json = jsoniter.ConfigCompatibleWithStandardLibrary
    json.Unmarshal(body,device)
    fmt.Println(device)

	diName := device.GetName()
	diNamespace := device.GetNamespace()
	fmt.Println(diName)
	fmt.Println(diNamespace)

    err := di.Client.Post().Namespace(diNamespace).Resource("devices").Body(device).Do().Error()
    if err != nil {
        fmt.Println(err)
        ctx.JSON(500, gin.H{
            "message": "Internal error",
            "code":  500,
            "reason":  err,
        })
        return
    }

    ctx.JSON(200, device)
}

func (di *Device) GetDevice(ctx *gin.Context){
    diName := ctx.Param("name")
    diNamespace := ctx.Param("namespace")
	fmt.Println(diName)
	fmt.Println(diNamespace)

	device := &v1alpha1.Device{}
	err := di.Client.Get().Namespace(diNamespace).Resource("devices").Name(diName).Body(&metav1.GetOptions{}).Do().Into(device)
	if err != nil {
        fmt.Println(err)
        ctx.JSON(404, gin.H{
            "message": "Namespace or Device not exist",
            "code":  404,
            "reason":  err,
        })
    }

    ctx.JSON(200, device)
}

func (di *Device) ListDevice(ctx *gin.Context){
	//diNamespace := ctx.Param("namespace")

	deviceList := &v1alpha1.DeviceList{}
    err := di.Client.Get().Resource("devices").Body(&metav1.GetOptions{}).Do().Into(deviceList)
    if err != nil {
        fmt.Println(err)
        ctx.JSON(404, gin.H{
            "message": "Device not exist",
            "code":  404,
            "reason":  err,
        })
    }

    ctx.JSON(200, deviceList)
}

func (di *Device) DeleteDevice(ctx *gin.Context){
    diName := ctx.Param("name")
    diNamespace := ctx.Param("namespace")
    fmt.Println(diName)
    fmt.Println(diNamespace)

    err := di.Client.Get().Namespace(diNamespace).Resource("devices").Name(diName).Body(&metav1.GetOptions{}).Do().Error()
    if err != nil {
        fmt.Println(err)
        ctx.JSON(404, gin.H{
            "message": "DeviceModel or Namespace not exist",
            "code":  404,
            "reason":  err,
        })
    }

	device := &v1alpha1.Device{}
    err = di.Client.Delete().Namespace(diNamespace).Resource("devices").Name(diName).Body(&metav1.DeleteOptions{}).Do().Into(device)
    if err != nil {
        fmt.Println(err)
        ctx.JSON(500, gin.H{
            "message": "Internal error",
            "code":  500,
            "reason":  err,
        })
    }
    ctx.JSON(200, device)
}

func (di *Device) UpdateDevice(ctx *gin.Context){

    body, _ := ioutil.ReadAll(ctx.Request.Body)
    device := &v1alpha1.Device{}
    var json = jsoniter.ConfigCompatibleWithStandardLibrary
    json.Unmarshal(body,device)
    fmt.Println(device)

	diName := device.GetName()
	diNamespace := device.GetNamespace()
	fmt.Println(diName)
	fmt.Println(diNamespace)

    err := di.Client.Get().Namespace(diNamespace).Resource("devices").Name(diName).Body(&metav1.GetOptions{}).Do().Into(device)
    if err != nil {
        ctx.JSON(404, gin.H{
            "message": "Device or Namespace not exist",
            "code":  404,
            "reason":  err,
        })
        return
    }

    device.SetLabels(device.GetLabels())

    err = di.Client.Put().Namespace(diNamespace).Resource("devices").Name(diName).Body(device).Do().Error()
    if err != nil {
        ctx.JSON(500, gin.H{
            "message": "Internal error",
            "code":  500,
            "reason":  err,
        })
    }
    ctx.JSON(200, device)
}
