package http

import (
	"angsuran-service/internal/controller/request"
	"angsuran-service/util"
	"encoding/json"
	"net/http"
)

func (a *AngsuranRoutes) CalculateAngsuranHandler(w http.ResponseWriter, r *http.Request) {
	var reqBody request.AngsuranRequest
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := util.ValidateRequestBody(reqBody); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	angsurans := a.au.CalculateAngsuran(reqBody.Plafond, reqBody.LamaPinjaman, reqBody.Bunga, util.ParseTanggal(reqBody.TanggalMulai))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"angsurans": angsurans})
}
