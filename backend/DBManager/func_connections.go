package DBManager

import (
	"diplom_server/backend/structs"
	"fmt"
	"github.com/lib/pq"
	"strconv"
	"strings"
	"time"
)

func (dbm *DBManger) saveNetworkChanges(agentId int, closed, opened []structs.Connection) {

	// first step. Mark closed connections

	if closed != nil {
		// database time MUST be SYNC with host, otherwise this code should be refactored
		closeStmt := "UPDATE data.changes_network SET (status, closedWhen) = (0, $1) WHERE id " + "in ("
		for _, row := range closed {
			closeStmt += strconv.Itoa(row.FakeId) + ", "
		}
		closeStmt = closeStmt[:len(closeStmt)-2] + ");"
		if _, err := dbm.postgres.Exec(closeStmt, time.Now()); err != nil {
			fmt.Println(err)
		}
	}

	if opened != nil {
		// second step. insert new connections
		updateStmt := "INSERT INTO data.changes_network(laddr, lport, raddr, rport, pid, procname, procowner, fk_connection_id, activeSince, status) values "

		vals := []interface{}{}
		for i, row := range opened {
			updateStmt += fmt.Sprintf("($%d, $%d, $%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d),",
				i*10+1, i*10+2, i*10+3, i*10+4, i*10+5, i*10+6, i*10+7, i*10+8, i*10+9, i*10+10)
			vals = append(vals,
				strings.Trim(strings.Replace(fmt.Sprint(row.LAddr), " ", ".", -1), "[]"), row.LPort,
				strings.Trim(strings.Replace(fmt.Sprint(row.RAddr), " ", ".", -1), "[]"), row.RPort,
				row.Pid, row.ProcName, row.ProcOwner, agentId, time.Now(), 1)
		}

		updateStmt = updateStmt[0 : len(updateStmt)-1] // trim the last ,
		stmt, err := dbm.postgres.Prepare(updateStmt)
		if err != nil {
			fmt.Println(err)
			return
		}
		//format all vals at once
		_, err = stmt.Exec(vals...)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func (dbm *DBManger) saveNetwork(agentId int, stats structs.Stats) {
	var tocId int64

	//
	createNewTocSTMT := `INSERT INTO data.toc_network(fk_connection_id, stamp) VALUES ($1, $2) RETURNING id;`
	err := dbm.postgres.QueryRow(createNewTocSTMT, agentId, time.Now()).Scan(&tocId)
	if err != nil {
		fmt.Println(err)
		return
	}

	sqlStr := "INSERT INTO data.network(laddr, lport, raddr, rport, pid, procname, procowner, fk_network_id) values "

	vals := []interface{}{}

	for i, row := range stats.Connections {
		//prepStr += "(?, ?, ?, ?, ?, ?, ?, (select id from tocid)),"
		sqlStr += fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d),",
			i*8+1, i*8+2, i*8+3, i*8+4, i*8+5, i*8+6, i*8+7, i*8+8)
		vals = append(vals,
			strings.Trim(strings.Replace(fmt.Sprint(row.LAddr), " ", ".", -1), "[]"), row.LPort,
			strings.Trim(strings.Replace(fmt.Sprint(row.RAddr), " ", ".", -1), "[]"), row.RPort,
			row.Pid, row.ProcName,
			row.ProcOwner, tocId)
	}
	sqlStr = sqlStr[0 : len(sqlStr)-1] // trim the last ,
	stmt, err := dbm.postgres.Prepare(sqlStr)
	if err != nil {
		fmt.Println(err)
		return
	}
	//format all vals at once
	_, err = stmt.Exec(vals...)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (dbm *DBManger) GetConnection(id int64) (conn structs.Connection, agentIP string) {

	connectionId := int(id)

	var (
		laddr       string
		lport       int
		raddr       string
		rport       int
		pid         int
		procname    string
		procowner   string
		status      int
		activeSince pq.NullTime
		closedWhen  pq.NullTime
	)

	stmt := `SELECT c.ip, cn.laddr, cn.lport, cn.raddr, cn.rport, cn.pid, cn.procname, cn.procowner, cn.activesince, cn.closedwhen, cn.status
			 FROM data.changes_network cn LEFT JOIN data.connections c ON cn.fk_connection_id = c.id
    		 WHERE cn.id = $1`

	row := dbm.postgres.QueryRow(stmt, connectionId)
	err := row.Scan(&agentIP, &laddr, &lport, &raddr, &rport, &pid, &procname, &procowner, &activeSince, &closedWhen, &status)
	if err != nil {
		fmt.Println(err)
	}

	laddrTmp := strings.Split(laddr, ".")
	laddrResult := make([]int, len(laddrTmp))
	for i, elem := range laddrTmp {
		laddrResult[i], err = strconv.Atoi(elem)
		if err != nil {
			fmt.Println(err)
		}
	}

	raddrTmp := strings.Split(raddr, ".")
	raddrResult := make([]int, len(raddrTmp))
	for i, elem := range raddrTmp {
		raddrResult[i], err = strconv.Atoi(elem)
		if err != nil {
			fmt.Println(err)
		}
	}

	conn = structs.Connection{
		FakeId:      connectionId,
		LAddr:       laddrResult,
		LPort:       lport,
		RAddr:       raddrResult,
		RPort:       rport,
		Pid:         pid,
		ProcName:    procname,
		ProcOwner:   procowner,
		ActiveSince: activeSince.Time,
		ClosedWhen:  closedWhen.Time,
		Status:      status,
	}

	return
}
func (dbm *DBManger) CloseConnection(fakeId interface{}) {
	fmt.Println("DBM::CloseConnection::FakeID:", fakeId)
}

func (dbm *DBManger) GetChangesOfConnectionsByID(agentId int) (connections []structs.Connection) {
	var err error
	var (
		fakeId    int
		laddr     string
		lport     int
		raddr     string
		rport     int
		pid       int
		procname  string
		procowner string
	)

	stmt := `
	SELECT n.id, n.laddr, n.lport, n.raddr, n.rport, n.pid, n.procname, n.procowner
	FROM data.changes_network n
	WHERE n.fk_network_id = (
    	SELECT tn.id FROM data.toc_network tn
    	WHERE tn.fk_connection_id = $1
    	ORDER BY tn.stamp DESC LIMIT 1
	);`

	rows, err := dbm.postgres.Query(stmt, agentId)
	if err != nil {
		fmt.Println(err)
	}
	for rows.Next() {
		err = rows.Scan(&fakeId, &laddr, &lport, &raddr, &rport, &pid, &procname, &procowner)
		if err != nil {
			fmt.Println(err)
		}
		laddrTmp := strings.Split(laddr, ".")
		laddrResult := make([]int, len(laddrTmp))
		for i, elem := range laddrTmp {
			laddrResult[i], err = strconv.Atoi(elem)
			if err != nil {
				fmt.Println(err)
			}
		}

		raddrTmp := strings.Split(raddr, ".")
		raddrResult := make([]int, len(raddrTmp))
		for i, elem := range raddrTmp {
			raddrResult[i], err = strconv.Atoi(elem)
			if err != nil {
				fmt.Println(err)
			}
		}

		connections = append(connections, structs.Connection{
			FakeId:    fakeId,
			LAddr:     laddrResult,
			LPort:     lport,
			RAddr:     raddrResult,
			RPort:     rport,
			Pid:       pid,
			ProcName:  procname,
			ProcOwner: procowner,
		})

	}
	return
}

func (dbm *DBManger) GetAgentConnectionsByID(agentId int) (connections []structs.Connection) {
	var err error
	var (
		fakeId      int
		laddr       string
		lport       int
		raddr       string
		rport       int
		pid         int
		procname    string
		procowner   string
		activeSince time.Time
		closedWhen  pq.NullTime
		status      int
	)

	stmt := `
		SELECT n.id, n.laddr, n.lport, 
		       n.raddr, n.rport, n.pid, 
		       n.procname, n.procowner, 
		       n.activesince, n.closedwhen, n.status
		FROM data.changes_network n
		WHERE n.fk_connection_id = $1 AND n.status > 0;`

	rows, err := dbm.postgres.Query(stmt, agentId)
	if err != nil {
		fmt.Println(err)
	}
	for rows.Next() {
		err = rows.Scan(&fakeId, &laddr, &lport, &raddr, &rport, &pid, &procname, &procowner, &activeSince, &closedWhen, &status)
		if err != nil {
			fmt.Println(err)
		}
		laddrTmp := strings.Split(laddr, ".")
		laddrResult := make([]int, len(laddrTmp))
		for i, elem := range laddrTmp {
			laddrResult[i], err = strconv.Atoi(elem)
			if err != nil {
				fmt.Println(err)
			}
		}

		raddrTmp := strings.Split(raddr, ".")
		raddrResult := make([]int, len(raddrTmp))
		for i, elem := range raddrTmp {
			raddrResult[i], err = strconv.Atoi(elem)
			if err != nil {
				fmt.Println(err)
			}
		}

		connections = append(connections, structs.Connection{
			FakeId:      fakeId,
			LAddr:       laddrResult,
			LPort:       lport,
			RAddr:       raddrResult,
			RPort:       rport,
			Pid:         pid,
			ProcName:    procname,
			ProcOwner:   procowner,
			ActiveSince: activeSince,
			ClosedWhen:  closedWhen.Time,
			Status:      status,
		})

	}
	return
}
