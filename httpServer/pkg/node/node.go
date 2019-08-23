package node

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

type Node struct{
    ClientSet *kubernetes.Clientset
}

func (n *Node) AddNode(ctx *gin.Context){
    body, _ := ioutil.ReadAll(ctx.Request.Body)
    nodeName := gjson.GetBytes(body, "metadata.name").String()
	nodeLabels := gjson.GetBytes(body, "metadata.Labels.name").String()
    fmt.Println(nodeName)
	fmt.Println(nodeLabels)

    node := &corev1.Node{
        ObjectMeta: metav1.ObjectMeta{
            Name: nodeName,
			Labels:map[string]string{
				"name": nodeLabels,
			},
        },
    }

    node, err := n.ClientSet.CoreV1().Nodes().Create(node)
    if err != nil {
        fmt.Println(err)
        ctx.JSON(400, gin.H{
            "message": "Internal error",
            "code":  400,
            "reason":  err,
        })
    }

    ctx.JSON(200, node)
}

func (n *Node) DeleteNode(ctx *gin.Context){
    body, _ := ioutil.ReadAll(ctx.Request.Body)
    nodeName := gjson.GetBytes(body, "metadata.name").String()
    fmt.Println(nodeName)

    node, err := n.ClientSet.CoreV1().Nodes().Get(nodeName, metav1.GetOptions{})
    if err != nil{
        ctx.JSON(404, gin.H{
            "message": "Node not exist",
            "code":  404,
            "reason":  err,
        })
        return
    }

    deletePolicy := metav1.DeletePropagationForeground
    err_d := n.ClientSet.CoreV1().Nodes().Delete(nodeName, &metav1.DeleteOptions{
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

    ctx.JSON(200, node)
}

func (n *Node) ListNode(ctx *gin.Context){
    nodeList, err := n.ClientSet.CoreV1().Nodes().List(metav1.ListOptions{})
    if err != nil {
        ctx.JSON(404, gin.H{
            "message": "Node not exist",
            "code":  404,
            "reason":  err,
        })
        return
    }

    ctx.JSON(200, nodeList)
}

func (n *Node) GetNode(ctx *gin.Context){
    body, _ := ioutil.ReadAll(ctx.Request.Body)
    nodeName := gjson.GetBytes(body, "metadata.name").String()

    nodeSrc, err := n.ClientSet.CoreV1().Nodes().Get(nodeName, metav1.GetOptions{})
    if err != nil {
        ctx.JSON(404, gin.H{
            "message": "Node not exist",
            "code":  404,
            "reason":  err,
        })
        return
    }

    ctx.JSON(200, nodeSrc)
}

func (n *Node) UpdateNode(ctx *gin.Context){
    body, _ := ioutil.ReadAll(ctx.Request.Body)
    nodeName := gjson.GetBytes(body, "metadata.name").String()
    label_json := gjson.GetBytes(body, "metadata.Labels").String()
    label_map := make(map[string]string)
    err := json.Unmarshal([]byte(label_json), &label_map)
    if err != nil{
        fmt.Println("JsonToMapDemo err:", err)
    }

    node, err := n.ClientSet.CoreV1().Nodes().Get(nodeName, metav1.GetOptions{})
    if err != nil{
        ctx.JSON(404, gin.H{
            "message": "Node not exist",
            "code":  404,
            "reason":  err,
        })
        return
    }

    node.SetLabels(label_map)

    nodeNew, err := n.ClientSet.CoreV1().Nodes().Update(node)
    if err != nil {
        ctx.JSON(400, gin.H{
                "message": "Internal error",
                "code":  400,
                "reason":  err,
            })
        }
    ctx.JSON(200, nodeNew)
}

