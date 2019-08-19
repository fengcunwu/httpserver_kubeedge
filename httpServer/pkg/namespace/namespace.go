package namespace

import (
    "github.com/gin-gonic/gin"
    "encoding/json"
    //"net/http"
    "fmt"
    "github.com/tidwall/gjson"
    "io/ioutil"
    corev1 "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/kubernetes"
    //"Management_Portal_API/pkg/signals"
)

type Namespace struct{
    ClientSet *kubernetes.Clientset
}

func (n *Namespace) AddNamespace(ctx *gin.Context){
    body, _ := ioutil.ReadAll(ctx.Request.Body)
    nsName := gjson.GetBytes(body, "metadata.name").String()
    fmt.Println(nsName)

    namespace := &corev1.Namespace{
        ObjectMeta: metav1.ObjectMeta{
            Name: nsName,
        },
    }

    namespace, err := n.ClientSet.CoreV1().Namespaces().Create(namespace)
    if err != nil {
        fmt.Println(err)
        ctx.JSON(400, gin.H{
            "message": "Internal error",
            "code":  400,
            "reason":  err,
        })
    }

    ctx.JSON(200, namespace)
}

func (n *Namespace) DeleteNamespace(ctx *gin.Context){
    body, _ := ioutil.ReadAll(ctx.Request.Body)
    nsName := gjson.GetBytes(body, "metadata.name").String()
    fmt.Println(nsName)

    namespace, err := n.ClientSet.CoreV1().Namespaces().Get(nsName, metav1.GetOptions{})
    if err != nil{
        ctx.JSON(404, gin.H{
            "message": "Namespace not exist",
            "code":  404,
            "reason":  err,
        })
        return
    }

    deletePolicy := metav1.DeletePropagationForeground
    err_d := n.ClientSet.CoreV1().Namespaces().Delete(nsName, &metav1.DeleteOptions{
        PropagationPolicy: &deletePolicy,
    })

    if err_d != nil {
        fmt.Println(err_d)
        ctx.JSON(400, gin.H{
            "message": "Internal error",
            "code":  400,
            "reason":  err_d,
        })
    }

    ctx.JSON(200, namespace)
}

func (n *Namespace) ListNamespace(ctx *gin.Context){
    nameSpaceList, err := n.ClientSet.CoreV1().Namespaces().List(metav1.ListOptions{})
    if err != nil {
        ctx.JSON(404, gin.H{
            "message": "Namespace not exist",
            "code":  404,
            "reason":  err,
        })
        return
    }

    ctx.JSON(200, nameSpaceList)
}

func (n *Namespace) GetNamespace(ctx *gin.Context){
    body, _ := ioutil.ReadAll(ctx.Request.Body)
    nsName := gjson.GetBytes(body, "metadata.name").String()

    nameSpaceSrc, err := n.ClientSet.CoreV1().Namespaces().Get(nsName, metav1.GetOptions{})
    if err != nil {
        ctx.JSON(404, gin.H{
            "message": "Namespace not exist",
            "code":  404,
            "reason":  err,
        })
        return
    }

    ctx.JSON(200, nameSpaceSrc)
}

func (n *Namespace) UpdateNamespace(ctx *gin.Context){
    body, _ := ioutil.ReadAll(ctx.Request.Body)
    nsName := gjson.GetBytes(body, "metadata.name").String()
    label_json := gjson.GetBytes(body, "metadata.Labels").String()
    label_map := make(map[string]string)
    err := json.Unmarshal([]byte(label_json), &label_map)
    if err != nil{
        fmt.Println("JsonToMapDemo err:", err)
    }

    namespace, err := n.ClientSet.CoreV1().Namespaces().Get(nsName, metav1.GetOptions{})
    if err != nil{
        ctx.JSON(404, gin.H{
            "message": "Namespace not exist",
            "code":  404,
            "reason":  err,
        })
        return
    }

    namespace.SetLabels(label_map)

    nameNew, err := n.ClientSet.CoreV1().Namespaces().Update(namespace)
    if err != nil {
        ctx.JSON(400, gin.H{
                "message": "Internal error",
                "code":  400,
                "reason":  err,
            })
        }
    ctx.JSON(200, nameNew)
}
