package device_model

import (
    "github.com/gin-gonic/gin"
	"github.com/json-iterator/go"
    "fmt"
    //"github.com/tidwall/gjson"
    "io/ioutil"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	rest "k8s.io/client-go/rest"
	v1alpha1 "httpServer/pkg/apis/devices/v1alpha1"
)

type DeviceModel struct{
    Client *rest.RESTClient
}

func (dm *DeviceModel) AddDeviceModel(ctx *gin.Context){
    body, _ := ioutil.ReadAll(ctx.Request.Body)

	devicemodel := &v1alpha1.DeviceModel{}
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	json.Unmarshal(body,devicemodel)

    dmName := devicemodel.GetName()
    dmNamespace := devicemodel.GetNamespace()
    fmt.Println(dmName)
    fmt.Println(dmNamespace)

	err := dm.Client.Post().Namespace(dmNamespace).Resource("devicemodels").Body(devicemodel).Do().Into(devicemodel)
    if err != nil {
        fmt.Println(err)
        ctx.JSON(400, gin.H{
            "message": "Internal error",
            "code":  400,
            "reason":  err,
        })
    }
    ctx.JSON(200, devicemodel)
}

func (dm *DeviceModel) GetDeviceModel(ctx *gin.Context){
    dmName := ctx.Param("name")
    dmNamespace := ctx.Param("namespace")
    fmt.Println(dmName)
    fmt.Println(dmNamespace)

    result := &v1alpha1.DeviceModel{}
	err := dm.Client.Get().Namespace(dmNamespace).Resource("devicemodels").Name(dmName).Body(&metav1.GetOptions{}).Do().Into(result)
    if err != nil {
        fmt.Println(err)
        ctx.JSON(404, gin.H{
            "message": "DeviceModel or Namespace not exist",
            "code":  404,
            "reason":  err,
        })
    }

    ctx.JSON(200, result)
}

func (dm *DeviceModel) ListDeviceModel(ctx *gin.Context){
    //dmNamespace := ctx.Param("namespace")
    //fmt.Println(dmNamespace)

	result := &v1alpha1.DeviceModelList{}
    err := dm.Client.Get().Resource("devicemodels").Body(&metav1.GetOptions{}).Do().Into(result)
    if err != nil {
        fmt.Println(err)
        ctx.JSON(404, gin.H{
            "message": "DeviceModel or Namespace not exist",
            "code":  404,
            "reason":  err,
        })
    }

    ctx.JSON(200, result)
}

func (dm *DeviceModel) DeleteDeviceModel(ctx *gin.Context){
    dmName := ctx.Param("name")
    dmNamespace := ctx.Param("namespace")
    fmt.Println(dmName)
    fmt.Println(dmNamespace)

    devicemodel := &v1alpha1.DeviceModel{}
    err := dm.Client.Get().Namespace(dmNamespace).Resource("devicemodels").Name(dmName).Body(&metav1.GetOptions{}).Do().Error()
    if err != nil {
        fmt.Println(err)
        ctx.JSON(404, gin.H{
            "message": "DeviceModel or Namespace not exist",
            "code":  404,
            "reason":  err,
        })
    }

    err = dm.Client.Delete().Namespace(dmNamespace).Resource("devicemodels").Name(dmName).Body(&metav1.DeleteOptions{}).Do().Error()
    if err != nil {
        fmt.Println(err)
        ctx.JSON(500, gin.H{
            "message": "Internal error",
            "code":  500,
            "reason":  err,
        })
    }

    ctx.JSON(200, devicemodel)
}

func (dm *DeviceModel) UpdateDeviceModel(ctx *gin.Context){
    body, _ := ioutil.ReadAll(ctx.Request.Body)

	devicemodel := &v1alpha1.DeviceModel{}
    var json = jsoniter.ConfigCompatibleWithStandardLibrary
    json.Unmarshal(body,devicemodel)
    fmt.Println(devicemodel)

	dmName := devicemodel.GetName()
	dmNamespace := devicemodel.GetNamespace()

    err := dm.Client.Get().Namespace(dmNamespace).Resource("devicemodels").Name(dmName).Body(&metav1.GetOptions{}).Do().Error()
    if err != nil {
        fmt.Println(err)
        ctx.JSON(404, gin.H{
            "message": "DeviceModel or Namespace not exist",
            "code":  404,
            "reason":  err,
        })
		return
    }

    devicemodel.SetLabels(devicemodel.GetLabels())

    err = dm.Client.Put().Namespace(dmNamespace).Resource("devicemodels").Name(dmName).Body(devicemodel).Do().Error()
    if err != nil {
        ctx.JSON(400, gin.H{
                "message": "Internal error",
                "code":  400,
                "reason":  err,
            })
        }

    ctx.JSON(200, devicemodel)
}
