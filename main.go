//
// main.go
// Copyright (C) 2019 Yongwen Zhuang <zeoman@163.com>
//
// Distributed under terms of the MIT license.
//

package main

import (
	"fmt"
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func confirmCallback(response bool) {
	fmt.Println("Responded with", response)
}

func main() {
	a := app.New()
	w := a.NewWindow("Fyne Demo")

	logo := canvas.NewImageFromResource(theme.FyneLogo())
	logo.SetMinSize(fyne.NewSize(64, 64))

	name := widget.NewEntry()
	name.SetPlaceHolder("[格式]姓 名 字")
	minAttr := widget.NewEntry()
	story := widget.NewMultiLineEntry()

	form := &widget.Form{
		OnCancel: func() {
			w.Close()
		},
		OnSubmit: func() {
			fmt.Println("Form submitted")
			fmt.Println("Name:", name.Text)
			fmt.Println("Email:", minAttr.Text)
		},
	}
	gen := widget.NewCheck("自动生成", func(on bool) { fmt.Println("checked", on) })
	form.Append("姓名", name)
	form.Append("最低能力值", minAttr)
	form.Append("", gen)
	form.Append("列传", story)
	w.SetContent(widget.NewVBox(
		widget.NewToolbar(widget.NewToolbarAction(theme.MailComposeIcon(), func() { fmt.Println("New") }),
			widget.NewToolbarSeparator(),
			widget.NewToolbarSpacer(),
			widget.NewToolbarAction(theme.ContentCutIcon(), func() { fmt.Println("Cut") }),
			widget.NewToolbarAction(theme.ContentCopyIcon(), func() { fmt.Println("Copy") }),
			widget.NewToolbarAction(theme.ContentPasteIcon(), func() { fmt.Println("Paste") }),
		),
		form,
	))
	w.ShowAndRun()
}
