package usecase

import (
	"angsuran-service/config"
	"angsuran-service/internal/controller/request"
	"angsuran-service/internal/entity"
	"angsuran-service/pkg/logger"
	"angsuran-service/util"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type IAngsuranUsecase interface {
	CalculateAngsuran(data *request.AngsuranRequest) ([]*entity.Angsuran, error)
}

type AngsuranUsecase struct {
	l   *logger.Logger
	cfg *config.Config
}

func NewAngsuranUsecase(l *logger.Logger, cfg *config.Config) *AngsuranUsecase {
	return &AngsuranUsecase{l, cfg}
}

func (a *AngsuranUsecase) CalculateAngsuran(data *request.AngsuranRequest) ([]*entity.Angsuran, error) {
	var angsurans []*entity.Angsuran

	startTime := time.Now()
	req := map[string]interface{}{"request": data}
	jsonReq, err := json.Marshal(req)
	if err != nil {
		a.l.CreateLog(&logger.Log{
			Event:        util.USECASE_ANGSURAN + "|POST",
			Method:       "POST",
			StatusCode:   http.StatusBadRequest,
			Request:      string(jsonReq),
			Response:     err.Error(),
			ResponseTime: time.Since(startTime),
			Message:      util.FAIL_CALCULATE_ANGSURAN,
		}, logger.LVL_ERROR)
		return nil, fmt.Errorf("failed to marshal request data: %v", err)
	}

	if data == nil {
		a.l.CreateLog(&logger.Log{
			Event:        util.USECASE_ANGSURAN + "|POST",
			Method:       "POST",
			StatusCode:   http.StatusBadRequest,
			Request:      string(jsonReq),
			Response:     util.ErrDataIsNil,
			ResponseTime: time.Since(startTime),
			Message:      util.FAIL_CALCULATE_ANGSURAN,
		}, logger.LVL_ERROR)
		return nil, util.ErrDataIsNil
	}

	if data.LamaPinjaman <= 0 || data.Bunga <= 0 {
		a.l.CreateLog(&logger.Log{
			Event:        util.USECASE_ANGSURAN + "|POST",
			Method:       "POST",
			StatusCode:   http.StatusBadRequest,
			Request:      string(jsonReq),
			Response:     util.ErrInvalidLoanDurationAndInterestRate,
			ResponseTime: time.Since(startTime),
			Message:      util.FAIL_CALCULATE_ANGSURAN,
		}, logger.LVL_ERROR)
		return nil, util.ErrInvalidLoanDurationAndInterestRate
	}

	ratePerMonth := data.Bunga / 12
	totalAngsuran := (float64(data.Plafond) * ratePerMonth) * (util.Pow(1+ratePerMonth, float64(data.LamaPinjaman))) /
		(util.Pow(1+ratePerMonth, float64(data.LamaPinjaman)) - 1)

	sisaPinjaman := float64(data.Plafond)
	tanggalMulai := util.ParseTanggal(data.TanggalMulai)

	for i := 1; i <= data.LamaPinjaman; i++ {
		angsuranBunga := (data.Bunga / 360.0) * 30 * sisaPinjaman
		angsuranPokok := totalAngsuran - angsuranBunga
		tanggalAngsuran := tanggalMulai.AddDate(0, i-1, 0).Format("2006-01-02")

		angsurans = append(angsurans, &entity.Angsuran{
			AngsuranKe:    i,
			Tanggal:       tanggalAngsuran,
			TotalAngsuran: util.RoundToTwoDecimal(totalAngsuran),
			AngsuranPokok: util.RoundToTwoDecimal(angsuranPokok),
			AngsuranBunga: util.RoundToTwoDecimal(angsuranBunga),
			SisaPinjaman:  util.RoundToTwoDecimal(sisaPinjaman - angsuranPokok),
		})

		sisaPinjaman -= angsuranPokok
	}

	jsonRes, err := json.Marshal(angsurans)
	if err != nil {
		a.l.CreateLog(&logger.Log{
			Event:        util.USECASE_ANGSURAN + "|POST",
			Method:       "POST",
			StatusCode:   http.StatusBadRequest,
			Request:      string(jsonReq),
			Response:     err.Error(),
			ResponseTime: time.Since(startTime),
			Message:      util.FAIL_CALCULATE_ANGSURAN,
		}, logger.LVL_ERROR)
		return nil, fmt.Errorf("failed to marshal request data: %v", err)
	}

	a.l.CreateLog(&logger.Log{
		Event:        util.USECASE_ANGSURAN + "|POST",
		Method:       "POST",
		StatusCode:   http.StatusOK,
		Request:      string(jsonReq),
		Response:     string(jsonRes),
		ResponseTime: time.Since(startTime),
		Message:      util.SUCCESS_CALCULATE_ANGSURAN,
	}, logger.LVL_INFO)

	return angsurans, nil
}
