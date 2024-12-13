package protected_controller

import (
	"strconv"
	"time"

	"github.com/Builderhummel/thesis-app/app/Models/db_model"
	view_protected_edit_supervision_request "github.com/Builderhummel/thesis-app/app/Views/handler/protected/edit_supervision_request"
)

func getSupervisorsSliceFromPersonalData(pd []db_model.PersonalData) []string {
	var supervisors []string
	for _, p := range pd {
		supervisors = append(supervisors, p.Name)
	}
	return supervisors
}

func convertDataPDtoViewEditPD(pd db_model.PersonalData) view_protected_edit_supervision_request.PersonalData {
	return view_protected_edit_supervision_request.PersonalData{
		PDUid: pd.PDUid,
		Name:  pd.Name,
	}
}

func convertSliceDataPDtoViewEditPD(pd []db_model.PersonalData) []view_protected_edit_supervision_request.PersonalData {
	var viewPDs []view_protected_edit_supervision_request.PersonalData
	for _, p := range pd {
		viewPDs = append(viewPDs, convertDataPDtoViewEditPD(p))
	}
	return viewPDs
}

func annotateSelectedPDs(allPDs, selectedPDs []view_protected_edit_supervision_request.PersonalData) []view_protected_edit_supervision_request.PersonalData {
	for i, pd := range allPDs {
		for _, selectedPD := range selectedPDs {
			if pd.PDUid == selectedPD.PDUid {
				allPDs[i].Selected = true
			}
		}
	}
	return allPDs
}

func processGradeString(gpa string) (float64, error) {
	gpa_float, err := strconv.ParseFloat(gpa, 64)
	if err != nil {
		if numErr, ok := err.(*strconv.NumError); ok && numErr.Err == strconv.ErrSyntax {
			return -1, nil
		} else {
			return 0, err
		}
	}

	if gpa_float < 1 {
		return -1, nil
	}

	return gpa_float, nil
}

func concatSemesterInfo(semester, year string) string {
	return semester + year
}

// TODO: Can never throw an error, check edge cases
func parseDateStringToGoDate(date string) (time.Time, error) {
	t, err := time.Parse("2006-01-02", date)
	if err != nil {
		return time.Time{}, nil
	}
	return t, nil
}
