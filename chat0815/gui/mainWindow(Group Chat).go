package gui

import (
	"chat0815/contivity"
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/container"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"log"
)

func BuildApp(cStatusC chan *contivity.ChatroomStatus, errorC chan contivity.ErrorMessage, refresh chan bool) fyne.App {
	a := app.New()

	go manageLogWindow(errorC, a)

	mainWin := a.NewWindow("chat 0815")
	mainWin.Resize(fyne.NewSize(1200, 600))
	mainWin.SetFixedSize(true)
	mainWin.SetMaster()
	mainWin.SetOnClosed(func() { contivity.GBXX(cStatusC) })
	startUpWin := BuildStartUp(cStatusC, errorC, a, mainWin)
	chatDisplay := widget.NewList(
		func() int {
			cStatus := <-cStatusC
			cStatusC <- cStatus
			return len(cStatus.ChatContent)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("Template")
		},
		func(i widget.ListItemID, obj fyne.CanvasObject) {
			cStatus := <-cStatusC
			contents := cStatus.ChatContent
			obj.(*widget.Label).SetText(contents[len(contents)-1-i])
			cStatusC <- cStatus
		},
	)
	//Manages Refresh requests
	go manageChatDisplayRefresh(refresh, chatDisplay)

	//Input Chat Console
	input := widget.NewEntry()
	inputEntryConfiguration(a, cStatusC, input)
	//Send Button
	send := widget.NewButton("Send it!", func() {
		sendButtonClicked(input, cStatusC, errorC)
	})

	chat := container.NewMax(chatDisplay)
	navigation := fyne.NewContainerWithLayout(layout.NewVBoxLayout(), input, send)
	content := container.NewHSplit(navigation, chat)
	mainWin.SetContent(content)
	startUpWin.Show()
	return a
}

func manageChatDisplayRefresh(refresh chan bool, chatDisplay *widget.List) {
	for {
		check := <-refresh
		if check {
			chatDisplay.Refresh()
		}
	}
}

//funktionen: Placeholder, TODO Cap Max Letters
func inputEntryConfiguration(a fyne.App, cStatusC chan *contivity.ChatroomStatus, input *widget.Entry) {
	input.SetPlaceHolder("Write a Message")
	input.OnChanged = func(typed string) {
		if len(typed) >= 50 {
			//input.Disable()
			input.SetText(input.Text[:49])
		}
		if input.Text == "/privateDebug" {
			input.SetText("")
			log.Println("DEBUG PRIVATE CHAT")
			go openRealPrivateWin(a, cStatusC, "", "NONAME")
		}
		if input.Text == "/privateChat" {
			input.SetText("")
			log.Println("Private Chat Please")
			go OpenPrivateWin(a, cStatusC)
		}
	}
}
func sendButtonClicked(input *widget.Entry, cStatusC chan *contivity.ChatroomStatus, errorC chan contivity.ErrorMessage) {
	if input.Text == "" {
		return
	}
	contivity.NGMX(input.Text, cStatusC, errorC)
	input.SetText("")
}

func manageLogWindow(errorC chan contivity.ErrorMessage, a fyne.App) {
	var logs contivity.ErrorMessage
	for {
		logs = <-errorC
		go showLog(logs, a)
	}
}
