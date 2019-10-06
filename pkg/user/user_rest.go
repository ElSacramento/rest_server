package user

import (
	"encoding/json"
	"net/http"
	"rest_server/pkg/database"
	"rest_server/pkg/rest"

	"github.com/sirupsen/logrus"

	"gopkg.in/go-playground/validator.v9"

	"golang.org/x/xerrors"
)

type Service struct {
	DB database.DataStore
}

func (s *Service) Get(w http.ResponseWriter, r *http.Request) {
	rawJSON, err := rest.GetRequestQuery(r)
	if err != nil {
		http.Error(w, xerrors.Errorf("Get request query failed: %v", err).Error(), http.StatusBadRequest)
		return
	}

	requestUser := struct {
		ID int64 `json:"user_id,string"`
	}{}
	if err := json.Unmarshal(rawJSON, &requestUser); err != nil {
		logrus.Infof("Parse request params failed: %v", err)
		http.Error(w, "bad user_id", http.StatusBadRequest)
		return
	}

	if requestUser.ID == 0 {
		http.Error(w, xerrors.Errorf("user_id is required").Error(), http.StatusBadRequest)
		return
	}

	dbUser, err := s.DB.GetUser(r.Context(), requestUser.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if dbUser == nil {
		http.Error(w, xerrors.Errorf("user not exists").Error(), http.StatusBadRequest)
		return
	}

	rawUser, err := json.Marshal(dbUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := rest.SendResponse(w, rawUser); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Service) Add(w http.ResponseWriter, r *http.Request) {
	data, err := rest.GetRequestBody(r)
	if err != nil {
		http.Error(w, xerrors.Errorf("Get request body failed: %v", err).Error(), http.StatusBadRequest)
		return
	}

	acc := database.Account{}
	if err := json.Unmarshal(data, &acc); err != nil {
		http.Error(w, xerrors.Errorf("Parse request body failed: %v", err).Error(), http.StatusBadRequest)
		return
	}

	validate := validator.New()
	if err := validate.Struct(&acc); err != nil {
		http.Error(w, xerrors.Errorf("some required params are empty").Error(), http.StatusBadRequest)
		return
	}

	userID, err := s.DB.InsertUser(r.Context(), &acc)
	if err != nil {
		if err, ok := err.(database.UserAlreadyExistsError); ok {
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
