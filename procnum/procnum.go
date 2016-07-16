package procnum

import (
	"time"

	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/control/plugin/cpolicy"
	"github.com/intelsdi-x/snap/core"
	"github.com/intelsdi-x/snap/core/ctypes"
)

const (
	// Name of plugin
	Name = "procnum"
	// Version of plugin
	Version = 1
	// Type of plugin
	Type = plugin.CollectorPluginType
)

var (
	// make sure that we actually satisify requierd interface
	_ plugin.CollectorPlugin = (*Procnum)(nil)

	metricNames = []string{
		"proc_num",
	}
)

type Procnum struct {
}

func New() *Procnum {
	return &Procnum{}
}

// CollectMetrics collects metrics for testing
func (p *Procnum) CollectMetrics(mts []plugin.MetricType) ([]plugin.MetricType, error) {
	var err error

	conf := mts[0].Config().Table()
	var statpath string
	statpathConf, ok := conf["statpath"]
	if !ok {
		statpath = statpathConf.(ctypes.ConfigValueStr).Value
	} else {
		statpath = "/proc/stat"
	}

	metrics, err := procnum(statpath, mts)
	if err != nil {
		return nil, err
	}

	return metrics, nil
}

func procnum(statpath string, mts []plugin.MetricType) ([]plugin.MetricType, error) {
	result, err := GatherProcInfo(statpath)
	if err != nil {
		return nil, err
	}
	runTime := time.Now()

	// leave room for expansion
	stats := make(map[string]float64)
	if result.Procnum != nil {
		stats["proc_num"] = *result.Procnum
	}

	metrics := make([]plugin.MetricType, 0, len(stats))
	for _, m := range mts {
		stat := m.Namespace()[2].Value
		if value, ok := stats[stat]; ok {
			mt := plugin.MetricType{
				Data_:      value,
				Namespace_: core.NewNamespace("raintank", "processes", stat),
				Timestamp_: runTime,
				Version_:   m.Version(),
			}
			metrics = append(metrics, mt)
		}
	}

	return metrics, nil
}

//GetMetricTypes returns metric types for testing
func (p *Procnum) GetMetricTypes(cfg plugin.ConfigType) ([]plugin.MetricType, error) {
	mts := []plugin.MetricType{}
	for _, metricName := range metricNames {
		mts = append(mts, plugin.MetricType{
			Namespace_: core.NewNamespace("raintank", "processes", metricName),
		})
	}
	return mts, nil
}

//GetConfigPolicy returns a ConfigPolicyTree for testing
func (p *Procnum) GetConfigPolicy() (*cpolicy.ConfigPolicy, error) {
	c := cpolicy.New()
	rule0, _ := cpolicy.NewStringRule("statpath", false, "/proc/stat")
	cp := cpolicy.NewPolicyNode()
	cp.Add(rule0)
	c.Add([]string{"raintank", "processes"}, cp)
	return c, nil
}

//Meta returns meta data for testing
func Meta() *plugin.PluginMeta {
	return plugin.NewPluginMeta(
		Name,
		Version,
		Type,
		[]string{plugin.SnapGOBContentType},
		[]string{plugin.SnapGOBContentType},
		plugin.Unsecure(true),
		plugin.ConcurrencyCount(5000),
	)
}
