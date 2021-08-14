package debugtools

import (
	"fmt"

	"github.com/oakmound/oak/v3/dlog"
	"github.com/oakmound/oak/v3/event"
	"github.com/oakmound/oak/v3/mouse"
	"github.com/oakmound/oak/v3/scene"
)

// DebugMouseRelease will print the position and button pressed of the mouse when the mouse is released, if the given
// key is held down at the time. If no key is given, it will always be printed
func DebugMouseRelease(ctx *scene.Context, k string) {
	event.GlobalBind(mouse.Release, func(_ event.CID, ev interface{}) int {
		mev, ok := ev.(*mouse.Event)
		if !ok {
			dlog.Error("got wrong event", fmt.Sprintf("%T", mev))
			return 0
		}
		if ctx.KeyState.IsDown(k) || k == "" {
			dlog.Info(mev)
		}
		return 0
	})
}
