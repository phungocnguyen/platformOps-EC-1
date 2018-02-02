package services

import (
	"database/sql"
	"fmt"
	"log"
	"platformOps-EC/models"
)

func InsertBaseline(db *sql.DB, baseline models.Baseline) (gen_id int) {
	return insertBaseline(db, baseline)
}

func InsertControl(db *sql.DB, control models.Control) (gen_id int) {
	return insertControl(db, control)
}

func ReadBaselineAll(db *sql.DB) {
	readBaselineAll(db)
}

func ReadControlByBaselineId(db *sql.DB, baseline_id int) {
	readControlByBaselineId(db, baseline_id)
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

func readBaselineById(db *sql.DB, baseline_id int) {

	rows, err := db.Query("SELECT name FROM baseline WHERE id = $1", baseline_id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s is %d\n", name, baseline_id)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

}

func insertBaseline(db *sql.DB, baseline models.Baseline) (gen_id int) {
	sqlStatement := "INSERT INTO baseline (name) VALUES ($1) RETURNING id"
	id := 0
	err := db.QueryRow(sqlStatement, baseline.GetName()).Scan(&id)
	if err != nil {
		panic(err)
	}
	fmt.Println("New record ID is:", id)
	baseline.SetId(id)
	return id
}

func readControlByBaselineId(db *sql.DB, baseline_id int) {
	sqlStatement := `SELECT id, req_id, cis_id, category,
                    requirement, discussion, check_text,
                    fix_text, row_desc, baseline_id
                    FROM control WHERE baseline_id=$1;`

	rows, err := db.Query(sqlStatement, baseline_id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var cis_id, category, requirement, discussion, check_text, fix_text, row_desc string
		var id, req_id, baseline_id int
		if err := rows.Scan(&id, &req_id, &cis_id, &category, &requirement, &discussion, &check_text, &fix_text, &row_desc, &baseline_id); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("result:  %d, %v, %v, %v, %v, %v, %v, %v\n",
			id, category, requirement, discussion, check_text, fix_text, row_desc, baseline_id)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

}

func insertControl(db *sql.DB, control models.Control) (gen_id int) {
	sqlStatement := `INSERT INTO control
                    (req_id, cis_id, category, requirement,
                    discussion, check_text, fix_text, row_desc, baseline_id)
                    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id;`
	id := 0
	err := db.QueryRow(sqlStatement, control.Req_id, control.Cis_id,
		control.Category, control.Requirement, control.Discussion,
		control.Check_text, control.Fix_text, control.Row_desc,
		control.Baseline_id).Scan(&id)
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
	return models.Control{Req_id: 2, Cis_id: "N/A", Category: "Test Category",
		Requirement: "Test Requirement", Discussion: "Test Discussion",
		Check_text: "Test Check_text", Fix_text: "Test Fix_text",
		Row_desc: "Test Row Desc", Baseline_id: 1}

}