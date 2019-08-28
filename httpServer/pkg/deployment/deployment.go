package deployment

import (
    "github.com/gin-gonic/gin"
	"github.com/json-iterator/go"
    "fmt"
    "github.com/tidwall/gjson"
    "io/ioutil"
	//corev1 "k8s.io/api/core/v1"
    appv1 "k8s.io/api/apps/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/kubernetes"
)

type Deployment struct{
    ClientSet *kubernetes.Clientset
}

func(n *Deployment) AddDeployment(ctx *gin.Context){
    body, _ := ioutil.ReadAll(ctx.Request.Body)
	dpNamespace := gjson.GetBytes(body, "metadata.namespace").String()

	deployment := &appv1.Deployment{}
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
    json.Unmarshal(body,deployment)

	deployment, err := n.ClientSet.AppsV1().Deployments(dpNamespace).Create(deployment)
    if err != nil {
        fmt.Println(err)
        ctx.JSON(400, gin.H{
            "message": "Internal error",
            "code":  400,
            "reason":  err,
        })
    }

    ctx.JSON(200, deployment)
}

func (n *Deployment) DeleteDeployment(ctx *gin.Context){
    dpName := ctx.Param("name")
    dpNamespace := ctx.Param("namespace")
	fmt.Println(dpName)
    fmt.Println(dpNamespace)

	deployment, err := n.ClientSet.AppsV1().Deployments(dpNamespace).Get(dpName, metav1.GetOptions{})
    if err != nil{
		fmt.Println(err)
        ctx.JSON(404, gin.H{
            "message": "Namespace or Deployment not exist",
            "code":  404,
            "reason":  err,
        })
        return
    }

    err = n.ClientSet.AppsV1().Deployments(dpNamespace).Delete(dpName, &metav1.DeleteOptions{})

    if err != nil {
        ctx.JSON(500, gin.H{
            "message": "Internal error",
            "code":  500,
            "reason":  err,
        })
    }

    ctx.JSON(200, deployment)
}

func (n *Deployment) ListDeployment(ctx *gin.Context){
	dpNamespace := ctx.Param("namespace")
    fmt.Println(dpNamespace)

	//deploymentList, err := n.ClientSet.AppsV1().Deployments(corev1.NamespaceAll).List(metav1.ListOptions{})
	deploymentList, err := n.ClientSet.AppsV1().Deployments(dpNamespace).List(metav1.ListOptions{})
    if err != nil {
        ctx.JSON(404, gin.H{
            "message": "Deployment or Namespace not exist",
            "code":  404,
            "reason":  err,
        })
        return
    }

    ctx.JSON(200, deploymentList)
}

func (n *Deployment) GetDeployment(ctx *gin.Context){
	dpName := ctx.Param("name")
	dpNamespace := ctx.Param("namespace")
    fmt.Println(dpName)
    fmt.Println(dpNamespace)

	deploymentSrc, err := n.ClientSet.AppsV1().Deployments(dpNamespace).Get(dpName, metav1.GetOptions{})
    if err != nil {
		fmt.Println(err)
        ctx.JSON(404, gin.H{
            "message": "Deployment or Namespace not exist",
            "code":  404,
            "reason":  err,
        })
        return
    }

    ctx.JSON(200, deploymentSrc)
}

//func (n *Deployment) UpdateDeployment(ctx *gin.Context){
//    body, _ := ioutil.ReadAll(ctx.Request.Body)
//    dpName := gjson.GetBytes(body, "metadata.name").String()
//	dpNamespace := gjson.GetBytes(body, "metadata.namespace").String()
//
//    label_json := gjson.GetBytes(body, "metadata.Labels").String()
//    label_map := make(map[string]string)
//    err := json.Unmarshal([]byte(label_json), &label_map)
//    if err != nil{
//        fmt.Println("JsonToMapDemo err:", err)
//    }
//
//    deployment, err := n.ClientSet.AppsV1().Deployments(dpNamespace).Get(dpName, metav1.GetOptions{})
//    if err != nil{
//        ctx.JSON(404, gin.H{
//            "message": "Deployment not exist",
//            "code":  404,
//            "reason":  err,
//        })
//        return
//    }
//
//    deployment.SetLabels(label_map)
//
//    newDeployment, err := n.ClientSet.AppsV1().Deployments(dpNamespace).Update(deployment)
//    if err != nil {
//        ctx.JSON(400, gin.H{
//                "message": "Internal error",
//                "code":  400,
//                "reason":  err,
//            })
//        }
//    ctx.JSON(200, newDeployment)
//}

func (n *Deployment) UpdateDeployment(ctx *gin.Context){

    body, _ := ioutil.ReadAll(ctx.Request.Body)
    deployment := &appv1.Deployment{}
    var json = jsoniter.ConfigCompatibleWithStandardLibrary
    json.Unmarshal(body,deployment)
    fmt.Println(deployment)

    deployment_Old, err := n.ClientSet.AppsV1().Deployments(deployment.GetNamespace()).Get(deployment.GetName(), metav1.GetOptions{})
    if err != nil {
        ctx.JSON(404, gin.H{
            "message": "Deployment or Namespace not exist",
            "code":  404,
            "reason":  err,
        })
        return
    }

    deployment_Old.SetLabels(deployment.GetLabels())

    deployment_New, err := n.ClientSet.AppsV1().Deployments(deployment_Old.GetNamespace()).Update(deployment_Old)
    if err != nil {
        ctx.JSON(500, gin.H{
            "message": "Internal error",
            "code":  500,
            "reason":  err,
        })
    }
    ctx.JSON(200, deployment_New)
}
