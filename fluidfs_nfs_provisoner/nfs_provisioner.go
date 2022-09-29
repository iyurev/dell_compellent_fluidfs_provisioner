package fluidfs_nfs_provisoner

import (
	"fmt"
	"github.com/iyurev/dell_compellent_fluidfs_provisioner/nfs_permissions"
	"github.com/iyurev/go_dell_compellent_api/compellent_api"
	"github.com/iyurev/sig-storage-lib-external-provisioner/controller"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const defaultPort = 3033
const tempMountPath = "/mnt"
const shareGid = 0
const permMask = 0775

type FluidfsNfsProvisioner struct {
	Username string
	Password string
	BaseUrl  string
	Port     int
	//ExportAddr           string
	FluidFsClusterName string
	//NasVolumeRootFoleder string
}

func (provisioner *FluidfsNfsProvisioner) CreateVolume(options controller.ProvisionOptions, volName string) error {
	port := defaultPort
	if provisioner.Port != 0 {
		port = provisioner.Port
	}
	pvcVolSize := options.PVC.Spec.Resources.Requests[v1.ResourceStorage]
	gbSize := pvcVolSize.Value()
	nfsAccess, ok := options.StorageClass.Parameters["nfs_access"]
	if !ok {
		return fmt.Errorf("Missing <nfs_access> storage class parameters!!")
	}
	nasVolFolder, ok := options.StorageClass.Parameters["nas_volume_folder"]
	if !ok {
		return fmt.Errorf("Missing <nas_volume_root_folder> storageclass parameter!!!")
	}

	nas, err := compellent_api.NewCompelentREST(provisioner.BaseUrl, provisioner.Username, provisioner.Password, port)
	if err != nil {

		return err
	}
	err = nas.CreateNfsPV(provisioner.FluidFsClusterName, volName, gbSize, nasVolFolder, nfsAccess)
	if err != nil {
		return err
	}

	return nil
}

func (provisioner *FluidfsNfsProvisioner) Provision(options controller.ProvisionOptions) (*v1.PersistentVolume, error) {
	volName := fmt.Sprintf("%s-%s-%s", options.PVC.Namespace, options.PVC.Name, options.PVName)
	err := provisioner.CreateVolume(options, volName)
	if err != nil {
		return &v1.PersistentVolume{}, err
	}
	exportAddr, ok := options.StorageClass.Parameters["nfs_export_addr"]
	//options.Parameters["nfs_export_addr"]
	if !ok {
		return &v1.PersistentVolume{}, fmt.Errorf("Missing <nfs_export_addr> storageclass parameter!!!")
	}
	exportPath := fmt.Sprintf("/%s", volName)
	annotation := make(map[string]string)
	annotation["storage-system"] = "FluidFs NFS volume"
	annotation["nas-volume-folder"] = options.StorageClass.Parameters["nas_volume_folder"]
	//Setup correct permissions on FS
	nfsShare, err := nfs_permissions.NewNFSShare(exportAddr, exportPath, tempMountPath, shareGid, permMask)
	if err != nil {
		return &v1.PersistentVolume{}, err
	}
	err = nfsShare.SetPermissions()
	if err != nil {
		return &v1.PersistentVolume{}, err
	}
	//Describe Pv object metadata
	metadata := metav1.ObjectMeta{
		Name:        volName,
		Labels:      map[string]string{},
		Annotations: annotation,
	}
	spec := v1.PersistentVolumeSpec{
		PersistentVolumeReclaimPolicy: *options.StorageClass.ReclaimPolicy,
		//options.PersistentVolumeReclaimPolicy,
		AccessModes: options.PVC.Spec.AccessModes,
		//StorageClassName:              *o.PVC.Spec.StorageClassName,
		Capacity: v1.ResourceList{
			v1.ResourceStorage: options.PVC.Spec.Resources.Requests[v1.ResourceStorage]},
		PersistentVolumeSource: v1.PersistentVolumeSource{
			NFS: &v1.NFSVolumeSource{
				Server:   exportAddr,
				Path:     exportPath,
				ReadOnly: false},
		},
	}
	return &v1.PersistentVolume{ObjectMeta: metadata, Spec: spec}, nil
}

func (provisioner *FluidfsNfsProvisioner) Delete(pv *v1.PersistentVolume) error {
	port := defaultPort
	if provisioner.Port != 0 {
		port = provisioner.Port
	}

	if nas, err := compellent_api.NewCompelentREST(provisioner.BaseUrl, provisioner.Username, provisioner.Password, port); err != nil {

		return err
	} else {
		err := nas.RemoveNfsPV(provisioner.FluidFsClusterName, pv.Name)
		if err != nil {
			return err
		}
	}
	return nil
}
