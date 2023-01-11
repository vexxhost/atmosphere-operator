package endpoints

import (
	"context"
	"errors"
	"fmt"
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

	CoreDNSClusterIP string

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
	OctaviaHeartbeatKey     string

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

// TODO: libvirt

func WithNeutron(ctx context.Context, c client.Client, neutron *openstackv1alpha1.Neutron) func(*EndpointConfig) error {
	return func(ec *EndpointConfig) error {
		databaseRef := neutron.Spec.DatabaseReference.WithNamespace(neutron.Namespace)
		if err := WithDatabase(ctx, c, &databaseRef)(ec); err != nil {
			return err
		}

		rabbitmqRef := neutron.Spec.RabbitmqReference.WithNamespace(neutron.Namespace)
		if err := WithRabbitmq(ctx, c, &rabbitmqRef)(ec); err != nil {
			return err
		}

		ec.RegionName = neutron.Spec.RegionName
		ec.NeutronHost = neutron.Spec.Ingress.Host

		secret := &corev1.Secret{}
		if err := c.Get(ctx, neutron.Spec.SecretsRef.WithNamespace(neutron.Namespace).NativeNamespacedName(), secret); err != nil {
			return err
		}

		ec.MemcacheSecretKey = string(secret.Data["memcache"])
		ec.NeutronDatabasePassword = string(secret.Data["database"])
		ec.NeutronRabbitmqPassword = string(secret.Data["rabbitmq"])
		ec.NeutronKeystonePassword = string(secret.Data["keystone"])

		coreDnsService := &corev1.Service{}
		if err := c.Get(ctx, neutron.Spec.CoreDNSRef.WithNamespace(neutron.Namespace).NativeNamespacedName(), coreDnsService); err != nil {
			return err
		}

		ec.CoreDNSClusterIP = coreDnsService.Spec.ClusterIP

		return nil
	}
}

func WithNeutronRef(ctx context.Context, c client.Client, ref *openstackv1alpha1.NamespacedName) func(*EndpointConfig) error {
	return func(ec *EndpointConfig) error {
		neutron := &openstackv1alpha1.Neutron{}
		if err := c.Get(ctx, ref.NativeNamespacedName(), neutron); err != nil {
			return err
		}

		return WithNeutron(ctx, c, neutron)(ec)
	}
}

func WithIronic(ctx context.Context, c client.Client, ironic *openstackv1alpha1.Ironic) func(*EndpointConfig) error {
	return func(ec *EndpointConfig) error {
		databaseRef := ironic.Spec.DatabaseReference.WithNamespace(ironic.Namespace)
		if err := WithDatabase(ctx, c, &databaseRef)(ec); err != nil {
			return err
		}

		rabbitmqRef := ironic.Spec.RabbitmqReference.WithNamespace(ironic.Namespace)
		if err := WithRabbitmq(ctx, c, &rabbitmqRef)(ec); err != nil {
			return err
		}

		ec.RegionName = ironic.Spec.RegionName
		ec.IronicHost = ironic.Spec.Ingress.Host

		secret := &corev1.Secret{}
		if err := c.Get(ctx, ironic.Spec.SecretsRef.WithNamespace(ironic.Namespace).NativeNamespacedName(), secret); err != nil {
			return err
		}

		ec.MemcacheSecretKey = string(secret.Data["memcache"])
		// ec.IronicDatabasePassword = string(secret.Data["database"])
		// ec.IronicRabbitmqPassword = string(secret.Data["rabbitmq"])
		ec.IronicKeystonePassword = string(secret.Data["keystone"])

		return nil
	}
}

func WithIronicRef(ctx context.Context, c client.Client, ref *openstackv1alpha1.NamespacedName) func(*EndpointConfig) error {
	return func(ec *EndpointConfig) error {
		ironic := &openstackv1alpha1.Ironic{}
		if err := c.Get(ctx, ref.NativeNamespacedName(), ironic); err != nil {
			return err
		}

		return WithIronic(ctx, c, ironic)(ec)
	}
}

func WithNova(ctx context.Context, c client.Client, nova *openstackv1alpha1.Nova) func(*EndpointConfig) error {
	return func(ec *EndpointConfig) error {
		databaseRef := nova.Spec.DatabaseReference.WithNamespace(nova.Namespace)
		if err := WithDatabase(ctx, c, &databaseRef)(ec); err != nil {
			return err
		}

		rabbitmqRef := nova.Spec.RabbitmqReference.WithNamespace(nova.Namespace)
		if err := WithRabbitmq(ctx, c, &rabbitmqRef)(ec); err != nil {
			return err
		}

		ec.RegionName = nova.Spec.RegionName
		ec.NovaHost = nova.Spec.Ingress.Host

		secret := &corev1.Secret{}
		if err := c.Get(ctx, nova.Spec.SecretsRef.WithNamespace(nova.Namespace).NativeNamespacedName(), secret); err != nil {
			return err
		}

		ec.MemcacheSecretKey = string(secret.Data["memcache"])
		ec.NovaDatabasePassword = string(secret.Data["database"])
		ec.NovaRabbitmqPassword = string(secret.Data["rabbitmq"])
		ec.NovaKeystonePassword = string(secret.Data["keystone"])
		ec.NovaMetadataSecret = string(secret.Data["metadata"])

		return nil
	}
}

func WithNovaRef(ctx context.Context, c client.Client, ref *openstackv1alpha1.NamespacedName) func(*EndpointConfig) error {
	return func(ec *EndpointConfig) error {
		nova := &openstackv1alpha1.Nova{}
		if err := c.Get(ctx, ref.NativeNamespacedName(), nova); err != nil {
			return err
		}

		return WithNova(ctx, c, nova)(ec)
	}
}

// TODO: senlin

func WithDesignate(ctx context.Context, c client.Client, designate *openstackv1alpha1.Designate) func(*EndpointConfig) error {
	return func(ec *EndpointConfig) error {
		databaseRef := designate.Spec.DatabaseReference.WithNamespace(designate.Namespace)
		if err := WithDatabase(ctx, c, &databaseRef)(ec); err != nil {
			return err
		}

		rabbitmqRef := designate.Spec.RabbitmqReference.WithNamespace(designate.Namespace)
		if err := WithRabbitmq(ctx, c, &rabbitmqRef)(ec); err != nil {
			return err
		}

		ec.RegionName = designate.Spec.RegionName
		ec.DesignateHost = designate.Spec.Ingress.Host

		secret := &corev1.Secret{}
		if err := c.Get(ctx, designate.Spec.SecretsRef.WithNamespace(designate.Namespace).NativeNamespacedName(), secret); err != nil {
			return err
		}

		ec.MemcacheSecretKey = string(secret.Data["memcache"])
		ec.DesignateDatabasePassword = string(secret.Data["database"])
		ec.DesignateRabbitmqPassword = string(secret.Data["rabbitmq"])
		ec.DesignateKeystonePassword = string(secret.Data["keystone"])

		return nil
	}
}

func WithDesignateRef(ctx context.Context, c client.Client, ref *openstackv1alpha1.NamespacedName) func(*EndpointConfig) error {
	return func(ec *EndpointConfig) error {
		designate := &openstackv1alpha1.Designate{}
		if err := c.Get(ctx, ref.NativeNamespacedName(), designate); err != nil {
			return err
		}

		return WithDesignate(ctx, c, designate)(ec)
	}
}

// TODO: heat

func WithOctavia(ctx context.Context, c client.Client, octavia *openstackv1alpha1.Octavia) func(*EndpointConfig) error {
	return func(ec *EndpointConfig) error {
		databaseRef := octavia.Spec.DatabaseReference.WithNamespace(octavia.Namespace)
		if err := WithDatabase(ctx, c, &databaseRef)(ec); err != nil {
			return err
		}

		rabbitmqRef := octavia.Spec.RabbitmqReference.WithNamespace(octavia.Namespace)
		if err := WithRabbitmq(ctx, c, &rabbitmqRef)(ec); err != nil {
			return err
		}

		ec.RegionName = octavia.Spec.RegionName
		ec.OctaviaHost = octavia.Spec.Ingress.Host

		secret := &corev1.Secret{}
		if err := c.Get(ctx, octavia.Spec.SecretsRef.WithNamespace(octavia.Namespace).NativeNamespacedName(), secret); err != nil {
			return err
		}

		ec.MemcacheSecretKey = string(secret.Data["memcache"])
		ec.OctaviaDatabasePassword = string(secret.Data["database"])
		ec.OctaviaRabbitmqPassword = string(secret.Data["rabbitmq"])
		ec.OctaviaKeystonePassword = string(secret.Data["keystone"])

		ec.OctaviaHeartbeatKey = string(secret.Data["heartbeat"])
		if ec.OctaviaHeartbeatKey == "" {
			return fmt.Errorf("octavia heartbeat key is empty")
		}

		return nil
	}
}

func WithOctaviaRef(ctx context.Context, c client.Client, ref *openstackv1alpha1.NamespacedName) func(*EndpointConfig) error {
	return func(ec *EndpointConfig) error {
		octavia := &openstackv1alpha1.Octavia{}
		if err := c.Get(ctx, ref.NativeNamespacedName(), octavia); err != nil {
			return err
		}

		return WithOctavia(ctx, c, octavia)(ec)
	}
}

// TODO: magnum

func WithHorizon(ctx context.Context, c client.Client, horizon *openstackv1alpha1.Horizon) func(*EndpointConfig) error {
	return func(ec *EndpointConfig) error {
		databaseRef := horizon.Spec.DatabaseReference.WithNamespace(horizon.Namespace)
		if err := WithDatabase(ctx, c, &databaseRef)(ec); err != nil {
			return err
		}

		ec.HorizonHost = horizon.Spec.Ingress.Host

		secret := &corev1.Secret{}
		if err := c.Get(ctx, horizon.Spec.SecretsRef.WithNamespace(horizon.Namespace).NativeNamespacedName(), secret); err != nil {
			return err
		}

		ec.HorizonDatabasePassword = string(secret.Data["database"])

		return nil
	}
}
