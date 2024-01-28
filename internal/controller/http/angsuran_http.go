package http

import (
	"angsuran-service/internal/controller/request"
	"angsuran-service/internal/controller/response"
	"angsuran-service/pkg/logger"
	"angsuran-service/util"
	"encoding/json"
	"net/http"
	"time"
)

func (a *AngsuranRoutes) CalculateAngsuranHandler(w http.ResponseWriter, r *http.Request) {
	var reqBody request.AngsuranRequest

	timeStart := time.Now()
	req := map[string]interface{}{"payload": reqBody}
	jsonReq, _ := json.Marshal(req)

	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		a.l.CreateLog(&logger.Log{
			Event:        util.HTTP_ANGSURAN + "|POST",
			Method:       "POST",
			StatusCode:   http.StatusBadRequest,
			Request:      string(jsonReq),
			Response:     err.Error(),
			ResponseTime: time.Since(timeStart),
			Message:      util.ErrInvalidPayload,
		}, logger.LVL_ERROR)
		response.HttpErrorResponse(w, false, http.StatusBadRequest, "400", err.Error())
		return
	}

	if err := util.ValidateRequestBody(reqBody); err != nil {
		a.l.CreateLog(&logger.Log{
			Event:        util.HTTP_ANGSURAN + "|POST",
			Method:       "POST",
			StatusCode:   http.StatusBadRequest,
			Request:      string(jsonReq),
			Response:     err.Error(),
			ResponseTime: time.Since(timeStart),
			Message:      util.ErrInvalidPayload,
		}, logger.LVL_ERROR)
		response.HttpErrorResponse(w, false, http.StatusBadRequest, "400", err.Error())
		return
	}

	angsurans, err := a.au.CalculateAngsuran(&reqBody)
	if err != nil {
		a.l.CreateLog(&logger.Log{
			Event:        util.HTTP_ANGSURAN + "|POST",
			Method:       "POST",
			StatusCode:   http.StatusBadRequest,
			Request:      string(jsonReq),
			Response:     err.Error(),
			ResponseTime: time.Since(timeStart),
			Message:      util.FAIL_CALCULATE_ANGSURAN,
		}, logger.LVL_ERROR)
		response.HttpErrorResponse(w, false, http.StatusBadRequest, "400", err.Error())
		return
	}

	response.HttpSuccessResponse(w, true, http.StatusOK, "200", util.SUCCESS_CALCULATE_ANGSURAN, angsurans)
}
