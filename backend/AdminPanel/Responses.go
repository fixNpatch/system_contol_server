package AdminPanel

import (
	"diplom_server/backend/structs"
	"fmt"
	"golang.org/x/net/websocket"
	"html/template"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

func Index(w http.ResponseWriter, req *http.Request) {
	fmt.Println("index called")
	tmpl := template.Must(template.ParseFiles("./frontend/index.html"))
	err := tmpl.Execute(w, nil)
	if err != nil {
		fmt.Println(err)
	}
}

func (r *Router) ManagementList(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Management List")
	hosts := r.dbm.GetHostsConfig()
	RenderJSON(w, req, hosts)
}

func GetData(w http.ResponseWriter, req *http.Request) {
	fmt.Println("data")

	test := structs.Stats{
		VmStat: structs.VmStat{1, 2, 3},
		Disk:   structs.Disk{},
		Cpu:    structs.Cpu{},
		Host:   structs.HostInfo{},
	}
	RenderJSON(w, req, test)
}

func CloseConnection(w http.ResponseWriter, req *http.Request) {
	fmt.Println("AdminPanel::Responses::closeConnection")
}

func UpdateTest(ws *websocket.Conn) {
	var err error
	var checkResponse string
	fmt.Println(ws.RemoteAddr())
	for {
		now := time.Now()
		str := strconv.Itoa(rand.Int())
		fmt.Println("try to send", now.Format(time.RFC822))

		if err = websocket.Message.Send(ws, str); err != nil {
			fmt.Println("Can't send")
			break
		}

		if err = websocket.Message.Receive(ws, &checkResponse); err != nil {
			fmt.Println(err)
			break
		}

	}
	if err = ws.Close(); err != nil {
		fmt.Println(err)
	}
}

func (r *Router) updateInfo(ws *websocket.Conn, agentIp string, data ...interface{}) {
	var err error
	var Data structs.Stats
	var checkResponse string

	id, err := r.dbm.GetAgentId(agentIp)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("updateInfo::getAgentId", id)

	for {
		info, disk, cpu := r.dbm.GetAgentStateByID(id)

		Data = structs.Stats{
			VmStat: structs.VmStat{
				Total:       20,
				Free:        0,
				UsedPercent: 0,
			},
			Host: info,
			Disk: disk,
			Cpu:  cpu,
		}

		if err = websocket.JSON.Send(ws, Data); err != nil {
			fmt.Println("Can't send")
			break
		}

		if err = websocket.Message.Receive(ws, &checkResponse); err != nil {
			fmt.Println(err)
			break
		}

		time.Sleep(2 * time.Second)
	}

	if err = ws.Close(); err != nil {
		fmt.Println(err)
	}

}

//todo update host
func (r *Router) updateHost(ws *websocket.Conn, data ...interface{}) {
	var err error
	var checkResponse string
	var Data []structs.Host

	for {
		Data = r.dbm.GetHostsStatus()

		//Data = []structs.Host{
		//	{1, "AlexPC", "192.168.0.1", rand.Intn(10)},
		//	{2, "AdamPC", "192.168.0.2", rand.Intn(10)},
		//	{3, "BretPC", "192.168.1.1", rand.Intn(10)},
		//	{4, "AllyPC", "192.168.1.2", rand.Intn(10)},
		//}

		if err = websocket.JSON.Send(ws, Data); err != nil {
			fmt.Println("Can't send")
			break
		}

		if err = websocket.Message.Receive(ws, &checkResponse); err != nil {
			fmt.Println(err)
			break
		}

		time.Sleep(2 * time.Second)
	}

	if err = ws.Close(); err != nil {
		fmt.Println(err)
	}
}

//todo update plot
func UpdatePlot(ws *websocket.Conn) {
	var err error

	for {
		now := time.Now()

		str := strconv.Itoa(rand.Int())
		fmt.Println("try to send", now.Format(time.RFC822))

		if err = websocket.Message.Send(ws, str); err != nil {
			fmt.Println("Can't send")
			break
		}
	}
}

//todo updatenetwork
func (r *Router) updateNetwork(ws *websocket.Conn, agentIp string, data ...interface{}) {
	var err error
	var Data structs.Stats
	var checkResponse string

	id, err := r.dbm.GetAgentId(agentIp)
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		connections := r.dbm.GetAgentConnectionsByID(id)
		//connections := r.dbm.GetAgentConnectionsByID(id)

		Data = structs.Stats{
			Connections: connections,
		}

		if err = websocket.JSON.Send(ws, Data); err != nil {
			fmt.Println("Can't send")
			break
		}

		if err = websocket.Message.Receive(ws, &checkResponse); err != nil {
			fmt.Println(err)
			break
		}

		time.Sleep(2 * time.Second)
	}

	if err = ws.Close(); err != nil {
		fmt.Println(err)
	}
}
