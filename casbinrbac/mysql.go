package casbinrbac

import (
	"database/sql"
	"fmt"
	"github.com/Shanghai-Lunara/pkg/zaplogger"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"net/url"
	"os"
	"time"
)

const (
	DBDSNFormat = "%s:%s@tcp(%s:%d)/%s?timeout=5s&readTimeout=5s&writeTimeout=5s&parseTime=true&loc=Local&charset=utf8mb4,utf8"
	//dsn = fmt.Sprintf("%s:%s@%s/%s?timeout=30s", user, pass, netAddr, dbname)
)

type MysqlConfig struct {
	Host            string `json:"host" yaml:"host"`
	Port            int    `json:"port" yaml:"port"`
	User            string `json:"user" yaml:"user"`
	Password        string `json:"password" yaml:"password"`
	Database        string `json:"database" yaml:"database"`
	MaxIdleConns    int    `json:"max_idle_conns" yaml:"max_idle_conns"`
	MaxOpenConns    int    `json:"max_open_conns" yaml:"max_open_conns"`
	ConnMaxLifetime int    `json:"conn_max_lifetime" yaml:"conn_max_lifetime"`
}

type MysqlCluster struct {
	Master MysqlConfig `yaml:"master"`
	Slave  MysqlConfig `yaml:"slave"`
}

type Mysql struct {
	Mysql MysqlCluster `yaml:"Mysql"`
}

type MysqlClusterPool struct {
	Master *sql.DB
	Slave  *sql.DB
}

func newMysqlPool(v *MysqlConfig) *sql.DB {
	dsn := fmt.Sprintf(DBDSNFormat, v.User, v.Password, v.Host, v.Port, v.Database)
	db, err := sql.Open("mysql", setDSNTimezone(dsn))
	if err != nil {
		zaplogger.Sugar().Errorw("newMysqlPool failed", "err", err)
		return nil
	}
	db.SetMaxIdleConns(v.MaxIdleConns)
	db.SetMaxOpenConns(v.MaxOpenConns)
	db.SetConnMaxLifetime(time.Second * time.Duration(v.ConnMaxLifetime))
	return db
}

var mysql *Mysql

func LoadMysqlConf(configFile string) {
	mysql = &Mysql{}
	var data []byte
	var err error
	if data, err = ioutil.ReadFile(configFile); err != nil {
		zaplogger.Sugar().Fatal(err)
	}
	if err := yaml.Unmarshal(data, mysql); err != nil {
		zaplogger.Sugar().Fatal(err)
	}
}

func GetMysqlCluster() *MysqlClusterPool {
	if mysql == nil {
		zaplogger.Sugar().Fatal("error: nil Mysql, please call casbinrbac.LoadMysqlConf(configFile string) before")
	}
	return &MysqlClusterPool{
		Master: newMysqlPool(&mysql.Mysql.Master),
		Slave:  newMysqlPool(&mysql.Mysql.Slave),
	}
}

func getContainerTimezone() string {
	if tz := os.Getenv("TZ"); tz != "" {
		return tz
	}
	return "Local"
}

func setDSNTimezone(dsn string) string {
	return dsn + "&loc=" + url.QueryEscape(getContainerTimezone())
}

func MasterDsn() string {
	if mysql == nil {
		zaplogger.Sugar().Fatal("error: nil Mysql, please call LoadMysqlConf() before")
	}
	return fmt.Sprintf(DBDSNFormat, mysql.Mysql.Master.User, mysql.Mysql.Master.Password, mysql.Mysql.Master.Host, mysql.Mysql.Master.Port, mysql.Mysql.Master.Database)
}
