package provider

type untisImportProvider struct {
}

func (u *untisImportProvider) Name() string {
	return ""
}

func (u *untisImportProvider) ValidateLogin(username string, password string) error {
	return nil
}

func (u *untisImportProvider) NeedsLogin() bool {
	return false
}

func (u *untisImportProvider) Import(startDate string) ([]*Entry, error) {
	//TODO implement me
	panic("implement me")
}

func NewUntisProvider() ImportProvider {
	return &untisImportProvider{}
}
