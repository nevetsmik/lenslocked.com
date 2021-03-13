package services

import (
	"lenslocked.com/hash"
	"lenslocked.com/interfaces"
	"lenslocked.com/models"
	"lenslocked.com/rand"
)

type pwResetValidator struct {
	interfaces.PwResetDBInt
	hmac hash.HMAC
}

func NewPwResetValidator(db interfaces.PwResetDBInt, hmac hash.HMAC) *pwResetValidator {
	return &pwResetValidator{
		PwResetDBInt: db,
		hmac:         hmac,
	}
}

func (pwrv *pwResetValidator) ByToken(token string) (*models.PwReset, error) {
	pwr := models.PwReset{Token: token}
	err := runPwResetValFns(&pwr, pwrv.hmacToken)
	if err != nil {
		return nil, err
	}
	return pwrv.PwResetDBInt.ByToken(pwr.TokenHash)
}

func (pwrv *pwResetValidator) CreatePwResetToken(pwr *models.PwReset) error {
	err := runPwResetValFns(pwr,
		pwrv.requireUserID,
		pwrv.setTokenIfUnset,
		pwrv.hmacToken,
	)
	if err != nil {
		return err
	}
	return pwrv.PwResetDBInt.CreatePwResetToken(pwr)
}

func (pwrv *pwResetValidator) DeletePwResetToken(id uint) error {
	if id <= 0 {
		return models.ErrIDInvalid
	}
	return pwrv.PwResetDBInt.DeletePwResetToken(id)
}

type pwResetValFn func(*models.PwReset) error

func runPwResetValFns(pwr *models.PwReset, fns ...pwResetValFn) error {
	for _, fn := range fns {
		if err := fn(pwr); err != nil {
			return err
		}
	}
	return nil
}

func (pwrv *pwResetValidator) requireUserID(pwr *models.PwReset) error {
	if pwr.UserID <= 0 {
		return models.ErrUserIDRequired
	}
	return nil
}

func (pwrv *pwResetValidator) setTokenIfUnset(pwr *models.PwReset) error {
	if pwr.Token != "" {
		return nil
	}
	token, err := rand.RememberToken()
	if err != nil {
		return err
	}
	pwr.Token = token
	return nil
}

func (pwrv *pwResetValidator) hmacToken(pwr *models.PwReset) error {
	if pwr.Token == "" {
		return nil
	}
	pwr.TokenHash = pwrv.hmac.Hash(pwr.Token)
	return nil
}
