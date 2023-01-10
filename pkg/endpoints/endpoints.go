package endpoints

import (
	"fmt"

	"golang.org/x/exp/slices"
	"helm.sh/helm/v3/pkg/chart"
)

func endpointAuthUsers(endpoint interface{}) []string {
	skipList := []string{"test"}
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
				"public": 443,
			},
		},
	}, nil
}

func ForChart(chart *chart.Chart, config *EndpointConfig) (map[string]interface{}, error) {
	endpoints := make(map[string]interface{})
	chartEndpoints := chart.Values["endpoints"].(map[string]interface{})

	for endpointName, endpointValues := range chartEndpoints {
		switch endpointName {
		case "cluster_domain_suffix":
		case "ingress":
		case "kube_dns":
		case "ldap":
		case "local_image_registry":
		case "oci_image_registry":
		case "fluentd":
			continue
		case "container_infra":
			endpoint, err := basicEndpoint(config.MagnumAPIHost)
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
		case "orchestration":
			endpoint, err := basicEndpoint(config.HeatHost)
			if err != nil {
				return nil, err
			}
			endpoints[endpointName] = endpoint
		case "oslo_cache":
			// TODO
			endpoints[endpointName] = map[string]interface{}{
				"auth": map[string]interface{}{
					"memcache_secret_key": config.MemcacheSecretKey,
				},
			}
		case "oslo_db":
			auth, err := endpointAuth(endpointValues, map[string]string{
				"admin":    config.DatabaseRootPassword,
				"keystone": config.KeystoneDatabasePassword,
				"magnum":   config.MagnumDatabasePassword,
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
				"admin":    config.RabbitmqAdminPassword,
				"keystone": config.KeystoneRabbitmqPassword,
				"magnum":   config.MagnumRabbitmqPassword,
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
				"statefulset": nil,
				"namespace":   config.RabbitmqNamespace,
				"hosts": map[string]interface{}{
					"default": config.RabbitmqServiceName,
				},
			}
		case "identity":
			auth, err := endpointAuth(endpointValues, map[string]string{
				"admin":             config.KeystoneAdminPassword,
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
		default:
			return nil, fmt.Errorf("endpoint %s not supported", endpointName)
		}
	}

	return endpoints, nil
}
