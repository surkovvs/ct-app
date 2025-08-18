package zorro

import "sync/atomic"

type Zorro struct {
	status *uint64
}

type (
	Status uint64
	Mask   uint64
)

func New() Zorro {
	var status uint64
	return Zorro{
		status: &status,
	}
}

func (z Zorro) GetStatus() Status {
	return Status(atomic.LoadUint64(z.status))
}

// Concurrently safe setup bits
func (c Zorro) SetStatus(status Status, mask Mask) {
	for {
		if c.TrySetStatus(status, mask) {
			return
		}
	}
}

func (c Zorro) TrySetStatus(status Status, mask Mask) bool {
	cur := atomic.LoadUint64(c.status)
	new := Status(cur).SetWithMask(status, mask)
	return atomic.CompareAndSwapUint64(c.status, cur, new)
}

// compare and swap but masked
func (c Zorro) TryChangeStatus(prev, next Status, mask Mask) bool {
	cur := atomic.LoadUint64(c.status)
	if !prev.CompareMasked(Status(cur), mask) {
		return false
	}
	new := Status(cur).SetWithMask(next, mask)
	return atomic.CompareAndSwapUint64(c.status, cur, new)
}

// status 1010 mask 0011 result 0010
func (s Status) Querying(m Mask) uint64 {
	return uint64(s) & uint64(m)
}

// status 1010 mask 0011 result 1011
func (s Status) MaskedOn(m Mask) uint64 {
	return uint64(s) | uint64(m)
}

// status 1010 mask 0011 result 1000
func (s Status) MaskedOff(m Mask) uint64 {
	return uint64(s) &^ uint64(m)
}

// status 1010 mask 0011 set 0101 result 1001
func (s Status) SetWithMask(set Status, m Mask) uint64 {
	return s.MaskedOff(m) | set.Querying(m)
}

func (s Status) CompareMasked(is Status, m Mask) bool {
	return s.Querying(m) == is.Querying(m)
}
