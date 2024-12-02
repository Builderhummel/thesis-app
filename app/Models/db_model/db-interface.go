package db_model

var dbSession *DBController

// TODO: Check for better error handling in whole file
func Init() {
	dbSession = &DBController{}
	dbSession.OpenConnection()

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

func GetDataThesisTableOpenRequests() ([]map[string]string, error) {
	data, err := dbSession.GtDataTblOpenReq()
	if err != nil {
		return nil, err
	}
	return data, nil
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
