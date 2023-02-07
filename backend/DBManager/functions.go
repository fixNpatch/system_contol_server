package DBManager

import (
	"diplom_server/backend/structs"
	"fmt"
	"math"
	"regexp"
	"time"
)

// upsert procedure === >
// create or replace procedure upsert_status(_id integer, _ip varchar(22), _status integer) language 'plpgsql'
// as $$
//    begin
//        insert INTO data.connections(id, ip, status) VALUES (_id, _ip, _status)
//        ON CONFLICT (ip)
//        DO UPDATE SET status = _status;
//    end;
//    $$;
func (dbm *DBManger) SaveStatus(statusMap map[string]int) (err error) {
	statement := `call upsert_status($1, $2, $3);`
	r, _ := regexp.Compile(`\d+.\d+.\d+.\d+.\d+`)

	counter := 1
	for i, status := range statusMap {
		// EXEC MUST BE USED IF UPDATE OR INSERT
		// https://aloksinhanov.medium.com/query-vs-exec-vs-prepare-in-golang-e7c49212c36c

		_, err = dbm.postgres.Exec(statement, counter, r.FindString(i), status)
		if err != nil {
			fmt.Println(err)
			return
		}
		counter++
	}

	return nil
}

func (dbm *DBManger) GetAgentId(agentIp string) (int, error) {
	var id int
	stmt := `SELECT c.id FROM data.connections c WHERE c.ip = $1;`

	rows, err := dbm.postgres.Query(stmt, agentIp)
	if err != nil {
		fmt.Println(err)
		return -1, err
	}

	for rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			fmt.Println(err)
			return -1, err
		}
	}

	return id, nil
}

// TODO
func (dbm *DBManger) saveProcs(agentId int, stats structs.Stats) {

	// TODO
	fmt.Println(stats.Processes)

	//fmt.Printf("%#v\n", stats.Processes[0])

	//var tocId int64
	//createNewTocSTMT := `INSERT INTO data.toc_procs(fk_connection_id, stamp) VALUES ($1, $2) RETURNING id;`
	//err := dbm.postgres.QueryRow(createNewTocSTMT, agentId, time.Now()).Scan(&tocId)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//
	//sqlStr := "INSERT INTO data.network(pid, procname, status, rport, " +
	//	"ppid, uids, gids, _groups, numthreads, createtime, fk_network_id) values "
	//
	//vals := []interface{}{}
	//
	//for i, row := range stats.Processes {
	//	sqlStr += fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d),",
	//		i*11+1, i*11+2, i*11+3,
	//		i*11+4, i*11+5, i*11+6,
	//		i*11+7, i*11+8, i*11+9,
	//		i*11+10, i*11+11)
	//	vals = append(vals,
	//		row.Pid, row.Name, row.Status
	//		strings.Trim(strings.Replace(fmt.Sprint(row.LAddr), " ", ".", -1), "[]"), row.LPort,
	//		strings.Trim(strings.Replace(fmt.Sprint(row.RAddr), " ", ".", -1), "[]"), row.RPort,
	//		row.Pid, row.ProcName,
	//		row.ProcOwner, tocId)
	//}
	//sqlStr = sqlStr[0 : len(sqlStr)-1] // trim the last ,
	//stmt, err := dbm.postgres.Prepare(sqlStr)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	////format all vals at once
	//_, err = stmt.Exec(vals...)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
}

//func (dbm *DBManger) saveNetwork(agentId int, stats structs.Stats, changesOnly bool) {
//	var tocId int64
//	//
//	createNewTocSTMT := `INSERT INTO data.toc_network(fk_connection_id, stamp) VALUES ($1, $2) RETURNING id;`
//	err := dbm.postgres.QueryRow(createNewTocSTMT, agentId, time.Now()).Scan(&tocId)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//
//	sqlStr := "INSERT INTO data.network(laddr, lport, raddr, rport, pid, procname, procowner, fk_network_id) values "
//
//	vals := []interface{}{}
//
//	for i, row := range stats.Connections {
//		//prepStr += "(?, ?, ?, ?, ?, ?, ?, (select id from tocid)),"
//		sqlStr += fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d),",
//			i*8+1, i*8+2, i*8+3, i*8+4, i*8+5, i*8+6, i*8+7, i*8+8)
//		vals = append(vals,
//			strings.Trim(strings.Replace(fmt.Sprint(row.LAddr), " ", ".", -1), "[]"), row.LPort,
//			strings.Trim(strings.Replace(fmt.Sprint(row.RAddr), " ", ".", -1), "[]"), row.RPort,
//			row.Pid, row.ProcName,
//			row.ProcOwner, tocId)
//	}
//	sqlStr = sqlStr[0 : len(sqlStr)-1] // trim the last ,
//	stmt, err := dbm.postgres.Prepare(sqlStr)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	//format all vals at once
//	_, err = stmt.Exec(vals...)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//}

func (dbm *DBManger) GetHostsConfig() (hosts []structs.Host) {
	var (
		id   int
		name string
		ip   string
	)

	stmt := `SELECT c.id, c.name, c.ip FROM config.hosts c;`

	rows, err := dbm.postgres.Query(stmt)
	if err != nil {
		fmt.Println(err)
	}
	for rows.Next() {
		err = rows.Scan(&id, &name, &ip)
		if err != nil {
			fmt.Println(err)
		}
		hosts = append(hosts, structs.Host{
			Id:     id,
			Name:   name,
			IP:     ip,
			Status: -1,
		})
	}

	return
}

func (dbm *DBManger) GetAgentStateByID(agentId int) (structs.HostInfo, structs.Disk, structs.Cpu) {
	var (
		info structs.HostInfo
		disk structs.Disk
		cpu  structs.Cpu
	)
	var (
		cpu_usage           int
		disk_free           int
		disk_total          int
		disk_used           int
		cpu_cores           int
		cpu_model           string
		os                  string
		os_platform         string
		os_platform_version string
	)

	stmt := `SELECT 
			 d.cpu_usage, d.disk_free, d.disk_total, d.disk_used, 
			 i.cpu_cores, i.cpu_model, i.os_os, i.os_platform, i.os_platform_version
			 FROM data.state d
			 LEFT JOIN data.info i ON i.id = $1
			 WHERE d.fk_connection_id = i.id
			 ORDER BY d.stamp DESC limit 1;`
	rows, err := dbm.postgres.Query(stmt, agentId)
	for rows.Next() {
		err = rows.Scan(&cpu_usage,
			&disk_free, &disk_total,
			&disk_used, &cpu_cores,
			&cpu_model, &os, &os_platform,
			&os_platform_version)
		if err != nil {
			fmt.Println(err)
		}
	}
	info = structs.HostInfo{
		Procs:           0,
		OS:              os,
		PlatformVersion: os_platform_version,
		Platform:        os_platform,
	}
	disk = structs.Disk{
		Total:       uint64(disk_total),
		Free:        uint64(disk_free),
		Used:        uint64(disk_used),
		UsedPercent: 0,
	}
	cpu = structs.Cpu{
		Percentage: []float64{float64(cpu_usage)},
		Model:      cpu_model,
		Cores:      cpu_cores,
	}
	return info, disk, cpu
}

func (dbm *DBManger) saveState(agentId int, stats structs.Stats) {

	stmt := `INSERT INTO data.state(fk_connection_id, cpu_usage, disk_total, disk_free, disk_used, stamp) 
			 values ($1, $2, $3, $4, $5, $6);`

	_, err := dbm.postgres.Exec(stmt, agentId, int(math.Round(stats.Cpu.Percentage[0])), stats.Disk.Total/1024/1024, stats.Disk.Free/1024/1024, stats.Disk.Used/1024/1024, time.Now())
	if err != nil {
		fmt.Println(err)
		return
	}

}

func (dbm *DBManger) GetDataByIP(agent string) (result structs.Stats) {
	r, _ := regexp.Compile(`\d+.\d+.\d+.\d+.\d+`)
	agent = r.FindString(agent)

	// получаем id агента
	id, err := dbm.GetAgentId(agent)
	if err != nil {
		fmt.Println(err)
		return
	}

	info, disk, cpu := dbm.GetAgentStateByID(id)
	conns := dbm.GetAgentConnectionsByID(id)

	result = structs.Stats{
		VmStat:      structs.VmStat{},
		Disk:        disk,
		Cpu:         cpu,
		Host:        info,
		Processes:   nil,
		Connections: conns,
		HostTime:    time.Time{},
	}

	return
}

func (dbm *DBManger) SaveInfo(agent string, data structs.Stats) {

}

func (dbm *DBManger) SaveChangesByIP(agent string, received structs.Stats, closedConnections, openedConnections []structs.Connection) {
	r, _ := regexp.Compile(`\d+.\d+.\d+.\d+.\d+`)
	agent = r.FindString(agent)

	// получаем id агента
	id, err := dbm.GetAgentId(agent)
	if err != nil {
		fmt.Println(err)
		return
	}

	dbm.saveNetworkChanges(id, closedConnections, openedConnections)
	dbm.saveState(id, received)
}

func (dbm *DBManger) SaveData(agent string, stats structs.Stats) {
	r, _ := regexp.Compile(`\d+.\d+.\d+.\d+.\d+`)
	agent = r.FindString(agent)

	// получаем id агента
	id, err := dbm.GetAgentId(agent)
	if err != nil {
		fmt.Println(err)
		return
	}

	dbm.saveState(id, stats)
	dbm.saveNetwork(id, stats)

	//dbm.saveProcs(id, stats) // processes not followed

	return
}

func (dbm *DBManger) GetHostsStatus() (data []structs.Host) {
	var (
		id     int
		name   string
		ip     string
		status int
	)
	stmt := `SELECT c.id, c.ip, c.status, h.name 
			 FROM data.connections c JOIN config.hosts h
			 ON h.ip = c.ip ORDER BY c.ip;`

	rows, err := dbm.postgres.Query(stmt)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	for rows.Next() {
		err = rows.Scan(&id, &ip, &status, &name)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		data = append(data, structs.Host{
			Id:     id,
			Name:   name,
			IP:     ip,
			Status: status,
		})
	}

	return
}

func (dbm *DBManger) LoadInitialSettings() (ipStack []string) {
	var ip string

	statement := `SELECT h.ip FROM config.hosts h;`
	rows, err := dbm.postgres.Query(statement)
	if err != nil {
		fmt.Println(err)
		return
	}

	for rows.Next() {
		err = rows.Scan(&ip)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Follow:", ip)
		ipStack = append(ipStack, ip)
	}

	if ipStack == nil {
		fmt.Println("ipStack is empty")
	}
	return
}

func (dbm *DBManger) Test() (err error) {
	var test string

	rows, err := dbm.postgres.Query("SELECT a.value FROM data.test a WHERE a.id = $1", 1)
	if rows == nil {
		fmt.Println("rows == nil")
		return
	}

	for rows.Next() {
		if err = rows.Scan(&test); err != nil {
			fmt.Println(err)
			return err
		}
		fmt.Println(test)
	}
	return err
}
