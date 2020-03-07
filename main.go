package main

import (
	"fmt"
	"github.com/AllenDang/giu"
)

type ScReader struct {
	Index int
	Error error
}

type CardStatusFunc func (index int) error

var firstName string
var lastName string

func onClickMe() {
	fmt.Println("Hello world!")
}

func onImSoCute() {
	fmt.Println("Im sooooooo cute!!")
}

func loop() {
	giu.SingleWindow("hello world", giu.Layout{
		giu.Label("Hello world from giu"),
		giu.Line(
			giu.Button("Click Me", onClickMe),
			giu.Button("I'm so cute", onImSoCute)),
			giu.Line(
				giu.InputText("Ime", 50.0, &firstName),
				giu.InputText("Prezime", 50.0, &lastName)),
	})
}

func main() {
	reader := NewReaderDevice()
	defer reader.Release()

	wnd := giu.NewMasterWindow("Identity", 400, 200, 0, nil)
	wnd.Main(loop)
}

func die(err error) {
	fmt.Println(err)
}
