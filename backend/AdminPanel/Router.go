package AdminPanel

import (
	"bufio"
	"diplom_server/backend/DBManager"
	"diplom_server/backend/structs"
	"encoding/json"
	"fmt"
	"golang.org/x/net/websocket"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"time"
)

var regexpComment = regexp.MustCompile("#.+")
var regexpStaticPath = regexp.MustCompile("(StaticPath[ ]+)(.+)")
var regexpRs = regexp.MustCompile("(GET|POST|PUT)[ ]+(.+)[ ]+(.+)")
var regexpParams = regexp.MustCompile(":([a-zA-Z]+)")
var regexpFaviconPath = regexp.MustCompile("(FaviconPath[ ]+)(.+)")

var StaticPath = ""
var FaviconPath = ""

type Router struct {
	dbm             *DBManager.DBManger
	f               *Filter
	dispatchChannel chan structs.Event
}

type Route struct {
	Method   string
	Path     string // Request
	BindPath string // thing to response
	Params   map[string]int
}

var Config []Route

var routeBinging = map[string]func(w http.ResponseWriter, req *http.Request){
	"index":           Index,
	"management":      Index,
	"settings":        Index,
	"management/list": nil,
	"ClearAllData":    nil,
	"GetData":         GetData,
	"CloseConnection": CloseConnection,
}

func RenderJSON(w http.ResponseWriter, req *http.Request, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	// get a payload p := Payload{d}
	if err := json.NewEncoder(w).Encode(data); err != nil {
		fmt.Println(err)
	}
}

func (r *Router) Init(c chan structs.Event) {
	r.dbm = new(DBManager.DBManger)
	r.dbm.Init(DBManager.POSTGRES)
	r.f = new(Filter)
	r.dispatchChannel = c
	r.parseConfig() // initialization
}

func (r *Router) mainRouting(w http.ResponseWriter, req *http.Request) {

	for _, route := range Config {
		if route.Method == req.Method {
			if route.Path == req.URL.Path {

				switch route.BindPath {
				case "CloseConnection":
					fmt.Println("Yo-ho-ho. That's what I've expected")

					connectionIdAsString := req.FormValue("fakeId")
					connectionId, err := strconv.Atoi(connectionIdAsString)
					if err != nil {
						fmt.Println(err)
						break
					}

					fmt.Println("connection id:", connectionId)

					event := structs.Event{
						Name:  structs.CloseEvent,
						Data:  connectionId,
						Delay: time.Time{},
					}
					r.dispatchChannel <- event
					break
				case "management/list":
					r.ManagementList(w, req)
					return
				case "ClearAllData":
					r.dbm.ClearAllData()
					return
				}

				routeBinging[route.BindPath](w, req)
			}
		}
	}
}

func (r *Router) parseParams(url string) map[string]int {
	tmpMap := make(map[string]int)
	tmp := regexpParams.FindAllStringSubmatch(url, -1)
	for i, match := range tmp {
		if len(match) > 1 {
			if _, exist := tmpMap[match[1]]; exist {
				log.Panic("you cannot use same params in one URI")
			} else {
				tmpMap[match[1]] = i
			}
		}
	}
	return tmpMap
}

func (r *Router) parseConfig() {
	Config = []Route{}
	file, err := os.Open("./backend/AdminPanel/routes.conf")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	// optionally, resize scanner's capacity for lines over 64K, see next example
	for scanner.Scan() {
		rs := scanner.Text()
		if idx := regexpComment.FindIndex([]byte(scanner.Text())); idx != nil {
			rs = rs[:idx[0]]
		}
		if len(rs) < 1 {
			continue
		}
		switch true {
		case len(regexpRs.FindString(rs)) > 0:
			{
				groups := regexpRs.FindStringSubmatch(rs)
				tmpMap := r.parseParams(groups[2])
				tmp := Route{
					Method:   groups[1],
					Path:     groups[2],
					BindPath: groups[3],
					Params:   tmpMap,
				}
				Config = append(Config, tmp)
			}
		case len(regexpStaticPath.FindString(rs)) > 0:
			{
				StaticPath = regexpStaticPath.FindStringSubmatch(rs)[2]
				// add new case for some variables
			}
		case len(regexpFaviconPath.FindString(rs)) > 0:
			{
				FaviconPath = regexpFaviconPath.FindStringSubmatch(rs)[2]
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func (r *Router) websocketRouter(ws *websocket.Conn) {
	fmt.Println(ws.Request().URL.Path)

	switch ws.Request().URL.Path {
	case "/update/test":
		UpdateTest(ws)
		break
	case "/update/hosts":
		fmt.Println("ws ask to update hosts")
		r.updateHost(ws)
		break
	case "/update/network":
		fmt.Println("get update network")
		keys, err := url.ParseQuery(ws.Request().URL.RawQuery)
		if err != nil {
			fmt.Println(err)
			break
		}
		r.updateNetwork(ws, keys["ip"][0])
		break
	case "/update/info":
		keys, err := url.ParseQuery(ws.Request().URL.RawQuery)
		if err != nil {
			fmt.Println(err)
			break
		}
		fmt.Println("get update info", keys["ip"][0])
		r.updateInfo(ws, keys["ip"][0])
		break
	default:
		fmt.Println("Try to access to unhandled route")
		break
	}
}

func (r *Router) switcher(filter *Filter) http.Handler {
	return filter.Manage(r.mainRouting)
}

func (r *Router) Manage() {

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(StaticPath))))
	http.Handle("/favicon.ico", http.FileServer(http.Dir(FaviconPath)))
	//http.Handle("/", r.switcher(f))
	http.Handle("/update/test", websocket.Handler(r.websocketRouter))
	http.Handle("/update/hosts", websocket.Handler(r.websocketRouter))
	http.Handle("/update/network", websocket.Handler(r.websocketRouter))
	http.Handle("/update/info", websocket.Handler(r.websocketRouter))

	http.HandleFunc("/", r.f.Manage(r.mainRouting)) // handle all incoming requests
	// "/" means every path. Filter incoming request then Route.
}
