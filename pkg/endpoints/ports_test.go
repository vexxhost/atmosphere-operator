package endpoints

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"helm.sh/helm/v3/pkg/chart/loader"
)

func TestGetPortFromChart(t *testing.T) {
	testCases := []struct {
		chart    string
		endpoint string
		port     string
		expected int
	}{
		{
			chart:    "keystone",
			endpoint: "identity",
			port:     "api",
			expected: 5000,
		},
		{
			chart:    "barbican",
			endpoint: "key_manager",
			port:     "api",
			expected: 9311,
		},
		{
			chart:    "octavia",
			endpoint: "load_balancer",
			port:     "api",
			expected: 9876,
		},
		{
			chart:    "nova",
			endpoint: "compute",
			port:     "api",
			expected: 8774,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.chart, func(t *testing.T) {
			chart, err := loader.Load(fmt.Sprintf("testdata/helm-charts/%s", tc.chart))
			require.NoError(t, err)

			port := GetPortFromChart(chart, tc.endpoint, tc.port)
			require.Equal(t, int32(tc.expected), port)
		})
	}
}
