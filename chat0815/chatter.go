package main

import (
	"chat0815/contivity"
	"chat0815/gui"
	"log"
	"net"
)

func main() {
	outboundAddr := contivity.TcpAddr(contivity.GetOutboundIP())
	chatContent := make([]string, 0)

	chatContent = append(chatContent, "Take care of each other and watch your drink")
	chatContent = append(chatContent, "Welcome to chat0815")

	cStatus := contivity.ChatroomStatus{
		ChatContent: chatContent,
		UserAddr:    []net.Addr{},
		BlockedAddr: []net.Addr{},
	}
	refresh := make(chan bool)
	cStatusC := make(chan *contivity.ChatroomStatus)
	go manageCStatus(&cStatus, cStatusC)

	//__________________________________________________________________
	l, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Println("Listener died")
		log.Fatal(err)
	}
	defer l.Close()
	//start server
	go contivity.RunServer(l, cStatusC, refresh)
	//Sync with yourself first
	go contivity.GetStatusUpdate(&outboundAddr, cStatusC, refresh)
	//FYNE STUFF
	log.Println("LOL")
	a := gui.BuildApp(cStatusC, refresh)
	a.Run()

}

//Provides the pointer to the current ChatRoomStatus, always waits for a pointer in return.
//Should run in own goroutine
func manageCStatus(cStatus *contivity.ChatroomStatus, c chan *contivity.ChatroomStatus) {
	for {
		c <- cStatus
		cStatus = <-c
	}
}
