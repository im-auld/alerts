package server

import (
	"crypto/sha1"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
	pb "github.com/im-auld/alerts/proto"
	"hash"
	"math/big"
	"time"
)

// Alerts have a default TTL of 7 days
const ALERT_DEFAULT_TTL int64 = 60 * 60 * 24 * 7

type Alert struct {
	RecipientId int64   `dynamo:",hash"`
	Uniq        int64   `dynamo:",range"`
	ThreadId    int64   `dynamo:",omitempty"`
	Message     string  `dynamo:",omitempty"`
	Timestamp   int64   `dynamo:",omitempty"`
	ActionPath  string  `dynamo:",omitempty"`
	Ttl         int64   `dynamo:",omitempty"`
	Seen        bool    `dynamo:",omitempty"`
}

func AlertFromProto(alert *pb.Alert) *Alert {
	return &Alert{
		RecipientId: alert.RecipientId,
		Uniq:        alert.Uniq,
		ThreadId:    alert.ThreadId,
		Message:     alert.Message,
		Timestamp:   alert.Timestamp,
		ActionPath:  alert.ActionPath,
		Ttl:         alert.Ttl,
		Seen:        alert.Seen,
	}
}

func AlertToProto(alert Alert) *pb.Alert {
	return &pb.Alert{
		RecipientId: alert.RecipientId,
		Uniq:        alert.Uniq,
		ThreadId:    alert.ThreadId,
		Message:     alert.Message,
		Timestamp:   alert.Timestamp,
		ActionPath:  alert.ActionPath,
		Ttl:         alert.Ttl,
		Seen:        alert.Seen,
	}
}

type DB struct {
	db    *dynamo.DB
	table dynamo.Table
}

func NewDB() DB {
	db := getDB()
	table := getTable(db)
	if err := table.Describe(); err != nil {
		db.CreateTable("alerts", Alert{}).Provision(2, 2).Run()
	}
	return DB{db: db, table: table}
}

func (db DB) SaveAlert(alert *Alert) error {
    alert.Timestamp = int64(time.Now().Unix())
	alert.Ttl = alert.Timestamp + ALERT_DEFAULT_TTL
	return saveAlertToTable(alert, db.table)
}

func (db DB) GetAlertsForRecipient(recipientId int64) ([]Alert, error) {
	return getUserAlerts(recipientId, db.table)
}

func (db DB) GetAlert(recipientId, uniq int64) (*Alert, error) {
    return getAlert(recipientId, uniq, db.table)
}

func (db DB) MarkAlertSeen(recipientId, uniq int64) error {
    alert, err := db.GetAlert(recipientId, uniq)
    if err != nil {
        return err
    }
    alert.Seen = true
    return db.SaveAlert(alert)
}

func saveAlertToTable(alert *Alert, table dynamo.Table) error {
	if alert.Uniq == 0 {
		alert.Uniq = getUniq(*alert)
	}
	err := table.Put(alert).Run()
	return err
}

func getUserAlerts(userId int64, table dynamo.Table) ([]Alert, error) {
	var results []Alert
	err := table.Get("RecipientId", userId).All(&results)
	return results, err
}

func getAlert(recipientId, uniq int64, table dynamo.Table) (*Alert, error) {
	var alert Alert
	err := table.Get("RecipientId", recipientId).Range("Uniq", dynamo.Equal, uniq).One(&alert)
	return &alert, err
}

func getTable(db *dynamo.DB) dynamo.Table {
	table := db.Table("alerts")
	return table
}

func getDB() *dynamo.DB {
	conf := &aws.Config{
		Endpoint: aws.String("http://localstack:4569"),
		Region:   aws.String("us-east-1"),
	}
	db := dynamo.New(session.New(), conf)
	return db
}

func computeKey(ha hash.Hash) int64 {
	n := new(big.Int).SetBytes(ha.Sum(nil))
	x := new(big.Int).Exp(big.NewInt(2), big.NewInt(32), nil)
	y := n.Mod(n, x)
	return y.Int64()
}

func getAlertHash(alert Alert) hash.Hash {
	ha := sha1.New()
	ha.Write([]byte(fmt.Sprintf("recipient_id:%d", alert.RecipientId)))
	ha.Write([]byte(fmt.Sprintf("thread_id:%s", alert.ThreadId)))
	ha.Write([]byte(fmt.Sprintf("action_path:%d", alert.ActionPath)))
	return ha
}

// GetUniq ...
func getUniq(alert Alert) int64 {
	ha := getAlertHash(alert)
	uniq := computeKey(ha)
	return uniq
}
