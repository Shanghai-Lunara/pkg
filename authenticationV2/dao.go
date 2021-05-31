package authentication

import (
	"database/sql"
	"github.com/Shanghai-Lunara/pkg/zaplogger"
	"time"
)

func appendPasswordSalt(pwd string) string {
	return pwd
}

func Query(db *sql.DB, account, password string) (Account, error) {
	ac := &Account{}
	if err := db.QueryRow("SELECT * FROM accounts WHERE account = ? AND password = ?", account, appendPasswordSalt(password)).
		Scan(&ac.Id, &ac.Account, &ac.Password, &ac.Routers, &ac.CreateTime, &ac.Status); err != nil {
		zaplogger.Sugar().Error(err)
		return *ac, err
	}
	return *ac, nil
}

func List(db *sql.DB) ([]Account, error) {
	rows, err := db.Query("SELECT * FROM accounts WHERE 1 ORDER BY id ASC")
	if err != nil {
		zaplogger.Sugar().Error(err)
		return nil, err
	}
	res := make([]Account, 0)
	for rows.Next() {
		ac := &Account{}
		if err := rows.Scan(&ac.Id, &ac.Account, &ac.Password, &ac.Routers, &ac.CreateTime, &ac.Status); err != nil {
			zaplogger.Sugar().Error(err)
			return nil, err
		}
		res = append(res, *ac)
	}
	return res, nil
}

func Add(db *sql.DB, account, password string) error {
	if _, err := db.Exec("INSERT INTO accounts (`account`,`password`,`createTime`,`status`) values (?,?,?,?)",
		account, appendPasswordSalt(password), time.Now().Unix(), Active); err != nil {
		zaplogger.Sugar().Error(err)
		return err
	}
	return nil
}

func ResetPassword(db *sql.DB, account, password string) error {
	if _, err := db.Query("UPDATE accounts SET password = ? WHERE account = ?",
		appendPasswordSalt(password), account); err != nil {
		zaplogger.Sugar().Error(err)
		return err
	}
	return nil
}

func Operator(db *sql.DB, account string, status Status) error {
	if _, err := db.Query("UPDATE accounts SET status = ? WHERE account = ?",
		status, account); err != nil {
		zaplogger.Sugar().Error(err)
		return err
	}
	return nil
}
