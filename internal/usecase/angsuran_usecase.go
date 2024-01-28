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
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/xuri/excelize/v2"
)

type IAngsuranUsecase interface {
	CalculateAngsuran(data *request.AngsuranRequest) ([]*entity.Angsuran, error)
	CreateExcelFile(data []*entity.Angsuran) (*excelize.File, error)
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

func (a *AngsuranUsecase) CreateExcelFile(data []*entity.Angsuran) (*excelize.File, error) {
	startTime := time.Now()
	xlsx := excelize.NewFile()
	defer func() {
		if err := xlsx.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	sheetName := "AngsuranTable"
	index, err := xlsx.NewSheet(sheetName)
	if err != nil {
		a.l.CreateLog(&logger.Log{
			Event:        util.USECASE_GENERATEXLSX_ANGSURAN + "|POST|NEWSHEET",
			Method:       "POST",
			StatusCode:   http.StatusInternalServerError,
			Request:      "Set NewSheet XLSX Angsuran",
			Response:     err.Error(),
			ResponseTime: time.Since(startTime),
			Message:      util.FAIL_GENERATE_ANGSURAN_XLSX,
		}, logger.LVL_ERROR)
		return nil, err
	}

	// Set header
	xlsx.SetCellValue(sheetName, "A1", "Angsuran Ke")
	xlsx.SetCellValue(sheetName, "B1", "Tanggal")
	xlsx.SetCellValue(sheetName, "C1", "Total Angsuran")
	xlsx.SetCellValue(sheetName, "D1", "Angsuran Pokok")
	xlsx.SetCellValue(sheetName, "E1", "Angsuran Bunga")
	xlsx.SetCellValue(sheetName, "F1", "Sisa Pinjaman")

	var wg sync.WaitGroup

	// Fill data
	for i, val := range data {
		wg.Add(1)
		go func(i int, val *entity.Angsuran) {
			defer wg.Done()
			row := i + 2
			xlsx.SetCellValue(sheetName, fmt.Sprintf("A%d", row), val.AngsuranKe)
			xlsx.SetCellValue(sheetName, fmt.Sprintf("B%d", row), val.Tanggal)
			xlsx.SetCellValue(sheetName, fmt.Sprintf("C%d", row), val.TotalAngsuran)
			xlsx.SetCellValue(sheetName, fmt.Sprintf("D%d", row), val.AngsuranPokok)
			xlsx.SetCellValue(sheetName, fmt.Sprintf("E%d", row), val.AngsuranBunga)
			xlsx.SetCellValue(sheetName, fmt.Sprintf("F%d", row), val.SisaPinjaman)
		}(i, val)
	}

	wg.Wait()

	xlsx.SetActiveSheet(index)

	filePath := a.cfg.Excelize.Path
	err = xlsx.SaveAs(filePath)
	if err != nil {
		a.l.CreateLog(&logger.Log{
			Event:        util.USECASE_GENERATEXLSX_ANGSURAN + "|POST|SAVE",
			Method:       "POST",
			StatusCode:   http.StatusInternalServerError,
			Request:      "Save XLSX Angsuran",
			Response:     err.Error(),
			ResponseTime: time.Since(startTime),
			Message:      util.FAIL_GENERATE_ANGSURAN_XLSX,
		}, logger.LVL_ERROR)
		return nil, errors.Wrap(err, "failed to save Excel file")
	}

	return xlsx, nil
}
