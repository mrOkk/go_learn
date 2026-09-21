package limiter

const (
	fixedWindow = iota
	slidingWindow
	tokenWindow
)

type Limiter interface {
	Allow() bool
}

func GetLimiter() Limiter {
	return nil
}
