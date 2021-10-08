package gwt

func (gwt *Gwt) ForceLogoutUser(userId string) error {
	return gwt.settings.storage.deleteAllTokens(userId)
}
