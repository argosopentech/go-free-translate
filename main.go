package main

import (
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/snakesel/libretranslate"
)

func main() {
	go func() {
		window := new(app.Window)
		err := run(window)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func run(window *app.Window) error {
	// Use the public LibreTranslate instance.
	// For a production app, you would host your own instance.
	translator := libretranslate.New(libretranslate.Config{
		Url: "https://libretranslate.com",
	})

	theme := material.NewTheme()

	var (
		ops          op.Ops
		inputEditor  widget.Editor
		outputEditor widget.Editor
		translateBtn widget.Clickable
		status       string
	)

	inputEditor.SetText("Hello world")
	outputEditor.ReadOnly = true

	for {
		switch e := window.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			// Handle button clicks.
			if translateBtn.Clicked(gtx) {
				status = "Translating..."
				// Run translation in a separate goroutine to avoid blocking the UI.
				go func() {
					inputText := inputEditor.Text()

					// Translate from English to Spanish.
					translatedText, err := translator.Translate(inputText, "en", "es")
					if err != nil {
						log.Printf("Translation error: %v", err)
						status = "Error: " + err.Error()
						window.Invalidate() // Request a redraw to show the error.
						return
					}

					status = "" // Clear status on success.
					outputEditor.SetText(translatedText)
					window.Invalidate() // Request a new frame to show the result.
				}()
			}

			// Define the layout.
			layout.UniformInset(16).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{
					Axis:    layout.Vertical,
					Spacing: layout.SpaceBetween,
				}.Layout(gtx,
					layout.Flexed(0.5, material.Editor(theme, &inputEditor, "Text to translate (English)").Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Inset{Top: 8, Bottom: 8}.Layout(gtx, material.Button(theme, &translateBtn, "Translate to Spanish").Layout)
					}),
					layout.Rigid(material.Body1(theme, status).Layout),
					layout.Flexed(0.5, material.Editor(theme, &outputEditor, "Translation").Layout),
				)
			})

			e.Frame(gtx.Ops)
		}
	}
}
