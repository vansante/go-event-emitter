package eventemitter

import (
	"sync"
)

// Emitter is the base struct which manages event subscriptions and calls all registered handlers on event emits.
type Emitter struct {
	mu            sync.RWMutex
	async         bool
	capturers     []*Capturer
	listeners     map[EventType][]*Listener
	listenersOnce map[EventType][]*Listener
}

// NewEmitter creates a new event emitter that implements the Observable interface.
// Async determines whether events listeners fire in separate goroutines or not.
func NewEmitter(async bool) (em *Emitter) {
	em = &Emitter{
		async:         async,
		listeners:     make(map[EventType][]*Listener),
		listenersOnce: make(map[EventType][]*Listener),
	}
	return em
}

// EmitEvent emits the given event to all listeners and capturers
func (em *Emitter) EmitEvent(event EventType, arguments ...interface{}) {
	// If we have no single listeners for this event, skip
	em.mu.RLock()
	if len(em.listenersOnce) > 0 {
		em.mu.RUnlock()
		// Get a full lock, we are changing a map
		em.mu.Lock()
		// Copy the slice
		listenersOnce := em.listenersOnce[event]
		// Create new empty slice
		em.listenersOnce[event] = make([]*Listener, 0)
		em.mu.Unlock()

		// No lock needed, we are working with an inaccessible copy
		em.emitListenerEvents(listenersOnce, arguments)
	} else {
		em.mu.RUnlock()
	}

	em.mu.RLock()
	// If we have no listeners for this event, skip
	if len(em.listeners[event]) > 0 {
		em.emitListenerEvents(em.listeners[event], arguments)
	}
	em.mu.RUnlock()

	em.mu.RLock()
	// If we have no capturers, skip
	if len(em.capturers) > 0 {
		em.emitCapturerEvents(em.capturers, event, arguments)
	}
	em.mu.RUnlock()
}

func (em *Emitter) emitListenerEvents(listeners []*Listener, arguments []interface{}) {
	for _, listener := range listeners {
		if em.async {
			go listener.handler(arguments...)
			continue
		}
		listener.handler(arguments...)
	}
}

func (em *Emitter) emitCapturerEvents(capturers []*Capturer, event EventType, arguments []interface{}) {
	for _, capturer := range capturers {
		if em.async {
			go capturer.handler(event, arguments...)
			continue
		}
		capturer.handler(event, arguments...)
	}
}

// AddListener adds a listener for the given event type
func (em *Emitter) AddListener(event EventType, handler HandleFunc) (listener *Listener) {
	em.mu.Lock()
	defer em.mu.Unlock()

	listener = &Listener{
		handler: handler,
	}
	em.listeners[event] = append(em.listeners[event], listener)
	return listener
}

// ListenOnce adds a listener for the given event type that removes itself after it has been fired once
func (em *Emitter) ListenOnce(event EventType, handler HandleFunc) (listener *Listener) {
	em.mu.Lock()
	defer em.mu.Unlock()

	listener = &Listener{
		handler: handler,
	}
	em.listenersOnce[event] = append(em.listenersOnce[event], listener)
	return listener
}

// AddCapturer adds an event capturer for all events
func (em *Emitter) AddCapturer(handler CaptureFunc) (capturer *Capturer) {
	em.mu.Lock()
	defer em.mu.Unlock()

	capturer = &Capturer{
		handler: handler,
	}
	em.capturers = append(em.capturers, capturer)
	return capturer
}

// RemoveListener removes the registered given listener for the given event
func (em *Emitter) RemoveListener(event EventType, listener *Listener) {
	em.mu.Lock()
	defer em.mu.Unlock()

	for index, list := range em.listeners[event] {
		if list == listener {
			em.removeListenerAt(event, index)
			return
		}
	}

	// If it hasnt been found yet, remove from listeners once if present there
	for index, list := range em.listenersOnce[event] {
		if list == listener {
			em.removeOnceListenerAt(event, index)
			return
		}
	}
}

// RemoveAllListenersForEvent removes all registered listeners for a given event type
func (em *Emitter) RemoveAllListenersForEvent(event EventType) {
	em.mu.Lock()
	defer em.mu.Unlock()

	em.listeners[event] = make([]*Listener, 0)
}

// RemoveAllListeners removes all registered listeners for all event types
func (em *Emitter) RemoveAllListeners() {
	em.mu.Lock()
	defer em.mu.Unlock()

	em.listeners = make(map[EventType][]*Listener)
	em.listenersOnce = make(map[EventType][]*Listener)
}

// RemoveCapturer removes the given capturer
func (em *Emitter) RemoveCapturer(capturer *Capturer) {
	em.mu.Lock()
	defer em.mu.Unlock()

	for index, capt := range em.capturers {
		if capt == capturer {
			em.removeCapturerAt(index)
			return
		}
	}
}

// RemoveAllCapturers removes all registered capturers
func (em *Emitter) RemoveAllCapturers() {
	em.mu.Lock()
	defer em.mu.Unlock()

	em.capturers = make([]*Capturer, 0)
}

func (em *Emitter) removeListenerAt(event EventType, index int) {
	copy(em.listeners[event][index:], em.listeners[event][index+1:])
	em.listeners[event][len(em.listeners[event])-1] = nil
	em.listeners[event] = em.listeners[event][:len(em.listeners[event])-1]
}

func (em *Emitter) removeOnceListenerAt(event EventType, index int) {
	copy(em.listenersOnce[event][index:], em.listenersOnce[event][index+1:])
	em.listenersOnce[event][len(em.listenersOnce[event])-1] = nil
	em.listenersOnce[event] = em.listenersOnce[event][:len(em.listenersOnce[event])-1]
}

func (em *Emitter) removeCapturerAt(index int) {
	copy(em.capturers[index:], em.capturers[index+1:])
	em.capturers[len(em.capturers)-1] = nil
	em.capturers = em.capturers[:len(em.capturers)-1]
}
