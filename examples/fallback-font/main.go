package main

import (
	"fmt"
	"image/color"

	"image"

	findfont "github.com/flopp/go-findfont"
	oak "github.com/oakmound/oak/v3"
	"github.com/oakmound/oak/v3/render"
	"github.com/oakmound/oak/v3/scene"
)

func main() {
	oak.AddScene("demo",
		scene.Scene{Start: func(*scene.Context) {

			const fontHeight = 16

			fg := render.DefFontGenerator
			fg.Color = image.NewUniform(color.RGBA{255, 0, 0, 255})
			fg.FontOptions.Size = fontHeight
			font, _ := fg.Generate()

			fallbackFonts := []string{
				"Arial.ttf",
				"Yumin.ttf",
				// TODO: support multi-color glyphs
				"Seguiemj.ttf",
			}

			for _, fontname := range fallbackFonts {
				fontPath, err := findfont.Find(fontname)
				if err != nil {
					fmt.Println("Do you have ", fontPath, "installed?")
					continue
				}
				fg := render.FontGenerator{
					File:  fontPath,
					Color: image.NewUniform(color.RGBA{255, 0, 0, 255}),
					FontOptions: render.FontOptions{
						Size: fontHeight,
					},
				}
				fallbackFont, err := fg.Generate()
				if err != nil {
					panic(err)
				}
				font.Fallbacks = append(font.Fallbacks, fallbackFont)
			}

			strs := []string{
				"Latin-lower: abcdefghijklmnopqrstuvwxyz",
				"Latin-upper: ABCDEFGHIJKLMNOPQRSTUVWXYZ",
				"Greek-lower: αβγδεζηθικλμνχοπρσςτυφψω",
				"Greek-upper: ΑΒΓΔΕΖΗΘΙΚΛΜΝΧΟΠΡΣΤΥΦΨΩ",
				"Japanese-kana: あいえおうかきけこくはひへほふさしせそすまみめもむ",
				"Kanji: 茂僕私華花日本英雄の時",
				"Emoji: 😀😃😄😁😆😅😂🤣🐶🐱🐭🐹🐰🦊🐻🐼",
			}

			y := 0.0
			for _, str := range strs {
				render.Draw(font.NewText(str, 10, y), 0)
				y += fontHeight
			}
		},
		})
	render.SetDrawStack(
		render.NewCompositeR(),
	)
	oak.Init("demo")
}
