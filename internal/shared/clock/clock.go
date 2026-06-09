// Package clock makes "now" injectable so domain code that depends on time
// (e.g. Incident duration — RN-012) stays testable.
package clock

import "time"

type Clock func() time.Time

func System() Clock {
	return time.Now
}
