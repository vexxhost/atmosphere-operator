package endpoints

import (
	"context"
	"errors"
	"strings"

	pxcv1 "github.com/percona/percona-xtradb-cluster-operator/pkg/apis/pxc/v1"
	infrav1alpha1 "github.com/vexxhost/atmosphere-operator/apis/infra/v1alpha1"
	openstackv1alpha1 "github.com/vexxhost/atmosphere-operator/apis/openstack/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type EndpointConfig struct {
	RegionName string

	MemcacheSecretKey string

	DatabaseNamespace    string
	DatabaseServiceName  string
	DatabaseRootPassword string

	RabbitmqNamespace     string
	RabbitmqServiceName   string
	RabbitmqAdminUsername string
	RabbitmqAdminPassword string

	KeystoneHost             string
	KeystoneDatabasePassword string
	KeystoneRabbitmqPassword string
	KeystoneAdminPassword    string

	BarbicanHost             string
	BarbicanDatabasePassword string
	BarbicanRabbitmqPassword string
	BarbicanKeystonePassword string
	BarbicanKeyEncryptionKey string

	GlanceHost             string
	GlanceDatabasePassword string
	GlanceRabbitmqPassword string
	GlanceKeystonePassword string

	CinderHost             string
	CinderDatabasePassword string
	CinderRabbitmqPassword string
	CinderKeystonePassword string

	PlacementHost             string
	PlacementDatabasePassword string
	PlacementKeystonePassword string

	NeutronHost             string
	NeutronDatabasePassword string
	NeutronRabbitmqPassword string
	NeutronKeystonePassword string

	IronicHost             string
	IronicKeystonePassword string

	NovaHost             string
	NovaNovncHost        string
	NovaDatabasePassword string
	NovaRabbitmqPassword string
	NovaKeystonePassword string
	NovaMetadataSecret   string

	SenlinHost             string
	SenlinDatabasePassword string
	SenlinRabbitmqPassword string
	SenlinKeystonePassword string

	DesignateHost             string
	DesignateDatabasePassword string
	DesignateRabbitmqPassword string
	DesignateKeystonePassword string

	HeatHost                      string
	HeatCloudFormationHost        string
	HeatDatabasePassword          string
	HeatRabbitmqPassword          string
	HeatKeystonePassword          string
	HeatTrusteeKeystonePassword   string
	HeatStackUserKeystonePassword string

	OctaviaHost             string
	OctaviaDatabasePassword string
	OctaviaRabbitmqPassword string
	OctaviaKeystonePassword string

	MagnumHost             string
	MagnumDatabasePassword string
	MagnumRabbitmqPassword string
	MagnumKeystonePassword string

	HorizonHost             string
	HorizonDatabasePassword string
}

func NewConfig(options ...func(*EndpointConfig) error) (*EndpointConfig, error) {
	config := &EndpointConfig{}

	for _, o := range options {
		err := o(config)
		if err != nil {
			return nil, err
		}
	}

	return config, nil
}

func WithNamespace(namespace string) func(*EndpointConfig) error {
	return func(ec *EndpointConfig) error {
		ec.DatabaseNamespace = namespace
		ec.RabbitmqNamespace = namespace
		return nil
	}
}

func WithDatabase(ctx context.Context, c client.Client, ref *openstackv1alpha1.NamespacedName) func(*EndpointConfig) error {
	return func(ec *EndpointConfig) error {
		pxc := &pxcv1.PerconaXtraDBCluster{}
		if err := c.Get(ctx, ref.NativeNamespacedName(), pxc); err != nil {
			return err
		}

		secret := &corev1.Secret{}
		if err := c.Get(ctx, client.ObjectKey{Namespace: pxc.Namespace, Name: pxc.Spec.SecretsName}, secret); err != nil {
			return err
		}

		ec.DatabaseNamespace = pxc.Namespace
		ec.DatabaseServiceName = strings.Split(pxc.Status.Host, ".")[0]
		ec.DatabaseRootPassword = string(secret.Data["root"])

		return nil
	}
}

func WithRabbitmq(ctx context.Context, c client.Client, ref *openstackv1alpha1.NamespacedName) func(*EndpointConfig) error {
	return func(ec *EndpointConfig) error {
		rabbitmq := &infrav1alpha1.RabbitmqCluster{}
		if err := c.Get(ctx, ref.NativeNamespacedName(), rabbitmq); err != nil {
			return err
		}

		if rabbitmq.Status.DefaultUser.ServiceReference == nil || rabbitmq.Status.DefaultUser.SecretReference == nil {
			return errors.New("rabbitmq is not ready")
		}

		secret := &corev1.Secret{}
		if err := c.Get(ctx, client.ObjectKey{Namespace: rabbitmq.Status.DefaultUser.SecretReference.Namespace, Name: rabbitmq.Status.DefaultUser.SecretReference.Name}, secret); err != nil {
			return err
		}

		ec.RabbitmqNamespace = rabbitmq.Status.DefaultUser.ServiceReference.Namespace
		ec.RabbitmqServiceName = rabbitmq.Status.DefaultUser.ServiceReference.Name
		ec.RabbitmqAdminUsername = string(secret.Data["username"])
		ec.RabbitmqAdminPassword = string(secret.Data["password"])

		return nil
	}
}

func WithKeystone(ctx context.Context, c client.Client, keystone *openstackv1alpha1.Keystone) func(*EndpointConfig) error {
	return func(ec *EndpointConfig) error {
		databaseRef := keystone.Spec.DatabaseReference.WithNamespace(keystone.Namespace)
		if err := WithDatabase(ctx, c, &databaseRef)(ec); err != nil {
			return err
		}

		rabbitmqRef := keystone.Spec.RabbitmqReference.WithNamespace(keystone.Namespace)
		if err := WithRabbitmq(ctx, c, &rabbitmqRef)(ec); err != nil {
			return err
		}

		ec.RegionName = keystone.Spec.RegionName
		ec.KeystoneHost = keystone.Spec.Ingress.Host

		secret := &corev1.Secret{}
		if err := c.Get(ctx, keystone.Spec.SecretsRef.WithNamespace(keystone.Namespace).NativeNamespacedName(), secret); err != nil {
			return err
		}

		ec.MemcacheSecretKey = string(secret.Data["memcache"])
		ec.KeystoneDatabasePassword = string(secret.Data["database"])
		ec.KeystoneRabbitmqPassword = string(secret.Data["rabbitmq"])
		ec.KeystoneAdminPassword = string(secret.Data["admin"])

		return nil
	}
}

func WithKeystoneRef(ctx context.Context, c client.Client, ref *openstackv1alpha1.NamespacedName) func(*EndpointConfig) error {
	return func(ec *EndpointConfig) error {
		keystone := &openstackv1alpha1.Keystone{}
		if err := c.Get(ctx, ref.NativeNamespacedName(), keystone); err != nil {
			return err
		}

		return WithKeystone(ctx, c, keystone)(ec)
	}
}

func WithBarbican(ctx context.Context, c client.Client, barbican *openstackv1alpha1.Barbican) func(*EndpointConfig) error {
	return func(ec *EndpointConfig) error {
		databaseRef := barbican.Spec.DatabaseReference.WithNamespace(barbican.Namespace)
		if err := WithDatabase(ctx, c, &databaseRef)(ec); err != nil {
			return err
		}

		rabbitmqRef := barbican.Spec.RabbitmqReference.WithNamespace(barbican.Namespace)
		if err := WithRabbitmq(ctx, c, &rabbitmqRef)(ec); err != nil {
			return err
		}

		ec.RegionName = barbican.Spec.RegionName
		ec.BarbicanHost = barbican.Spec.Ingress.Host

		secret := &corev1.Secret{}
		if err := c.Get(ctx, barbican.Spec.SecretsRef.WithNamespace(barbican.Namespace).NativeNamespacedName(), secret); err != nil {
			return err
		}

		ec.MemcacheSecretKey = string(secret.Data["memcache"])
		ec.BarbicanDatabasePassword = string(secret.Data["database"])
		ec.BarbicanRabbitmqPassword = string(secret.Data["rabbitmq"])
		ec.BarbicanKeystonePassword = string(secret.Data["keystone"])
		ec.BarbicanKeyEncryptionKey = string(secret.Data["kek"])

		return nil
	}
}

// TODO: ceph
// TODO: glance
// TODO: cinder

func WithPlacement(ctx context.Context, c client.Client, placement *openstackv1alpha1.Placement) func(*EndpointConfig) error {
	return func(ec *EndpointConfig) error {
		databaseRef := placement.Spec.DatabaseReference.WithNamespace(placement.Namespace)
		if err := WithDatabase(ctx, c, &databaseRef)(ec); err != nil {
			return err
		}

		ec.RegionName = placement.Spec.RegionName
		ec.PlacementHost = placement.Spec.Ingress.Host

		secret := &corev1.Secret{}
		if err := c.Get(ctx, placement.Spec.SecretsRef.WithNamespace(placement.Namespace).NativeNamespacedName(), secret); err != nil {
			return err
		}

		ec.MemcacheSecretKey = string(secret.Data["memcache"])
		ec.PlacementDatabasePassword = string(secret.Data["database"])
		ec.PlacementKeystonePassword = string(secret.Data["keystone"])

		return nil
	}
}

// TODO: ovs
// TODO: libvirt
// TODO: neutron
// TODO: nova
// TODO: senlin
// TODO: designate
// TODO: heat
// TODO: octavia
// TODO: magnum
// TODO: horizon
