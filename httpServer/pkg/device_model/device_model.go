package device_model

import (
    "github.com/gin-gonic/gin"
    "encoding/json"
    //"net/http"
    "fmt"
    "github.com/tidwall/gjson"
    "io/ioutil"
    //corev1 "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    //"k8s.io/client-go/kubernetes"
	//scheme "k8s.io/client-go/kubernetes/scheme"
	rest "k8s.io/client-go/rest"
	v1alpha1 "httpServer/pkg/apis/devices/v1alpha1"
)

type DeviceModel struct{
    Client *rest.RESTClient
}

func (dm *DeviceModel) AddDeviceModel(ctx *gin.Context){
    body, _ := ioutil.ReadAll(ctx.Request.Body)
    dmName := gjson.GetBytes(body, "metadata.name").String()
	dmNamespace := gjson.GetBytes(body, "metadata.namespace").String() 
    fmt.Println(dmName)
	fmt.Println(dmNamespace)

    devicemodel := &v1alpha1.DeviceModel{
        ObjectMeta: metav1.ObjectMeta{
            Name: dmName,
        },
    }

	result := &v1alpha1.DeviceModel{}
    err := dm.Client.Post().Namespace(dmNamespace).Resource("devicemodels").Body(devicemodel).Do().Into(result)
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
    body, _ := ioutil.ReadAll(ctx.Request.Body)
    dmName := gjson.GetBytes(body, "metadata.name").String()
    dmNamespace := gjson.GetBytes(body, "metadata.namespace").String()
    fmt.Println(dmName)
    fmt.Println(dmNamespace)

    result := &v1alpha1.DeviceModel{}
	err := dm.Client.Get().Namespace(dmNamespace).Resource("devicemodels").Name(dmName).Body(&metav1.GetOptions{}).Do().Into(result)
    if err != nil {
        fmt.Println(err)
        ctx.JSON(404, gin.H{
            "message": "DeviceModel not exist",
            "code":  404,
            "reason":  err,
        })
    }

    ctx.JSON(200, result)
}

func (dm *DeviceModel) ListDeviceModel(ctx *gin.Context){
    result := &v1alpha1.DeviceModelList{}
    err := dm.Client.Get().Resource("devicemodels").Body(&metav1.GetOptions{}).Do().Into(result)
    if err != nil {
        fmt.Println(err)
        ctx.JSON(404, gin.H{
            "message": "DeviceModel not exist",
            "code":  404,
            "reason":  err,
        })
    }

    ctx.JSON(200, result)
}

func (dm *DeviceModel) DeleteDeviceModel(ctx *gin.Context){
    body, _ := ioutil.ReadAll(ctx.Request.Body)
    dmName := gjson.GetBytes(body, "metadata.name").String()
    dmNamespace := gjson.GetBytes(body, "metadata.namespace").String()
    fmt.Println(dmName)
    fmt.Println(dmNamespace)

    result := &v1alpha1.DeviceModel{}
    err := dm.Client.Delete().Namespace(dmNamespace).Resource("devicemodels").Name(dmName).Body(&metav1.DeleteOptions{}).Do().Into(result)
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

func (dm *DeviceModel) UpdateDeviceModel(ctx *gin.Context){
    body, _ := ioutil.ReadAll(ctx.Request.Body)
    dmName := gjson.GetBytes(body, "metadata.name").String()
    dmNamespace := gjson.GetBytes(body, "metadata.namespace").String()
    label_json := gjson.GetBytes(body, "metadata.Labels").String()
    label_map := make(map[string]string)
    err := json.Unmarshal([]byte(label_json), &label_map)
    if err != nil{
        fmt.Println("JsonToMapDemo err:", err)
    }

	devicemodel := &v1alpha1.DeviceModel{}
    err_g := dm.Client.Get().Namespace(dmNamespace).Resource("devicemodels").Name(dmName).Body(&metav1.GetOptions{}).Do().Into(devicemodel)
    if err_g != nil {
        fmt.Println(err)
        ctx.JSON(404, gin.H{
            "message": "DeviceModel or Namespace not exist",
            "code":  404,
            "reason":  err,
        })
		return
    }
    devicemodel.SetLabels(label_map)

	result := &v1alpha1.DeviceModel{}
    err_d := dm.Client.Put().Namespace(dmNamespace).Resource("devicemodels").Name(dmName).Body(devicemodel).Do().Into(result)
    if err_d != nil {
        ctx.JSON(400, gin.H{
                "message": "Internal error",
                "code":  400,
                "reason":  err,
            })
        }
    ctx.JSON(200, result)
}
