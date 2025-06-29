package main

import (
	"fmt"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/snakesel/libretranslate"
)

// Language defines the structure for a supported language.
type Language struct {
	Code string
	Name string
}

// supportedLanguages is a hardcoded list of languages for selection.
// TODO call https://libretranslate.com/languages
// https://github.com/SnakeSel/libretranslate/issues/2
var supportedLanguages = []Language{
	{Code: "en", Name: "English"},
	{Code: "es", Name: "Spanish"},
	{Code: "fr", Name: "French"},
	{Code: "de", Name: "German"},
	{Code: "zh", Name: "Chinese"},
	{Code: "ru", Name: "Russian"},
}

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
		Url: "http://localhost:5000",
	})

	theme := material.NewTheme()

	var (
		ops             op.Ops
		inputEditor     widget.Editor
		outputEditor    widget.Editor
		translateBtn    widget.Clickable
		status          string
		sourceLangEnum  widget.Enum
		targetLangEnum  widget.Enum
		sourceLangList  widget.List
		targetLangList  widget.List
	)

	// Set default selections and list orientation.
	sourceLangEnum.Value = "en"
	targetLangEnum.Value = "es"
	sourceLangList.Axis = layout.Horizontal
	targetLangList.Axis = layout.Horizontal

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
					sourceCode := sourceLangEnum.Value
					targetCode := targetLangEnum.Value

					// Translate using selected languages.
					translatedText, err := translator.Translate(inputText, sourceCode, targetCode)
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

			// Find the full names of the selected languages for the button label.
			var sourceName, targetName string
			for _, lang := range supportedLanguages {
				if lang.Code == sourceLangEnum.Value {
					sourceName = lang.Name
				}
				if lang.Code == targetLangEnum.Value {
					targetName = lang.Name
				}
			}
			buttonText := fmt.Sprintf("Translate from %s to %s", sourceName, targetName)

			// Define the layout.
			layout.UniformInset(16).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{
					Axis: layout.Vertical,
				}.Layout(gtx,
					// Source Language Selector
					layout.Rigid(material.Body1(theme, "From:").Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Inset{Top: 4, Bottom: 8}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return material.List(theme, &sourceLangList).Layout(gtx, len(supportedLanguages), func(gtx layout.Context, i int) layout.Dimensions {
								lang := supportedLanguages[i]
								return layout.Inset{Right: 12}.Layout(gtx, material.RadioButton(theme, &sourceLangEnum, lang.Code, lang.Name).Layout)
							})
						})
					}),
					layout.Flexed(1, material.Editor(theme, &inputEditor, "Text to translate").Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Inset{Top: 8, Bottom: 8}.Layout(gtx, material.Button(theme, &translateBtn, buttonText).Layout)
					}),
					// Target Language Selector
					layout.Rigid(material.Body1(theme, "To:").Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Inset{Top: 4, Bottom: 8}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return material.List(theme, &targetLangList).Layout(gtx, len(supportedLanguages), func(gtx layout.Context, i int) layout.Dimensions {
								lang := supportedLanguages[i]
								return layout.Inset{Right: 12}.Layout(gtx, material.RadioButton(theme, &targetLangEnum, lang.Code, lang.Name).Layout)
							})
						})
					}),
					layout.Rigid(material.Body1(theme, status).Layout),
					layout.Flexed(1, material.Editor(theme, &outputEditor, "Translation").Layout),
				)
			})

			e.Frame(gtx.Ops)
		}
	}
}
