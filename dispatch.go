package models

// ToDTO converts a TaskDispatchPayload into a TaskDTO, mapping the flat
// dispatch fields into the nested App/AppVersion pointers that the engine's
// internal pipeline expects.
func (d TaskDispatchPayload) ToDTO() TaskDTO {
	return TaskDTO{
		BaseModelDTO:       BaseModelDTO{ID: d.ID, ShortID: d.ShortID},
		PermissionModelDTO: PermissionModelDTO{UserID: d.UserID, TeamID: d.TeamID},
		Status:             d.Status,
		AppID:              d.AppID,
		AppVersionID:       d.AppVersionID,
		AppVariant:         d.AppVariant,
		Function:           d.Function,
		Input:              d.Input,
		Setup:              d.Setup,
		WorkerID:           d.WorkerID,
		SessionID:          d.SessionID,
		SessionTimeout:     d.SessionTimeout,
		App: &AppDTO{
			Name: d.AppName,
		},
		AppVersion: &AppVersionDTO{
			Repository:        d.Repository,
			Kernel:            d.Kernel,
			Env:               d.AppEnv,
			RequiredResources: AppResources{GPU: AppGPUResource{Count: d.GPUCount}},
		},
	}
}
