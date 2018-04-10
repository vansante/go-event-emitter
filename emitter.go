package eventemitter

import (
	"sync"
)

type Emitter struct {
	async         bool
	capturers     []*Capturer
	listeners     map[EventType][]*Listener
	listenerMutex sync.Mutex
}

// NewEmitter creates a new event emitter that implements the Observable interface.
// Async determines whether events listeners fire in separate goroutines or not.
func NewEmitter(async bool) (em *Emitter) {
	em = &Emitter{
		async:     async,
		listeners: make(map[EventType][]*Listener),
	}
	return em
}

// EmitEvent emits the given event to all listeners and capturers
func (em *Emitter) EmitEvent(event EventType, arguments ...interface{}) {
	em.listenerMutex.Lock()
	defer em.listenerMutex.Unlock()

	removed := 0
	var adjustedIndex int
	for index := range em.listeners[event] {
		adjustedIndex = index - removed

		if em.async {
			go em.listeners[event][adjustedIndex].handler(arguments...)
		} else {
			em.listeners[event][adjustedIndex].handler(arguments...)
		}

		if em.listeners[event][adjustedIndex].once {
			em.removeListenerAtIndex(event, adjustedIndex)
			removed++
		}
	}

	removed = 0
	for index := range em.capturers {
		adjustedIndex = index - removed

		if em.async {
			go em.capturers[adjustedIndex].handler(event, arguments...)
		} else {
			em.capturers[adjustedIndex].handler(event, arguments...)
		}

		if em.capturers[adjustedIndex].once {
			em.removeCapturerAtIndex(adjustedIndex)
			removed++
		}
	}
}

// AddListener adds a listener for the given event type
func (em *Emitter) AddListener(event EventType, handler HandleFunc) (listener *Listener) {
	em.listenerMutex.Lock()
	defer em.listenerMutex.Unlock()

	listener = &Listener{
		handler: handler,
		once:    false,
	}
	em.listeners[event] = append(em.listeners[event], listener)
	return listener
}

// ListenOnce adds a listener for the given event type that removes itself after it has been fired once
func (em *Emitter) ListenOnce(event EventType, handler HandleFunc) (listener *Listener) {
	em.listenerMutex.Lock()
	defer em.listenerMutex.Unlock()

	listener = &Listener{
		handler: handler,
		once:    true,
	}
	em.listeners[event] = append(em.listeners[event], listener)
	return listener
}

// AddCapturer adds an event capturer for all events
func (em *Emitter) AddCapturer(handler CaptureFunc) (capturer *Capturer) {
	em.listenerMutex.Lock()
	defer em.listenerMutex.Unlock()

	capturer = &Capturer{
		handler: handler,
		once:    false,
	}
	em.capturers = append(em.capturers, capturer)
	return capturer
}

// CaptureOnce adds an event capturer for all events that removes itself after it has been fired once
func (em *Emitter) CaptureOnce(handler CaptureFunc) (capturer *Capturer) {
	em.listenerMutex.Lock()
	defer em.listenerMutex.Unlock()

	capturer = &Capturer{
		handler: handler,
		once:    true,
	}
	em.capturers = append(em.capturers, capturer)
	return capturer
}

// RemoveListener removes the registered given listener for the given event
func (em *Emitter) RemoveListener(event EventType, listener *Listener) {
	em.listenerMutex.Lock()
	defer em.listenerMutex.Unlock()

	for index, list := range em.listeners[event] {
		if list == listener {
			em.removeListenerAtIndex(event, index)
			return
		}
	}
}

// RemoveAllListenersForEvent removes all registered listeners for a given event type
func (em *Emitter) RemoveAllListenersForEvent(event EventType) {
	em.listenerMutex.Lock()
	defer em.listenerMutex.Unlock()

	em.listeners[event] = make([]*Listener, 0)
}

// RemoveAllListeners removes all registered listeners for all event types
func (em *Emitter) RemoveAllListeners() {
	em.listenerMutex.Lock()
	defer em.listenerMutex.Unlock()

	em.listeners = make(map[EventType][]*Listener)
}

// RemoveCapturer removes the given capturer
func (em *Emitter) RemoveCapturer(capturer *Capturer) {
	em.listenerMutex.Lock()
	defer em.listenerMutex.Unlock()

	for index, capt := range em.capturers {
		if capt == capturer {
			em.removeCapturerAtIndex(index)
			return
		}
	}
}

// RemoveAllCapturers removes all registered capturers
func (em *Emitter) RemoveAllCapturers() {
	em.listenerMutex.Lock()
	defer em.listenerMutex.Unlock()

	em.capturers = make([]*Capturer, 0)
}

func (em *Emitter) removeListenerAtIndex(event EventType, index int) {
	copy(em.listeners[event][index:], em.listeners[event][index+1:])
	em.listeners[event][len(em.listeners[event])-1] = nil
	em.listeners[event] = em.listeners[event][:len(em.listeners[event])-1]
}

func (em *Emitter) removeCapturerAtIndex(index int) {
	copy(em.capturers[index:], em.capturers[index+1:])
	em.capturers[len(em.capturers)-1] = nil
	em.capturers = em.capturers[:len(em.capturers)-1]
}
