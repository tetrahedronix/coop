package googol

import (
	"sync"
	"time"

	"github.com/tetrahedronix/coop/googol/ecs"
)

type Loop struct {
	mu      sync.RWMutex
	world   *World
	systems []System
	ticker  *time.Ticker
	stopCh  chan struct{}
	running bool // Loop flag
}

func (l *Loop) Add(sys System) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.systems = append(l.systems, sys)
}

func (l *Loop) Update() {
	// Ottiene le entità (copia)
	entities := l.world.GetEntities()

	// Esegue i sistemi
	for _, sys := range l.systems {
		for _, e := range entities {
			if sys.Match(e) {
				sys.Process(e)
			}
		}
	}

	// Swap dei buffer (double buffering)
	l.world.SwapBuffers()

	// Rimuove le entità contrassegnate
	l.world.PurgeRemoved()
}

func (l *Loop) Remove(sys ecs.System) {
	l.mu.Lock()
	defer l.mu.Unlock()

	for i, s := range l.systems {
		if s == sys {
			l.systems = append(l.systems[:i], l.systems[i+1:]...)
			break
		}
	}
}

// Run avvia il loop a 60 FPS
func (l *Loop) Run() {
	l.mu.Lock()

	if l.running {
		l.mu.Unlock()
		return
	}

	l.running = true
	l.ticker = time.NewTicker(time.Second / 60)
	l.mu.Unlock()

	go func() {
		for {
			select {
			case <-l.ticker.C:
				l.Update()
			case <-l.stopCh:
				l.ticker.Stop()
				return
			}
		}
	}()
}

func NewLoop(w *World) *Loop {
	return &Loop{
		world:   w,
		systems: make([]ecs.System, 0),
		ticker:  nil,
		stopCh:  make(chan struct{}),
	}
}
