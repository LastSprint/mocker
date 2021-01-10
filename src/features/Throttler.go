package features

import "time"

/// Throttler provides operation to throttle processing
/// This implementation provides throttling for current goroutine
type Throttler struct{}

/// This func connected to `Throttler` especially for improve testability
func (t Throttler) Throttle(seconds int) {
	time.Sleep(time.Second * time.Duration(seconds))
}
