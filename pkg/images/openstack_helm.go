package images

import (
	"golang.org/x/exp/slices"
	"helm.sh/helm/v3/pkg/chart"
)

func GetImageTagsForOpenstackHelmChart(chart *chart.Chart, registry string) (map[string]interface{}, error) {
	tags := make(map[string]interface{})

	// Loop over list of all image tags
	images := chart.Values["images"]
	chartTags := images.(map[string]interface{})["tags"].(map[string]interface{})

	for key, _ := range chartTags {
		if slices.Contains(SKIP_LIST, key) {
			continue
		}

		ref, err := GetImageReference(key)
		if err != nil {
			return nil, err
		}

		if registry == "" {
			tags[key] = ref.Remote()
			continue
		}

		image, err := OverrideRegistry(ref, registry)
		if err != nil {
			return nil, err
		}

		tags[key] = image.Remote()
	}

	return tags, nil
}
