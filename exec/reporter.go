package exec

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/x-mod/routine"
)

const maxNum = 1000000

type ReportData struct {
	AvgTotal float64
	Fastest  float64
	Slowest  float64
	Average  float64
	Ops      float64

	AvgDelay float64
	DelayMax float64
	DelayMin float64

	Lats    []float64
	Offsets []float64
	Codes   []int

	Uptime time.Duration
	Total  time.Duration

	ErrorDist map[string]int
	CodeDist  map[int]int
	SizeTotal int64
	NumRes    int64

	LatencyDistribution []LatencyDistribution
	Histogram           []Bucket
}

type LatencyDistribution struct {
	Percentage int
	Latency    float64
}

type Bucket struct {
	Mark      float64
	Count     int
	Frequency float64
}

const (
	barChar = "■"
)

type Reporter struct {
	upTime  time.Time
	results chan *routine.Result
	//stats
	avgTotal float64
	fastest  float64
	slowest  float64
	average  float64
	ops      float64

	offsets []float64
	codes   []int

	done  chan bool
	total time.Duration

	errorDist map[string]int
	lats      []float64
	sizeTotal int64
	reqTotal  int64
	respTotal int64
	numRes    int64
}

func NewReporter(n int, res chan *routine.Result) routine.Executor {
	cap := min(n, maxNum)
	return &Reporter{
		upTime:    time.Now(),
		results:   res,
		errorDist: make(map[string]int),
		lats:      make([]float64, 0, cap),
		codes:     make([]int, 0, cap),
	}
}

func (r *Reporter) Execute(ctx context.Context) error {
	close(r.results)
	// Loop will continue until channel is closed
	for res := range r.results {
		r.numRes++
		if res.Err != nil {
			r.errorDist[res.Err.Error()]++
			r.codes = append(r.codes, res.Code)
		} else {
			r.avgTotal += res.Duration.Seconds()
			if len(r.lats) < maxNum {
				r.lats = append(r.lats, res.Duration.Seconds())
				r.offsets = append(r.offsets, res.Begin.Sub(r.upTime).Seconds())
			}
			if res.ContentLength > 0 {
				r.sizeTotal += int64(res.ContentLength)
			}
		}
	}
	r.Finalize()
	return nil
}

func min(a, b int) int {
	if a < b {
		if a != 0 {
			return a
		}
	}
	return b
}

func (r *Reporter) Finalize() {
	r.total = time.Since(r.upTime)
	r.ops = float64(r.numRes) / r.total.Seconds()
	r.average = r.avgTotal / float64(len(r.lats))
	r.print()
}

func (r *Reporter) print() {
	buf := &bytes.Buffer{}
	if err := newTemplate().Execute(buf, r.snapshot()); err != nil {
		log.Println("error:", err.Error())
		return
	}
	r.printf(buf.String())

	r.printf("\n")
}

func (r *Reporter) printf(s string, v ...interface{}) {
	fmt.Printf(s, v...)
}

func (r *Reporter) snapshot() ReportData {
	snapshot := ReportData{
		AvgTotal:  r.avgTotal,
		Average:   r.average,
		Ops:       r.ops,
		SizeTotal: r.sizeTotal,
		Total:     r.total,
		Uptime:    time.Since(r.upTime),
		ErrorDist: r.errorDist,
		NumRes:    r.numRes,
		Lats:      make([]float64, len(r.lats)),
		Offsets:   make([]float64, len(r.lats)),
		Codes:     make([]int, len(r.lats)),
	}

	if len(r.lats) == 0 {
		return snapshot
	}

	copy(snapshot.Lats, r.lats)
	copy(snapshot.Codes, r.codes)
	copy(snapshot.Offsets, r.offsets)

	sort.Float64s(r.lats)
	r.fastest = r.lats[0]
	r.slowest = r.lats[len(r.lats)-1]

	snapshot.Histogram = r.histogram()
	snapshot.LatencyDistribution = r.latencies()

	snapshot.Fastest = r.fastest
	snapshot.Slowest = r.slowest

	CodeDist := make(map[int]int, len(snapshot.Codes))
	for _, code := range snapshot.Codes {
		CodeDist[code]++
	}
	snapshot.CodeDist = CodeDist

	return snapshot
}

func (r *Reporter) latencies() []LatencyDistribution {
	pctls := []int{10, 25, 50, 75, 90, 95, 99}
	data := make([]float64, len(pctls))
	j := 0
	for i := 0; i < len(r.lats) && j < len(pctls); i++ {
		current := i * 100 / len(r.lats)
		if current >= pctls[j] {
			data[j] = r.lats[i]
			j++
		}
	}
	res := make([]LatencyDistribution, len(pctls))
	for i := 0; i < len(pctls); i++ {
		if data[i] > 0 {
			res[i] = LatencyDistribution{Percentage: pctls[i], Latency: data[i]}
		}
	}
	return res
}

func (r *Reporter) histogram() []Bucket {
	bc := 10
	buckets := make([]float64, bc+1)
	counts := make([]int, bc+1)
	bs := (r.slowest - r.fastest) / float64(bc)
	for i := 0; i < bc; i++ {
		buckets[i] = r.fastest + bs*float64(i)
	}
	buckets[bc] = r.slowest
	var bi int
	var max int
	for i := 0; i < len(r.lats); {
		if r.lats[i] <= buckets[bi] {
			i++
			counts[bi]++
			if max < counts[bi] {
				max = counts[bi]
			}
		} else if bi < len(buckets)-1 {
			bi++
		}
	}
	res := make([]Bucket, len(buckets))
	for i := 0; i < len(buckets); i++ {
		res[i] = Bucket{
			Mark:      buckets[i],
			Count:     counts[i],
			Frequency: float64(counts[i]) / float64(len(r.lats)),
		}
	}
	return res
}
