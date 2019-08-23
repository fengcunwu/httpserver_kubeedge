package deployment

import (
    "github.com/gin-gonic/gin"
    "encoding/json"
    //"net/http"
    "fmt"
    "github.com/tidwall/gjson"
    "io/ioutil"
	corev1 "k8s.io/api/core/v1"
    appv1 "k8s.io/api/apps/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/kubernetes"
)

type Deployment struct{
    ClientSet *kubernetes.Clientset
}

func (n *Deployment) AddDeployment(ctx *gin.Context){
	body, _ := ioutil.ReadAll(ctx.Request.Body)
    dpName := gjson.GetBytes(body, "metadata.name").String()
    dpNamespace := gjson.GetBytes(body, "metadata.namespace").String()
    fmt.Println(dpName)
	fmt.Println(dpNamespace)

	var a int32
	a = 1
    deployment := &appv1.Deployment{
        ObjectMeta: metav1.ObjectMeta{
            Name: dpName,
			Namespace: dpNamespace,
        },
		Spec: appv1.DeploymentSpec{
            Replicas: &a,
            Selector: &metav1.LabelSelector{
                MatchLabels: map[string]string{
                    "app": dpName,
                },
            },
            Template: corev1.PodTemplateSpec{
                ObjectMeta: metav1.ObjectMeta{
                    Labels: map[string]string{
                        "app": dpName,
                        "version": "V1",
                    },
                },
                Spec: corev1.PodSpec{
                    Containers: []corev1.Container{
                        {
                            Name: dpName,
                            Image: "nginx:1.13.5-alpine",
                            ImagePullPolicy: "IfNotPresent",
                            Ports: []corev1.ContainerPort{
                                {
                                    Name: "http",
                                    Protocol: corev1.ProtocolTCP,
									ContainerPort: 80,
                                    HostPort: 80,
                                },
                            },
                        },
                    },
					NodeSelector: map[string]string{
						"name": "edge-node",
					},
                },
            },
		},
    }

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
    body, _ := ioutil.ReadAll(ctx.Request.Body)
    dpName := gjson.GetBytes(body, "metadata.name").String()
    dpNamespace := gjson.GetBytes(body, "metadata.namespace").String()
    fmt.Println(dpName)
    fmt.Println(dpNamespace)

	deployment, err := n.ClientSet.AppsV1().Deployments(dpNamespace).Get(dpName, metav1.GetOptions{})
    if err != nil{
		fmt.Println(err)
        ctx.JSON(404, gin.H{
            "message": "Deployment or Namespace not exist",
            "code":  404,
            "reason":  err,
        })
        return
    }


    deletePolicy := metav1.DeletePropagationForeground
    err_d := n.ClientSet.AppsV1().Deployments(dpNamespace).Delete(dpName, &metav1.DeleteOptions{
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

    ctx.JSON(200, deployment)
}

func (n *Deployment) ListDeployment(ctx *gin.Context){
    deploymentList, err := n.ClientSet.AppsV1().Deployments(corev1.NamespaceAll).List(metav1.ListOptions{})
    if err != nil {
        ctx.JSON(404, gin.H{
            "message": "Deployment not exist",
            "code":  404,
            "reason":  err,
        })
        return
    }

    ctx.JSON(200, deploymentList)
}

func (n *Deployment) GetDeployment(ctx *gin.Context){
    body, _ := ioutil.ReadAll(ctx.Request.Body)
    dpName := gjson.GetBytes(body, "metadata.name").String()
    dpNamespace := gjson.GetBytes(body, "metadata.namespace").String()

	deploymentSrc, err := n.ClientSet.AppsV1().Deployments(dpNamespace).Get(dpName, metav1.GetOptions{})
    if err != nil {
		fmt.Println(err)
        ctx.JSON(404, gin.H{
            "message": "Deployment not exist",
            "code":  404,
            "reason":  err,
        })
        return
    }

    ctx.JSON(200, deploymentSrc)
}

func (n *Deployment) UpdateDeployment(ctx *gin.Context){
    body, _ := ioutil.ReadAll(ctx.Request.Body)
    dpName := gjson.GetBytes(body, "metadata.name").String()
	dpNamespace := gjson.GetBytes(body, "metadata.namespace").String()

    label_json := gjson.GetBytes(body, "metadata.Labels").String()
    label_map := make(map[string]string)
    err := json.Unmarshal([]byte(label_json), &label_map)
    if err != nil{
        fmt.Println("JsonToMapDemo err:", err)
    }

    deployment, err := n.ClientSet.AppsV1().Deployments(dpNamespace).Get(dpName, metav1.GetOptions{})
    if err != nil{
        ctx.JSON(404, gin.H{
            "message": "Deployment not exist",
            "code":  404,
            "reason":  err,
        })
        return
    }

    deployment.SetLabels(label_map)

    newDeployment, err := n.ClientSet.AppsV1().Deployments(dpNamespace).Update(deployment)
    if err != nil {
        ctx.JSON(400, gin.H{
                "message": "Internal error",
                "code":  400,
                "reason":  err,
            })
        }
    ctx.JSON(200, newDeployment)
}
