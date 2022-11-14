package fps

import "time"

func New() *FPS {
	ticker := time.NewTicker(time.Second)
	fps := &FPS{
		ticker: ticker,
		count:  0,
		quit:   make(chan struct{}),
	}
	go func() {
		for {
			select {
			case <-ticker.C:
				fps.fps = fps.count
				fps.count = 0
			case <-fps.quit:
				return
			}
		}
	}()
	return fps
}

type FPS struct {
	fps    uint
	ticker *time.Ticker
	count  uint
	quit   chan struct{}
}

func (f *FPS) Tick() {
	f.count += 1
}

func (f *FPS) Close() {
	close(f.quit)
	f.ticker.Stop()
}

func (f *FPS) FPS() uint {
	return f.fps
}
