package app

import (
	"github.com/tuomasvapaavuori/site_installer/app/models"
)

type InstallTemplate struct {
	MysqlUser     models.RandomValue
	MysqlPassword models.RandomValue
	DatabaseName  models.RandomValue
}
