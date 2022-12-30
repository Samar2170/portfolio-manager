package main

import (
	"log"

	"github.com/Samar2170/portfolio-manager/account"
	"github.com/Samar2170/portfolio-manager/utils"
)

type CreateGeneralAccountJob struct {
	UserId uint
}

func (cgaj CreateGeneralAccountJob) Do() {
	accCode := utils.CreateRandomString(GACodeLen)
	ga := account.GeneralAccount{
		Code:   accCode,
		UserId: cgaj.UserId,
	}
	var err error
	ga, err = ga.GetOrCreate()
	if err != nil {
		log.Println(err)
	}

}

type CreateUserAccountStatusJob struct {
	UserId uint
}

func (cuas CreateUserAccountStatusJob) Do() {
	uas := account.UserAccountStatus{
		UserId:           cuas.UserId,
		GaAccountCreated: false,
	}
	uas.GetOrCreate()
}

const GACodeLen = 8
