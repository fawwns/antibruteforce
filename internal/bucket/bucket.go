package bucket

import "time"

type Bucket struct {
	Count      int
	Limit      int
	LastRefill time.Time
}

// Allow проверяет - можно ли выполнить ещё 1 операцию.
// Если можно, уменьшает токены на 1 и возвращает true.
// Также обновляет LastRefill.
func (b *Bucket) Allow() bool {
	if time.Since(b.LastRefill) >= time.Minute {
		b.Count = 0
		b.LastRefill = time.Now()
		return true
	}

	if b.Count < b.Limit {
		b.Count++
		return true
	}
	return false
}

// Reset полностью сбрасывает бакет.
// Вызывается после успешной авторизации, чтобы очистить счётчик попыток для логина/пароля.
func (b *Bucket) Reset() {
	b.Count = 0
	b.LastRefill = time.Now()
}

// LastSeen возвращает время последнего обращения (чтобы Store мог удалить неактивные).
func (b *Bucket) LastSeen() time.Time {
	return b.LastRefill
}
