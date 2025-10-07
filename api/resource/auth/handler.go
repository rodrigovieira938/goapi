package auth

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rodrigovieira938/goapi/api/resource/users"
	"github.com/rodrigovieira938/goapi/config"
	"github.com/rodrigovieira938/goapi/util"
	"github.com/rodrigovieira938/goapi/util/db"
)

type API struct {
	db  *sql.DB
	cfg *config.AuthConfig
}

func New(db *sql.DB, cfg *config.AuthConfig) *API {
	return &API{db: db, cfg: cfg}
}

func createJWT(userID int, cfg *config.AuthConfig) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Hour * 24).Unix(), //TODO: get this time from AuthConfig
	})
	return token.SignedString([]byte(cfg.Secret))
}

func (api *API) Login(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		util.JsonError(w, "{\"error\":\"Content-Type must be application/json\"}", http.StatusBadRequest)
		return
	}
	var loginInfo Login
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&loginInfo)
	if err != nil {
		util.JsonError(w, "{\"error\":\"Invalid JSON!\"}", http.StatusBadRequest)
		return
	}
	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(loginInfo)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		if validationErrors != nil {
			util.JsonError(w, "{\"error\":\""+validationErrors.Error()+"\"}", http.StatusBadRequest)
			return
		}
	}
	row := db.GetUserByEmail(api.db, loginInfo.Email)
	var user users.User
	if err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password); err != nil {
		if err == sql.ErrNoRows {
			util.JsonError(w, "{\"error\":\"Invalid email or password\"}", http.StatusBadRequest)
			return
		}
		util.JsonError(w, "{\"error\":\"Internal Server Error\"}", http.StatusInternalServerError)
		return
	}
	if user.Password != loginInfo.Password {
		util.JsonError(w, "{\"error\":\"Invalid email or password\"}", http.StatusBadRequest)
		return
	}
	token, err := createJWT(user.ID, api.cfg)
	if err != nil {
		util.JsonError(w, "{\"error\":\"Internal Server Error\"}", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(Response{token})
}
