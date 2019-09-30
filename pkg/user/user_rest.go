package user

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
	"net/http"
	"rest_server/pkg/database"
	"rest_server/pkg/parse"
	"strconv"
)

type Service struct {
	DB database.DataStore
}

func (s *Service) Get(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Get request to user")
	rawJSON, err := parse.RequestQueryJSON(r)
	if err != nil {
		http.Error(w, xerrors.Errorf("Parse request query failed: %v", err).Error(), http.StatusBadRequest)
		return
	}

	requestUser := struct {
		ID string `json:"user_id"`
	}{}
	if err := json.Unmarshal(rawJSON, &requestUser); err != nil {
		logrus.Debugf("request params: %s", rawJSON)
		http.Error(w, xerrors.Errorf("Parse request params failed: %v", err).Error(), http.StatusBadRequest)
		return
	}
	userID, err := strconv.Atoi(requestUser.ID)
	if err != nil {
		logrus.Debugf("request params user id: %s", requestUser.ID)
		http.Error(w, xerrors.Errorf("Bad user id: %v", err).Error(), http.StatusBadRequest)
		return
	}

	dbUser, err := s.DB.GetUser(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if dbUser == nil {
		logrus.Debugf("request params user id: %s", requestUser.ID)
		http.Error(w, xerrors.Errorf("user not exists").Error(), http.StatusBadRequest)
		return
	}

	rawUser, err := json.Marshal(dbUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := parse.ResponseJSON(w, rawUser); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Service) Add(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Get request to user/add")
	data, err := parse.RequestBodyJSON(r)
	if err != nil {
		http.Error(w, xerrors.Errorf("Parse request body failed: %v", err).Error(), http.StatusBadRequest)
		return
	}
	logrus.Infof("Request body: %s", string(data))

	rawUser, err := json.Marshal(database.User{ID: 2})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := parse.ResponseJSON(w, rawUser); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
