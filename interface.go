package eventemitter

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

type Observable interface {
	AddListener(event string, handler HandleFunc) (listener *Listener)
	ListenOnce(event string, handler HandleFunc) (listener *Listener)
	AddCapturer(handler CaptureFunc) (capturer *Capturer)
	CaptureOnce(handler CaptureFunc) (capturer *Capturer)
	RemoveListener(event string, listener *Listener)
	RemoveCapturer(capturer *Capturer)
}

type EventEmitter interface {
	EmitEvent(event string, arguments ...interface{})
}