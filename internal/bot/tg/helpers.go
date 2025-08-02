package tg

import "time"

func (b *Bot) checkRascladRateLimit(chatID int64) bool {
	now := time.Now()
	if lastRequest, ok := b.lastRascladRequestTimes.Load(chatID); ok {
		lastRequestTime := lastRequest.(time.Time)
		if now.Sub(lastRequestTime) < cooldownDuration {
			return false
		}
	}

	b.lastRascladRequestTimes.Store(chatID, now)
	return true
}
