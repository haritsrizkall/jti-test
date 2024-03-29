package http

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/haritsrizkall/jti-test/constant"
	"github.com/haritsrizkall/jti-test/domain"
	phone_ws "github.com/haritsrizkall/jti-test/phone/websocket"
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
		utils.NewResponse(resp, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.NewResponse(resp, http.StatusOK, "Success", phoneResponse)
}

func (h *PhoneHandler) GetByID(resp http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	id := mux.Vars(req)["id"]
	if id == "" {
		utils.NewResponse(resp, http.StatusBadRequest, "ID cannot be empty", nil)
		return
	}
	idInt, err := strconv.Atoi(id)
	if err != nil {
		utils.NewResponse(resp, http.StatusBadRequest, "ID must be integer", nil)
		return
	}

	phone, err := h.PhoneUsecase.GetByID(ctx, idInt)
	if err != nil {
		if err.Error() == constant.ErrNoRowsInResultSet {
			utils.NewResponse(resp, http.StatusNotFound, err.Error(), nil)
			return
		}
	}

	utils.NewResponse(resp, http.StatusOK, "Success", phone)
}

func (h *PhoneHandler) Create(resp http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	var request domain.CreatePhoneRequest
	err := utils.DecodeBody(req, &request)
	if err != nil {
		utils.NewResponse(resp, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	// validate
	err = pkg.Validate.Struct(request)
	if err != nil {
		utils.NewResponse(resp, http.StatusBadRequest, err.Error(), nil)
		return
	}

	phone, err := h.PhoneUsecase.Create(ctx, request)
	if err != nil {
		utils.NewResponse(resp, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.NewResponse(resp, http.StatusOK, "Success", phone)
}

func (h *PhoneHandler) Update(resp http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	var request domain.UpdatePhoneRequest
	err := utils.DecodeBody(req, &request)
	if err != nil {
		utils.NewResponse(resp, http.StatusBadRequest, err.Error(), nil)
		return
	}

	// validate
	err = pkg.Validate.Struct(request)
	if err != nil {
		utils.NewResponse(resp, http.StatusBadRequest, err.Error(), nil)
		return
	}

	id := mux.Vars(req)["id"]
	if id == "" {
		utils.NewResponse(resp, http.StatusBadRequest, "ID cannot be empty", nil)
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		utils.NewResponse(resp, http.StatusBadRequest, "ID must be integer", nil)
		return
	}
	request.ID = idInt

	phone, err := h.PhoneUsecase.Update(ctx, request)
	if err != nil {
		utils.NewResponse(resp, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.NewResponse(resp, http.StatusOK, "Success", phone)
}

func (h *PhoneHandler) Delete(resp http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	id := mux.Vars(req)["id"]
	if id == "" {
		utils.NewResponse(resp, http.StatusBadRequest, "ID cannot be empty", nil)
		return
	}
	idInt, err := strconv.Atoi(id)
	if err != nil {
		utils.NewResponse(resp, http.StatusBadRequest, "ID must be integer", nil)
		return
	}

	phone, err := h.PhoneUsecase.Delete(ctx, idInt)
	if err != nil {
		utils.NewResponse(resp, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.NewResponse(resp, http.StatusOK, "Success", phone)
}

func (h *PhoneHandler) AutoGenerate(resp http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	err := h.PhoneUsecase.AutoGenerate(ctx)
	if err != nil {
		utils.NewResponse(resp, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.NewResponse(resp, http.StatusOK, "Success", nil)
}

func (p *PhoneHandler) ServeWs(hub *phone_ws.Hub, w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{}
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("Client Connected")

	client := &phone_ws.Client{
		Hub:  hub,
		Conn: conn,
		Send: make(chan []byte, 256),
	}
	client.Hub.Register <- client

	go client.WritePump()
}
