package golang

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

// UserRest is REST controller used to handler request about the user entities.
type UserRest struct {
	wg        sync.WaitGroup
	buildInfo BuildInfo
	hostname  string
	startedAt time.Time
	db        *sql.DB
	mu        sync.Mutex
	healthy   bool
}

// NewUserRest constructs new UserRest controller used to handler request about the user entities.
func NewUserRest(buildInfo BuildInfo, dbInfo DBInfo) (*UserRest, error) {
	db, err := sql.Open("mysql", dbInfo.ConnectionString())
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
		healthy:   true,
	}, nil
}

// Unhealthy sets server health to false - used in gracefull shutdown mode
func (r *UserRest) Unhealthy() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.healthy = false
}

// Users handler responds with list of all users in system.
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

// GetUser handler responds with particular user based on the id of the user.
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

// AddUser handler is responsible for adding new user to system.
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

// Health handler is responsible for serving resonse with current health status of service. Healthz concept is used to leverage the regular health status pattern.
func (r *UserRest) Health(w http.ResponseWriter, req *http.Request) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.healthy {
		resp := HealthStatus{
			Hostname:  r.hostname,
			StartedAt: r.startedAt.Format("2006-01-02_15:04:05"),
			Uptime:    time.Now().UTC().Sub(r.startedAt).String(),
			Build:     r.buildInfo,
		}

		WriteJSON(w, resp)
	} else {
		WriteErr(w, fmt.Errorf("Server in graceful  shutdown mode"), http.StatusInternalServerError)
	}
}

// Err handler is dummy simulator of error which occurs in out service.
func (r *UserRest) Err(w http.ResponseWriter, req *http.Request) {
	bytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		WriteErr(w, errors.New("can't read request body"), http.StatusBadRequest)
		return
	}

	WriteErr(w, errors.New(string(bytes)), http.StatusInternalServerError)
}

// User holds basic info about the User entity.
type User struct {
	ID         int       `json:"id"`
	FirstName  string    `json:"firstName"`
	SecondName string    `json:"secondName"`
	BirthDate  time.Time `json:"birthDate"`
}

// HealthStatus holds basic info about the Health status of the application.
type HealthStatus struct {
	Build     BuildInfo `json:"buildInfo"`
	Hostname  string    `json:"hostname"`
	Uptime    string    `json:"uptime"`
	StartedAt string    `json:"startedAt"`
}

// BuildInfo holds basic info about the build based on the git statistics.
type BuildInfo struct {
	Version    string `json:"version"`
	GitVersion string `json:"gitVersion"`
	BuildTime  string `json:"buildTime"`
	LastCommit Commit `json:"lastCommit"`
}

// Commit holds info about the git commit.
type Commit struct {
	ID     string `json:"id"`
	Time   string `json:"time"`
	Author string `json:"author"`
}

// DBInfo holds configuration how to access to database.
type DBInfo struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Name     string `json:"name"`
}

// ConnectionString returns connection string based on the DBInfo configuration.
func (db *DBInfo) ConnectionString() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=true", db.User, db.Password, db.Host, db.Name)
}
