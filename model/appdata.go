package model

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/syhlion/go-common"
	"github.com/syhlion/gopusher/module/log"
	"strings"
	"time"
)

var (
	appData *AppData = nil
)

type AppDataResult struct {
	AppKey       string `json:"app_key"`
	AppName      string `json:"app_name"`
	AuthAccount  string `json:"auth_account"`
	AuthPassword string `json:"auth_password"`
	RequestIP    string `json:"request_ip"`
	Date         string `json:"date"`
	Timestamp    string `json:"timestamp"`
}
type AppData struct {
	sqlDir string
}

func NewAppData(sqlDir string) *AppData {
	return &AppData{sqlDir}
}

func (d *AppData) IsExist(app_key string) bool {
	db, err := sql.Open("sqlite3", d.sqlDir)
	defer db.Close()
	if err != nil {
		log.Logger.Fatal(err)
	}
	sql := "SELECT EXISTS(SELECT 1  FROM  `appdata` WHERE `app_key`= $1)"
	var result int
	err = db.QueryRow(sql, app_key).Scan(&result)
	if err != nil {
		log.Logger.Debug(app_key, " ", err)
		return false
	}

	if result == 0 {
		log.Logger.Debug(app_key, " no exist")
		return false
	}
	return true

}

func (d *AppData) Delete(app_key string) (err error) {
	db, err := sql.Open("sqlite3", d.sqlDir)
	defer db.Close()
	if err != nil {
		log.Logger.Fatal(err)
	}

	sql := "DELETE FROM `appdata` where app_key = ?"

	tx, err := db.Begin()
	if err != nil {
		log.Logger.Debug(err)
		return
	}
	stmt, err := tx.Prepare(sql)
	if err != nil {
		log.Logger.Debug(app_key, " ", err)
		return
	}

	_, err = stmt.Exec(app_key)
	if err != nil {
		log.Logger.Debug(app_key, " ", err)
		return
	}
	err = tx.Commit()
	if err != nil {
		log.Logger.Debug(app_key, " ", err)
		return
	}
	return

}

func (d *AppData) Get(app_key string) (r AppDataResult, err error) {

	db, err := sql.Open("sqlite3", d.sqlDir)
	defer db.Close()
	if err != nil {
		log.Logger.Fatal(err)
	}
	sql := "SELECT * FROM `appdata` WHERE app_key = ?"
	rows, err := db.Query(sql)
	if err != nil {
		log.Logger.Debug(err)
		return
	}
	for rows.Next() {
		err = rows.Scan(&r.AppName, &r.AuthAccount, &r.AuthPassword, &r.RequestIP, &r.AppKey, &r.Timestamp, &r.Date)
		if err != nil {
			log.Logger.Debug(err)
			return
		}
	}
	return
}

func (d *AppData) GetAll() (r []AppDataResult, err error) {

	db, err := sql.Open("sqlite3", d.sqlDir)
	defer db.Close()
	if err != nil {
		log.Logger.Fatal(err)
	}
	sql := "SELECT * FROM `appdata`"
	rows, err := db.Query(sql)
	if err != nil {
		log.Logger.Debug(err)
		return
	}
	var apps AppDataResult
	for rows.Next() {
		err = rows.Scan(&apps.AppName, &apps.AuthAccount, &apps.AuthPassword, &apps.RequestIP, &apps.AppKey, &apps.Timestamp, &apps.Date)
		if err != nil {
			log.Logger.Debug(err)
			return
		}
		r = append(r, apps)
	}
	return

}

func (d *AppData) Register(app_name string, auth_account string, auth_password string, request_ip string) (app_key string, err error) {
	db, err := sql.Open("sqlite3", d.sqlDir)
	defer db.Close()
	if err != nil {
		log.Logger.Fatal(err)
	}
	cmd := "INSERT INTO appdata(app_name,auth_account,auth_password,request_ip,app_key,timestamp,date) VALUES (?,?,?,?,?,?,?)"
	tx, err := db.Begin()
	if err != nil {
		log.Logger.Debug(err)
		return
	}
	stmt, err := tx.Prepare(cmd)
	if err != nil {
		log.Logger.Debug(app_name, " ", request_ip, " ", err)
		return
	}
	date := time.Now().Format("2006/01/02 15:04:05")

	seeds := []string{app_name, auth_account, auth_password, request_ip, common.TimeToString(), date}
	seed := strings.Join(seeds, ",")
	app_key = common.EncodeMd5(seed)

	log.Logger.Info(app_key)
	_, err = stmt.Exec(app_name, auth_account, auth_password, request_ip, app_key, common.Time(), date)
	if err != nil {
		log.Logger.Debug(app_name, " ", request_ip, " ", err)
		return
	}
	err = tx.Commit()
	if err != nil {
		log.Logger.Debug(app_name, " ", request_ip, " ", err)
		return
	}
	defer stmt.Close()
	return
}