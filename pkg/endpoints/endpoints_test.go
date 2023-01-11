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
	endpoint, err := basicEndpoint("test", "cloud.atmosphere.vexxhost.com")
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
				"public": float32(443),
			},
		},
	}, endpoint)
}

func TestTestBasicEndpointWithEmptyHost(t *testing.T) {
	_, err := basicEndpoint("foo", "")
	require.Error(t, err)
}

func assertEndpointHostFQDNOverride(t *testing.T, expected string, endpoints interface{}) {
	assert.Equal(t, expected, endpoints.(map[string]interface{})["host_fqdn_override"].(map[string]interface{})["public"].(map[string]interface{})["host"])
}

func TestForChart(t *testing.T) {
	folders, err := os.ReadDir("testdata/helm-charts")
	require.NoError(t, err)

	config := &EndpointConfig{
		RegionName:                    "RegionOne",
		MemcacheSecretKey:             "memcache",
		DatabaseNamespace:             "openstack",
		DatabaseServiceName:           "percona-xtradb",
		DatabaseRootPassword:          "db-root",
		RabbitmqNamespace:             "openstack",
		RabbitmqServiceName:           "keystone-rabbitmq",
		RabbitmqAdminUsername:         "rabbitmq-root",
		RabbitmqAdminPassword:         "rabbitmq-root",
		KeystoneHost:                  "cloud.atmosphere.vexxhost.com",
		KeystoneDatabasePassword:      "db-keystone",
		KeystoneRabbitmqPassword:      "rabbitmq-keystone",
		KeystoneAdminPassword:         "keystone-admin",
		BarbicanHost:                  "key-manager.atmosphere.vexxhost.com",
		BarbicanDatabasePassword:      "db-barbican",
		BarbicanRabbitmqPassword:      "rabbitmq-barbican",
		BarbicanKeystonePassword:      "barbican-keystone",
		GlanceHost:                    "image.atmosphere.vexxhost.com",
		GlanceDatabasePassword:        "db-glance",
		GlanceRabbitmqPassword:        "rabbitmq-glance",
		GlanceKeystonePassword:        "glance-keystone",
		CinderHost:                    "block-storage.atmosphere.vexxhost.com",
		CinderDatabasePassword:        "db-cinder",
		CinderRabbitmqPassword:        "rabbitmq-cinder",
		CinderKeystonePassword:        "cinder-keystone",
		PlacementHost:                 "placement.atmosphere.vexxhost.com",
		PlacementDatabasePassword:     "db-placement",
		PlacementKeystonePassword:     "placement-keystone",
		NeutronHost:                   "networking.atmosphere.vexxhost.com",
		NeutronDatabasePassword:       "db-neutron",
		NeutronRabbitmqPassword:       "rabbitmq-neutron",
		NeutronKeystonePassword:       "neutron-keystone",
		IronicHost:                    "baremetal.atmosphere.vexxhost.com",
		IronicKeystonePassword:        "ironic-keystone",
		NovaHost:                      "compute.atmosphere.vexxhost.com",
		NovaNovncHost:                 "vnc.atmosphere.vexxhost.com",
		NovaDatabasePassword:          "db-nova",
		NovaRabbitmqPassword:          "rabbitmq-nova",
		NovaKeystonePassword:          "nova-keystone",
		NovaMetadataSecret:            "metadata-secret",
		SenlinHost:                    "clustering.atmosphere.vexxhost.com",
		SenlinDatabasePassword:        "db-senlin",
		SenlinRabbitmqPassword:        "rabbitmq-senlin",
		SenlinKeystonePassword:        "senlin-keystone",
		DesignateHost:                 "dns.atmosphere.vexxhost.com",
		DesignateDatabasePassword:     "db-designate",
		DesignateRabbitmqPassword:     "rabbitmq-designate",
		DesignateKeystonePassword:     "designate-keystone",
		HeatHost:                      "orchestration.atmosphere.vexxhost.com",
		HeatCloudFormationHost:        "cloudformation.atmosphere.vexxhost.com",
		HeatDatabasePassword:          "db-heat",
		HeatRabbitmqPassword:          "rabbitmq-heat",
		HeatKeystonePassword:          "heat-keystone",
		HeatTrusteeKeystonePassword:   "heat-trustee-keystone",
		HeatStackUserKeystonePassword: "heat-stack-user-keystone",
		OctaviaHost:                   "load-balancer.atmosphere.vexxhost.com",
		OctaviaDatabasePassword:       "db-octavia",
		OctaviaRabbitmqPassword:       "rabbitmq-octavia",
		OctaviaKeystonePassword:       "octavia-keystone",
		MagnumHost:                    "container-infra.atmosphere.vexxhost.com",
		MagnumDatabasePassword:        "db-magnum",
		MagnumRabbitmqPassword:        "rabbitmq-magnum",
		MagnumKeystonePassword:        "magnum-keystone",
		HorizonHost:                   "dashboard.atmosphere.vexxhost.com",
		HorizonDatabasePassword:       "db-horizon",
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
				case "baremetal":
					assertEndpointHostFQDNOverride(t, config.IronicHost, endpoints[endpoint])
				case "clustering":
					assertEndpointHostFQDNOverride(t, config.SenlinHost, endpoints[endpoint])
				case "cloudformation":
					assertEndpointHostFQDNOverride(t, config.HeatCloudFormationHost, endpoints[endpoint])
				case "compute":
					assertEndpointHostFQDNOverride(t, config.NovaHost, endpoints[endpoint])
				case "compute_metadata":
					assert.Equal(t, config.NovaMetadataSecret, endpoints[endpoint].(map[string]interface{})["secret"])
					assert.Equal(t, "nova-metadata", endpoints[endpoint].(map[string]interface{})["hosts"].(map[string]interface{})["public"])
					assert.Equal(t, 8775, endpoints[endpoint].(map[string]interface{})["port"].(map[string]interface{})["metadata"].(map[string]interface{})["public"])
				case "compute_novnc_proxy":
					assertEndpointHostFQDNOverride(t, config.NovaNovncHost, endpoints[endpoint])
					assert.Equal(t, "/vnc_lite.html", endpoints[endpoint].(map[string]interface{})["path"].(map[string]interface{})["default"])
				case "container_infra":
					assertEndpointHostFQDNOverride(t, config.MagnumHost, endpoints[endpoint])
				case "dashboard":
					assertEndpointHostFQDNOverride(t, config.HorizonHost, endpoints[endpoint])
				case "dns":
					assertEndpointHostFQDNOverride(t, config.DesignateHost, endpoints[endpoint])
				case "identity":
					auth := endpoints[endpoint].(map[string]interface{})["auth"].(map[string]interface{})

					assert.Equal(t, config.RegionName, auth["admin"].(map[string]interface{})["region_name"])
					assert.Equal(t, fmt.Sprintf("admin-%s", config.RegionName), auth["admin"].(map[string]interface{})["username"])
					assert.Equal(t, config.KeystoneAdminPassword, auth["admin"].(map[string]interface{})["password"])

					if tc.Name() == "barbican" {
						assert.Equal(t, config.BarbicanKeystonePassword, auth["barbican"].(map[string]interface{})["password"])
					} else if tc.Name() == "glance" {
						assert.Equal(t, config.GlanceKeystonePassword, auth["glance"].(map[string]interface{})["password"])
					} else if tc.Name() == "cinder" {
						assert.Equal(t, config.CinderKeystonePassword, auth["cinder"].(map[string]interface{})["password"])
					} else if tc.Name() == "placement" {
						assert.Equal(t, config.PlacementKeystonePassword, auth["placement"].(map[string]interface{})["password"])
					} else if tc.Name() == "neutron" {
						assert.Equal(t, config.NeutronKeystonePassword, auth["neutron"].(map[string]interface{})["password"])
					} else if tc.Name() == "ironic" {
						assert.Equal(t, config.IronicKeystonePassword, auth["ironic"].(map[string]interface{})["password"])
					} else if tc.Name() == "nova" {
						assert.Equal(t, config.NovaKeystonePassword, auth["nova"].(map[string]interface{})["password"])
					} else if tc.Name() == "senlin" {
						assert.Equal(t, config.SenlinKeystonePassword, auth["senlin"].(map[string]interface{})["password"])
					} else if tc.Name() == "designate" {
						assert.Equal(t, config.DesignateKeystonePassword, auth["designate"].(map[string]interface{})["password"])
					} else if tc.Name() == "heat" {
						assert.Equal(t, config.HeatKeystonePassword, auth["heat"].(map[string]interface{})["password"])
					} else if tc.Name() == "octavia" {
						assert.Equal(t, config.OctaviaKeystonePassword, auth["octavia"].(map[string]interface{})["password"])
					} else if tc.Name() == "magnum" {
						assert.Equal(t, config.MagnumKeystonePassword, auth["magnum"].(map[string]interface{})["password"])
					} else if tc.Name() != "keystone" && tc.Name() != "horizon" {
						t.Errorf("untested keystone configuration for %s", tc.Name())
					}

					assert.Equal(t, "keystone-api", endpoints[endpoint].(map[string]interface{})["hosts"].(map[string]interface{})["default"])
					assert.Equal(t, float32(5000), endpoints[endpoint].(map[string]interface{})["port"].(map[string]interface{})["api"].(map[string]interface{})["default"])
					assert.Equal(t, float32(443), endpoints[endpoint].(map[string]interface{})["port"].(map[string]interface{})["api"].(map[string]interface{})["public"])

					assertEndpointHostFQDNOverride(t, config.KeystoneHost, endpoints[endpoint])
				case "image":
					assertEndpointHostFQDNOverride(t, config.GlanceHost, endpoints[endpoint])
				case "key_manager":
					assertEndpointHostFQDNOverride(t, config.BarbicanHost, endpoints[endpoint])
				case "load_balancer":
					assertEndpointHostFQDNOverride(t, config.OctaviaHost, endpoints[endpoint])
				case "mdns":
					assertEndpointHostFQDNOverride(t, config.DesignateHost, endpoints[endpoint])
				case "network":
					assertEndpointHostFQDNOverride(t, config.NeutronHost, endpoints[endpoint])
				case "orchestration":
					assertEndpointHostFQDNOverride(t, config.HeatHost, endpoints[endpoint])
				case "oslo_cache":
					assert.Equal(t, config.MemcacheSecretKey, endpoints[endpoint].(map[string]interface{})["auth"].(map[string]interface{})["memcache_secret_key"])
				case "oslo_db", "oslo_db_api", "oslo_db_cell0":
					auth := endpoints[endpoint].(map[string]interface{})["auth"].(map[string]interface{})

					assert.Equal(t, config.DatabaseRootPassword, auth["admin"].(map[string]interface{})["password"])

					if tc.Name() == "keystone" {
						assert.Equal(t, config.KeystoneDatabasePassword, auth["keystone"].(map[string]interface{})["password"])
					} else if tc.Name() == "barbican" {
						assert.Equal(t, config.BarbicanDatabasePassword, auth["barbican"].(map[string]interface{})["password"])
					} else if tc.Name() == "glance" {
						assert.Equal(t, config.GlanceDatabasePassword, auth["glance"].(map[string]interface{})["password"])
					} else if tc.Name() == "cinder" {
						assert.Equal(t, config.CinderDatabasePassword, auth["cinder"].(map[string]interface{})["password"])
					} else if tc.Name() == "placement" {
						assert.Equal(t, config.PlacementDatabasePassword, auth["placement"].(map[string]interface{})["password"])
					} else if tc.Name() == "neutron" {
						assert.Equal(t, config.NeutronDatabasePassword, auth["neutron"].(map[string]interface{})["password"])
					} else if tc.Name() == "nova" {
						assert.Equal(t, config.NovaDatabasePassword, auth["nova"].(map[string]interface{})["password"])
					} else if tc.Name() == "senlin" {
						assert.Equal(t, config.SenlinDatabasePassword, auth["senlin"].(map[string]interface{})["password"])
					} else if tc.Name() == "designate" {
						assert.Equal(t, config.DesignateDatabasePassword, auth["designate"].(map[string]interface{})["password"])
					} else if tc.Name() == "heat" {
						assert.Equal(t, config.HeatDatabasePassword, auth["heat"].(map[string]interface{})["password"])
					} else if tc.Name() == "octavia" {
						assert.Equal(t, config.OctaviaDatabasePassword, auth["octavia"].(map[string]interface{})["password"])
					} else if tc.Name() == "magnum" {
						assert.Equal(t, config.MagnumDatabasePassword, auth["magnum"].(map[string]interface{})["password"])
					} else if tc.Name() == "horizon" {
						assert.Equal(t, config.HorizonDatabasePassword, auth["horizon"].(map[string]interface{})["password"])
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
					} else if tc.Name() == "barbican" {
						assert.Equal(t, config.BarbicanRabbitmqPassword, auth["barbican"].(map[string]interface{})["password"])
					} else if tc.Name() == "glance" {
						assert.Equal(t, config.GlanceRabbitmqPassword, auth["glance"].(map[string]interface{})["password"])
					} else if tc.Name() == "cinder" {
						assert.Equal(t, config.CinderRabbitmqPassword, auth["cinder"].(map[string]interface{})["password"])
					} else if tc.Name() == "neutron" {
						assert.Equal(t, config.NeutronRabbitmqPassword, auth["neutron"].(map[string]interface{})["password"])
					} else if tc.Name() == "nova" {
						assert.Equal(t, config.NovaRabbitmqPassword, auth["nova"].(map[string]interface{})["password"])
					} else if tc.Name() == "senlin" {
						assert.Equal(t, config.SenlinRabbitmqPassword, auth["senlin"].(map[string]interface{})["password"])
					} else if tc.Name() == "designate" {
						assert.Equal(t, config.DesignateRabbitmqPassword, auth["designate"].(map[string]interface{})["password"])
					} else if tc.Name() == "heat" {
						assert.Equal(t, config.HeatRabbitmqPassword, auth["heat"].(map[string]interface{})["password"])
					} else if tc.Name() == "octavia" {
						assert.Equal(t, config.OctaviaRabbitmqPassword, auth["octavia"].(map[string]interface{})["password"])
					} else if tc.Name() == "magnum" {
						assert.Equal(t, config.MagnumRabbitmqPassword, auth["magnum"].(map[string]interface{})["password"])
					} else {
						t.Errorf("untested rabbitmq configuration for %s", tc.Name())
					}

					assert.Equal(t, false, endpoints[endpoint].(map[string]interface{})["statefulset"])
					assert.Equal(t, config.RabbitmqNamespace, endpoints[endpoint].(map[string]interface{})["namespace"])
					assert.Equal(t, config.RabbitmqServiceName, endpoints[endpoint].(map[string]interface{})["hosts"].(map[string]interface{})["default"])
				case "placement":
					assertEndpointHostFQDNOverride(t, config.PlacementHost, endpoints[endpoint])
				case "volumev3":
					assertEndpointHostFQDNOverride(t, config.CinderHost, endpoints[endpoint])
				default:
					t.Errorf("unexpected endpoint: %s", endpoint)
				}
			}
		})
	}

}
