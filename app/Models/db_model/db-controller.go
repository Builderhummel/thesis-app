package db_model

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Builderhummel/thesis-app/app/Constants/roles"
	"github.com/Builderhummel/thesis-app/app/config"

	_ "github.com/go-sql-driver/mysql"
)

var Config *config.Configuration

type DBController struct {
	db *sql.DB
}

func (dbc *DBController) OpenConnection() error {
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", Config.DBUsername, Config.DBPassword, Config.DBIP, Config.DBPort, Config.DBName)
	dbc.db, err = sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	return nil
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
			Booked BOOLEAN,
			ThesisType VARCHAR(255),
			ThesisStatus VARCHAR(255),
			Semester VARCHAR(255),
			ThesisTitle VARCHAR(255),
			GPA FLOAT,
			RequestDate DATE,
			ResponseDate DATE,
			RegisteredDate DATE,
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
			Email VARCHAR(255),
			IsSupervisor BOOLEAN DEFAULT FALSE,
			IsExaminer BOOLEAN DEFAULT FALSE
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
			Role VARCHAR(255),
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

	_, err = dbc.db.Exec(`
		CREATE TABLE ThesisFiles (
			FUID INT AUTO_INCREMENT PRIMARY KEY,
			TUID INT NOT NULL,
			FileName VARCHAR(255) NOT NULL,
			OriginalFileName VARCHAR(255) NOT NULL,
			FileSize BIGINT NOT NULL,
			UploadDate TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			PDUID INT,
			FOREIGN KEY (TUID) REFERENCES Thesis(TUID) ON DELETE CASCADE,
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

func (dbc *DBController) GtUsrRleByLgnHndle(userid string) (roles.Role, error) {
	var roleStr string
	err := dbc.db.QueryRow("SELECT Role FROM Account WHERE LoginHandle = ?", userid).Scan(&roleStr)
	if err != nil {
		return "", err
	}
	role := roles.Role(roleStr)
	return role, nil
}

func (dbc *DBController) ChkUserActive(handle string) (bool, error) {
	var active bool
	err := dbc.db.QueryRow("SELECT Active FROM Account WHERE LoginHandle = ?", handle).Scan(&active)
	if err != nil {
		return false, err
	}
	return active, nil
}

// TODO: Refactor (only used when logging in)
func (dbc *DBController) UpdtUser(handle, name, email string) error {
	_, err := dbc.db.Exec("UPDATE PersonalData pd JOIN Account a ON pd.PDUID = a.PDUID SET pd.Name = ?, pd.Email = ? WHERE a.LoginHandle = ?", name, email, handle)
	if err != nil {
		return err
	}
	return nil
}

func (dbc *DBController) GtAllUsrs() ([]PersonalData, error) {
	query := `
	SELECT 
		COALESCE(pd.PDUID, ''),
		COALESCE(pd.Name, ''),
		COALESCE(pd.Email, ''),
		COALESCE(a.Role, ''),
		COALESCE(a.LoginHandle, ''),
		COALESCE(a.Active, FALSE),
		COALESCE(pd.IsSupervisor, FALSE),
		COALESCE(pd.IsExaminer, FALSE)
    FROM 
        PersonalData pd
    LEFT JOIN 
        Account a ON pd.PDUID = a.PDUID
    `

	rows, err := dbc.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %v", err)
	}
	defer rows.Close()

	var users []PersonalData
	for rows.Next() {
		var user PersonalData
		err := rows.Scan(
			&user.PDUid,
			&user.Name,
			&user.Email,
			&user.Role,
			&user.Handle,
			&user.IsActive,
			&user.IsSupervisor,
			&user.IsExaminer,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user row: %v", err)
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating user rows: %v", err)
	}

	return users, nil
}

func (dbc *DBController) GtUsrByPUID(puid string) (PersonalData, error) {
	query := `
	SELECT 
		COALESCE(pd.PDUID, ''),
		COALESCE(pd.Name, ''),
		COALESCE(pd.Email, ''),
		COALESCE(a.Role, ''),
		COALESCE(a.LoginHandle, ''),
		COALESCE(a.Active, FALSE),
		COALESCE(pd.IsSupervisor, FALSE),
		COALESCE(pd.IsExaminer, FALSE)
    FROM 
        PersonalData pd
    LEFT JOIN 
        Account a ON pd.PDUID = a.PDUID
	WHERE
		pd.PDUID = ?
    `

	var user PersonalData
	err := dbc.db.QueryRow(query, puid).Scan(&user.PDUid, &user.Name, &user.Email, &user.Role, &user.Handle, &user.IsActive, &user.IsSupervisor, &user.IsExaminer)
	if err != nil {
		return PersonalData{}, fmt.Errorf("error fetching user: %v", err)
	}

	return user, nil
}

func (dbc *DBController) GtAllSupervisors() ([]PersonalData, error) {
	query := "SELECT PDUID, Name, Email FROM PersonalData WHERE IsSupervisor = 1"

	rows, err := dbc.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var supervisors []PersonalData

	for rows.Next() {
		var supervisor PersonalData
		if err := rows.Scan(&supervisor.PDUid, &supervisor.Name, &supervisor.Email); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		supervisors = append(supervisors, supervisor)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return supervisors, nil
}

func (dbc *DBController) GtAllExaminers() ([]PersonalData, error) {
	query := "SELECT PDUID, Name, Email FROM PersonalData WHERE IsExaminer = 1"

	rows, err := dbc.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var examiners []PersonalData

	for rows.Next() {
		var examiner PersonalData
		if err := rows.Scan(&examiner.PDUid, &examiner.Name, &examiner.Email); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		examiners = append(examiners, examiner)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return examiners, nil
}

func (dbc *DBController) GtUsrPuidFromUserId(user_id string) (string, error) {
	var puid string
	err := dbc.db.QueryRow("SELECT PDUID FROM Account WHERE LoginHandle = ?", user_id).Scan(&puid)
	if err != nil {
		return "", err
	}
	return puid, nil
}

func (dbc *DBController) InsrtNwUsr(name, email, loginHandle, role string, active, isSupervisor, isExaminer bool) error {
	tx, err := dbc.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	query1 := `
		INSERT INTO PersonalData (Name, Email, IsSupervisor, IsExaminer)
		VALUES (?, ?, ?, ?);
	`

	res, err := tx.Exec(query1, name, email, isSupervisor, isExaminer)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("failed to insert into PersonalData: %w", err)
	}

	lastInsertID, err := res.LastInsertId()
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("failed to retrieve last insert ID: %w", err)
	}

	query2 := `
		INSERT INTO Account (PDUID, LoginHandle, Role, Active)
		VALUES (?, ?, ?, ?);
	`

	_, err = tx.Exec(query2, lastInsertID, loginHandle, role, active)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("failed to insert into Account: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (dbc *DBController) UptFullUsr(puid, name, email, loginHandle, role string, active, isSupervisor, isExaminer bool) error {
	tx, err := dbc.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	query1 := `
		UPDATE PersonalData
		SET 
			Name = ?,
			Email = ?,
			IsSupervisor = ?,
			IsExaminer = ?
		WHERE 
			PDUID = ?;
	`

	_, err = tx.Exec(query1, name, email, isSupervisor, isExaminer, puid)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("failed to update PersonalData: %w", err)
	}

	query2 := `
		UPDATE Account
		SET 
			LoginHandle = ?,
			Role = ?,
			Active = ?
		WHERE 
			PDUID = ?;
	`

	_, err = tx.Exec(query2, loginHandle, role, active, puid)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("failed to update Account: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// Get Amount of request, contacted(=contacted, registered), working
func (dbc *DBController) GtHomepageRCW() (map[string]string, error) {
	// SQL query to select requested thesis information
	query := `
	SELECT 
		COUNT(ThesisStatus = 'request' OR NULL) AS Requested,
		COUNT(ThesisStatus = 'contacted' OR NULL) AS Contacted,
		COUNT(ThesisStatus = 'registered' OR NULL) AS Registered,
		COUNT(ThesisStatus = 'working' OR NULL) AS Working
	FROM 
		Thesis;
	`

	// Execute the query
	rows, err := dbc.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %v", err)
	}
	defer rows.Close()

	// Slice to store the results
	var results map[string]string

	// Iterate through the rows
	for rows.Next() {
		var (
			requested, contacted, registered, working string
		)

		// Scan the row values
		err := rows.Scan(
			&requested,
			&contacted,
			&registered,
			&working,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}

		intContacted, err := strconv.Atoi(contacted)
		if err != nil {
			return nil, fmt.Errorf("error converting contacted to int: %v", err)
		}

		intRegistered, err := strconv.Atoi(registered)
		if err != nil {
			return nil, fmt.Errorf("error converting registered to int: %v", err)
		}

		newContacted := strconv.Itoa(intContacted + intRegistered)

		//Map results
		results = map[string]string{
			"requested": requested,
			"contacted": newContacted,
			"working":   working,
		}
	}

	// Check for any errors encountered during iteration
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during row iteration: %v", err)
	}

	return results, nil
}

func (dbc *DBController) GtDataTblOpenReq() ([]map[string]string, error) {
	// SQL query to select requested thesis information
	query := `
    SELECT 
        TUID,
        ThesisType,
		studyProgram,
        Name,
        DATE_FORMAT(RequestDate, '%Y-%m-%d') AS RequestDate,
        ThesisStatus,
        Email,
		COALESCE(GPA, '-1') as GPA
    FROM 
        Thesis
    WHERE 
        ThesisStatus = 'request'
    ORDER BY 
        RequestDate ASC;
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
			tuid, thesisType, courseOfStudy, name, requestDate, status, email string
			gpa                                                               float64
		)

		// Scan the row values
		err := rows.Scan(
			&tuid,
			&thesisType,
			&courseOfStudy,
			&name,
			&requestDate,
			&status,
			&email,
			&gpa,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}

		//Map results
		row_data := map[string]string{
			"tuid":          tuid,
			"thesisType":    thesisType,
			"courseOfStudy": courseOfStudy,
			"name":          name,
			"requestDate":   requestDate,
			"status":        status,
			"email":         email,
			"gpa":           fmt.Sprintf("%.2f", gpa),
		}

		results = append(results, row_data)
	}

	// Check for any errors encountered during iteration
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during row iteration: %v", err)
	}

	return results, nil
}

func (dbc *DBController) GtDataTblAllSupervisions() ([]map[string]string, error) {
	// SQL query to select all thesis information
	query := `
		SELECT 
			TUID,
			ThesisType,
			ThesisTitle,
			Name,
			COALESCE(DATE_FORMAT(Deadline, '%Y-%m-%d'), '') AS Deadline,
			ThesisStatus,
			Email,
			COALESCE(Semester, '') as Semester
		FROM 
			Thesis
		ORDER BY 
			RequestDate DESC;

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
			tuid, thesisType, thesisTitle, name, deadline, status, email, semester string
		)

		// Scan the row values
		err := rows.Scan(
			&tuid,
			&thesisType,
			&thesisTitle,
			&name,
			&deadline,
			&status,
			&email,
			&semester,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}

		//Get supervisors for this thesis
		supervisorQuery := `
			SELECT pd.Name
			FROM SupervisorJunction sj
			JOIN PersonalData pd ON sj.PDUID = pd.PDUID
			WHERE sj.TUID = ?
		`

		supervisorRows, err := dbc.db.Query(supervisorQuery, tuid)
		if err != nil {
			return nil, fmt.Errorf("error fetching supervisors: %v", err)
		}
		defer supervisorRows.Close()

		var supervisors []string
		for supervisorRows.Next() {
			var supervisor string
			if err := supervisorRows.Scan(&supervisor); err != nil {
				return nil, fmt.Errorf("error scanning supervisor: %v", err)
			}
			supervisors = append(supervisors, supervisor)
		}

		//Map results
		row_data := map[string]string{
			"tuid":        tuid,
			"thesisType":  thesisType,
			"thesisTitle": thesisTitle,
			"name":        name,
			"deadline":    deadline,
			"status":      status,
			"email":       email,
			"supervisor":  strings.Join(supervisors, ", "),
			"semester":    semester,
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
            t.Email,
			COALESCE(t.Semester, '') as Semester
        FROM Thesis t
        JOIN SupervisorJunction sj ON t.TUID = sj.TUID
        JOIN PersonalData pd ON pd.PDUID = ?
        WHERE sj.PDUID = ?
		ORDER BY Deadline DESC`

	rows, err := dbc.db.Query(query, supervisor_puid, supervisor_puid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []map[string]string

	for rows.Next() {
		var tuid, thesisType, name, thesisTitle, deadline, thesisStatus, email, semester string

		err := rows.Scan(&tuid, &thesisType, &name, &thesisTitle, &deadline, &thesisStatus, &email, &semester)
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
			"semester":     semester,
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
        TUID, Name, Email, StudyProgram,
		COALESCE(Booked, FALSE) AS Booked,
		COALESCE(GPA, -1) as GPA, 
		ThesisType, ThesisTitle, 
        ThesisStatus,
		COALESCE(Semester, '') as Semester,
		COALESCE(FinalGrade, -1) as FinalGrade, 
		CAST(COALESCE(RequestDate, '0001-01-01') AS DATE) as RequestDate, 
		CAST(COALESCE(ResponseDate, '0001-01-01') AS DATE) as ResponseDate, 
		CAST(COALESCE(RegisteredDate, '0001-01-01') AS DATE) as RegisteredDate, 
		CAST(COALESCE(Deadline, '0001-01-01') AS DATE) as Deadline, 
		CAST(COALESCE(SubmitDate, '0001-01-01') AS DATE) as SubmitDate, 
		Notes
    FROM Thesis 
    WHERE TUID = ?`

	result := &ThesisFullData{}
	err := dbc.db.QueryRow(mainQuery, thesisID).Scan(
		&result.TUID, &result.Name, &result.Email, &result.StudyProgram, &result.Booked,
		&result.GPA, &result.ThesisType, &result.ThesisTitle,
		&result.ThesisStatus, &result.Semester, &result.FinalGrade, &result.RequestDate,
		&result.ResponseDate, &result.RegisteredDate, &result.Deadline, &result.SubmitDate,
		&result.Notes,
	)
	if err != nil {
		return nil, fmt.Errorf("error fetching thesis data: %v", err)
	}

	// Get supervisors
	supervisorQuery := `
	SELECT pd.PDUID, pd.Name, pd.Email
	FROM PersonalData pd
	JOIN SupervisorJunction sj ON pd.PDUID = sj.PDUID
	WHERE sj.TUID = ?`
	supervisorRows, err := dbc.db.Query(supervisorQuery, thesisID)
	if err != nil {
		return nil, fmt.Errorf("error fetching supervisors: %v", err)
	}
	defer supervisorRows.Close()

	var supervisors []PersonalData
	for supervisorRows.Next() {
		var supervisor PersonalData
		if err := supervisorRows.Scan(&supervisor.PDUid, &supervisor.Name, &supervisor.Email); err != nil {
			return nil, fmt.Errorf("error scanning supervisor: %v", err)
		}
		supervisors = append(supervisors, supervisor)
	}
	result.Supervisors = supervisors

	// Get examiners
	examinerQuery := `
	SELECT pd.PDUID, pd.Name, pd.Email
	FROM PersonalData pd
	JOIN ExaminerJunction ej ON pd.PDUID = ej.PDUID
	WHERE ej.TUID = ?`
	examinerRows, err := dbc.db.Query(examinerQuery, thesisID)
	if err != nil {
		return nil, fmt.Errorf("error fetching examiners: %v", err)
	}
	defer examinerRows.Close()

	var examiners []PersonalData
	for examinerRows.Next() {
		var examiner PersonalData
		if err := examinerRows.Scan(&examiner.PDUid, &examiner.Name, &examiner.Email); err != nil {
			return nil, fmt.Errorf("error scanning examiner: %v", err)
		}
		examiners = append(examiners, examiner)
	}
	result.Examiners = examiners

	return result, nil
}

func (dbc *DBController) GtAllDataFullSupervision() ([]*ThesisFullData, error) {
	// 1. Query all thesis main data
	mainQuery := `
		SELECT 
			TUID, Name, Email, StudyProgram,
			COALESCE(Booked, FALSE) AS Booked,
			COALESCE(GPA, -1) as GPA, 
			ThesisType, ThesisTitle, 
			ThesisStatus,
			COALESCE(Semester, '') as Semester,
			COALESCE(FinalGrade, -1) as FinalGrade, 
			CAST(COALESCE(RequestDate, '0001-01-01') AS DATE) as RequestDate, 
			CAST(COALESCE(ResponseDate, '0001-01-01') AS DATE) as ResponseDate, 
			CAST(COALESCE(RegisteredDate, '0001-01-01') AS DATE) as RegisteredDate, 
			CAST(COALESCE(Deadline, '0001-01-01') AS DATE) as Deadline, 
			CAST(COALESCE(SubmitDate, '0001-01-01') AS DATE) as SubmitDate, 
			Notes
		FROM Thesis
	`
	rows, err := dbc.db.Query(mainQuery)
	if err != nil {
		return nil, fmt.Errorf("error fetching thesis data: %v", err)
	}
	defer rows.Close()

	// Use map for quick lookups by TUID
	thesisMap := make(map[string]*ThesisFullData)
	var thesisList []*ThesisFullData

	for rows.Next() {
		result := &ThesisFullData{}
		err := rows.Scan(
			&result.TUID, &result.Name, &result.Email, &result.StudyProgram, &result.Booked,
			&result.GPA, &result.ThesisType, &result.ThesisTitle,
			&result.ThesisStatus, &result.Semester, &result.FinalGrade, &result.RequestDate,
			&result.ResponseDate, &result.RegisteredDate, &result.Deadline, &result.SubmitDate,
			&result.Notes,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning thesis: %v", err)
		}
		thesisMap[result.TUID] = result
		thesisList = append(thesisList, result)
	}

	// 2. Query all supervisors for all theses
	supervisorQuery := `
		SELECT sj.TUID, pd.PDUID, pd.Name, pd.Email
		FROM SupervisorJunction sj
		JOIN PersonalData pd ON pd.PDUID = sj.PDUID
	`
	supRows, err := dbc.db.Query(supervisorQuery)
	if err != nil {
		return nil, fmt.Errorf("error fetching supervisors: %v", err)
	}
	defer supRows.Close()

	for supRows.Next() {
		var tuid, pduid, name, email string
		if err := supRows.Scan(&tuid, &pduid, &name, &email); err != nil {
			return nil, fmt.Errorf("error scanning supervisor: %v", err)
		}
		if thesis, ok := thesisMap[tuid]; ok {
			thesis.Supervisors = append(thesis.Supervisors, PersonalData{
				PDUid: pduid,
				Name:  name,
				Email: email,
			})
		}
	}

	// 3. Query all examiners for all theses
	examinerQuery := `
		SELECT ej.TUID, pd.PDUID, pd.Name, pd.Email
		FROM ExaminerJunction ej
		JOIN PersonalData pd ON pd.PDUID = ej.PDUID
	`
	examRows, err := dbc.db.Query(examinerQuery)
	if err != nil {
		return nil, fmt.Errorf("error fetching examiners: %v", err)
	}
	defer examRows.Close()

	for examRows.Next() {
		var tuid, pduid, name, email string
		if err := examRows.Scan(&tuid, &pduid, &name, &email); err != nil {
			return nil, fmt.Errorf("error scanning examiner: %v", err)
		}
		if thesis, ok := thesisMap[tuid]; ok {
			thesis.Examiners = append(thesis.Examiners, PersonalData{
				PDUid: pduid,
				Name:  name,
				Email: email,
			})
		}
	}

	return thesisList, nil
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

func (dbc *DBController) UpdtThesisInfo(td *ThesisFullData) error {
	tx, err := dbc.db.Begin()
	if err != nil {
		return fmt.Errorf("begin transaction: %v", err)
	}
	defer tx.Rollback()

	_, err = tx.Exec(`
        UPDATE Thesis 
        SET Name=?, Email=?, StudyProgram=?, Booked=?, 
            ThesisType=?, ThesisStatus=?, Semester=?, 
            ThesisTitle=?, GPA=?, RequestDate=?, ResponseDate=?, RegisteredDate=?,
            SubmitDate=?, Deadline=?, FinalGrade=?, Notes=?
        WHERE TUID = ?`,
		td.Name, td.Email, td.StudyProgram, td.Booked,
		td.ThesisType, td.ThesisStatus, td.Semester,
		td.ThesisTitle,
		convertGradeToNullFloat(td.GPA),
		convertGoDateToSqlNullDate(td.RequestDate),
		convertGoDateToSqlNullDate(td.ResponseDate),
		convertGoDateToSqlNullDate(td.RegisteredDate),
		convertGoDateToSqlNullDate(td.SubmitDate),
		convertGoDateToSqlNullDate(td.Deadline),
		convertGradeToNullFloat(td.FinalGrade),
		td.Notes,
		td.TUID)
	if err != nil {
		return fmt.Errorf("update thesis: %v", err)
	}

	if err := dbc.updtJunction(tx, td.TUID, td.Supervisors, "SupervisorJunction"); err != nil {
		return fmt.Errorf("update supervisors: %v", err)
	}

	if err := dbc.updtJunction(tx, td.TUID, td.Examiners, "ExaminerJunction"); err != nil {
		return fmt.Errorf("update examiners: %v", err)
	}

	err = tx.Commit()

	return err
}

func (dbc *DBController) DelThesisRequest(thesisID string) error {
	tx, err := dbc.db.Begin()
	if err != nil {
		return err
	}
	defer func() error {
		if p := recover(); p != nil {
			tx.Rollback()
			return fmt.Errorf("transaction rollback due to panic: %v", p)
		}
		return nil
	}()

	// Delete from junction tables first (to avoid foreign key errors)
	_, err = tx.Exec("DELETE FROM SupervisorJunction WHERE TUID = ?", thesisID)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec("DELETE FROM ExaminerJunction WHERE TUID = ?", thesisID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// If you have other relations referencing Thesis by TUID, add their deletion here

	// Delete the thesis itself
	_, err = tx.Exec("DELETE FROM Thesis WHERE TUID = ?", thesisID)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (dbc *DBController) ChkIfThesisIsBooked(thesisID string) (bool, error) {
	var booked bool
	err := dbc.db.QueryRow("SELECT Booked FROM Thesis WHERE TUID = ?", thesisID).Scan(&booked)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil // Thesis not found, so not booked
		}
		return false, fmt.Errorf("error checking if thesis is booked: %v", err)
	}
	return booked, nil
}

func (dbc *DBController) updtJunction(tx *sql.Tx, thesisID string, people []PersonalData, junctionTable string) error {
	_, err := tx.Exec(fmt.Sprintf("DELETE FROM %s WHERE TUID = ?", junctionTable), thesisID) // Not insecure, bc fixed variable
	if err != nil {
		return err
	}

	for _, person := range people {
		var exists bool
		err := tx.QueryRow("SELECT EXISTS(SELECT 1 FROM PersonalData WHERE PDUID = ?)", person.PDUid).Scan(&exists)
		if err != nil {
			return fmt.Errorf("pduid validation error: %v", err)
		}
		if !exists {
			return fmt.Errorf("invalid PDUID: %s", person.PDUid)
		}

		_, err = tx.Exec(fmt.Sprintf("INSERT INTO %s (TUID, PDUID) VALUES (?, ?)", junctionTable), // Not insecure, bc fixed variable
			thesisID, person.PDUid)
		if err != nil {
			return fmt.Errorf("junction insert error: %v", err)
		}
	}

	return nil
}

// File management methods
func (dbc *DBController) InsrtFileRecord(tuid, fileName, originalFileName string, fileSize int64, pduid string) (int64, error) {
	result, err := dbc.db.Exec(
		"INSERT INTO ThesisFiles (TUID, FileName, OriginalFileName, FileSize, PDUID) VALUES (?, ?, ?, ?, ?)",
		tuid, fileName, originalFileName, fileSize, pduid,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to insert file record: %v", err)
	}

	fileID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get file ID: %v", err)
	}

	return fileID, nil
}

func (dbc *DBController) GtFilesByThesis(tuid string) ([]ThesisFile, error) {
	query := `
		SELECT FUID, TUID, FileName, OriginalFileName, FileSize, UploadDate, COALESCE(PDUID, '') as PDUID
		FROM ThesisFiles
		WHERE TUID = ?
		ORDER BY UploadDate DESC
	`

	rows, err := dbc.db.Query(query, tuid)
	if err != nil {
		return nil, fmt.Errorf("failed to query files: %v", err)
	}
	defer rows.Close()

	var files []ThesisFile
	for rows.Next() {
		var file ThesisFile
		err := rows.Scan(&file.FUID, &file.TUID, &file.FileName, &file.OriginalFileName, &file.FileSize, &file.UploadDate, &file.PDUID)
		if err != nil {
			return nil, fmt.Errorf("failed to scan file: %v", err)
		}
		files = append(files, file)
	}

	return files, nil
}

func (dbc *DBController) GtFileByID(fuid string) (*ThesisFile, error) {
	var file ThesisFile
	query := `
		SELECT FUID, TUID, FileName, OriginalFileName, FileSize, UploadDate, COALESCE(PDUID, '') as PDUID
		FROM ThesisFiles
		WHERE FUID = ?
	`

	err := dbc.db.QueryRow(query, fuid).Scan(&file.FUID, &file.TUID, &file.FileName, &file.OriginalFileName, &file.FileSize, &file.UploadDate, &file.PDUID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("file not found")
		}
		return nil, fmt.Errorf("failed to get file: %v", err)
	}

	return &file, nil
}

func (dbc *DBController) DelFileRecord(fuid string) error {
	result, err := dbc.db.Exec("DELETE FROM ThesisFiles WHERE FUID = ?", fuid)
	if err != nil {
		return fmt.Errorf("failed to delete file record: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("file not found")
	}

	return nil
}

func convertGoDateToSqlNullDate(date time.Time) sql.NullTime {
	if date.IsZero() {
		return sql.NullTime{Valid: false}
	}
	return sql.NullTime{Time: date, Valid: true}
}

func convertGradeToNullFloat(fg float64) sql.NullFloat64 {
	if fg == -1 {
		return sql.NullFloat64{Valid: false}
	}
	return sql.NullFloat64{Float64: fg, Valid: true}
}
