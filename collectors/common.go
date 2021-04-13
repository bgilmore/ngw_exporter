package collectors

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bgilmore/ngw_exporter/models"
	"github.com/prometheus/client_golang/prometheus"
)

func metricName(module, metric string) string {
	return prometheus.BuildFQName("ngw", module, metric)
}

func scrape(target string, m models.Model) error {
	resp, err := http.Get(models.URL(m, target))
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("http: scape returned %s", resp.Status)
	}
	if err := json.NewDecoder(resp.Body).Decode(&m); err != nil {
		return fmt.Errorf("json: %v", err)
	}
	return nil
}
