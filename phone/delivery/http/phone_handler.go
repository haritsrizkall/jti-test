package http

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/haritsrizkall/jti-test/domain"
	"github.com/haritsrizkall/jti-test/pkg"
	"github.com/haritsrizkall/jti-test/utils"
)

type PhoneHandler struct {
	PhoneUsecase domain.PhoneUsecase
}

func NewPhoneHandler(phoneUsecase domain.PhoneUsecase) *PhoneHandler {
	return &PhoneHandler{
		PhoneUsecase: phoneUsecase,
	}
}

func (h *PhoneHandler) GetAll(resp http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	phoneResponse, err := h.PhoneUsecase.GetAll(ctx)
	if err != nil {
		// handle error
		response := utils.Response{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		}

		utils.NewResponse(resp, response)
		return
	}

	// handle response
	response := utils.Response{
		Status:  http.StatusOK,
		Message: "Success",
		Data:    phoneResponse,
	}

	utils.NewResponse(resp, response)
}

func (h *PhoneHandler) Create(resp http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	var request domain.CreatePhoneRequest
	err := utils.DecodeBody(req, &request)
	if err != nil {
		// handle error
		response := utils.Response{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		}

		utils.NewResponse(resp, response)
		return
	}

	// validate
	err = pkg.Validate.Struct(request)
	if err != nil {
		// handle error
		response := utils.Response{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
			Data:    nil,
		}

		utils.NewResponse(resp, response)
		return
	}

	phone, err := h.PhoneUsecase.Create(ctx, request)
	if err != nil {
		// handle error
		response := utils.Response{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		}
		utils.NewResponse(resp, response)
		return
	}

	// handle response
	response := utils.Response{
		Status:  http.StatusOK,
		Message: "Success",
		Data:    phone,
	}

	utils.NewResponse(resp, response)
}

func (h *PhoneHandler) Update(resp http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	var request domain.UpdatePhoneRequest
	err := utils.DecodeBody(req, &request)
	if err != nil {
		// handle error
		response := utils.Response{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		}

		utils.NewResponse(resp, response)
		return
	}

	// validate
	err = pkg.Validate.Struct(request)
	if err != nil {
		// handle error
		fmt.Println(err)
		response := utils.Response{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
			Data:    nil,
		}

		utils.NewResponse(resp, response)
		return
	}

	id := mux.Vars(req)["id"]
	if id == "" {
		// handle error
		response := utils.Response{
			Status:  http.StatusInternalServerError,
			Message: "ID cannot be empty",
			Data:    nil,
		}

		utils.NewResponse(resp, response)
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		// handle error
		response := utils.Response{
			Status:  http.StatusInternalServerError,
			Message: "ID must be integer",
			Data:    nil,
		}

		utils.NewResponse(resp, response)
		return
	}
	request.ID = idInt

	phone, err := h.PhoneUsecase.Update(ctx, request)
	if err != nil {
		// handle error
		response := utils.Response{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		}
		utils.NewResponse(resp, response)
		return
	}

	// handle response
	response := utils.Response{
		Status:  http.StatusOK,
		Message: "Success",
		Data:    phone,
	}

	utils.NewResponse(resp, response)
}

func (h *PhoneHandler) Delete(resp http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	id := mux.Vars(req)["id"]
	if id == "" {
		// handle error
		response := utils.Response{
			Status:  http.StatusInternalServerError,
			Message: "ID cannot be empty",
			Data:    nil,
		}

		utils.NewResponse(resp, response)
		return
	}
	idInt, err := strconv.Atoi(id)
	if err != nil {
		// handle error
		response := utils.Response{
			Status:  http.StatusInternalServerError,
			Message: "ID must be integer",
			Data:    nil,
		}

		utils.NewResponse(resp, response)
		return
	}

	phone, err := h.PhoneUsecase.Delete(ctx, idInt)
	if err != nil {
		// handle error
		response := utils.Response{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
			Data:    nil,
		}
		utils.NewResponse(resp, response)
		return
	}

	// handle response
	response := utils.Response{
		Status:  http.StatusOK,
		Message: "Success",
		Data:    phone,
	}

	utils.NewResponse(resp, response)
}
