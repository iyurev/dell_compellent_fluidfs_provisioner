package main

import (
	"flag"
	"github.com/iyurev/dell_compellent_fluidfs_provisioner/fluidfs_nfs_provisoner"
	"github.com/iyurev/sig-storage-lib-external-provisioner/controller"
	"github.com/prometheus/common/log"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/klog"
	"os"
)

var (
	provName string = "fluidfs-nfs-provisioner"
	k8sConf  *rest.Config
)

func main() {
	klog.InitFlags(nil)
	err := flag.Set("stderrthreshold", "INFO")
	if err != nil {
		log.Fatal(err)
	}
	flag.Parse()
	k8sHost, k8sToken := os.Getenv("K8S_HOST"), os.Getenv("K8S_TOKEN")

	username := os.Getenv("USERNAME")
	if username == "" {
		log.Fatal("You must give username for FluidFs rest api!!")
	}
	password := os.Getenv("PASSWORD")
	if password == "" {
		log.Fatal("You must give password for FluidFs rest api!!")
	}
	baseurl := os.Getenv("BASE_URL")
	if baseurl == "" {
		log.Fatal("You must give base url for FluidFs rest api!!")
	}
	clusterName := os.Getenv("CLUSTER_NAME")
	if clusterName == "" {
		log.Fatal("You must set FluidFs cluster name!!!")
	}

	if k8sToken != "" && k8sHost != "" {
		tlsClientConfig := rest.TLSClientConfig{}
		k8sConf = &rest.Config{
			Host:            k8sHost,
			TLSClientConfig: tlsClientConfig,
			BearerToken:     k8sToken,
		}

	} else {
		var err error
		k8sConf, err = rest.InClusterConfig()
		if err != nil {
			log.Fatal(err)
		}
	}

	k8sClient, err := kubernetes.NewForConfig(k8sConf)
	if err != nil {
		log.Fatal(err)
	}

	k8sVer, err := k8sClient.Discovery().ServerVersion()
	if err != nil {
		log.Fatal(err)
	}

	fluidfsNfsProvisioner := &fluidfs_nfs_provisoner.FluidfsNfsProvisioner{
		Username:           username,
		Password:           password,
		BaseUrl:            baseurl,
		FluidFsClusterName: clusterName,
	}
	pc := controller.NewProvisionController(k8sClient, provName, fluidfsNfsProvisioner, k8sVer.GitVersion)
	pc.Run(wait.NeverStop)

}
