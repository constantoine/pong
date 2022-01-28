package main

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
)

func BenchmarkUpdate(b *testing.B) {
	game := NewGame()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		game.Update()
	}
}

func BenchmarkRefresh(b *testing.B) {
	a := app.New()
	game := NewGame()
	w := a.NewWindow("Hello World")
	w.SetFixedSize(true)
	w.Resize(fyne.NewSize(768, 512))
	c := canvas.NewImageFromImage(game.GetImg())
	c.FillMode = canvas.ImageFillOriginal
	w.SetContent(c)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		game.Update()
		canvas.Refresh(c)
	}
	w.Close()
}
