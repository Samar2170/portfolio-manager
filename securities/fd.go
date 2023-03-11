package securities

import (
	"time"

	"github.com/Samar2170/portfolio-manager/utils"
)

func GetFDByID(id uint) (FixedDeposit, error) {
	var fd FixedDeposit
	err := db.First(&fd, "id = ?", id).Error
	return fd, err
}

func getFdsByIds(ids []uint) ([]FixedDeposit, error) {
	var fds []FixedDeposit
	err := db.Where("id IN ?", ids).Find(&fds).Error
	return fds, err
}

func CalculateNextIPDate(fd FixedDeposit) time.Time {
	var nextIPDate time.Time
	today := time.Now()
	ipdate := fd.IPDate
	var currIpdate time.Time
	if ipdate.Before(today) {
		switch fd.IPFreq {
		case "A":
			currIpdate = time.Date(today.Year(), ipdate.Month(), ipdate.Day(), 0, 0, 0, 0, time.Local)
			if currIpdate.Before(today) {
				nextIPDate = currIpdate.AddDate(1, 0, 0)
			} else {
				nextIPDate = currIpdate
			}

		case "M":
			currIpdate = time.Date(today.Year(), today.Month(), ipdate.Day(), 0, 0, 0, 0, time.Local)
			nextIPDate = currIpdate.AddDate(0, 1, 0)
		case "Q":
			nextIPDate = utils.GetNextQuarter(today)
		case "MT":
			nextIPDate = fd.MtDate
		case "SA":
			nextIPDate = utils.GetNextHY(today)
			// case "QAD":
			// 	nextIPDate = utils.GetNextQuarter(today)
			// case "SAD":
			// 	nextIPDate = utils.GetNextHY(today)
		}
	} else {
		nextIPDate = currIpdate
	}
	return nextIPDate
}

func UpdateNextIPDatesFDs() error {
	var fds []FixedDeposit

	err := db.Where("next_ip_date < CURRENT_DATE").Find(&fds).Error
	if err != nil {
		return err
	}

	for _, fd := range fds {
		nipd := CalculateNextIPDate(fd)
		err = db.Model(&fd).Where("id = ?", fd.ID).Update("next_ip_date", nipd).Error
	}
	return err
}

func (fd *FixedDeposit) Create() error {
	return db.Create(&fd).Error
}

func CalculateAccruedInterest(fd FixedDeposit) float64 {
	t := time.Now()
	if fd.IPFreq == "MT" {
		totalInterest := fd.MtAmount - fd.Amount
		totalTime := fd.MtDate.Sub(fd.StartDate) / 24
		dailyInterest := totalInterest / float64(totalTime)
		timeDone := t.Sub(fd.StartDate) / 24
		return float64(timeDone) * dailyInterest
	} else {
		var daysDiff int
		var intMultiple float64
		annualInterest := fd.IPRate * fd.Amount
		switch fd.IPFreq {
		case "M":
			daysDiff = int(t.Day())
			intMultiple = float64(daysDiff) / 30 * 0.0833

		case "Q":
			cqDate := utils.GetCurrentQuarterFirstDate(t)
			daysDiff = int(t.Sub(cqDate).Hours() / 24)
			intMultiple = float64(daysDiff) / 90 * 0.25
		case "SA":
			cqDate := utils.GetCurrentHYFirstDate(t)
			daysDiff = int(t.Sub(cqDate).Hours() / 24)
			intMultiple = float64(daysDiff) / 180 * 0.5
		case "A":
			ipDate := time.Date(t.Year(), fd.StartDate.Month(), fd.StartDate.Day(), 0, 0, 0, 0, time.UTC)
			if t.After(ipDate) {
				delta := ipDate.Sub(t)
				daysDiff = int(delta.Hours() / 24)
				intMultiple = float64(daysDiff) / 365
			} else {
				ipDate = time.Date(t.Year()-1, fd.StartDate.Month(), fd.StartDate.Day(), 0, 0, 0, 0, time.UTC)
				delta := t.Sub(ipDate)
				daysDiff := int(delta.Hours() / 24)
				intMultiple = float64(daysDiff) / 365
			}
		}
		accruedInt := annualInterest * intMultiple
		return accruedInt
	}
}

func CalculateAccruedInterestAllFDs() error {
	var ids []uint
	db.Raw("SELECT id FROM fixed_deposits").Scan(&ids)
	batches := len(ids) / BATCHSIZE

	for i := 0; i <= batches; i++ {
		newIds := ids[i*BATCHSIZE : i*BATCHSIZE+BATCHSIZE]
		fds, err := getFdsByIds(newIds)
		if err != nil {
			return err
		}
		for _, fd := range fds {
			accruedInterest := CalculateAccruedInterest(fd)
			db.Model(&fd).Where("id = ?", fd.ID).Update("accrued_interest", accruedInterest)
		}
	}
	return nil
}
