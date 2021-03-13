package services

import (
	"github.com/jinzhu/gorm"

	"lenslocked.com/models"
)

type pwResetGorm struct {
	db *gorm.DB
}

func (pwrg *pwResetGorm) ByToken(tokenHash string) (*models.PwReset, error) {
	var pwr models.PwReset
	err := first(pwrg.db.Where("token_hash = ?", tokenHash), &pwr)
	if err != nil {
		return nil, err
	}
	return &pwr, nil
}

func (pwrg *pwResetGorm) CreatePwResetToken(pwr *models.PwReset) error {
	return pwrg.db.Create(pwr).Error
}

func (pwrg *pwResetGorm) DeletePwResetToken(id uint) error {
	pwr := models.PwReset{Model: gorm.Model{ID: id}}
	return pwrg.db.Delete(&pwr).Error
}
