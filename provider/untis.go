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

func (u *untisImportProvider) Import() ([]*Entry, error) {
	return nil, nil
}

func NewUntisProvider() ImportProvider {
	return &untisImportProvider{}
}
