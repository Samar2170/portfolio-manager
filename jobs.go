package main

import (
	"github.com/Samar2170/portfolio-manager/account"
	"github.com/Samar2170/portfolio-manager/portfolio"
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

type ParseFDFileJob struct {
	FileId uint
}
type ParseStockFileJob struct {
	FileId uint
}
type ParseMFFileJob struct {
	FileId uint
}

func (pffj ParseFDFileJob) Do() error {
	err := portfolio.ParseFDFile(pffj.FileId)
	return err
}
func (psfj ParseStockFileJob) Do() error {
	err := portfolio.ParseStockFile(psfj.FileId)
	return err
}
func (pmfj ParseMFFileJob) Do() error {
	err := portfolio.ParseStockFile(pmfj.FileId)
	return err
}
