package model

import (
	"github.com/pkg/errors"
	"isucon8/isubank"
	"isucon8/isulogger"
	"sync"
)

const (
	BankEndpoint = "bank_endpoint"
	BankAppid    = "bank_appid"
	LogEndpoint  = "log_endpoint"
	LogAppid     = "log_appid"
)

//go:generate scanner
type Setting struct {
	Name string
	Val  string
}

func SetSetting(d QueryExecutor, k, v string) error {
	_, err := d.Exec(`INSERT INTO setting (name, val) VALUES (?, ?) ON DUPLICATE KEY UPDATE val = VALUES(val)`, k, v)
	return err
}

func GetSetting(d QueryExecutor, k string) (string, error) {
	s, err := scanSetting(d.Query(`SELECT * FROM setting WHERE name = ?`, k))
	if err != nil {
		return "", err
	}
	return s.Val, nil
}

func Isubank(d QueryExecutor) (*isubank.Isubank, error) {
	ep, err := GetSetting(d, BankEndpoint)
	if err != nil {
		return nil, errors.Wrapf(err, "getSetting failed. %s", BankEndpoint)
	}
	id, err := GetSetting(d, BankAppid)
	if err != nil {
		return nil, errors.Wrapf(err, "getSetting failed. %s", BankAppid)
	}
	return isubank.NewIsubank(ep, id)
}

func Logger(d QueryExecutor) (*isulogger.Isulogger, error) {
	ep, err := GetSetting(d, LogEndpoint)
	if err != nil {
		return nil, errors.Wrapf(err, "getSetting failed. %s", LogEndpoint)
	}
	id, err := GetSetting(d, LogAppid)
	if err != nil {
		return nil, errors.Wrapf(err, "getSetting failed. %s", LogAppid)
	}
	return isulogger.NewIsulogger(ep, id)
}

func sendLog(d QueryExecutor, tag string, v interface{}) {
	BufferedLogsMutex.Lock()
	BufferedLogs = append(BufferedLogs, LogPayload{
		Tag:   tag,
		Value: v,
	})
	BufferedLogsMutex.Unlock()
}

type LogPayload struct {
	Tag string
	Value interface{}
}

var BufferedLogs []LogPayload
var BufferedLogsMutex sync.RWMutex