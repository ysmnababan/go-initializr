package initializer

type service struct {
}

func NewService() *service {
	return &service{}
}

func (s *service) InitializeBoilerplate(req *BasicConfigRequest) (folderId string, err error) {
	return
}

func (s *service) DownloadProjectByFolderID(folderID string) (err error){
	return
}
