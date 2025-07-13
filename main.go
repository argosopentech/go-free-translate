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
	// Use a local LibreTranslate instance.
	translator := libretranslate.New(libretranslate.Config{
		Url: "http://localhost:5000",
	})

	theme := material.NewTheme()

	var (
		ops            op.Ops
		inputEditor    widget.Editor
		outputEditor   widget.Editor
		translateBtn   widget.Clickable
		status         string
		sourceLangEnum widget.Enum
		targetLangEnum widget.Enum
		sourceLangList widget.List
		targetLangList widget.List
	)

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

			if translateBtn.Clicked(gtx) {
				status = "Translating..."
				go func() {
					inputText := inputEditor.Text()
					sourceCode := sourceLangEnum.Value
					targetCode := targetLangEnum.Value

					translatedText, err := translator.Translate(inputText, sourceCode, targetCode)
					if err != nil {
						log.Printf("Translation error: %v", err)
						status = "Error: " + err.Error()
						window.Invalidate()
						return
					}

					status = ""
					outputEditor.SetText(translatedText)
					window.Invalidate()
				}()
			}

			layoutUI(gtx, theme, &inputEditor, &outputEditor, &translateBtn,
				&sourceLangEnum, &targetLangEnum, &sourceLangList, &targetLangList, status)

			e.Frame(gtx.Ops)
		}
	}
}

func layoutUI(
	gtx layout.Context,
	theme *material.Theme,
	inputEditor, outputEditor *widget.Editor,
	translateBtn *widget.Clickable,
	sourceLangEnum, targetLangEnum *widget.Enum,
	sourceLangList, targetLangList *widget.List,
	status string,
) layout.Dimensions {
	return layout.UniformInset(16).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{
			Axis: layout.Vertical,
		}.Layout(gtx,
			layout.Rigid(material.Body1(theme, "From:").Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return sourceLangSelector(gtx, theme, sourceLangEnum, sourceLangList)
			}),
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				return inputEditorSection(gtx, theme, inputEditor)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return translateButton(gtx, theme, translateBtn, sourceLangEnum.Value, targetLangEnum.Value)
			}),
			layout.Rigid(material.Body1(theme, "To:").Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return targetLangSelector(gtx, theme, targetLangEnum, targetLangList)
			}),
			layout.Rigid(material.Body1(theme, status).Layout),
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				return outputEditorSection(gtx, theme, outputEditor)
			}),
		)
	})
}


func sourceLangSelector(gtx layout.Context, theme *material.Theme, enum *widget.Enum, list *widget.List) layout.Dimensions {
	return layout.Inset{Top: 4, Bottom: 8}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return material.List(theme, list).Layout(gtx, len(supportedLanguages), func(gtx layout.Context, i int) layout.Dimensions {
			lang := supportedLanguages[i]
			return layout.Inset{Right: 12}.Layout(gtx, material.RadioButton(theme, enum, lang.Code, lang.Name).Layout)
		})
	})
}

func targetLangSelector(gtx layout.Context, theme *material.Theme, enum *widget.Enum, list *widget.List) layout.Dimensions {
	return layout.Inset{Top: 4, Bottom: 8}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return material.List(theme, list).Layout(gtx, len(supportedLanguages), func(gtx layout.Context, i int) layout.Dimensions {
			lang := supportedLanguages[i]
			return layout.Inset{Right: 12}.Layout(gtx, material.RadioButton(theme, enum, lang.Code, lang.Name).Layout)
		})
	})
}

func inputEditorSection(gtx layout.Context, theme *material.Theme, editor *widget.Editor) layout.Dimensions {
	return material.Editor(theme, editor, "Text to translate").Layout(gtx)
}

func outputEditorSection(gtx layout.Context, theme *material.Theme, editor *widget.Editor) layout.Dimensions {
	return material.Editor(theme, editor, "Translation").Layout(gtx)
}

func translateButton(gtx layout.Context, theme *material.Theme, btn *widget.Clickable, sourceCode string, targetCode string) layout.Dimensions {
	var sourceName, targetName string
	for _, lang := range supportedLanguages {
		if lang.Code == sourceCode {
			sourceName = lang.Name
		}
		if lang.Code == targetCode {
			targetName = lang.Name
		}
	}
	buttonText := fmt.Sprintf("Translate from %s to %s", sourceName, targetName)
	return layout.Inset{Top: 8, Bottom: 8}.Layout(gtx, material.Button(theme, btn, buttonText).Layout)
}
