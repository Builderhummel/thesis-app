package db_model

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/Builderhummel/thesis-app/app/config"

	_ "github.com/go-sql-driver/mysql"
)

var Config *config.Configuration

type DBController struct {
	db *sql.DB
}

func (dbc *DBController) OpenConnection() (*sql.DB, error) {
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", Config.DBUsername, Config.DBPassword, Config.DBIP, Config.DBPort, Config.DBName)
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

func (dbc *DBController) GetLoginHandleFromDB(handle string) (string, error) {
	var user string
	err := dbc.db.QueryRow("SELECT LoginHandle FROM Account WHERE LoginHandle = ?", handle).Scan(&user)
	if err != nil {
		return "", err
	}
	return user, nil
}

func (dbc *DBController) ChkUserActive(handle string) (bool, error) {
	var active bool
	err := dbc.db.QueryRow("SELECT Active FROM Account WHERE LoginHandle = ?", handle).Scan(&active)
	if err != nil {
		return false, err
	}
	return active, nil
}

func (dbc *DBController) UpdtUser(handle, name, email string) error {
	_, err := dbc.db.Exec("UPDATE PersonalData pd JOIN Account a ON pd.PDUID = a.PDUID SET pd.Name = ?, pd.Email = ? WHERE a.LoginHandle = ?", name, email, handle)
	if err != nil {
		return err
	}
	return nil
}

func (dbc *DBController) GtUsrPuidFromUserId(user_id string) (string, error) {
	var puid string
	err := dbc.db.QueryRow("SELECT PDUID FROM Account WHERE LoginHandle = ?", user_id).Scan(&puid)
	if err != nil {
		return "", err
	}
	return puid, nil
}

func (dbc *DBController) GtDataTblOpenReq() ([]map[string]string, error) {
	// SQL query to select requested thesis information
	query := `
    SELECT 
        TUID,
        ThesisType,
		ThesisTitle,
        Name,
        DATE_FORMAT(RequestDate, '%Y-%m-%d') AS RequestDate,
        ThesisStatus,
        Email
    FROM 
        Thesis
    WHERE 
        ThesisStatus = 'request'
    ORDER BY 
        RequestDate
    `

	// Execute the query
	rows, err := dbc.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %v", err)
	}
	defer rows.Close()

	// Slice to store the results
	var results []map[string]string

	// Iterate through the rows
	for rows.Next() {
		var (
			tuid, thesisType, thesisTitle, name, requestDate, status, email string
		)

		// Scan the row values
		err := rows.Scan(
			&tuid,
			&thesisType,
			&thesisTitle,
			&name,
			&requestDate,
			&status,
			&email,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}

		//Map results
		row_data := map[string]string{
			"tuid":        tuid,
			"thesisType":  thesisType,
			"thesisTitle": thesisTitle,
			"name":        name,
			"requestDate": requestDate,
			"status":      status,
			"email":       email,
		}

		results = append(results, row_data)
	}

	// Check for any errors encountered during iteration
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during row iteration: %v", err)
	}

	return results, nil
}

func (dbc *DBController) GtDataTblMySupervisions(supervisor_puid string) ([]map[string]string, error) {
	// Main query for thesis information
	query := `
        SELECT DISTINCT
            t.TUID,
            t.ThesisType,
            t.Name,
            t.ThesisTitle,
            COALESCE(t.Deadline, '') as Deadline,
            t.ThesisStatus,
            t.Email
        FROM Thesis t
        JOIN SupervisorJunction sj ON t.TUID = sj.TUID
        JOIN PersonalData pd ON pd.PDUID = ?
        WHERE sj.PDUID = ?`

	rows, err := dbc.db.Query(query, supervisor_puid, supervisor_puid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []map[string]string

	for rows.Next() {
		var tuid, thesisType, name, thesisTitle, deadline, thesisStatus, email string

		err := rows.Scan(&tuid, &thesisType, &name, &thesisTitle, &deadline, &thesisStatus, &email)
		if err != nil {
			return nil, err
		}

		// Get all supervisors for this thesis
		supervisorQuery := `
            SELECT pd.Name
            FROM SupervisorJunction sj
            JOIN PersonalData pd ON sj.PDUID = pd.PDUID
            WHERE sj.TUID = ?`

		supervisorRows, err := dbc.db.Query(supervisorQuery, tuid)
		if err != nil {
			return nil, err
		}
		defer supervisorRows.Close()

		var supervisors []string
		for supervisorRows.Next() {
			var supervisor string
			if err := supervisorRows.Scan(&supervisor); err != nil {
				return nil, err
			}
			supervisors = append(supervisors, supervisor)
		}

		row := map[string]string{
			"tuid":         tuid,
			"thesistype":   thesisType,
			"name":         name,
			"thesistitle":  thesisTitle,
			"deadline":     deadline,
			"supervisor":   strings.Join(supervisors, ", "),
			"thesisstatus": thesisStatus,
			"email":        email,
		}

		result = append(result, row)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (dbc *DBController) GtDataFullSupervision(thesisID string) (*ThesisFullData, error) {
	// Main thesis data query
	mainQuery := `
    SELECT 
        Name, Email, StudyProgram, GPA, ThesisType, ThesisTitle, 
        ThesisStatus,
		COALESCE(FinalGrade, -1) as FinalGrade, 
		CAST(COALESCE(RequestDate, '0001-01-01') AS DATE) as RequestDate, 
		CAST(COALESCE(ContactDate, '0001-01-01') AS DATE) as ContactDate, 
		CAST(COALESCE(Deadline, '0001-01-01') AS DATE) as Deadline, 
		CAST(COALESCE(SubmitDate, '0001-01-01') AS DATE) as SubmitDate, 
		Notes
    FROM Thesis 
    WHERE TUID = ?`

	result := &ThesisFullData{}
	err := dbc.db.QueryRow(mainQuery, thesisID).Scan(
		&result.Name, &result.Email, &result.StudyProgram,
		&result.GPA, &result.ThesisType, &result.ThesisTitle,
		&result.ThesisStatus, &result.FinalGrade, &result.RequestDate,
		&result.ContactDate, &result.Deadline, &result.SubmitDate,
		&result.Notes,
	)
	if err != nil {
		return nil, fmt.Errorf("error fetching thesis data: %v", err)
	}

	// Get supervisors
	supervisorQuery := `
    SELECT pd.Name 
    FROM PersonalData pd
    JOIN SupervisorJunction sj ON pd.PDUID = sj.PDUID
    WHERE sj.TUID = ?`

	supervisorRows, err := dbc.db.Query(supervisorQuery, thesisID)
	if err != nil {
		return nil, fmt.Errorf("error fetching supervisors: %v", err)
	}
	defer supervisorRows.Close()

	for supervisorRows.Next() {
		var name string
		if err := supervisorRows.Scan(&name); err != nil {
			return nil, fmt.Errorf("error scanning supervisor: %v", err)
		}
		result.Supervisors = append(result.Supervisors, name)
	}

	// Get examiners
	examinerQuery := `
    SELECT pd.Name 
    FROM PersonalData pd
    JOIN ExaminerJunction ej ON pd.PDUID = ej.PDUID
    WHERE ej.TUID = ?`

	examinerRows, err := dbc.db.Query(examinerQuery, thesisID)
	if err != nil {
		return nil, fmt.Errorf("error fetching examiners: %v", err)
	}
	defer examinerRows.Close()

	for examinerRows.Next() {
		var name string
		if err := examinerRows.Scan(&name); err != nil {
			return nil, fmt.Errorf("error scanning examiner: %v", err)
		}
		result.Examiners = append(result.Examiners, name)
	}

	return result, nil
}

func (dbc *DBController) InsrtNwThsisRequest(name, email, courseOfStudy, thesisType, thesisTitle, gpa, requestDate, notes string) error {
	_, err := dbc.db.Exec(`
		INSERT INTO Thesis (Name, Email, StudyProgram, ThesisType, ThesisStatus, ThesisTitle, GPA, RequestDate, Notes)
		VALUES (?, ?, ?, ?, 'request', ?, ?, ?, ?)
	`, name, email, courseOfStudy, thesisType, thesisTitle, gpa, requestDate, notes)
	if err != nil {
		return err
	}
	return nil
}
