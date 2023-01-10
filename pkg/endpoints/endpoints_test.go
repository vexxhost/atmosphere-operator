package endpoints

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"helm.sh/helm/v3/pkg/chart/loader"
)

func TestBasicEndpoint(t *testing.T) {
	endpoint, err := basicEndpoint("cloud.atmosphere.vexxhost.com")
	require.NoError(t, err)

	assert.Equal(t, map[string]interface{}{
		"scheme": map[string]interface{}{
			"public": "https",
		},
		"host_fqdn_override": map[string]interface{}{
			"public": map[string]interface{}{
				"host": "cloud.atmosphere.vexxhost.com",
			},
		},
		"port": map[string]interface{}{
			"api": map[string]interface{}{
				"public": 443,
			},
		},
	}, endpoint)
}

func TestTestBasicEndpointWithEmptyHost(t *testing.T) {
	_, err := basicEndpoint("")
	require.Error(t, err)
}

func assertEndpointHostFQDNOverride(t *testing.T, expected string, endpoints interface{}) {
	assert.Equal(t, expected, endpoints.(map[string]interface{})["host_fqdn_override"].(map[string]interface{})["public"].(map[string]interface{})["host"])
}

func TestForChart(t *testing.T) {
	folders, err := os.ReadDir("testdata/helm-charts")
	require.NoError(t, err)

	config := &EndpointConfig{
		RegionName:               "RegionOne",
		MemcacheSecretKey:        "memcache",
		DatabaseNamespace:        "openstack",
		DatabaseServiceName:      "percona-xtradb",
		DatabaseRootPassword:     "db-root",
		RabbitmqNamespace:        "openstack",
		RabbitmqServiceName:      "keystone-rabbitmq",
		RabbitmqAdminUsername:    "rabbitmq-root",
		RabbitmqAdminPassword:    "rabbitmq-root",
		KeystoneHost:             "cloud.atmosphere.vexxhost.com",
		KeystoneDatabasePassword: "db-keystone",
		KeystoneRabbitmqPassword: "rabbitmq-keystone",
		KeystoneAdminPassword:    "keystone-admin",
		BarbicanHost:             "key-manager.atmosphere.vexxhost.com",
		HeatHost:                 "orchestration.atmosphere.vexxhost.com",
		MagnumAPIHost:            "container-infra.atmosphere.vexxhost.com",
		MagnumDatabasePassword:   "db-magnum",
		MagnumRabbitmqPassword:   "rabbitmq-magnum",
		MagnumKeystonePassword:   "magnum-keystone",
	}

	for _, tc := range folders {
		if !tc.IsDir() || tc.Name() == "helm-toolkit" {
			continue
		}

		t.Run(tc.Name(), func(t *testing.T) {
			chart, err := loader.Load(fmt.Sprintf("testdata/helm-charts/%s.tgz", tc.Name()))
			require.NoError(t, err)

			endpoints, err := ForChart(chart, config)
			require.NoError(t, err)

			for endpoint := range endpoints {
				switch endpoint {
				case "container_infra":
					assertEndpointHostFQDNOverride(t, config.MagnumAPIHost, endpoints["container_infra"])
				case "key_manager":
					assertEndpointHostFQDNOverride(t, config.BarbicanHost, endpoints["key_manager"])
				case "orchestration":
					assertEndpointHostFQDNOverride(t, config.HeatHost, endpoints["orchestration"])
				case "oslo_cache":
					assert.Equal(t, config.MemcacheSecretKey, endpoints["oslo_cache"].(map[string]interface{})["auth"].(map[string]interface{})["memcache_secret_key"])
				case "oslo_db":
					auth := endpoints[endpoint].(map[string]interface{})["auth"].(map[string]interface{})

					assert.Equal(t, config.DatabaseRootPassword, auth["admin"].(map[string]interface{})["password"])

					if tc.Name() == "keystone" {
						assert.Equal(t, config.KeystoneDatabasePassword, auth["keystone"].(map[string]interface{})["password"])
					} else if tc.Name() == "magnum" {
						assert.Equal(t, config.MagnumDatabasePassword, auth["magnum"].(map[string]interface{})["password"])
					} else {
						t.Errorf("untested database configuration for %s", tc.Name())
					}

					assert.Equal(t, config.DatabaseNamespace, endpoints[endpoint].(map[string]interface{})["namespace"])
					assert.Equal(t, config.DatabaseServiceName, endpoints[endpoint].(map[string]interface{})["hosts"].(map[string]interface{})["default"])
				case "oslo_messaging":
					auth := endpoints[endpoint].(map[string]interface{})["auth"].(map[string]interface{})

					assert.Equal(t, config.RabbitmqAdminUsername, auth["admin"].(map[string]interface{})["username"])
					assert.Equal(t, config.RabbitmqAdminPassword, auth["admin"].(map[string]interface{})["password"])

					if tc.Name() == "keystone" {
						assert.Equal(t, config.KeystoneRabbitmqPassword, auth["keystone"].(map[string]interface{})["password"])
					} else if tc.Name() == "magnum" {
						assert.Equal(t, config.MagnumRabbitmqPassword, auth["magnum"].(map[string]interface{})["password"])
					} else {
						t.Errorf("untested rabbitmq configuration for %s", tc.Name())
					}

					assert.Equal(t, nil, endpoints[endpoint].(map[string]interface{})["statefulset"])
					assert.Equal(t, config.RabbitmqNamespace, endpoints[endpoint].(map[string]interface{})["namespace"])
					assert.Equal(t, config.RabbitmqServiceName, endpoints[endpoint].(map[string]interface{})["hosts"].(map[string]interface{})["default"])
				case "identity":
					auth := endpoints[endpoint].(map[string]interface{})["auth"].(map[string]interface{})

					assert.Equal(t, config.RegionName, auth["admin"].(map[string]interface{})["region_name"])
					assert.Equal(t, fmt.Sprintf("admin-%s", config.RegionName), auth["admin"].(map[string]interface{})["username"])
					assert.Equal(t, config.KeystoneAdminPassword, auth["admin"].(map[string]interface{})["password"])

					if tc.Name() == "magnum" {
						assert.Equal(t, config.MagnumKeystonePassword, auth["magnum"].(map[string]interface{})["password"])
					} else if tc.Name() != "keystone" {
						t.Errorf("untested keystone configuration for %s", tc.Name())
					}

					assert.Equal(t, "keystone-api", endpoints[endpoint].(map[string]interface{})["hosts"].(map[string]interface{})["default"])
					assert.Equal(t, 5000, endpoints[endpoint].(map[string]interface{})["port"].(map[string]interface{})["api"].(map[string]interface{})["default"])
					assert.Equal(t, 443, endpoints[endpoint].(map[string]interface{})["port"].(map[string]interface{})["api"].(map[string]interface{})["public"])

					assertEndpointHostFQDNOverride(t, config.KeystoneHost, endpoints[endpoint])
				default:
					t.Errorf("unexpected endpoint: %s", endpoint)
				}
			}
		})
	}

}
