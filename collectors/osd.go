package collectors

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/prometheus/client_golang/prometheus"
)

//OsdCollector sample comment
type OsdCollector struct {
	conn Conn

	//CrushWeight is a persistent setting, and it affects how CRUSH assigns data to OSDs.
	//It displays the CRUSH weight for the OSD
	CrushWeight *prometheus.GaugeVec

	//Depth displays the OSD's level of hierarchy in the CRUSH map
	Depth *prometheus.GaugeVec

	//Reweight sets an override weight on the OSD.
	//It displays value within 0 to 1.
	Reweight *prometheus.GaugeVec

	//Bytes displays the total bytes available in the OSD
	Bytes *prometheus.GaugeVec

	//UsedBytes displays the total used bytes in the OSD
	UsedBytes *prometheus.GaugeVec

	//AvailBytes displays the total available bytes in the OSD
	AvailBytes *prometheus.GaugeVec

	//Utilization displays current utilization of the OSD
	Utilization *prometheus.GaugeVec

	//Pgs displays total no. of placement groups in the OSD.
	//Available in Ceph Jewel version.
	Pgs *prometheus.GaugeVec

	//CommitLatency displays in seconds how long it takes for an operation to be applied to disk
	CommitLatency *prometheus.GaugeVec

	//ApplyLatency displays in seconds how long it takes to get applied to the backing filesystem
	ApplyLatency *prometheus.GaugeVec

	//OsdsIn displays the In state of the OSD
	OsdIn *prometheus.GaugeVec

	//OsdsUP displays the Up state of the OSD
	OsdUp *prometheus.GaugeVec

	//TotalBytes displays total bytes in all OSDs
	TotalBytes prometheus.Gauge

	//TotalUsedBytes displays total used bytes in all OSDs
	TotalUsedBytes prometheus.Gauge

	//TotalAvailBytes displays total available bytes in all OSDs
	TotalAvailBytes prometheus.Gauge

	//AverageUtil displays average utilization in all OSDs
	AverageUtil prometheus.Gauge
}

//NewOsdCollector creates an instance of the OsdCollector and instantiates
// the individual metrics that show information about the osd.
func NewOsdCollector(conn Conn) *OsdCollector {
	return &OsdCollector{
		conn: conn,

		CrushWeight: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: cephNamespace,
				Name:      "osd_crush_weight",
				Help:      "OSD Crush Weight",
			},
			[]string{"osd"},
		),

		Depth: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: cephNamespace,
				Name:      "osd_depth",
				Help:      "OSD Depth",
			},
			[]string{"osd"},
		),

		Reweight: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: cephNamespace,
				Name:      "osd_reweight",
				Help:      "OSD Reweight",
			},
			[]string{"osd"},
		),

		Bytes: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: cephNamespace,
				Name:      "osd_bytes",
				Help:      "OSD Total Bytes",
			},
			[]string{"osd"},
		),

		UsedBytes: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: cephNamespace,
				Name:      "osd_used_bytes",
				Help:      "OSD Used Storage in Bytes",
			},
			[]string{"osd"},
		),

		AvailBytes: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: cephNamespace,
				Name:      "osd_avail_bytes",
				Help:      "OSD Available Storage in Bytes",
			},
			[]string{"osd"},
		),

		Utilization: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: cephNamespace,
				Name:      "osd_utilization",
				Help:      "OSD Utilization",
			},
			[]string{"osd"},
		),

		Pgs: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: cephNamespace,
				Name:      "osd_pgs",
				Help:      "OSD Placement Group Count",
			},
			[]string{"osd"},
		),

		TotalBytes: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: cephNamespace,
				Name:      "osd_total_bytes",
				Help:      "OSD Total Storage Bytes",
			},
		),
		TotalUsedBytes: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: cephNamespace,
				Name:      "osd_total_used_bytes",
				Help:      "OSD Total Used Storage Bytes",
			},
		),

		TotalAvailBytes: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: cephNamespace,
				Name:      "osd_total_avail_bytes",
				Help:      "OSD Total Available Storage Bytes ",
			},
		),

		AverageUtil: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: cephNamespace,
				Name:      "osd_average_utilization",
				Help:      "OSD Average Utilization",
			},
		),

		CommitLatency: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: cephNamespace,
				Name:      "osd_perf_commit_latency_seconds",
				Help:      "OSD Perf Commit Latency",
			},
			[]string{"osd"},
		),

		ApplyLatency: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: cephNamespace,
				Name:      "osd_perf_apply_latency_seconds",
				Help:      "OSD Perf Apply Latency",
			},
			[]string{"osd"},
		),

		OsdIn: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: cephNamespace,
				Name:      "osd_in",
				Help:      "OSD In Status",
			},
			[]string{"osd"},
		),

		OsdUp: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: cephNamespace,
				Name:      "osd_up",
				Help:      "OSD Up Status",
			},
			[]string{"osd"},
		),
	}
}

func (o *OsdCollector) collectorList() []prometheus.Collector {
	return []prometheus.Collector{
		o.CrushWeight,
		o.Depth,
		o.Reweight,
		o.Bytes,
		o.UsedBytes,
		o.AvailBytes,
		o.Utilization,
		o.Pgs,
		o.TotalBytes,
		o.TotalUsedBytes,
		o.TotalAvailBytes,
		o.AverageUtil,
		o.CommitLatency,
		o.ApplyLatency,
		o.OsdIn,
		o.OsdUp,
	}
}

type cephOsdDf struct {
	OsdNodes []struct {
		Name        string      `json:"name"`
		CrushWeight json.Number `json:"crush_weight"`
		Depth       json.Number `json:"depth"`
		Reweight    json.Number `json:"reweight"`
		KB          json.Number `json:"kb"`
		UsedKB      json.Number `json:"kb_used"`
		AvailKB     json.Number `json:"kb_avail"`
		Utilization json.Number `json:"utilization"`
		Pgs         json.Number `json:"pgs"`
	} `json:"nodes"`

	Summary struct {
		TotalKB      json.Number `json:"total_kb"`
		TotalUsedKB  json.Number `json:"total_kb_used"`
		TotalAvailKB json.Number `json:"total_kb_avail"`
		AverageUtil  json.Number `json:"average_utilization"`
	} `json:"summary"`
}

type cephPerfStat struct {
	PerfInfo []struct {
		ID    json.Number `json:"id"`
		Stats struct {
			CommitLatency json.Number `json:"commit_latency_ms"`
			ApplyLatency  json.Number `json:"apply_latency_ms"`
		} `json:"perf_stats"`
	} `json:"osd_perf_infos"`
}

type cephOsdDump struct {
	Osds []struct {
		Osd json.Number `json:"osd"`
		Up  json.Number `json:"up"`
		In  json.Number `json:"in"`
	} `json:"osds"`
}

func (o *OsdCollector) collect() error {
	cmd := o.cephOSDDfCommand()

	buf, _, err := o.conn.MonCommand(cmd)
	if err != nil {
		log.Println("[ERROR] Unable to collect data from ceph osd df", err)
		return err
	}

	osdDf := &cephOsdDf{}
	if err := json.Unmarshal(buf, osdDf); err != nil {
		return err
	}

	for _, node := range osdDf.OsdNodes {

		crushWeight, err := node.CrushWeight.Float64()
		if err != nil {
			return err
		}

		o.CrushWeight.WithLabelValues(node.Name).Set(crushWeight)

		depth, err := node.Depth.Float64()
		if err != nil {

			return err
		}

		o.Depth.WithLabelValues(node.Name).Set(depth)

		reweight, err := node.Reweight.Float64()
		if err != nil {
			return err
		}

		o.Reweight.WithLabelValues(node.Name).Set(reweight)

		osdKB, err := node.KB.Float64()
		if err != nil {
			return nil
		}

		o.Bytes.WithLabelValues(node.Name).Set(osdKB * 1e3)

		usedKB, err := node.UsedKB.Float64()
		if err != nil {
			return err
		}

		o.UsedBytes.WithLabelValues(node.Name).Set(usedKB * 1e3)

		availKB, err := node.AvailKB.Float64()
		if err != nil {
			return err
		}

		o.AvailBytes.WithLabelValues(node.Name).Set(availKB * 1e3)

		util, err := node.Utilization.Float64()
		if err != nil {
			return err
		}

		o.Utilization.WithLabelValues(node.Name).Set(util)

		pgs, err := node.Pgs.Float64()
		if err != nil {
			continue
		}

		o.Pgs.WithLabelValues(node.Name).Set(pgs)

	}

	totalKB, err := osdDf.Summary.TotalKB.Float64()
	if err != nil {
		return err
	}

	o.TotalBytes.Set(totalKB * 1e3)

	totalUsedKB, err := osdDf.Summary.TotalUsedKB.Float64()
	if err != nil {
		return err
	}

	o.TotalUsedBytes.Set(totalUsedKB * 1e3)

	totalAvailKB, err := osdDf.Summary.TotalAvailKB.Float64()
	if err != nil {
		return err
	}

	o.TotalAvailBytes.Set(totalAvailKB * 1e3)

	averageUtil, err := osdDf.Summary.AverageUtil.Float64()
	if err != nil {
		return err
	}

	o.AverageUtil.Set(averageUtil)

	return nil

}

func (o *OsdCollector) collectOsdPerf() error {
	osdPerfCmd := o.cephOSDPerfCommand()
	buf, _, err := o.conn.MonCommand(osdPerfCmd)
	if err != nil {
		log.Println("[ERROR] Unable to collect data from ceph osd perf", err)
		return err
	}

	osdPerf := &cephPerfStat{}
	if err := json.Unmarshal(buf, osdPerf); err != nil {
		return err
	}

	for _, perfStat := range osdPerf.PerfInfo {
		osdID, err := perfStat.ID.Int64()
		if err != nil {
			return err
		}
		osdName := fmt.Sprintf("osd.%v", osdID)

		commitLatency, err := perfStat.Stats.CommitLatency.Float64()
		if err != nil {
			return err
		}
		o.CommitLatency.WithLabelValues(osdName).Set(commitLatency / 1e3)

		applyLatency, err := perfStat.Stats.ApplyLatency.Float64()
		if err != nil {
			return err
		}
		o.ApplyLatency.WithLabelValues(osdName).Set(applyLatency / 1e3)
	}

	return nil
}

func (o *OsdCollector) collectOsdDump() error {
	osdDumpCmd := o.cephOsdDump()
	buff, _, err := o.conn.MonCommand(osdDumpCmd)
	if err != nil {
		log.Println("[ERROR] Unable to collect data from ceph osd dump", err)
		return err
	}

	osdDump := &cephOsdDump{}
	if err := json.Unmarshal(buff, osdDump); err != nil {
		return err
	}

	for _, dumpInfo := range osdDump.Osds {
		osdID, err := dumpInfo.Osd.Int64()
		if err != nil {
			return err
		}
		osdName := fmt.Sprintf("osd.%v", osdID)

		in, err := dumpInfo.In.Float64()
		if err != nil {
			return err
		}

		o.OsdIn.WithLabelValues(osdName).Set(in)

		up, err := dumpInfo.Up.Float64()
		if err != nil {
			return err
		}

		o.OsdUp.WithLabelValues(osdName).Set(up)
	}

	return nil

}

func (o *OsdCollector) cephOsdDump() []byte {
	cmd, err := json.Marshal(map[string]interface{}{
		"prefix": "osd dump",
		"format": "json",
	})
	if err != nil {
		panic(err)
	}
	return cmd
}

func (o *OsdCollector) cephOSDDfCommand() []byte {
	cmd, err := json.Marshal(map[string]interface{}{
		"prefix": "osd df",
		"format": "json",
	})
	if err != nil {
		panic(err)
	}
	return cmd
}

func (o *OsdCollector) cephOSDPerfCommand() []byte {
	cmd, err := json.Marshal(map[string]interface{}{
		"prefix": "osd perf",
		"format": "json",
	})
	if err != nil {
		panic(err)
	}
	return cmd
}

// Describe sends the descriptors of each OsdCollector related metrics we have defined
// to the provided prometheus channel.
func (o *OsdCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, metric := range o.collectorList() {
		metric.Describe(ch)
	}

}

// Collect sends all the collected metrics to the provided prometheus channel.
// It requires the caller to handle synchronization.
func (o *OsdCollector) Collect(ch chan<- prometheus.Metric) {

	if err := o.collectOsdPerf(); err != nil {
		log.Println("failed collecting cluster osd perf stats:", err)
	}

	if err := o.collectOsdDump(); err != nil {
		log.Println("failed collecting cluster osd dump", err)
	}

	if err := o.collect(); err != nil {
		log.Println("failed collecting osd metrics:", err)
	}

	for _, metric := range o.collectorList() {
		metric.Collect(ch)
	}

}
