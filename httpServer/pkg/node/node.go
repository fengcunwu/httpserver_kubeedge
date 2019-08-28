package node

import (
	"github.com/gin-gonic/gin"
	"fmt"
	"io/ioutil"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/klog"
	_"reflect"
	"github.com/json-iterator/go"
)

type Node struct{
    ClientSet *kubernetes.Clientset
}

func (n *Node) AddNode(ctx *gin.Context){
	body, _ := ioutil.ReadAll(ctx.Request.Body)

    node := &corev1.Node{}
    var json = jsoniter.ConfigCompatibleWithStandardLibrary
    json.Unmarshal(body,node)
    fmt.Println(node)

    node, err := n.ClientSet.CoreV1().Nodes().Create(node)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(500, gin.H{
			"message": "Internal error",
			"code":  500,
			"reason":  err,
		})
		return
	}

	ctx.JSON(200, node)
}

func (n *Node) DeleteNode(ctx *gin.Context){
    nodeName := ctx.Param("name")
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
    err = n.ClientSet.CoreV1().Nodes().Delete(nodeName, &metav1.DeleteOptions{
        PropagationPolicy: &deletePolicy,
    })

    if err != nil {
        ctx.JSON(500, gin.H{
            "message": "Internal error",
            "code":  500,
            "reason":  err,
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
    nodeName := ctx.Param("name")
	fmt.Println(nodeName)

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
	node := &corev1.Node{}
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	json.Unmarshal(body,node)
	fmt.Println(node)

	nodeOld, err := n.ClientSet.CoreV1().Nodes().Get(node.GetName(), metav1.GetOptions{})
	if err != nil {
		ctx.JSON(404, gin.H{
			"message": "Node not exist",
			"code":  404,
			"reason":  err,
		})
		return
	}

	nodeOld.SetLabels(node.GetLabels())

	nodeNew, err := n.ClientSet.CoreV1().Nodes().Update(nodeOld)
	if err != nil {
		ctx.JSON(500, gin.H{
			"message": "Internal error",
			"code":  500,
			"reason":  err,
		})
    }
	ctx.JSON(200, nodeNew)
}

