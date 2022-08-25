package provider

type testProvider struct {
}

func (t *testProvider) Name() string {
	return "test"
}

func (t *testProvider) ValidateLogin(username string, password string) error {
	return nil
}

func (t *testProvider) NeedsLogin() bool {
	return true
}

func (t *testProvider) Export(report *Report) error {
	return nil
}

func NewTestProvider() ExportProvider {
	return &testProvider{}
}
