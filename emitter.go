package eventemitter

import "sync"

type Emitter struct {
	capturers     []*Capturer
	listeners     map[string][]*Listener
	listenerMutex sync.Mutex
}

type HandleFunc func(arguments ...interface{})
type Listener struct {
	handler HandleFunc
	once    bool
}

type CaptureFunc func(eventName string, arguments ...interface{})
type Capturer struct {
	handler CaptureFunc
	once    bool
}

func NewEmitter() (em *Emitter) {
	em = &Emitter{
		listeners: make(map[string][]*Listener),
	}
	return
}

func (em *Emitter) EmitEvent(event string, arguments ...interface{}) {
	em.listenerMutex.Lock()
	defer em.listenerMutex.Unlock()

	removed := 0
	var adjustedIndex int
	for index := range em.listeners[event] {
		adjustedIndex = index - removed

		go em.listeners[event][adjustedIndex].handler(arguments...)
		if em.listeners[event][adjustedIndex].once {
			em.removeListenerAtIndex(event, adjustedIndex)
			removed++
		}
	}

	removed = 0
	for index := range em.capturers {
		adjustedIndex = index - removed

		go em.capturers[adjustedIndex].handler(event, arguments...)
		if em.capturers[adjustedIndex].once {
			em.removeCapturerAtIndex(adjustedIndex)
			removed++
		}
	}
}

func (em *Emitter) AddListener(event string, handler HandleFunc) {
	em.listenerMutex.Lock()
	defer em.listenerMutex.Unlock()

	em.listeners[event] = append(em.listeners[event], &Listener{
		handler: handler,
		once:    false,
	})
}

func (em *Emitter) ListenOnce(event string, handler HandleFunc) {
	em.listenerMutex.Lock()
	defer em.listenerMutex.Unlock()

	em.listeners[event] = append(em.listeners[event], &Listener{
		handler: handler,
		once:    true,
	})
}

func (em *Emitter) AddCapturer(handler CaptureFunc) {
	em.listenerMutex.Lock()
	defer em.listenerMutex.Unlock()

	em.capturers = append(em.capturers, &Capturer{
		handler: handler,
		once:    false,
	})
}

func (em *Emitter) CaptureOnce(handler CaptureFunc) {
	em.listenerMutex.Lock()
	defer em.listenerMutex.Unlock()

	em.capturers = append(em.capturers, &Capturer{
		handler: handler,
		once:    true,
	})
}

func (em *Emitter) RemoveListener(event string, handler HandleFunc) {
	em.listenerMutex.Lock()
	defer em.listenerMutex.Unlock()

	for index, listener := range em.listeners[event] {
		if &listener.handler == &handler {
			em.removeListenerAtIndex(event, index)
			return
		}
	}
}

func (em *Emitter) RemoveAllListenersForEvent(event string) {
	em.listenerMutex.Lock()
	defer em.listenerMutex.Unlock()

	em.listeners[event] = make([]*Listener, 0)
}

func (em *Emitter) RemoveAllListeners() {
	em.listenerMutex.Lock()
	defer em.listenerMutex.Unlock()

	em.listeners = make(map[string][]*Listener)
}

func (em *Emitter) RemoveCapturer(handler CaptureFunc) {
	em.listenerMutex.Lock()
	defer em.listenerMutex.Unlock()

	for index, capturer := range em.capturers {
		if &capturer.handler == &handler {
			em.removeCapturerAtIndex(index)
			return
		}
	}
}

func (em *Emitter) RemoveAllCapturers() {
	em.listenerMutex.Lock()
	defer em.listenerMutex.Unlock()

	em.capturers = make([]*Capturer, 0)
}

func (em *Emitter) removeListenerAtIndex(event string, index int) {
	copy(em.listeners[event][index:], em.listeners[event][index+1:])
	em.listeners[event][len(em.listeners[event])-1] = nil
	em.listeners[event] = em.listeners[event][:len(em.listeners[event])-1]
}

func (em *Emitter) removeCapturerAtIndex(index int) {
	copy(em.capturers[index:], em.capturers[index+1:])
	em.capturers[len(em.capturers)-1] = nil
	em.capturers = em.capturers[:len(em.capturers)-1]
}
