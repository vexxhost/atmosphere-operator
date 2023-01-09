package monitoring

import (
	"embed"
	"encoding/json"

	"github.com/google/go-jsonnet"
	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
)

//go:embed jsonnet/vendor
var jsonnetFiles embed.FS

func GetJsonnetVM() (*jsonnet.VM, error) {
	vm := jsonnet.MakeVM()
	vm.Importer(&EmbedImporter{
		FS:     jsonnetFiles,
		JPaths: []string{"jsonnet/vendor"},
	})

	return vm, nil
}

func GetMemcachedPrometheusRules() ([]monitoringv1.RuleGroup, error) {
	vm, err := GetJsonnetVM()
	if err != nil {
		return nil, err
	}

	jsonStr, err := vm.EvaluateAnonymousSnippet("rules.jsonnet", `
		local memcached = import 'github.com/grafana/jsonnet-libs/memcached-mixin/mixin.libsonnet';
		memcached.prometheusAlerts.groups
	`)
	if err != nil {
		return nil, err
	}

	groups := []monitoringv1.RuleGroup{}
	json.Unmarshal([]byte(jsonStr), &groups)

	return groups, nil
}
