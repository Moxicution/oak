package event

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	updateCh := make(chan bool)
	assert.Nil(t, UpdateLoop(60, updateCh))
	triggers := 0
	Bind(func(int, interface{}) int {
		triggers++
		return 0
	}, Enter, 0)
	sleep()
	assert.Equal(t, 1, triggers)
	<-updateCh
	sleep()
	assert.Equal(t, 2, triggers)
	assert.Equal(t, 2, FramesElapsed())
	assert.Nil(t, SetTick(1))
	<-updateCh
	assert.Nil(t, Stop())
	sleep()
	sleep()
	select {
	case <-updateCh:
		t.Fatal("Handler should be closed")
	default:
	}
	assert.Nil(t, Update())
	sleep()
	assert.Equal(t, 3, triggers)
	assert.Nil(t, Flush())

	BindPriority(func(int, interface{}) int {
		triggers = 100
		return 0
	}, BindingOption{
		Event: Event{
			Name:     Enter,
			CallerID: 0,
		},
		Priority: 4,
	})

	BindPriority(func(int, interface{}) int {
		if triggers != 100 {
			t.Fatal("Wrong call order")
		}
		return 0
	}, BindingOption{
		Event: Event{
			Name:     Enter,
			CallerID: 0,
		},
		Priority: 3,
	})

	Flush()
	sleep()
	assert.Nil(t, Update())
	sleep()
	sleep()
	Reset()
}

func BenchmarkHandler(b *testing.B) {
	triggers := 0
	entities := 10
	go DefaultBus.ResolvePending()
	for i := 0; i < entities; i++ {
		DefaultBus.GlobalBind(func(int, interface{}) int {
			triggers++
			return 0
		}, Enter)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		<-DefaultBus.TriggerBack(Enter, DefaultBus.framesElapsed)
	}
}

func TestPauseAndResume(t *testing.T) {
	b := NewBus()
	b.ResolvePending()
	triggerCt := 0
	b.Bind(func(int, interface{}) int {
		triggerCt++
		return 0
	}, "EnterFrame", 0)
	ch := make(chan bool, 1000)
	b.UpdateLoop(60, ch)
	time.Sleep(1 * time.Second)
	b.Pause()
	time.Sleep(1 * time.Second)
	oldCt := triggerCt
	time.Sleep(1 * time.Second)
	if oldCt != triggerCt {
		t.Fatalf("pause did not stop enter frame from triggering: expected %v got %v", oldCt, triggerCt)
	}

	b.Resume()
	time.Sleep(1 * time.Second)
	newCt := triggerCt
	if newCt == oldCt {
		t.Fatalf("resume did not resume enter frame triggering: expected %v got %v", oldCt, newCt)
	}
}
