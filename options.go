package qmg

import (
	"time"

	"github.com/zeromicro/go-zero/core/syncx"
)

const (
	defaultTimeout       = time.Second
	defaultSlowThreshold = time.Millisecond * 500
)

var slowThreshold = syncx.ForAtomicDuration(defaultSlowThreshold)

type (
	opt struct {
		timeout time.Duration
	}

	// Option defines the method to customize a mongo model.
	Option func(opts *opt)
)

// SetSlowThreshold sets the slow threshold.
func SetSlowThreshold(threshold time.Duration) {
	slowThreshold.Set(threshold)
}

func defaultOptions() *opt {
	return &opt{
		timeout: defaultTimeout,
	}
}
