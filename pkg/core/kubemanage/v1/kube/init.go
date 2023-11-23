package kube

import (
	"flag"
	"github.com/noovertime7/kubemanage/pkg/logger"
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var K8s k8s

type k8s struct {
	// Config 集群配置对象
	Config *rest.Config
	// ClientSet kubernetes的客户端集合
	ClientSet *kubernetes.Clientset
}

// Init 初始化 k8s客户端+配置
func (k *k8s) Init() error {
	var err error
	// 记录 集群配置
	var config *rest.Config
	// 记录 .kubeconfig 的路径
	var kubeConfig *string

	// homeDir()取home目录的路径，拼接 kubeconfig 文件路径
	if home := homeDir(); home != "" {
		// 添加命令行参数，并给出默认值
		kubeConfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		// 如果没有找到HOME路径，则默认配置文件路径设置为空，让用户自己从命令行设置kubeconfig文件的路径
		kubeConfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// TODO 使用 ServiceAccount 创建集群配置（InCluster模式）
	if config, err = rest.InClusterConfig(); err != nil {
		// 使用 KubeConfig 文件创建集群配置
		if config, err = clientcmd.BuildConfigFromFlags("", *kubeConfig); err != nil {
			return err
		}
	}

	// 创建 clientSet
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}
	log := logger.New()
	log.Info("获取k8s clientSet 成功")
	k.ClientSet = clientSet
	k.Config = config
	return nil
}

// homeDir 从环境变量中，取出HOME或USERPROFILE的值，即为 用户的目录
func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	// 没有HOME，说明是Windows系统，取 USERPROFILE 的值
	return os.Getenv("USERPROFILE") // windows
}
