package eventemitter

type EventType string

type HandleFunc func(arguments ...interface{})
type Listener struct {
	handler HandleFunc
	once    bool
}

type CaptureFunc func(eventName EventType, arguments ...interface{})
type Capturer struct {
	handler CaptureFunc
	once    bool
}

// Observable describes an object that can be listened to by event listeners and capturers
type Observable interface {
	// AddListener adds a listener for the given event type
	AddListener(event EventType, handler HandleFunc) (listener *Listener)
	// ListenOnce adds a listener for the given event type that removes itself after it has been fired once
	ListenOnce(event EventType, handler HandleFunc) (listener *Listener)
	// AddCapturer adds an event capturer for all events
	AddCapturer(handler CaptureFunc) (capturer *Capturer)
	// CaptureOnce adds an event capturer for all events that removes itself after it has been fired once
	CaptureOnce(handler CaptureFunc) (capturer *Capturer)
	// RemoveListener removes the registered given listener for the given event
	RemoveListener(event EventType, listener *Listener)
	// // RemoveCapturer removes the given capturer
	RemoveCapturer(capturer *Capturer)
}

type EventEmitter interface {
	// EmitEvent emits the given event to all listeners and capturers
	EmitEvent(event EventType, arguments ...interface{})
}