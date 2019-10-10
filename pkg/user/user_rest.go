package user

import (
	"encoding/json"
	"net/http"

	"rest_server/pkg/errors"
	"rest_server/pkg/rest"

	"github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
	"gopkg.in/go-playground/validator.v9"
)

type Handler struct {
	user service
}

type requestUser struct {
	Password string          `json:"password" validate:"required"`
	Email    string          `json:"email" validate:"required,email"`
	Name     string          `json:"name" validate:"required"`
	Phone    string          `json:"phone" validate:"required"`
	RegionID int64           `json:"region_id" validate:"required"`
	Meta     json.RawMessage `json:"meta,omitempty" validate:"omitempty"`
}

type responseUser struct {
	ID       int64           `json:"user_id"`
	Email    string          `json:"email"`
	Name     string          `json:"name"`
	Phone    string          `json:"phone"`
	RegionID int64           `json:"region_id"`
	Meta     json.RawMessage `json:"meta"`
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	rawJSON, err := rest.GetRequestQuery(r)
	if err != nil {
		http.Error(w, xerrors.Errorf("Get request query failed: %v", err).Error(), http.StatusBadRequest)
		return
	}

	requestData := struct {
		ID int64 `json:"user_id,string"`
	}{}
	if err := json.Unmarshal(rawJSON, &requestData); err != nil {
		logrus.Infof("Parse request params failed: %v", err)
		http.Error(w, "bad user_id", http.StatusBadRequest)
		return
	}

	if requestData.ID == 0 {
		http.Error(w, xerrors.Errorf("user_id is required").Error(), http.StatusBadRequest)
		return
	}

	usr, err := h.user.Get(r.Context(), requestData.ID)
	if err != nil {
		if err, ok := err.(errors.UserNotExistsError); ok {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := responseUser{
		ID:       usr.ID,
		Email:    usr.Email,
		Name:     usr.Name,
		Phone:    usr.Phone,
		RegionID: usr.RegionID,
		Meta:     usr.Meta,
	}

	rawUser, err := json.Marshal(&response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := rest.SendResponse(w, rawUser); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) Add(w http.ResponseWriter, r *http.Request) {
	data, err := rest.GetRequestBody(r)
	if err != nil {
		http.Error(w, xerrors.Errorf("Get request body failed: %v", err).Error(), http.StatusBadRequest)
		return
	}

	usr := requestUser{}
	if err := json.Unmarshal(data, &usr); err != nil {
		http.Error(w, xerrors.Errorf("Parse request body failed: %v", err).Error(), http.StatusBadRequest)
		return
	}

	validate := validator.New()
	if err := validate.Struct(&usr); err != nil {
		http.Error(w, xerrors.Errorf("some required params are empty").Error(), http.StatusBadRequest)
		return
	}

	dbUser := serviceUser{
		Email:    usr.Email,
		Name:     usr.Name,
		Phone:    usr.Phone,
		RegionID: usr.RegionID,
		Meta:     usr.Meta,
		Password: usr.Password,
	}

	userID, err := h.user.Insert(r.Context(), &dbUser)
	if err != nil {
		if err, ok := err.(errors.UserAlreadyExistsError); ok {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if userID == 0 {
		http.Error(w, xerrors.Errorf("impossible error").Error(), http.StatusInternalServerError)
		return
	}

	response := struct {
		ID int64 `json:"user_id"`
	}{}
	response.ID = userID

	rawResponse, err := json.Marshal(&response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := rest.SendResponse(w, rawResponse); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
