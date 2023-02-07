package HostForAgents

import (
	"bytes"
	"diplom_server/backend/DBManager"
	"diplom_server/backend/structs"
	"encoding/json"
	"fmt"
	"golang.org/x/net/websocket"
	"io"
	"log"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	_ "github.com/lib/pq" // загружаем драйвер postgres
)

type Conn struct {
	conn              *websocket.Conn
	status            int
	ConnectionsToKill []structs.Connection
}

type ParamsToCloseConnection struct {
}

// HostStatus
const (
	NotConnected = iota
	Active
	Lost
)

type config struct {
	clientUpdateInterval time.Duration // seconds
	mode                 int           // look at const above
	dbm                  *DBManager.DBManger
}

type Server2 struct {
	config
	ipStack           []string
	wsStack           map[string]Conn
	eventListenerChan chan structs.Event
}

func (s *Server2) Init(c chan structs.Event) {
	s.wsStack = make(map[string]Conn)
	s.dbm = new(DBManager.DBManger)
	s.dbm.Init(DBManager.POSTGRES)
	s.ipStack = s.dbm.LoadInitialSettings()
	s.eventListenerChan = c
	s.clientUpdateInterval = time.Second * 2
	s.mode = structs.SaveOnlyChanges
}

func (s *Server2) saveStatus() error {
	statusMap := make(map[string]int)

	for i, elem := range s.wsStack {
		statusMap[i] = elem.status
	}

	if err := s.dbm.SaveStatus(statusMap); err != nil {
		fmt.Println("Server2::saveStatus::", err)
		return err
	}

	return nil
}

func (s *Server2) update() {
	var wg sync.WaitGroup

	for idx := range s.wsStack {
		go func(agentIP string, stack map[string]Conn) {
			wg.Add(1)
			defer wg.Done()

			conn := stack[agentIP].conn

			switch stack[agentIP].status {
			case NotConnected:
				return
			case Active:
				{
					var tmp []structs.Connection
					for _, row := range stack[agentIP].ConnectionsToKill {
						var bodyString string
						r, _ := regexp.Compile(`\d+.\d+.\d+.\d+.\d+`)
						agent := r.FindString(agentIP)

						values := map[string]string{
							"laddr": strings.Trim(strings.Replace(fmt.Sprint(row.LAddr), " ", ".", -1), "[]"),
							"raddr": strings.Trim(strings.Replace(fmt.Sprint(row.RAddr), " ", ".", -1), "[]"),
							"lport": strconv.Itoa(row.LPort),
							"rport": strconv.Itoa(row.RPort),
							"pid":   strconv.Itoa(row.Pid),
						}

						bytez, err := json.Marshal(values)
						if err != nil {
							fmt.Println(err)
						}

						urlstring := "http://" + agent + "/close"
						resp, err := http.Post(urlstring, "application/json", bytes.NewBuffer(bytez))
						if err != nil {
							fmt.Println(err)
						}

						if resp.StatusCode == http.StatusOK {
							bodyBytes, err := io.ReadAll(resp.Body)
							if err != nil {
								log.Fatal(err)
							}
							bodyString = string(bodyBytes)
							fmt.Println(bodyString)
						}

						if err = resp.Body.Close(); err != nil {
							fmt.Println(err)
						}

						if bodyString != "ok" {
							tmp = append(tmp, row)
						}
					}

					// переприсваивание только ради обновления ConnectionsToKill
					s.wsStack[agentIP] = Conn{
						conn:              conn,
						status:            stack[agentIP].status,
						ConnectionsToKill: tmp,
					}

					if err := websocket.Message.Send(conn, []byte("fupd")); err != nil {
						stack[agentIP] = Conn{
							conn:   nil,
							status: Lost,
						}
						return
					}
					break
				}
			case Lost:
				{
					if wsRenew, err := websocket.Dial(agentIP, "", agentIP[:9]); err == nil {
						stack[agentIP] = Conn{
							conn:   wsRenew,
							status: Active,
						}
						conn = wsRenew
					} else {
						fmt.Println(err)
						return
					}
					break
				}
			}

			var ans structs.Stats
			if err := websocket.JSON.Receive(conn, &ans); err != nil {
				fmt.Println(err)
				return
			}
			s.saveData(agentIP, ans)
		}(idx, s.wsStack)
	}

	wg.Wait()
}

func (s *Server2) saveData(agentIP string, received structs.Stats) {
	switch s.mode {
	case structs.SaveOnlyChanges:
		previousData := s.dbm.GetDataByIP(agentIP)
		closedConnections, openedConnections := s.findNetworkChanges(previousData, received)
		s.dbm.SaveChangesByIP(agentIP, received, closedConnections, openedConnections)
		s.dbm.SaveInfo(agentIP, received)
		break
	case structs.SaveFullData:
		s.dbm.SaveData(agentIP, received)
		break
	}
}

func compareConnections(prev, next structs.Connection) bool {
	if prev.LPort == next.LPort &&
		prev.RPort == next.RPort &&
		prev.Pid == next.Pid &&
		prev.ProcName == next.ProcName &&
		prev.ProcOwner == next.ProcOwner &&
		reflect.DeepEqual(prev.LAddr, next.LAddr) &&
		reflect.DeepEqual(prev.RAddr, next.RAddr) {
		return true
	}
	return false
}

func (s *Server2) findNetworkChanges(prev, next structs.Stats) (closed, opened []structs.Connection) {
	var found bool
	idxArray := make(map[int]bool, len(next.Connections))
	//fmt.Println("\n-------------------")
	//fmt.Println("prev", prev)
	//fmt.Println("next", next)

	for _, savedElem := range prev.Connections {
		found = false
		for i, receivedElem := range next.Connections {
			if compareConnections(savedElem, receivedElem) {
				idxArray[i] = true
				found = true
				break
			}
		}
		if !found {
			closed = append(closed, savedElem)
		}
	}

	for i, receivedElem := range next.Connections {
		if _, ok := idxArray[i]; !ok {
			opened = append(opened, receivedElem)
		}
	}

	return
}

func (s *Server2) makeConnections() {
	for _, ip := range s.ipStack {
		// устанавливаем соединение
		status := Active

		url := fmt.Sprint("ws://" + ip + "/")
		ws, err := websocket.Dial(url, "", "ws://"+ip[:9])
		if err != nil {
			ws = nil
			status = NotConnected
			fmt.Println(err)
			continue
		}

		s.wsStack[url] = Conn{
			conn:              ws,
			status:            status,
			ConnectionsToKill: nil,
		}
	}
}

func (s *Server2) makeTaskOnClose(event structs.Event) {
	fmt.Println("server::makeTaskOnClose::", event)

	value := reflect.ValueOf(event.Data)
	conn, agent := s.dbm.GetConnection(value.Int())
	if agent == "" {
		fmt.Println("Cannot find agentIP to kill connection")
		return
	}

	agentURL := "ws://" + agent + "/"

	fmt.Println("server::makeTaskOnClose::", s.wsStack[agentURL])
	// получаем копию entry
	if entry, Found := s.wsStack[agentURL]; Found {
		fmt.Println("ok")
		entry.ConnectionsToKill = append(entry.ConnectionsToKill, conn)
		// переприсваиваем
		s.wsStack[agentURL] = entry
	}

	fmt.Println("server::makeTaskOnClose::", s.wsStack[agentURL].ConnectionsToKill)
}

func (s *Server2) Manage() {
	flagInitiated := false // указываем на то, что соединения не инициализированы

	for {
		if !flagInitiated { // инициализируем соединения
			s.makeConnections()
			flagInitiated = true
		}

		if err := s.saveStatus(); err != nil {
			fmt.Println("saveStatus::", err)
			return
		}

		s.update()

		if len(s.eventListenerChan) > 0 { // если пришел event из web-server'а
			event := <-s.eventListenerChan
			fmt.Println("server2::", event)

			switch event.Name {
			case structs.CloseEvent:
				fmt.Println("Get event to close connection")
				s.makeTaskOnClose(event)
				break
			case structs.RefreshConnectionListEvent:
				fmt.Println("Get event to refresh connections")
				flagInitiated = false
				break
			}
		}

		time.Sleep(s.clientUpdateInterval)
	}
}
