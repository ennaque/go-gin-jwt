package gwt

type Service struct {
	settings *Settings
}

func (service *Service) ForceLogoutUser(userId string) error {
	return service.settings.Storage.DeleteAllTokens(userId)
}
