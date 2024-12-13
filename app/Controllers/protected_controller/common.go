package protected_controller

import (
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
