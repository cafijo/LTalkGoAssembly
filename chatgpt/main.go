package main

import (
	"context"
	"fmt"
	gogpt "github.com/sashabaranov/go-gpt3"
	"os"
	"syscall/js"
)

var document = js.Global().Get("document")

func getElementByID(id string) js.Value {
	return document.Call("getElementById", id)
}

func renderEditor(parent js.Value) js.Value {
	editorMarkup := `
		<div id="editor" style="display: flex; flex-flow: row wrap;">
			<textarea id="markdown" style="width: 50%; height: 400px"></textarea>
			<div id="preview" style="width: 50%;"></div>
			<button id="render">ChatGPT</button>
		</div>
	`
	parent.Call("insertAdjacentHTML", "beforeend", editorMarkup)
	return getElementByID("editor")
}

func main() {
	quit := make(chan struct{}, 0)

	// See example 2: Enable the stop button
	stopButton := getElementByID("stop")
	stopButton.Set("disabled", false)
	stopButton.Set("onclick", js.FuncOf(func(js.Value, []js.Value) interface{} {
		println("stopping")
		stopButton.Set("disabled", true)
		quit <- struct{}{}
		return nil
	}))

	c := gogpt.NewClient(os.Getenv("APIKEY_CHATGPT"))

	// Simple markdown editor
	editor := renderEditor(document.Get("body"))
	preview := getElementByID("preview")
	markdown := getElementByID("markdown")

	renderButton := getElementByID("render")
	renderButton.Set("onclick", js.FuncOf(func(js.Value, []js.Value) interface{} {
		go func() {
			md := markdown.Get("value").String()
			ctx := context.Background()
			req := gogpt.CompletionRequest{
				Model:            gogpt.GPT3TextDavinci003,
				MaxTokens:        100,
				Temperature:      0.6,
				TopP:             1,
				BestOf:           1,
				FrequencyPenalty: 0,
				PresencePenalty:  0,
				Prompt:           md,
			}
			resp, err := c.CreateCompletion(ctx, req)
			if err != nil {
				fmt.Printf("%+v", err)
			}
			preview.Set("innerHTML", resp.Choices[0].Text)
		}()
		return nil
	}))
	<-quit

	editor.Call("remove")
}
