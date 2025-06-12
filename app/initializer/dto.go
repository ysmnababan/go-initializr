package initializer

import "go-initializr/app/common"

type BasicConfigRequest struct {
	ProjectName string               `json:"project_name" validate:"required"`
	JWT         bool                 `json:"jwt"`
	Swagger     bool                 `json:"swagger"`
	Redis       bool                 `json:"redis"`
	Validator   bool                 `json:"validator"`
	DB          common.DBType        `json:"db" validate:"required,oneof=postgres mysql sqlserver"`
	FrameWork   common.FrameWorkType `json:"framework" validate:"required,oneof=echo gin"`
	ModInit     bool                 `json:"mod_init"`
}
