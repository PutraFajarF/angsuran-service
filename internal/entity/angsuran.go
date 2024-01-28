package entity

type Angsuran struct {
	AngsuranKe    int     `json:"angsuran_ke,omitempty"`
	Tanggal       string  `json:"tanggal,omitempty"`
	TotalAngsuran float64 `json:"total_angsuran,omitempty"`
	AngsuranPokok float64 `json:"angsuran_pokok,omitempty"`
	AngsuranBunga float64 `json:"angsuran_bunga,omitempty"`
	SisaPinjaman  float64 `json:"sisa_pinjaman,omitempty"`
}
