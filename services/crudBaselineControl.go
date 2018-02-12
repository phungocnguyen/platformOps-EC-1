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
	sqlStatement := "INSERT INTO baseline (name) VALUES ($1) RETURNING id"
	id := 0
	err := db.QueryRow(sqlStatement, baseline.Name).Scan(&id)
	if err != nil {
		panic(err)
	}
	fmt.Println("New record ID is:", id)
	baseline.SetId(id)
	return id
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

func insertControl(db *sql.DB, control models.Control) (genId int) {
	sqlStatement := `INSERT INTO control
                    (req_id, cis_id, category, requirement,
                    discussion, check_text, fix_text, row_desc, baseline_id)
                    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id;`
	id := 0
	err := db.QueryRow(sqlStatement, control.ReqId, control.CisId,
		control.Category, control.Requirement, control.Discussion,
		control.CheckText, control.FixText, control.RowDesc,
		control.BaselineId).Scan(&id)
	if err != nil {
		panic(err)
	}
	control.SetId(id)
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

func populateControl() (control models.Control) {
	return models.Control{ReqId: 2, CisId: "N/A", Category: "Test Category",
		Requirement: "Test Requirement", Discussion: "Test Discussion",
		CheckText: "Test CheckText", FixText: "Test FixText",
		RowDesc: "Test Row Desc", BaselineId: 1}

}