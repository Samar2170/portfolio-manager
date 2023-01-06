package main

import (
	"github.com/Samar2170/portfolio-manager/account"
	"github.com/Samar2170/portfolio-manager/utils"
)

type CreateGeneralAccountJob struct {
	UserId uint
}

func (cgaj CreateGeneralAccountJob) Do() error {
	accCode := utils.CreateRandomString(GACodeLen)
	ga := account.GeneralAccount{
		Code:   accCode,
		UserId: cgaj.UserId,
	}
	var err error
	ga, err = ga.Create()
	if err != nil {
		return err
	}
	return nil
}

type CreateUserAccountStatusJob struct {
	UserId uint
}

func (cuas CreateUserAccountStatusJob) Do() error {
	uas := account.UserAccountStatus{
		UserId:           cuas.UserId,
		GaAccountCreated: false,
	}
	uas.GetOrCreate()
	return nil
}

const GACodeLen = 8
