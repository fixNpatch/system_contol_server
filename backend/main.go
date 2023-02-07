package main

import (
	"diplom_server/backend/AdminPanel"
	"diplom_server/backend/HostForAgents"
	"diplom_server/backend/structs"
	"net/http"
)

var (
	router  *AdminPanel.Router
	server2 *HostForAgents.Server2
)

func main() {

	ch := make(chan structs.Event, 1)

	go func() {
		router = new(AdminPanel.Router)
		router.Init(ch)
		router.Manage()
	}()

	go func() {
		server2 = new(HostForAgents.Server2)
		server2.Init(ch)
		server2.Manage()
	}()

	_ = http.ListenAndServe(":80", nil)
}
