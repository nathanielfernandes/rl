package ratelimit

import "time"

type RatelimitManager struct {
	Rate  int16
	Per   int64
	c     uint8
	rlMap map[string]*Ratelimit
}

func NewRatelimitManager(r int16, p int64) *RatelimitManager {
	return &RatelimitManager{Rate: r, Per: p, c: 0, rlMap: make(map[string]*Ratelimit)}
}

func (rm *RatelimitManager) cleanseMapping() {
	ct := time.Now().UnixMilli()
	stale := []*string{}
	for id, rl := range rm.rlMap {
		if rl.isStale(ct) {
			stale = append(stale, &id)
		}
	}
	for _, id := range stale {
		delete(rm.rlMap, *id)
	}
}

func (rm *RatelimitManager) intervalCleanse() {
	if rm.c == 255 {
		rm.cleanseMapping()
		rm.c = 0
	}
	rm.c++
}

func (rm *RatelimitManager) IsRatelimited(id string) bool {
	defer rm.intervalCleanse()

	rl, ok := rm.rlMap[id]
	if ok {
		return rl.updateCalls()
	}

	newRl := newRatelimit(rm.Rate, rm.Per)
	rm.rlMap[id] = &newRl

	return newRl.updateCalls()
}

// TimeUntilReset returns the duration until the rate limit will reset for the given ID.
// If the ID is not rate limited or doesn't exist, it returns 0.
func (rm *RatelimitManager) TimeUntilReset(id string) time.Duration {
	rl, ok := rm.rlMap[id]
	if !ok {
		return 0
	}

	currentTime := time.Now().UnixMilli()
	if rl.isStale(currentTime) {
		return 0
	}

	resetTime := rl.tf + rl.per
	remainingTime := resetTime - currentTime
	if remainingTime <= 0 {
		return 0
	}

	return time.Duration(remainingTime) * time.Millisecond
}

// RemainingCalls returns the number of calls remaining for the given ID.
// If the ID doesn't exist, it returns the maximum rate.
func (rm *RatelimitManager) RemainingCalls(id string) int16 {
	rl, ok := rm.rlMap[id]
	if !ok {
		return rm.Rate
	}

	return rl.remainingCalls(time.Now().UnixMilli())
}