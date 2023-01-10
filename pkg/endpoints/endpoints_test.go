package endpoints

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"helm.sh/helm/v3/pkg/chart/loader"
)

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
				case "oslo_cache":
					assert.Equal(t, config.MemcacheSecretKey, endpoints["oslo_cache"].(map[string]interface{})["auth"].(map[string]interface{})["memcache_secret_key"])
				case "oslo_db":
					assert.Equal(t, config.DatabaseRootPassword, endpoints[endpoint].(map[string]interface{})["auth"].(map[string]interface{})["admin"].(map[string]interface{})["password"])
					assert.Equal(t, config.KeystoneDatabasePassword, endpoints[endpoint].(map[string]interface{})["auth"].(map[string]interface{})["keystone"].(map[string]interface{})["password"])
					assert.Equal(t, config.DatabaseNamespace, endpoints[endpoint].(map[string]interface{})["namespace"])
					assert.Equal(t, config.DatabaseServiceName, endpoints[endpoint].(map[string]interface{})["hosts"].(map[string]interface{})["default"])
				case "oslo_messaging":
					assert.Equal(t, config.RabbitmqAdminUsername, endpoints[endpoint].(map[string]interface{})["auth"].(map[string]interface{})["admin"].(map[string]interface{})["username"])
					assert.Equal(t, config.RabbitmqAdminPassword, endpoints[endpoint].(map[string]interface{})["auth"].(map[string]interface{})["admin"].(map[string]interface{})["password"])
					assert.Equal(t, config.KeystoneRabbitmqPassword, endpoints[endpoint].(map[string]interface{})["auth"].(map[string]interface{})["keystone"].(map[string]interface{})["password"])
					assert.Equal(t, nil, endpoints[endpoint].(map[string]interface{})["statefulset"])
					assert.Equal(t, config.RabbitmqNamespace, endpoints[endpoint].(map[string]interface{})["namespace"])
					assert.Equal(t, config.RabbitmqServiceName, endpoints[endpoint].(map[string]interface{})["hosts"].(map[string]interface{})["default"])
				case "identity":
					assert.Equal(t, config.RegionName, endpoints[endpoint].(map[string]interface{})["auth"].(map[string]interface{})["admin"].(map[string]interface{})["region_name"])
					assert.Equal(t, fmt.Sprintf("admin-%s", config.RegionName), endpoints[endpoint].(map[string]interface{})["auth"].(map[string]interface{})["admin"].(map[string]interface{})["username"])
					assert.Equal(t, config.KeystoneAdminPassword, endpoints[endpoint].(map[string]interface{})["auth"].(map[string]interface{})["admin"].(map[string]interface{})["password"])
					assert.Equal(t, "keystone-api", endpoints[endpoint].(map[string]interface{})["hosts"].(map[string]interface{})["default"])
					assert.Equal(t, "https", endpoints[endpoint].(map[string]interface{})["scheme"].(map[string]interface{})["public"])
					assert.Equal(t, config.KeystoneHost, endpoints[endpoint].(map[string]interface{})["host_fqdn_override"].(map[string]interface{})["public"].(map[string]interface{})["host"])
					assert.Equal(t, 5000, endpoints[endpoint].(map[string]interface{})["port"].(map[string]interface{})["api"].(map[string]interface{})["default"])
					assert.Equal(t, 443, endpoints[endpoint].(map[string]interface{})["port"].(map[string]interface{})["api"].(map[string]interface{})["public"])
				default:
					t.Errorf("unexpected endpoint: %s", endpoint)
				}
			}
		})
	}

}
