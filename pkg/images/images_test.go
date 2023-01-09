package images

import (
	"testing"

	dockerparser "github.com/novln/docker-parser"
	"github.com/stretchr/testify/assert"
)

func TestGetImageRef(t *testing.T) {
	image, err := GetImageReference("dep_check")

	assert.NoError(t, err)
	assert.Equal(t, "quay.io/vexxhost/kubernetes-entrypoint:latest", image.Remote())
}

func TestOverrideRegistry(t *testing.T) {
	testCases := []struct {
		name     string
		expected string
	}{
		{
			name:     "docker.io/libary/memcached:1.6.17",
			expected: "localhost:3000/memcached:1.6.17",
		},
		{
			name:     "us-docker.pkg.dev/vexxhost-infra/openstack/senlin:wallaby",
			expected: "localhost:3000/senlin:wallaby",
		},
		{
			name:     "quay.io/vexxhost/magnum@sha256:46e7c910780864f4532ecc85574f159a36794f37aac6be65e4b48c46040ced17",
			expected: "localhost:3000/magnum@sha256:46e7c910780864f4532ecc85574f159a36794f37aac6be65e4b48c46040ced17",
		},
		{
			name:     "k8s.gcr.io/ingress-nginx/controller:v1.1.1@sha256:0bc88eb15f9e7f84e8e56c14fa5735aaa488b840983f87bd79b1054190e660de",
			expected: "localhost:3000/controller@sha256:0bc88eb15f9e7f84e8e56c14fa5735aaa488b840983f87bd79b1054190e660de",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ref, err := dockerparser.Parse(tc.name)
			assert.NoError(t, err)

			overriddenRef, err := OverrideRegistry(ref, "localhost:3000")
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, overriddenRef.Remote())
		})
	}
}

func TestGetImageRefWithInvalidImageName(t *testing.T) {
	_, err := GetImageReference("invalid")

	assert.Error(t, err)
}
