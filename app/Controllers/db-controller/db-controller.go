package dbcontroller

import (
	"database/sql"
	"fmt"

	"github.com/Builderhummel/thesis-app/app/config"

	_ "github.com/go-sql-driver/mysql"
)

var Config *config.Configuration

type DBController struct {
	db *sql.DB
}

func (dbc *DBController) OpenConnection() (*sql.DB, error) {
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", Config.DBUsername, Config.DBPassword, Config.DBIP, Config.DBPort, Config.DBName)
	fmt.Println(dsn)
	dbc.db, err = sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	return dbc.db, nil
}

func (dbc *DBController) CloseConnection() error {
	return dbc.db.Close()
}

func (dbc *DBController) InitDatabase() error {
	fmt.Println("Try to init table Thesis")
	_, err := dbc.db.Exec(`
		CREATE TABLE Thesis (
			TUID INT AUTO_INCREMENT PRIMARY KEY,
			Name VARCHAR(255),
			Email VARCHAR(255),
			StudyProgram VARCHAR(255),
			ThesisType VARCHAR(255),
			ThesisStatus VARCHAR(255),
			ThesisTitle VARCHAR(255),
			GPA FLOAT,
			RequestDate DATE,
			ContactDate DATE,
			SubmitDate DATE,
			Deadline DATE,
			FinalGrade FLOAT,
			Notes TEXT
		)
	`)
	if err != nil {
		return err
	}

	_, err = dbc.db.Exec(`
		CREATE TABLE PersonalData (
			PDUID INT AUTO_INCREMENT PRIMARY KEY,
			Name VARCHAR(255),
			Email VARCHAR(255)
		)
	`)
	if err != nil {
		return err
	}

	_, err = dbc.db.Exec(`
		CREATE TABLE Account (
			AUID INT AUTO_INCREMENT PRIMARY KEY,
			PDUID INT,
			LoginHandle VARCHAR(255),
			Active BOOLEAN,
			FOREIGN KEY (PDUID) REFERENCES PersonalData(PDUID)
		)
	`)
	if err != nil {
		return err
	}

	_, err = dbc.db.Exec(`
		CREATE TABLE SupervisorJunction (
			TUID INT,
			PDUID INT,
			FOREIGN KEY (TUID) REFERENCES Thesis(TUID),
			FOREIGN KEY (PDUID) REFERENCES PersonalData(PDUID)
		)
	`)
	if err != nil {
		return err
	}

	_, err = dbc.db.Exec(`
		CREATE TABLE ExaminerJunction (
			TUID INT,
			PDUID INT,
			FOREIGN KEY (TUID) REFERENCES Thesis(TUID),
			FOREIGN KEY (PDUID) REFERENCES PersonalData(PDUID)
		)
	`)
	if err != nil {
		return err
	}

	return nil
}

func (dbc *DBController) CheckIfDatabaseIsInitialized() (bool, error) {
	result, err := dbc.db.Query("SHOW TABLES LIKE 'Account'")
	if err != nil {
		return false, err
	}
	defer result.Close()
	for result.Next() {
		var tableName string
		if err := result.Scan(&tableName); err != nil {
			return false, err
		}
		if tableName == "Account" {
			return true, nil
		}
	}
	return false, nil
}
