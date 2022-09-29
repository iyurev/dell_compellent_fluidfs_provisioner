package nfs_permissions

import (
	"fmt"
	"k8s.io/klog"
	"k8s.io/kubernetes/pkg/util/mount"
	"os"
)

type NFSShare struct {
	Server        string
	Path          string
	DestMountPath string
	FSGroup       int
	PermMask      int
	Mounter       mount.Interface
}

func NewNFSShare(nfsServer, path, destPath string, groupUid, permMask int) (*NFSShare, error) {
	if permMask == 0 {
		return &NFSShare{}, fmt.Errorf("You must set unix permissions mask!!! ")
	}
	mounter := mount.New("")

	return &NFSShare{
		Server:        nfsServer,
		Path:          path,
		DestMountPath: destPath,
		FSGroup:       groupUid,
		PermMask:      permMask,
		Mounter:       mounter,
	}, nil
}

func (nfsShare *NFSShare) SetPermissions() error {
	klog.Infof("Mount NFS share: %s:%s  to %s\n", nfsShare.Server, nfsShare.Path, nfsShare.DestMountPath)
	err := nfsShare.Mount()
	if err != nil {
		return err
	}
	klog.Infof("Change group owner for NFS share: %s:%s to GID %d\n", nfsShare.Server, nfsShare.Path, nfsShare.FSGroup)
	if err := nfsShare.Chown(-1, 0); err != nil {
		return err
	}
	klog.Infof("Change unix acl for NFS share: %s:%s to %s\n", nfsShare.Server, nfsShare.Path, os.FileMode(nfsShare.PermMask).String())
	if err := nfsShare.Chmod(); err != nil {
		return err
	}
	klog.Infof("Unmount NFS share: %s:%s, mount point  %s\n", nfsShare.Server, nfsShare.Path, nfsShare.DestMountPath)
	if err := nfsShare.UmountNFS(); err != nil {
		return err
	}
	return nil
}

func (nfsShare *NFSShare) Mount() error {
	src := fmt.Sprintf("%s:%s", nfsShare.Server, nfsShare.Path)
	err := nfsShare.Mounter.Mount(src, nfsShare.DestMountPath, "nfs", []string{})
	if err != nil {
		return err
	}
	return nil
}

func (nfsShare *NFSShare) UmountNFS() error {
	err := nfsShare.Mounter.Unmount(nfsShare.DestMountPath)
	if err != nil {
		return err
	}
	return nil
}

func (nfsShare *NFSShare) Chmod() error {
	fileMode := os.FileMode(nfsShare.PermMask)
	err := os.Chmod(nfsShare.DestMountPath, fileMode)
	if err != nil {
		return err
	}
	return nil
}
func (nfsShare *NFSShare) Chown(uid, gid int) error {
	err := os.Chown(nfsShare.DestMountPath, uid, gid)
	if err != nil {
		return err
	}
	return nil
}
