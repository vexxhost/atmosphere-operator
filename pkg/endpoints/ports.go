package endpoints

import "helm.sh/helm/v3/pkg/chart"

func GetPortFromChart(chart *chart.Chart, endpointName string, portName string) int32 {
	endpoint := chart.Values["endpoints"].(map[string]interface{})[endpointName].(map[string]interface{})
	portConfig := endpoint["port"].(map[string]interface{})
	portGroup := portConfig[portName].(map[string]interface{})

	port, ok := portGroup["service"].(float64)
	if !ok {
		port = portGroup["default"].(float64)
	}

	return int32(port)
}
