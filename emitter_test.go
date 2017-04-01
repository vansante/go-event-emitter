package eventemitter

import (
	"testing"
	"time"
)

func TestEmitter(t *testing.T) {
	var em EventEmitter
	var ob Observable

	e := NewEmitter()
	em = e
	ob = e

	var ASingle, AListener, capture, captureOnce int

	listener := ob.AddListener("test event A", func(args ...interface{}) {
		verifyArgs(t, args)
		AListener++
	})

	ob.ListenOnce("test event A", func(args ...interface{}){
		verifyArgs(t, args)
		ASingle++
	})

	capturer := ob.AddCapturer(func(event string, args ...interface{}) {
		verifyArgs(t, args)
		capture++
	})

	ob.CaptureOnce(func(event string, args ...interface{}) {
		verifyArgs(t, args)

		captureOnce++
		if event != "test event A" {
			t.Log("wrong event for captureOnce")
			t.Fail()
		}
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

	// Events are async, so wait a bit for them to finish
	time.Sleep(time.Second)

	if ASingle != 1 {
		t.Log("Single A event not triggered right", ASingle)
		t.Fail()
	}
	if AListener != 3 {
		t.Log("A event not triggered right", AListener)
		t.Fail()
	}
	if capture != 5 {
		t.Log("Capture all not triggered right", capture)
		t.Fail()
	}
	if captureOnce != 1 {
		t.Log("Capture once not triggered right", captureOnce)
		t.Fail()
	}
}

func verifyArgs(t *testing.T, args []interface{}) {
	if len(args) != 3 {
		t.Log("Too few arguments", args)
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