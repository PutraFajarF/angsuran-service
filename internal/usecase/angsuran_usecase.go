package usecase

import (
	"angsuran-service/config"
	"angsuran-service/internal/entity"
	"angsuran-service/pkg/logger"
	"angsuran-service/util"
	"time"
)

type IAngsuranUsecase interface {
	CalculateAngsuran(plafond, lamaPinjaman int, bunga float64, tanggalMulai time.Time) []*entity.Angsuran
}

type AngsuranUsecase struct {
	l   *logger.Logger
	cfg *config.Config
}

func NewAngsuranUsecase(l *logger.Logger, cfg *config.Config) *AngsuranUsecase {
	return &AngsuranUsecase{l, cfg}
}

func (a *AngsuranUsecase) CalculateAngsuran(plafond, lamaPinjaman int, bunga float64, tanggalMulai time.Time) []*entity.Angsuran {
	var angsurans []*entity.Angsuran

	ratePerMonth := bunga / 12
	totalAngsuran := (float64(plafond) * ratePerMonth) * (util.Pow(1+ratePerMonth, float64(lamaPinjaman))) /
		(util.Pow(1+ratePerMonth, float64(lamaPinjaman)) - 1)

	sisaPinjaman := float64(plafond)

	for i := 1; i <= lamaPinjaman; i++ {
		angsuranBunga := (bunga / 360.0) * 30 * sisaPinjaman
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

	return angsurans
}
