package test

type MockCfg struct {
	CurrentUser string
	SetUserErr  error
}

func (m *MockCfg) SetUser(name string) error {
	if m.SetUserErr != nil {
		return m.SetUserErr
	}
	m.CurrentUser = name
	return nil
}

func (m *MockCfg) GetCurrentUser() string {
	return m.CurrentUser
}
