package golang

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type UserRest struct {
	wg        sync.WaitGroup
	buildInfo BuildInfo
	hostname  string
	startedAt time.Time
	db        *sql.DB
}

func NewUserRest(buildInfo BuildInfo) (*UserRest, error) {
	db, err := sql.Open("mysql", "root:7J52xZ0B9V@tcp(mysql-mysql:3306)/users?charset=utf8&parseTime=true")
	if err != nil {
		return nil, err
	}

	host, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	return &UserRest{
		db:        db,
		buildInfo: buildInfo,
		hostname:  host,
		startedAt: time.Now().UTC(),
	}, nil
}

func (r *UserRest) Users(w http.ResponseWriter, req *http.Request) {
	rows, err := r.db.Query("SELECT * FROM users")
	if err != nil {
		WriteErr(w, err, http.StatusInternalServerError)
		return
	}

	users := make([]User, 0, 0)
	for rows.Next() {
		var uid int
		var firstName string
		var secondName string
		var birthDate time.Time
		err = rows.Scan(&uid, &firstName, &secondName, &birthDate)
		if err != nil {
			WriteErr(w, err, http.StatusInternalServerError)
			return
		}

		users = append(users, User{ID: uid, FirstName: firstName, SecondName: secondName, BirthDate: birthDate})
	}

	if rows.Err() != nil {
		WriteErr(w, rows.Err(), http.StatusInternalServerError)
		return
	}

	WriteJSON(w, users)
}

func (r *UserRest) GetUser(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	if id == "" {
		WriteErr(w, errors.New("please provide User id"), http.StatusBadRequest)
		return
	}

	row := r.db.QueryRow("SELECT * FROM users WHERE id = ?", id)
	var uid int
	var firstName string
	var secondName string
	var birthDate time.Time
	err := row.Scan(&uid, &firstName, &secondName, &birthDate)
	if err != nil {
		if err == sql.ErrNoRows {
			WriteErr(w, errors.New("can't find user with id: "+id), http.StatusNotFound)
			return
		}

		WriteErr(w, err, http.StatusInternalServerError)
		return
	}

	user := User{ID: uid, FirstName: firstName, SecondName: secondName, BirthDate: birthDate}

	WriteJSON(w, user)
}

func (r *UserRest) AddUser(w http.ResponseWriter, req *http.Request) {
	var user User
	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		WriteErr(w, errors.New("can't deserialize json with user"), http.StatusBadRequest)
		return
	}

	_, err = r.db.Exec("INSERT INTO users(firstName, secondName, birthDate) values(?,?,?)", user.FirstName, user.SecondName, user.BirthDate)
	if err != nil {
		WriteErr(w, err, http.StatusInternalServerError)
		return
	}
}

func (r *UserRest) Health(w http.ResponseWriter, req *http.Request) {
	resp := HealthStatus{
		Hostname:  r.hostname,
		StartedAt: r.startedAt.Format("2006-01-02_15:04:05"),
		Uptime:    time.Now().UTC().Sub(r.startedAt).String(),
		Build:     r.buildInfo,
	}

	WriteJSON(w, resp)
}

func (r *UserRest) Err(w http.ResponseWriter, req *http.Request) {
	bytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		WriteErr(w, errors.New("can't read request body"), http.StatusBadRequest)
		return
	}

	WriteErr(w, errors.New(string(bytes)), http.StatusInternalServerError)
}

type User struct {
	ID         int       `json:"id"`
	FirstName  string    `json:"firstName"`
	SecondName string    `json:"secondName"`
	BirthDate  time.Time `json:"birthDate"`
}

type HealthStatus struct {
	Build     BuildInfo `json:"buildInfo"`
	Hostname  string    `json:"hostname"`
	Uptime    string    `json:"uptime"`
	StartedAt string    `json:"startedAt"`
}

type BuildInfo struct {
	Version    string `json:"version"`
	GitVersion string `json:"gitVersion"`
	BuildTime  string `json:"buildTime"`
	LastCommit Commit `json:"lastCommit"`
}

type Commit struct {
	Id     string `json:"id"`
	Time   string `json:"time"`
	Author string `json:"author"`
}
