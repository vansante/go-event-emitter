package eventemitter

import (
	"sync/atomic"
	"testing"
	"time"
)

func TestEmitterSynced(t *testing.T) {
	testEmitter(t, false)
}

func TestEmitterAsync(t *testing.T) {
	testEmitter(t, true)
}

func testEmitter(t *testing.T, async bool) {
	var em EventEmitter
	var ob Observable

	e := NewEmitter(async)
	em = e
	ob = e

	var ASingle, AListener, capture int32

	listener := ob.AddListener("test event A", func(args ...interface{}) {
		verifyArgs(t, args)
		atomic.AddInt32(&AListener, 1)
	})

	ob.ListenOnce("test event A", func(args ...interface{}) {
		verifyArgs(t, args)
		atomic.AddInt32(&ASingle, 1)
	})

	capturer := ob.AddCapturer(func(event EventType, args ...interface{}) {
		verifyArgs(t, args)
		atomic.AddInt32(&capture, 1)
	})

	em.EmitEvent("test event A", "test", 123, true)
	em.EmitEvent("test event B", "test", 123, true)
	em.EmitEvent("test event C", "test", 123, true)
	em.EmitEvent("test event A", "test", 123, true)
	em.EmitEvent("test event A", "test", 123, true)

	ob.RemoveListener("test event A", listener)
	ob.RemoveCapturer(capturer)

	em.EmitEvent("Testing 123", 1)
	em.EmitEvent("test event A", 1)
	em.EmitEvent("Wow", 2)

	if async {
		// Events are async, so wait a bit for them to finish
		time.Sleep(time.Millisecond * 200)
	}

	if atomic.LoadInt32(&ASingle) != 1 {
		t.Log("Single A event not triggered right", atomic.LoadInt32(&ASingle))
		t.Fail()
	}
	if atomic.LoadInt32(&AListener) != 3 {
		t.Log("A event not triggered right", atomic.LoadInt32(&AListener))
		t.Fail()
	}
	if atomic.LoadInt32(&capture) != 5 {
		t.Log("Capture all not triggered right", atomic.LoadInt32(&capture))
		t.Fail()
	}
}

func verifyArgs(t *testing.T, args []interface{}) {
	if len(args) != 3 {
		t.Logf("Too few arguments (%d) %#v", len(args), args)
		t.Fail()
		return
	}

	s, ok := args[0].(string)
	if !ok || s != "test" {
		t.Log("Wrong argument for 1:test!")
		t.Fail()
	}

	i, ok := args[1].(int)
	if !ok || i != 123 {
		t.Log("Wrong argument for 2:123!")
		t.Fail()
	}

	b, ok := args[2].(bool)
	if !ok || b != true {
		t.Log("Wrong argument for 3:true!")
		t.Fail()
	}
}

func TestEmitNonAsyncRecursive(t *testing.T) {
	e := NewEmitter(false)

	var rootFired int
	e.AddListener("rootevent", func(args ...interface{}) {
		rootFired++
		e.EmitEvent("subevent", 1, 2, 3)
		e.EmitEvent("subevent", 1, 2, 3)
	})

	var subFired int
	e.AddListener("subevent", func(args ...interface{}) {
		if len(args) != 3 {
			t.Logf("Too few arguments (%d) %#v", len(args), args)
			t.Fail()
			return
		}
		subFired++
	})

	e.EmitEvent("rootevent", "test")

	if rootFired != 1 {
		t.Log("Root event all not triggered right", rootFired)
		t.Fail()
	}

	if subFired != 2 {
		t.Log("Sub event all not triggered right", subFired)
		t.Fail()
	}
}
