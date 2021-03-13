package interfaces

import "lenslocked.com/models"

type PwResetDBInt interface {
	ByToken(token string) (*models.PwReset, error)
	CreatePwResetToken(pwr *models.PwReset) error
	DeletePwResetToken(id uint) error
}
