package domain

type TimeProvider interface {
	NowUnixMilli() int64
}
