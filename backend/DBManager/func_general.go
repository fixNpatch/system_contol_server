package DBManager

import "fmt"

func (dbm *DBManger) ClearAllData() {
	clearNetwork := `TRUNCATE data.network; SELECT SETVAL('data.network_id_seq', 1);`
	clearTocNetwork := `TRUNCATE data.toc_network cascade; SELECT SETVAL('data.toc_network_id_seq', 1);`
	clearState := `TRUNCATE data.state; SELECT setval('data.state_id_seq', 1);`
	clearNetworkChanges := `TRUNCATE data.changes_network; SELECT setval('data.changes_network_id_seq', 1);`

	if res, err := dbm.postgres.Exec(clearNetwork); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(res)
	}

	if res, err := dbm.postgres.Exec(clearTocNetwork); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(res)
	}

	if res, err := dbm.postgres.Exec(clearState); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(res)
	}

	if res, err := dbm.postgres.Exec(clearNetworkChanges); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(res)
	}

}
