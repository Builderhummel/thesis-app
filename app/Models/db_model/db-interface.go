package db_model

import (
	"fmt"
	"time"

	"github.com/Builderhummel/thesis-app/app/Constants/roles"
)

var dbSession *DBController

// TODO: Check for better error handling in whole file
func Init() {
	dbSession = &DBController{}
	for {
		err := dbSession.OpenConnection()
		if err != nil {
			fmt.Printf("Failed to connect to the database: %v. Retrying in 5s...\n", err)
		} else {
			if pingErr := dbSession.db.Ping(); pingErr != nil {
				fmt.Printf("Database ping failed: %v. Retrying in 5s...\n", pingErr)
			} else {
				fmt.Println("Database connection established.")
				break
			}
		}
		time.Sleep(5 * time.Second)
	}

	// Check if database is initialized
	check, err := dbSession.CheckIfDatabaseIsInitialized()
	if err != nil {
		panic(err)
	}
	if !check {
		err = dbSession.InitDatabase()
		if err != nil {
			panic(err)
		}
	}
}

func VerifyLoginUser(userid string) (bool, error) {
	uid, err := dbSession.GetLoginHandleFromDB(userid)
	if err != nil {
		return false, err
	}
	if uid == "" {
		return false, nil
	}
	return true, nil
}

func GetUserRoleByLoginHandle(userid string) (roles.Role, error) {
	role, err := dbSession.GtUsrRleByLgnHndle(userid)
	if err != nil {
		return "", err
	}
	return role, nil
}

func CheckUserActive(userid string) (bool, error) {
	authorized, err := dbSession.ChkUserActive(userid)
	if err != nil {
		return false, err
	}
	return authorized, nil
}

func UpdateUser(userid, name, email string) error {
	err := dbSession.UpdtUser(userid, name, email)
	if err != nil {
		return err
	}
	return nil
}

func GetAllUsers() ([]PersonalData, error) {
	data, err := dbSession.GtAllUsrs()
	if err != nil {
		return nil, err
	}
	return data, nil
}

func GetUserByPUID(puid string) (PersonalData, error) {
	data, err := dbSession.GtUsrByPUID(puid)
	if err != nil {
		return PersonalData{}, err
	}
	return data, nil
}

func GetUserByLoginHandle(loginHandle string) (PersonalData, error) {
	data, err := dbSession.GtUsrByLgnHndle(loginHandle)
	if err != nil {
		return PersonalData{}, err
	}
	return data, nil
}

func GetAllSupervisors() ([]PersonalData, error) {
	data, err := dbSession.GtAllSupervisors()
	if err != nil {
		return nil, err
	}
	return data, nil
}

func GetAllExaminers() ([]PersonalData, error) {
	data, err := dbSession.GtAllExaminers()
	if err != nil {
		return nil, err
	}
	return data, nil
}

func InsertNewUser(name, email, handle, role string, active, isSupervisor, isExaminer bool) error {
	err := dbSession.InsrtNwUsr(name, email, handle, role, active, isSupervisor, isExaminer)
	if err != nil {
		return err
	}
	return nil
}

func UpdateFullUser(puid, name, email, handle, role string, active, isSupervisor, isExaminer bool) error {
	err := dbSession.UptFullUsr(puid, name, email, handle, role, active, isSupervisor, isExaminer)
	if err != nil {
		return err
	}
	return nil
}

func GetHomepageRCW() (map[string]string, error) {
	result, err := dbSession.GtHomepageRCW()
	if err != nil {
		return nil, err
	}
	return result, nil
}

func GetDataThesisTableOpenRequests() ([]map[string]string, error) {
	data, err := dbSession.GtDataTblOpenReq()
	if err != nil {
		return nil, err
	}
	return data, nil
}

func GetDataTableAllSupervisions() ([]map[string]string, error) {
	data, err := dbSession.GtDataTblAllSupervisions()
	if err != nil {
		return nil, err
	}
	return data, nil
}

func GetUserPuidFromUserId(user_id string) (string, error) {
	return dbSession.GtUsrPuidFromUserId(user_id)
}

func GetDataTableMySupervisions(user_id string) ([]map[string]string, error) {
	puid, err := dbSession.GtUsrPuidFromUserId(user_id)
	if err != nil {
		return nil, err
	}
	data, err := dbSession.GtDataTblMySupervisions(puid)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func GetDataFullSupervision(thesisID string) (*ThesisFullData, error) {
	data, err := dbSession.GtDataFullSupervision(thesisID)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func GetAllDataFullSupervison() ([]*ThesisFullData, error) {
	data, err := dbSession.GtAllDataFullSupervision()
	if err != nil {
		return nil, err
	}
	return data, nil
}

func InsertNewThesisRequest(name, email, courseOfStudy, thesisType, thesisTitle, gpa, contactDate, notes string) error {
	err := dbSession.InsrtNwThsisRequest(name, email, courseOfStudy, thesisType, thesisTitle, gpa, contactDate, notes)
	if err != nil {
		return err
	}
	return nil
}

func UpdateThesisInfo(td *ThesisFullData) error {
	err := dbSession.UpdtThesisInfo(td)
	if err != nil {
		return err
	}
	return nil
}

func AddThesisSupervisor(tuid string, supervisor PersonalData) error {
	err := dbSession.addToJunction(tuid, supervisor, "SupervisorJunction")
	if err != nil {
		return err
	}
	return nil
}

func UpdateThesisSupervisors(tuid string, supervisors []PersonalData) error {
	err := dbSession.updtJunction(tuid, supervisors, "SupervisorJunction")
	if err != nil {
		return err
	}
	return nil
}

func DeleteThesisRequest(tuid string) error {
	err := dbSession.DelThesisRequest(tuid)
	if err != nil {
		return err
	}
	return nil
}

func CheckIfThesisIsBooked(thesisID string) (bool, error) {
	return dbSession.ChkIfThesisIsBooked(thesisID)
}

// File operations
func InsertFileRecord(tuid, fileName, originalFileName string, fileSize int64, pduid, category string) (int64, error) {
	return dbSession.InsrtFileRecord(tuid, fileName, originalFileName, fileSize, pduid, category)
}

func GetFilesByThesis(tuid string) ([]ThesisFile, error) {
	return dbSession.GtFilesByThesis(tuid)
}

func GetFileByID(fuid string) (*ThesisFile, error) {
	return dbSession.GtFileByID(fuid)
}

func GetLatestFileByCategory(tuid, category string) (*ThesisFile, error) {
	return dbSession.GtLatestFileByCategory(tuid, category)
}

func DeleteFileRecord(fuid string) error {
	return dbSession.DelFileRecord(fuid)
}
