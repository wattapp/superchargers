package metrics

import (
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"

	influx "github.com/influxdata/influxdb/client/v2"
)

var (
	client   influx.Client
	Database string
)

type MetricWriter interface {
	Write(influx.BatchPoints) error
}

func Connect() error {
	addr := os.Getenv("INFLUXDB_ADDR")
	if len(addr) == 0 {
		client = NullWriter{}
		return nil
	}
	Database = os.Getenv("INFLUXDB_DATABASE")
	var err error
	client, err = influx.NewHTTPClient(influx.HTTPConfig{
		Addr:     addr,
		Username: os.Getenv("INFLUXDB_USERNAME"),
		Password: os.Getenv("INFLUXDB_PASSWORD"),
	})

	return err
}

func RecordHTTPResponse(process func() (*http.Response, error)) (*http.Response, error) {
	start := time.Now()
	resp, err := process()

	tags := map[string]string{
		"path":   resp.Request.URL.Path,
		"host":   resp.Request.URL.Host,
		"method": resp.Request.Method,
		"status": strconv.Itoa(resp.StatusCode),
	}

	fields := map[string]interface{}{
		"length": resp.ContentLength,
		"took":   time.Since(start).Seconds(),
		"url":    resp.Request.URL.String(),
	}

	Write("http_request", tags, fields)

	return resp, err
}

func Decr(key string) {
	tags := map[string]string{
		"key": key,
	}
	fields := map[string]interface{}{
		"value": -1,
	}
	Write("event", tags, fields)
}

func Incr(key string) {
	tags := map[string]string{
		"key": key,
	}
	fields := map[string]interface{}{
		"value": 1,
	}
	Write("event", tags, fields)
}

func Write(name string, tags map[string]string, fields map[string]interface{}) error {
	bp, err := influx.NewBatchPoints(influx.BatchPointsConfig{
		Database: Database,
	})
	if err != nil {
		return err
	}

	pt, err := influx.NewPoint(name, tags, fields, time.Now())
	if err != nil {
		return err
	}

	bp.AddPoint(pt)
	err = client.Write(bp)
	if err != nil {
		return err
	}

	return nil
}

func Stats(interval time.Duration) {
	var lastPauseNs uint64 = 0
	memStats := &runtime.MemStats{}
	hostname, _ := os.Hostname()

	fmt.Printf("Memory stats goroutine started at %+v interval\n", interval)

	for {
		runtime.ReadMemStats(memStats)

		tags := map[string]string{"hostname": hostname}
		fields := map[string]interface{}{
			"pid":                        os.Getpid(),
			"proc.goroutines":            int64(runtime.NumGoroutine()),
			"proc.memory.allocated":      int64(memStats.Alloc),
			"proc.memory.mallocs":        int64(memStats.Mallocs),
			"proc.memory.frees":          int64(memStats.Frees),
			"proc.memory.gc.total_pause": int64(time.Duration(memStats.PauseTotalNs) / time.Millisecond),
			"proc.memory.heap":           int64(memStats.HeapAlloc),
			"proc.memory.stack":          int64(memStats.StackInuse),
		}

		if lastPauseNs > 0 {
			pauseSinceLastSample := int64(memStats.PauseTotalNs - lastPauseNs)
			fields["proc.memory.gc.pause_per_second"] = pauseSinceLastSample / int64(time.Millisecond) / int64(interval.Seconds())
		}

		lastPauseNs = memStats.PauseTotalNs
		err := Write("runtime.mem_stats", tags, fields)
		if err != nil {
			fmt.Println(err)
		}

		time.Sleep(interval)
	}
}

// The writer used for testing
type NullWriter struct{}

func (nw NullWriter) Write(bp influx.BatchPoints) error {
	fmt.Println(bp)
	return nil
}

// Ping checks that status of cluster, and will always return 0 time and no
// error for UDP clients
func (nw NullWriter) Ping(timeout time.Duration) (time.Duration, string, error) {
	return time.Duration(0), "", nil
}

func (nw NullWriter) Query(q influx.Query) (*influx.Response, error) {
	fmt.Println(q)
	return nil, nil
}

func (nw NullWriter) Close() error {
	return nil
}
