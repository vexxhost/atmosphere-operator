package monitoring

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetMemcachedPrometheusRules(t *testing.T) {
	groups, err := GetMemcachedPrometheusRules()
	require.NoError(t, err)

	assert.Equal(t, 1, len(groups))
	assert.Equal(t, "memcached", groups[0].Name)
	assert.Equal(t, 3, len(groups[0].Rules))
}
