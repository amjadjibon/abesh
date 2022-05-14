package iface

// Counter is metrics counter.
type Counter interface {
	With(levels ...string) Counter
	Inc()
	Add(delta float64)
}

// Gauge is metrics gauge.
type Gauge interface {
	With(levels ...string) Gauge
	Set(value float64)
	Add(delta float64)
	Sub(delta float64)
}

// Observer is metrics observer.
type Observer interface {
	With(levels ...string) Observer
	Observe(float64)
}
