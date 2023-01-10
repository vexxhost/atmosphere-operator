package builders

import (
	corev1 "k8s.io/api/core/v1"

	pxcv1 "github.com/percona/percona-xtradb-cluster-operator/pkg/apis/pxc/v1"
	"github.com/vexxhost/atmosphere-operator/pkg/images"
)

// PxcBuilder defines the interface to build a PXC
type PxcBuilder struct {
	obj *pxcv1.PerconaXtraDBCluster
	err error
}

// Namespace returns a new Namespace builder with default configurations
func defaultPxc(existing *pxcv1.PerconaXtraDBCluster) *PxcBuilder {
	pb = &PxcBuilder{
		obj: existing,
		err: nil,
	}
	pb.obj.spec.pxc.size = int32(3)
	pb.obj.spec.pxc.autoRecovery = true
	pb.obj.spec.pxc.configuration = "[mysqld]\nmax_connections=8192\n"
	pb.obj.spec.pxc.nodeSelector = &corev1.NodeSelector{
		NodeSelectorTerms: []corev1.NodeSelectorTerm{
			{
				MatchExpressions: []corev1.NodeSelectorRequirement{
					{
						Key:      "openstack-control-plane",
						Operator: corev1.NodeSelectorOpIn,
						Values:   []string{"enabled"},
					},
				},
			},
		},
	}
	q, _ := k8sresource.ParseQuantity("160Gi")
	pb.obj.spec.pxc.volumeSpec = &pxcv1.VolumeSpec{
		persistentVolumeClaim: &corev1.PersistentVolumeClaimSpec{
			resources: corev1.ResourceRequirements{
				requests: map[corev1.ResourceName]k8sresource.Quantity{
					corev1.ResourceStorage: q,
				},
			},
		},
	}
	pb.obj.spec.haproxy.enabled = true
	pb.obj.spec.haproxy.size = int32(3)
	pb.obj.spec.haproxy.nodeSelector = &corev1.NodeSelector{
	return pb
}

// Annotation sets one set of Annotation
func (pb *PxcBuilder) overrideImageRegistry(imageRegistry string) *PxcBuilder {
	if imageRegistry == "" {
		return pb
	}
	image, err := images.GetImageReference("percona_xtradb_cluster")
	if err != nil {
		pb.err = err
		return pb
	}
	image, err = images.OverrideRegistry(image, imageRegistry)
	if err != nil {
		pb.err = err
		return pb
	}
	pb.obj.spec.pxc.image = image.Remote()
	
	image, err = images.GetImageReference("percona_xtradb_cluster_haproxy")
	if err != nil {
		pb.err = err
		return pb
	}

	image, err = images.OverrideRegistry(image, imageRegistry)
	if err != nil {
		pb.err = err
		return pb
	}
	pb.obj.spec.haproxy.image = image.Remote()

	return pb
}

// Build returns a complete Pxc object
func (pb *PxcBuilder) Build() (*pxcv1.PerconaXtraDBCluster, error) {
	return pb.obj, pb,err
}
