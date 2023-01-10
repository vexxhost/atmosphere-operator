package endpoints

import (
	"fmt"

	"golang.org/x/exp/slices"
	"helm.sh/helm/v3/pkg/chart"
)

func endpointAuthUsers(endpoint interface{}) []string {
	skipList := []string{"nova_api", "test"}
	keys := []string{}

	for k := range endpoint.(map[string]interface{})["auth"].(map[string]interface{}) {
		if slices.Contains(skipList, k) {
			continue
		}

		keys = append(keys, k)
	}
	return keys
}

func endpointAuth(endpoint interface{}, passwords map[string]string) (map[string]interface{}, error) {
	auth := map[string]interface{}{}

	for _, user := range endpointAuthUsers(endpoint) {
		password, ok := passwords[user]
		if !ok || password == "" {
			return nil, fmt.Errorf("password for user %s missing", user)
		}

		auth[user] = map[string]interface{}{
			"password": password,
		}
	}

	return auth, nil
}

func basicEndpoint(host string) (map[string]interface{}, error) {
	if host == "" {
		return nil, fmt.Errorf("host is required")
	}

	return map[string]interface{}{
		"scheme": map[string]interface{}{
			"public": "https",
		},
		"host_fqdn_override": map[string]interface{}{
			"public": map[string]interface{}{
				"host": host,
			},
		},
		"port": map[string]interface{}{
			"api": map[string]interface{}{
				"public": float32(443),
			},
		},
	}, nil
}

func ForChart(chart *chart.Chart, config *EndpointConfig) (map[string]interface{}, error) {
	endpoints := make(map[string]interface{})
	chartEndpoints := chart.Values["endpoints"].(map[string]interface{})

	for endpointName, endpointValues := range chartEndpoints {
		switch endpointName {
		case "ceph_mon", "ceph_object_store", "cloudwatch", "cluster_domain_suffix", "compute_spice_proxy", "fluentd", "ingress", "kube_dns", "ldap", "libvirt_exporter", "local_image_registry", "monitoring", "object_store", "oci_image_registry", "powerdns", "prometheus_rabbitmq_exporter":
			continue
		case "baremetal":
			endpoint, err := basicEndpoint(config.IronicHost)
			if err != nil {
				return nil, err
			}
			endpoints[endpointName] = endpoint
		case "clustering":
			endpoint, err := basicEndpoint(config.SenlinHost)
			if err != nil {
				return nil, err
			}
			endpoints[endpointName] = endpoint
		case "cloudformation":
			endpoint, err := basicEndpoint(config.HeatCloudFormationHost)
			if err != nil {
				return nil, err
			}
			endpoints[endpointName] = endpoint
		case "compute":
			endpoint, err := basicEndpoint(config.NovaHost)
			if err != nil {
				return nil, err
			}
			endpoints[endpointName] = endpoint
		case "compute_metadata":
			if config.NovaMetadataSecret == "" {
				return nil, fmt.Errorf("nova metadata secret is required")
			}

			endpoints[endpointName] = map[string]interface{}{
				"secret": config.NovaMetadataSecret,
				"hosts": map[string]interface{}{
					"public": "nova-metadata",
				},
				"port": map[string]interface{}{
					"metadata": map[string]interface{}{
						"public": 8775,
					},
				},
			}
		case "compute_novnc_proxy":
			endpoint, err := basicEndpoint(config.NovaNovncHost)
			if err != nil {
				return nil, err
			}
			endpoint["path"] = map[string]interface{}{
				"default": "/vnc_lite.html",
			}
			endpoints[endpointName] = endpoint
		case "container_infra":
			endpoint, err := basicEndpoint(config.MagnumHost)
			if err != nil {
				return nil, err
			}
			endpoints[endpointName] = endpoint
		case "dashboard":
			endpoint, err := basicEndpoint(config.HorizonHost)
			if err != nil {
				return nil, err
			}
			endpoints[endpointName] = endpoint
		case "dns":
			endpoint, err := basicEndpoint(config.DesignateHost)
			if err != nil {
				return nil, err
			}
			endpoints[endpointName] = endpoint
		case "identity":
			auth, err := endpointAuth(endpointValues, map[string]string{
				"admin":             config.KeystoneAdminPassword,
				"barbican":          config.BarbicanKeystonePassword,
				"glance":            config.GlanceKeystonePassword,
				"cinder":            config.CinderKeystonePassword,
				"placement":         config.PlacementKeystonePassword,
				"neutron":           config.NeutronKeystonePassword,
				"ironic":            config.IronicKeystonePassword,
				"nova":              config.NovaKeystonePassword,
				"senlin":            config.SenlinKeystonePassword,
				"designate":         config.DesignateKeystonePassword,
				"heat":              config.HeatKeystonePassword,
				"heat_trustee":      config.HeatTrusteeKeystonePassword,
				"heat_stack_user":   config.HeatStackUserKeystonePassword,
				"octavia":           config.OctaviaKeystonePassword,
				"magnum":            config.MagnumKeystonePassword,
				"magnum_stack_user": config.MagnumKeystonePassword,
			})
			if err != nil {
				return nil, err
			}

			if config.RegionName == "" {
				return nil, fmt.Errorf("region name is required")
			}

			for _, user := range endpointAuthUsers(endpointValues) {
				auth[user].(map[string]interface{})["region_name"] = config.RegionName
				auth[user].(map[string]interface{})["username"] = fmt.Sprintf("%s-%s", user, config.RegionName)
			}

			endpoint, err := basicEndpoint(config.KeystoneHost)
			if err != nil {
				return nil, err
			}

			endpoint["auth"] = auth
			endpoint["hosts"] = map[string]interface{}{
				"default": "keystone-api",
			}
			endpoint["port"] = map[string]interface{}{
				"api": map[string]interface{}{
					"default": 5000,
					"public":  443,
				},
			}

			endpoints[endpointName] = endpoint
		case "image":
			endpoint, err := basicEndpoint(config.GlanceHost)
			if err != nil {
				return nil, err
			}
			endpoints[endpointName] = endpoint
		case "key_manager":
			endpoint, err := basicEndpoint(config.BarbicanHost)
			if err != nil {
				return nil, err
			}
			endpoints[endpointName] = endpoint
		case "load_balancer":
			endpoint, err := basicEndpoint(config.OctaviaHost)
			if err != nil {
				return nil, err
			}
			endpoints[endpointName] = endpoint
		case "mdns":
			endpoint, err := basicEndpoint(config.DesignateHost)
			if err != nil {
				return nil, err
			}
			endpoints[endpointName] = endpoint
		case "network":
			endpoint, err := basicEndpoint(config.NeutronHost)
			if err != nil {
				return nil, err
			}
			endpoints[endpointName] = endpoint
		case "orchestration":
			endpoint, err := basicEndpoint(config.HeatHost)
			if err != nil {
				return nil, err
			}
			endpoints[endpointName] = endpoint
		case "oslo_cache":
			if config.MemcacheSecretKey == "" {
				return nil, fmt.Errorf("memcache secret key is required")
			}

			endpoints[endpointName] = map[string]interface{}{
				"auth": map[string]interface{}{
					"memcache_secret_key": config.MemcacheSecretKey,
				},
			}
		case "oslo_db", "oslo_db_api", "oslo_db_cell0":
			auth, err := endpointAuth(endpointValues, map[string]string{
				"admin":     config.DatabaseRootPassword,
				"keystone":  config.KeystoneDatabasePassword,
				"barbican":  config.BarbicanDatabasePassword,
				"glance":    config.GlanceDatabasePassword,
				"cinder":    config.CinderDatabasePassword,
				"placement": config.PlacementDatabasePassword,
				"neutron":   config.NeutronDatabasePassword,
				"nova":      config.NovaDatabasePassword,
				"nova_api":  config.NovaDatabasePassword,
				"senlin":    config.SenlinDatabasePassword,
				"designate": config.DesignateDatabasePassword,
				"heat":      config.HeatDatabasePassword,
				"octavia":   config.OctaviaDatabasePassword,
				"magnum":    config.MagnumDatabasePassword,
				"horizon":   config.HorizonDatabasePassword,
			})
			if err != nil {
				return nil, err
			}

			if config.DatabaseNamespace == "" {
				return nil, fmt.Errorf("database namespace is required")
			}
			if config.DatabaseServiceName == "" {
				return nil, fmt.Errorf("database service name is required")
			}

			endpoints[endpointName] = map[string]interface{}{
				"auth":      auth,
				"namespace": config.DatabaseNamespace,
				"hosts": map[string]interface{}{
					"default": config.DatabaseServiceName,
				},
			}
		case "oslo_messaging":
			auth, err := endpointAuth(endpointValues, map[string]string{
				"admin":     config.RabbitmqAdminPassword,
				"keystone":  config.KeystoneRabbitmqPassword,
				"barbican":  config.BarbicanRabbitmqPassword,
				"glance":    config.GlanceRabbitmqPassword,
				"cinder":    config.CinderRabbitmqPassword,
				"neutron":   config.NeutronRabbitmqPassword,
				"nova":      config.NovaRabbitmqPassword,
				"senlin":    config.SenlinRabbitmqPassword,
				"designate": config.DesignateRabbitmqPassword,
				"heat":      config.HeatRabbitmqPassword,
				"octavia":   config.OctaviaRabbitmqPassword,
				"magnum":    config.MagnumRabbitmqPassword,
			})
			if err != nil {
				return nil, err
			}

			// NOTE(mnaser): RabbitMQ operator generates random usernames for the
			//               admin user.
			if config.RabbitmqAdminUsername == "" {
				return nil, fmt.Errorf("rabbitmq admin username is required")
			}
			auth["admin"].(map[string]interface{})["username"] = config.RabbitmqAdminUsername

			if config.RabbitmqNamespace == "" {
				return nil, fmt.Errorf("rabbitmq namespace is required")
			}
			if config.RabbitmqServiceName == "" {
				return nil, fmt.Errorf("rabbitmq service name is required")
			}

			endpoints[endpointName] = map[string]interface{}{
				"auth":        auth,
				"statefulset": false,
				"namespace":   config.RabbitmqNamespace,
				"hosts": map[string]interface{}{
					"default": config.RabbitmqServiceName,
				},
			}
		case "placement":
			endpoint, err := basicEndpoint(config.PlacementHost)
			if err != nil {
				return nil, err
			}
			endpoints[endpointName] = endpoint
		case "volumev3":
			endpoint, err := basicEndpoint(config.CinderHost)
			if err != nil {
				return nil, err
			}
			endpoints[endpointName] = endpoint
		default:
			return nil, fmt.Errorf("endpoint %s not supported", endpointName)
		}
	}

	return endpoints, nil
}
