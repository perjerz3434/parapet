package trace

import (
	"log"
	"net/http"
	"time"

	"contrib.go.opencensus.io/exporter/stackdriver"
	sdpropagation "go.opencensus.io/exporter/stackdriver/propagation"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/trace"
	"go.opencensus.io/trace/propagation"
)

// New creates new trace middleware
func New() *Trace {
	return &Trace{}
}

// Trace middleware
type Trace struct {
	ProjectID            string
	Propagation          propagation.HTTPFormat
	BundleCountThreshold int
	BundleDelayThreshold time.Duration
	IsPublicEndpoint     bool
	FormatSpanName       func(r *http.Request) string
	StartOptions         trace.StartOptions
}

// ServeHandler implements middleware interface
func (m *Trace) ServeHandler(h http.Handler) http.Handler {
	if m.Propagation == nil {
		m.Propagation = &sdpropagation.HTTPFormat{}
	}
	if m.FormatSpanName == nil {
		m.FormatSpanName = func(r *http.Request) string {
			proto := r.Header.Get("X-Forwarded-Proto")
			return proto + "://" + r.Host + r.RequestURI
		}
	}
	if m.StartOptions.Sampler == nil {
		m.StartOptions.Sampler = trace.AlwaysSample()
	}
	if m.StartOptions.SpanKind == trace.SpanKindUnspecified {
		m.StartOptions.SpanKind = trace.SpanKindServer
	}

	exporter, err := stackdriver.NewExporter(stackdriver.Options{
		ProjectID:            m.ProjectID,
		BundleCountThreshold: m.BundleCountThreshold,
		BundleDelayThreshold: m.BundleDelayThreshold,
	})
	if err != nil {
		log.Println("stackdriver/trace:", err)
		return h
	}

	trace.RegisterExporter(exporter)

	return &ochttp.Handler{
		Handler:          h,
		Propagation:      m.Propagation,
		FormatSpanName:   m.FormatSpanName,
		StartOptions:     m.StartOptions,
		IsPublicEndpoint: m.IsPublicEndpoint,
	}
}
