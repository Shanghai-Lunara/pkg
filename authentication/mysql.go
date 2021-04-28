package authentication

import (
	"database/sql"
	"github.com/Shanghai-Lunara/pkg/zaplogger"
	"time"
)

//[GIN-debug] GET    /rbac/account/login
//[GIN-debug] GET    /rbac/account/list
//[GIN-debug] POST   /rbac/account/:account/:pwd
//[GIN-debug] PUT    /rbac/account/:account/:pwd
//[GIN-debug] Delete /rbac/account/:account

//var db *sql.DB

func Query(db *sql.DB, account string) (Account, error) {
	ac := &Account{}
	if err := db.QueryRow("SELECT * FROM accounts WHERE account = ?", account).
		Scan(&ac.Account, &ac.Password, &ac.CreateTime, &ac.Status); err != nil {
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
		if err := rows.Scan(&ac.Account, &ac.Password, &ac.CreateTime, &ac.Status); err != nil {
			zaplogger.Sugar().Error(err)
			return nil, err
		}
		res = append(res, *ac)
	}
	return res, nil
}

func Add(db *sql.DB, account, password string) error {
	if _, err := db.Exec("INSERT INTO accounts (`account`,`password`,`createTime`,`status`) values (?,?,?,?)",
		account, password, time.Now().Unix(), Active); err != nil {
		zaplogger.Sugar().Error(err)
		return err
	}
	return nil
}

func ResetPassword(db *sql.DB, account, password string) error {
	if _, err := db.Query("UPDATE accounts SET password = ? WHERE account = ? AND status = ?",
		password, account, Active); err != nil {
		zaplogger.Sugar().Error(err)
		return err
	}
	return nil
}

func Disable(db *sql.DB, account string) error {
	if _, err := db.Query("UPDATE accounts SET status = ? WHERE account = ?",
		Inactive, account); err != nil {
		zaplogger.Sugar().Error(err)
		return err
	}
	return nil
}
