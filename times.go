package retry

type timesStrategy struct {
	maxtimes int
}

func (strategy *timesStrategy) Next(times int) bool {
	return strategy.maxtimes > times
}

type infiniteTimesStrategy struct {
}

func (strategy *infiniteTimesStrategy) Next(times int) bool {
	return true
}

// WithTimes .
func WithTimes(times int) Option {
	return func(action *actionImpl) {
		action.times = &timesStrategy{
			maxtimes: times,
		}
	}
}

// Infinite infinit retry
func Infinite() Option {
	return func(action *actionImpl) {
		action.times = &infiniteTimesStrategy{}
	}
}
