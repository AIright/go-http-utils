package listener

type Metrics interface {
	Increment(key string)
}

type stubMetrics struct{}

func (sm stubMetrics) Increment(_ string) {}
