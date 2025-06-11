package initializer

type service struct {
}

func NewService() *service {
	return &service{}
}

func (s *service) InitializeBoilerplate() (folderId string, err error) {
	return
}
