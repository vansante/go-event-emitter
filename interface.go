package eventemitter

type EventID string

type HandleFunc func(arguments ...interface{})
type Listener struct {
	handler HandleFunc
	once    bool
}

type CaptureFunc func(eventName EventID, arguments ...interface{})
type Capturer struct {
	handler CaptureFunc
	once    bool
}

type Observable interface {
	AddListener(event EventID, handler HandleFunc) (listener *Listener)
	ListenOnce(event EventID, handler HandleFunc) (listener *Listener)
	AddCapturer(handler CaptureFunc) (capturer *Capturer)
	CaptureOnce(handler CaptureFunc) (capturer *Capturer)
	RemoveListener(event EventID, listener *Listener)
	RemoveCapturer(capturer *Capturer)
}

type EventEmitter interface {
	EmitEvent(event EventID, arguments ...interface{})
}