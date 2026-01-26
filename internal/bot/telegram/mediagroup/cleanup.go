package telegram

import (
	"context"
	"time"
)

func (mgm *MediaGroupManager) Cleanup(maxAge time.Duration) {
	mgm.mutex.Lock()
	defer mgm.mutex.Unlock()

	now := time.Now()
	for id, group := range mgm.groups {
		if now.Sub(group.Created) > maxAge {
			if group.Timer != nil {
				group.Timer.Stop()
			}
			delete(mgm.groups, id)
			go mgm.onComplete(group)
		}
	}
}

// starts a background cleanup routine
func (mgm *MediaGroupManager) StartCleanupRoutine(ctx context.Context, cleanupInterval, maxAge time.Duration) {
	ticker := time.NewTicker(cleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			mgm.Cleanup(maxAge)
		}
	}
}
