package services

import (
	"database/sql"
	"fmt"
	"log"
	"platformOps-EC/models"
)

func InsertBaseline(db *sql.DB, baseline models.Baseline) (genId int) {
	return insertBaseline(db, baseline)
}

func InsertControl(db *sql.DB, control models.Control) (genId int) {
	return insertControl(db, control)
}

func ReadBaselineAll(db *sql.DB) {
	readBaselineAll(db)
}

func ReadControlByBaselineId(db *sql.DB, baselineId int) {
	readControlByBaselineId(db, baselineId)
}

func SetSearchPath(db *sql.DB, schema string) {

	sqlStatement := "SET search_path TO " + schema

	_, err := db.Exec(sqlStatement)
	if err != nil {
		panic(err)
	}
}

func readBaselineAll(db *sql.DB) {
	rows, err := db.Query("SELECT name, id FROM baseline")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		var id int
		if err := rows.Scan(&name, &id); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s is %d\n", name, id)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

}

func readBaselineById(db *sql.DB, baselineId int) {

	rows, err := db.Query("SELECT name FROM baseline WHERE id = $1", baselineId)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s is %d\n", name, baselineId)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

}

func insertBaseline(db *sql.DB, baseline models.Baseline) (genId int) {
	sqlStatement := `INSERT INTO baseline (name, baseline_uuid) 
					VALUES ($1, $2) RETURNING id`
	id := 0
	err := db.QueryRow(sqlStatement, baseline.Name, baseline.Uid).Scan(&id)
	if err != nil {
		panic(err)
	}
	fmt.Println("New record ID is:", id)
	baseline.Id = id
	return id
}

func GetManifestByBaselineId(db *sql.DB, baselineUid string) []models.ECManifest {
	SetSearchPath(db, "baseline")

	sqlStatement := `SELECT c.req_id, c.category, b.name, c.baseline_uuid, c.control_uuid
                    FROM control c, baseline b WHERE c.baseline_uuid=b.baseline_uuid AND b.baseline_uuid=$1;`

	rows, err := db.Query(sqlStatement, baselineUid)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var manifests []models.ECManifest
	for rows.Next() {
		var category, baselineName, controlUid string
		var reqId int
		if err := rows.Scan(&reqId, &category, &baselineName, &baselineUid, &controlUid); err != nil {
			log.Fatal(err)
		}

		commands := GetCommandByControlId(db, controlUid)
		manifest := models.ECManifest{reqId, category, getCommandStringArray(commands), baselineUid, controlUid}
		manifests = append(manifests, manifest)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return manifests
}

func getCommandStringArray(commands []models.Command) []string {
	var cmds []string
	for i := range commands {
		cmds = append(cmds, commands[i].Cmd)
	}

	return cmds
}

func GetCommandByControlId(db *sql.DB, controlId string) []models.Command {
	SetSearchPath(db, "baseline")
	sqlStatement := `SELECT id, cmd, exec_order
					 	FROM  command  
						WHERE control_uuid = $1 ORDER BY exec_order ASC;`

	rows, err := db.Query(sqlStatement, controlId)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var commands []models.Command

	for rows.Next() {
		var cmd string
		var id, exeOrder int
		if err := rows.Scan(&id, &cmd, &exeOrder); err != nil {
			log.Fatal(err)
		}
		command := models.Command{Id: id, Cmd: cmd, ExeOrder: exeOrder, ControlId: controlId}

		commands = append(commands, command)

	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return commands

}

func readControlByBaselineId(db *sql.DB, baselineId int) {
	sqlStatement := `SELECT id, req_id, cis_id, category,
                    requirement, discussion, check_text,
                    fix_text, row_desc, baselineId
                    FROM control WHERE baselineId=$1;`

	rows, err := db.Query(sqlStatement, baselineId)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var cisId, category, requirement, discussion, checkText, fixText, rowDesc string
		var id, reqId, baselineId int
		if err := rows.Scan(&id, &reqId, &cisId, &category, &requirement, &discussion, &checkText, &fixText, &rowDesc, &baselineId); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("result:  %d, %v, %v, %v, %v, %v, %v, %v\n",
			id, category, requirement, discussion, checkText, fixText, rowDesc, baselineId)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

}

func GetBaselineIdByName(db *sql.DB, name string) (baselineId int) {
	var id int
	sqlStatement := `SELECT id 
					FROM baseline 
					WHERE name=$1;`
	err := db.QueryRow(sqlStatement, name).Scan(&id)
	switch {
	case err == sql.ErrNoRows:
		log.Printf("No user with that ID.")
	case err != nil:
		log.Fatal(err)
	default:
		fmt.Printf("Searched baseline name %s has Id %s\n", name, id)
	}
	return id
}

func GetECResultById(db *sql.DB, batchUid string) models.BatchSubmision {
	var resultSubmit, dateSubmit, timeSubmit, userSubmit string
	var id int
	sqlStatement := `SELECT id, to_char(date_submit, 'YYYY-MM-DD') as date_submit, to_char(time_submit, 'HH24:MI:SS') as time_submit, user_submit, ec_result 
					FROM batch_submission 
					WHERE batch_uuid=$1;`
	err := db.QueryRow(sqlStatement, batchUid).Scan(&id, &dateSubmit, &timeSubmit, &userSubmit, &resultSubmit)
	switch {
	case err == sql.ErrNoRows:
		log.Printf("No ec result with that ID.")
	case err != nil:
		log.Fatal(err)
	default:
		fmt.Printf("Searched ec result  Id %v\n", batchUid)
	}

	var ecResults []models.ECResult

	models.ToObject(resultSubmit, &ecResults)

	return models.BatchSubmision{id, batchUid, dateSubmit, timeSubmit, userSubmit, ecResults}
}

func SaveECResult(db *sql.DB, ecResults []models.ECResult, submitDate string, submitUser string) (status int, batchUid string) {
	id, batchUid := saveBatchSubmission(db, ecResults, submitDate, submitUser)

	if id <= 0 {
		return 0, "Error posting batch to master"
	}

	sqlStatement := `INSERT INTO exec_result 
					(baseline_uuid, host_exec, exec_date, exec_time, commands, std_err, std_output, batch_uuid, control_uuid)
					VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`
	id = 0
	dateTime := ecResults[0].DateExe
	err := db.QueryRow(sqlStatement, ecResults[0].BaselineUid, ecResults[0].HostExec, dateTime, GetTimeZoneString(dateTime), models.ToJson(ecResults[0].Command), models.ToJson(ecResults[0].StdErrOutput), models.ToJson(ecResults[0].StdOutput), batchUid, ecResults[0].ControlUid).Scan(&id)

	if err != nil {
		log.Print(err)
	}

	return id, batchUid

}

func saveBatchSubmission(db *sql.DB, ecResults []models.ECResult, submitDate string, submitUser string) (id int, batchUid string) {
	sqlStatement := `INSERT INTO batch_submission 
					(batch_uuid, user_submit, date_submit, time_submit, ec_result)
					VALUES ($1, $2, $3, $4, $5) RETURNING id`

	id = 0
	batchUid = NewUUID()
	err := db.QueryRow(sqlStatement, batchUid, submitUser, submitDate, GetTimeZoneString(submitDate), models.ToJson(ecResults)).Scan(&id)

	if err != nil {
		log.Print(err)
	}

	return id, batchUid
}

func insertControl(db *sql.DB, control models.Control) (genId int) {
	sqlStatement := `INSERT INTO control
                    (req_id, cis_id, category, requirement,
                    discussion, check_text, fix_text, row_desc, baseline_id)
                    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id;`
	id := 0
	err := db.QueryRow(sqlStatement, control.ReqId, control.CisId,
		control.Category, control.Requirement, control.Discussion,
		control.CheckText, control.FixText, control.RowDesc,
		control.BaselineUid).Scan(&id)
	if err != nil {
		panic(err)
	}
	control.Id = id
	fmt.Println("New record ID is:", id)
	return id
}

func deleteControl(db *sql.DB) {
	id := 3
	sqlStatement := `
					DELETE FROM baseline.control
					WHERE id = $1;`
	_, err := db.Exec(sqlStatement, id)
	if err != nil {
		panic(err)
	}

}
