package golang

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

// WriteJSON writes JSON response struct to ResponseWriter.
func WriteJSON(w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json, err := json.Marshal(response)
	if err != nil {
		return err
	}

	if _, err := w.Write(json); err != nil {
		return err
	}

	return nil
}

// WriteErr writes error to ResponseWriter.
func WriteErr(w http.ResponseWriter, err error, httpCode int) {
	logrus.Error(err.Error())

	// write error to response
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	var errMap = map[string]interface{}{
		"httpStatus": httpCode,
		"error":      err.Error(),
	}

	errJSON, err := json.Marshal(errMap)
	if err != nil {
		logrus.Errorf("can't marshal error response. err: %s", err)
	}

	logrus.Error(string(errJSON))
	w.WriteHeader(httpCode)
	w.Write(errJSON)
}

// callStats holds various stats associated with HTTP request-response pair.
type callStats struct {
	w       http.ResponseWriter
	code    int // response status code
	resSize int64
}

func newCallStats(w http.ResponseWriter) *callStats {
	return &callStats{w: w}
}

func (r *callStats) Header() http.Header {
	return r.w.Header()
}

func (r *callStats) WriteHeader(code int) {
	r.w.WriteHeader(code)
	r.code = code
}

func (r *callStats) Write(p []byte) (n int, err error) {
	if r.code == 0 {
		r.code = http.StatusOK
	}
	n, err = r.w.Write(p)
	r.resSize += int64(n)
	return
}

func (r *callStats) StatusCode() int {
	return r.code
}

func (r *callStats) ResponseSize() int64 {
	return r.resSize
}

// NewLogginHandler creates and returns LoggingHandler with custom metrics.
func NewLogginHandler(root http.Handler) http.Handler {
	duration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "monitoring",
			Subsystem: "rest",
			Name:      "http_durations_histogram_seconds",
			Help:      "Request time duration.",
		},
		[]string{"code", "method", "endpoint"},
	)

	requests := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "monitoring",
			Subsystem: "rest",
			Name:      "http_requests_total",
			Help:      "Total number of requests received.",
		},
		[]string{"code", "method", "endpoint"},
	)

	lh := LoggingHandler{root: root, duration: duration, requests: requests}

	prometheus.DefaultRegisterer.Register(lh)

	return lh
}

// LoggingHandler writes basic information about each request and response to
// the log.
type LoggingHandler struct {
	root     http.Handler
	duration *prometheus.HistogramVec
	requests *prometheus.CounterVec
}

func (h LoggingHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	t := time.Now()
	stats := newCallStats(w)
	h.root.ServeHTTP(stats, req)

	elapsed := time.Since(t)

	buf := new(bytes.Buffer)
	buf.WriteString(getRemoteAddr(req))
	buf.WriteString(" - \"")
	buf.WriteString(req.Method)
	buf.WriteByte(' ')
	buf.WriteString(req.RequestURI)
	buf.WriteByte(' ')
	buf.WriteString(req.Proto)
	buf.WriteString("\" ")
	buf.WriteString(strconv.Itoa(stats.StatusCode()))
	buf.WriteByte(' ')
	buf.WriteString(strconv.FormatInt(stats.ResponseSize(), 10))
	buf.WriteString(" \"")
	buf.WriteString(req.UserAgent())
	buf.WriteString("\" Took: ")
	buf.WriteString(elapsed.String())

	h.requests.WithLabelValues(strconv.Itoa(stats.StatusCode()), strconv.Itoa(stats.StatusCode()), req.RequestURI).Inc()
	h.duration.WithLabelValues(strconv.Itoa(stats.StatusCode()), strconv.Itoa(stats.StatusCode()), req.RequestURI).Observe(elapsed.Seconds())

	logrus.Infof(buf.String())
}

func getRemoteAddr(r *http.Request) string {
	forwaredFor := r.Header.Get("X-Forwarded-For")
	if forwaredFor == "" {
		return r.RemoteAddr
	}

	return forwaredFor
}

// Describe implements prometheus.Collector interface.
func (h LoggingHandler) Describe(in chan<- *prometheus.Desc) {
	h.duration.Describe(in)
	h.requests.Describe(in)
}

// Collect implements prometheus.Collector interface.
func (h LoggingHandler) Collect(in chan<- prometheus.Metric) {
	h.duration.Collect(in)
	h.requests.Collect(in)
}
