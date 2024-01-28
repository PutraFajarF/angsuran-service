package request

type AngsuranRequest struct {
	Plafond      int     `json:"plafond"`
	LamaPinjaman int     `json:"lama_pinjaman"`
	Bunga        float64 `json:"bunga"`
	TanggalMulai string  `json:"tanggal_mulai"`
}
