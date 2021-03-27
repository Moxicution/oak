package main

import (
	"fmt"
	"image/color"

	"github.com/oakmound/oak/v2"
	"github.com/oakmound/oak/v2/event"
	"github.com/oakmound/oak/v2/mouse"
	"github.com/oakmound/oak/v2/render"
	"github.com/oakmound/oak/v2/scene"
)

func main() {
	// TODO: SetupConfig needs to be overhauled
	oak.SetupConfig.Debug.Level = "VERBOSE"
	oak.SetupConfig.DrawFrameRate = 1200
	oak.SetupConfig.FrameRate = 60
	c1 := oak.NewController()
	c1.InitialDrawStack = render.NewDrawStack(render.NewDynamicHeap())

	type GlobalBinder interface {
		GlobalBind(string, event.Bindable)
	}

	// Two windows cannot share the same logic handler
	b1 := event.NewBus()
	c1.SetLogicHandler(b1)
	c1.FirstSceneInput = color.RGBA{255, 0, 0, 255}
	c1.AddScene("scene1", scene.Scene{
		Start: func(ctx *scene.Context) {
			fmt.Println("Start scene 1")
			cb := render.NewColorBox(50, 50, ctx.SceneInput.(color.RGBA))
			cb.SetPos(50, 50)
			ctx.DrawStack.Draw(cb, 0)
			dFPS := render.NewDrawFPS(0.1, nil, 600, 10)
			ctx.DrawStack.Draw(dFPS, 1)
			ctx.EventHandler.(GlobalBinder).GlobalBind(mouse.Press, mouse.Binding(func(_ event.CID, me mouse.Event) int {
				cb.SetPos(me.X(), me.Y())
				return 0
			}))
		},
	})
	go c1.Init("scene1")

	c2 := oak.NewController()
	c2.InitialDrawStack = render.NewDrawStack(render.NewDynamicHeap())
	b2 := event.NewBus()
	c2.SetLogicHandler(b2)
	c2.FirstSceneInput = color.RGBA{0, 255, 0, 255}
	c2.AddScene("scene2", scene.Scene{
		Start: func(ctx *scene.Context) {
			fmt.Println("Start scene 2")
			cb := render.NewColorBox(50, 50, ctx.SceneInput.(color.RGBA))
			cb.SetPos(50, 50)
			ctx.DrawStack.Draw(cb, 0)
			dFPS := render.NewDrawFPS(0.1, nil, 600, 10)
			ctx.DrawStack.Draw(dFPS, 1)
			ctx.EventHandler.(GlobalBinder).GlobalBind(mouse.Press, mouse.Binding(func(_ event.CID, me mouse.Event) int {
				cb.SetPos(me.X(), me.Y())
				return 0
			}))
		},
	})
	c2.Init("scene2")

	//oak.Init() => oak.NewController(render.GlobalDrawStack, dlog.DefaultLogger ...).Init()
}
