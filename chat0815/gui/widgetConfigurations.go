package gui

import (
	"chat0815/contivity"
	"chat0815/errPopUps"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"log"
	"strings"
)

func groupChatNavigationConfiguration(chatC chan contivity.ChatStorage, gcStatusC chan *contivity.GroupChatStatus, a fyne.Window) *widget.List {
	list := widget.NewList(
		func() int {
			gcStatus := <-gcStatusC
			l := len(gcStatus.UserNames) - 1
			gcStatusC <- gcStatus
			return l
		},
		func() fyne.CanvasObject {
			return widget.NewButton("Template", func() {})
		},
		func(i widget.ListItemID, obj fyne.CanvasObject) {
			gcStatus := <-gcStatusC
			users := GetSortedKeyMap(gcStatus.UserNames, contivity.AddrWithoutPort(contivity.TcpAddr(contivity.GetOutboundIP())))
			addr := ""
			name := ""
			for j, userAddr := range users {
				if j == i {
					name = gcStatus.UserNames[userAddr]
					addr = userAddr
				}
			}
			obj.(*widget.Button).SetText(name)
			obj.(*widget.Button).OnTapped = func() {
				openPrivateTab(chatC, addr, name, a)
			}
			obj.(*widget.Button).Refresh()
			gcStatusC <- gcStatus

		},
	)
	return list
}

func newGroupChatDisplayConfiguration(gcStatusC chan *contivity.GroupChatStatus) *widget.List {
	mainChatDisplay := widget.NewList(
		func() int {
			gcStatus := <-gcStatusC
			gcStatusC <- gcStatus
			return len(gcStatus.ChatContent)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("Template")
		},
		func(i widget.ListItemID, obj fyne.CanvasObject) {
			gcStatus := <-gcStatusC
			contents := gcStatus.ChatContent
			msg := contents[len(contents)-1-i]
			obj.(*widget.Label).SetText(msg)
			if msg[0:6] == gcStatus.UserName {
				obj.(*widget.Label).Alignment = fyne.TextAlignTrailing
				obj.(*widget.Label).SetText(strings.Split(msg, ":")[1])
			} else if i == len(contents)-1 {
				obj.(*widget.Label).Alignment = fyne.TextAlignCenter
			} else if i == len(contents)-2 {
				obj.(*widget.Label).Alignment = fyne.TextAlignCenter
			} else {
				obj.(*widget.Label).Alignment = fyne.TextAlignLeading
			}
			gcStatusC <- gcStatus
		},
	)
	return mainChatDisplay
}

func newGroupInputEntry(gcStatusC chan *contivity.GroupChatStatus, errorC chan errPopUps.ErrorMessage) *groupInputEntry {
	entry := &groupInputEntry{}
	entry.ExtendBaseWidget(entry)
	entry.gcStatusC = gcStatusC
	entry.errorC = errorC

	entry.SetPlaceHolder("Write a Message")
	entry.OnChanged = func(typed string) {
		if len(typed) >= 50 {
			entry.SetText(entry.Text[:49])
		}
		if entry.Text == "/privateDebug" {
			entry.SetText("")
			log.Println("DEBUG PRIVATE CHAT")
			//TODO REMOVE THOSE WE HAVE NAVIGATION BUTTONS NOW
		}
		if entry.Text == "/privateChat" {
			entry.SetText("")
			log.Println("Private Chat Please")
			//TODO REMOVE THOSE WE HAVE NAVIGATION BUTTONS NOW
		}
	}
	return entry
}

type groupInputEntry struct {
	widget.Entry
	gcStatusC chan *contivity.GroupChatStatus
	errorC    chan errPopUps.ErrorMessage
}

func (e *groupInputEntry) onEnter() {
	if e.Entry.Text == "" {
		return
	}
	contivity.NGMX(e.Entry.Text, e.gcStatusC, e.errorC)
	e.Entry.SetText("")
}

func (e *groupInputEntry) KeyDown(key *fyne.KeyEvent) {
	switch key.Name {
	case fyne.KeyReturn:
		e.onEnter()
	}
}
