package images

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"helm.sh/helm/v3/pkg/chart/loader"
)

func TestGetImageTagsForOpenstackHelmChart(t *testing.T) {
	testCases := []struct {
		chart    string
		registry string
		expected map[string]interface{}
	}{
		{
			chart:    "testdata/helm-charts/memcached",
			registry: "",
			expected: map[string]interface{}{
				"dep_check":                     "quay.io/vexxhost/kubernetes-entrypoint:latest",
				"memcached":                     "docker.io/library/memcached:1.6.17",
				"prometheus_memcached_exporter": "quay.io/prometheus/memcached-exporter:v0.10.0",
			},
		},
		{
			chart:    "testdata/helm-charts/memcached",
			registry: "localhost:3000",
			expected: map[string]interface{}{
				"dep_check":                     "localhost:3000/kubernetes-entrypoint:latest",
				"memcached":                     "localhost:3000/memcached:1.6.17",
				"prometheus_memcached_exporter": "localhost:3000/memcached-exporter:v0.10.0",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.chart, func(t *testing.T) {
			chart, err := loader.Load(tc.chart)
			require.NoError(t, err)

			tags, err := GetImageTagsForOpenstackHelmChart(chart, tc.registry)
			require.NoError(t, err)

			assert.Equal(t, tc.expected, tags)
		})
	}
}
