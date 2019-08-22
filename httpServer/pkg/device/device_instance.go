package device

import (
    "github.com/gin-gonic/gin"
    "encoding/json"
    "fmt"
    "github.com/tidwall/gjson"
    "io/ioutil"
    v1 "k8s.io/api/core/v1"
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
    diName := gjson.GetBytes(body, "metadata.name").String()
    diNamespace := gjson.GetBytes(body, "metadata.namespace").String()
	diModel := gjson.GetBytes(body, "spec.deviceModelRef.name").String()
    fmt.Println(diName)
    fmt.Println(diNamespace)
	fmt.Println(diModel)

    device := &v1alpha1.Device{
		TypeMeta: metav1.TypeMeta {
			Kind: "Device",
			APIVersion: "devices.kubeedge.io/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta {
			Name: diName,
			Namespace: diNamespace,
		},
		Spec: v1alpha1.DeviceSpec {
			DeviceModelRef: &v1.LocalObjectReference{
				Name: "diModel",
			},
			Protocol:  v1alpha1.ProtocolConfig{
				Modbus: &v1alpha1.ProtocolConfigModbus{
					TCP: &v1alpha1.ProtocolConfigModbusTCP{
						IP: "192.168.226.139",
						Port: 5028,
						SlaveID: "1",
					},
				},
			},
			NodeSelector: &v1.NodeSelector{
				NodeSelectorTerms: []v1.NodeSelectorTerm{
					{
						MatchExpressions: []v1.NodeSelectorRequirement{
							{
								Key: "",
								Operator: "In",
								Values: []string{"fb4ebb70-2783-42b8-b3ef-63e2fd6d242e"},
							},
						},
					},
				},
			},
		},
		Status: v1alpha1.DeviceStatus{
			Twins: []v1alpha1.Twin{
				{
					PropertyName: "Temperature",
					Desired: v1alpha1.TwinProperty{
						Value: "12",
						Metadata:map[string]string{
							"timestamp": "1550049403598",
							"type": "string",
						},
					},
					Reported:v1alpha1.TwinProperty{
						Value: "12",
						Metadata:map[string]string{
							"timestamp": "1550049403598",
							"type": "string",
						},
					},
				},
			},
		},
    }

    //result := &v1alpha1.Device{}
    err := di.Client.Post().Namespace(diNamespace).Resource("devices").Body(device).Do().Error()
    if err != nil {
        fmt.Println(err)
        ctx.JSON(400, gin.H{
            "message": "Internal error",
            "code":  400,
            "reason":  err,
        })
    }
    ctx.JSON(200, device)
}

func (di *Device) GetDevice(ctx *gin.Context){
	body, _ := ioutil.ReadAll(ctx.Request.Body)
    diName := gjson.GetBytes(body, "metadata.name").String()
    diNamespace := gjson.GetBytes(body, "metadata.namespace").String()

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
    body, _ := ioutil.ReadAll(ctx.Request.Body)
    diName := gjson.GetBytes(body, "metadata.name").String()
    diNamespace := gjson.GetBytes(body, "metadata.namespace").String()
    fmt.Println(diName)
    fmt.Println(diNamespace)

    device := &v1alpha1.Device{}
    err := di.Client.Delete().Namespace(diNamespace).Resource("devices").Name(diName).Body(&metav1.DeleteOptions{}).Do().Into(device)
    if err != nil {
        fmt.Println(err)
        ctx.JSON(404, gin.H{
            "message": "DeviceModel or Namespace not exist",
            "code":  404,
            "reason":  err,
        })
    }

    ctx.JSON(200, device)
}

func (di *Device) UpdateDevice(ctx *gin.Context){
    body, _ := ioutil.ReadAll(ctx.Request.Body)
    diName := gjson.GetBytes(body, "metadata.name").String()
    diNamespace := gjson.GetBytes(body, "metadata.namespace").String()
    label_json := gjson.GetBytes(body, "metadata.Labels").String()
    label_map := make(map[string]string)
    err := json.Unmarshal([]byte(label_json), &label_map)
    if err != nil{
        fmt.Println("JsonToMapDemo err:", err)
    }

    device := &v1alpha1.Device{}
    err_g := di.Client.Get().Namespace(diNamespace).Resource("devices").Name(diName).Body(&metav1.GetOptions{}).Do().Into(device)
    if err_g != nil {
        fmt.Println(err)
        ctx.JSON(404, gin.H{
            "message": "DeviceModel or Namespace not exist",
            "code":  404,
            "reason":  err,
        })
        return
    }
    device.SetLabels(label_map)

    result := &v1alpha1.Device{}
    err_d := di.Client.Put().Namespace(diNamespace).Resource("devices").Name(diName).Body(device).Do().Into(result)
    if err_d != nil {
        ctx.JSON(400, gin.H{
                "message": "Internal error",
                "code":  400,
                "reason":  err,
            })
		}
    ctx.JSON(200, result)
}

