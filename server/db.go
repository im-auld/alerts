package server

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
	pb "github.com/im-auld/alerts/proto"
	"hash"
	"math/big"
	"crypto/sha1"
	"fmt"
)

type Alert struct {
	RecipientId      int64   `dynamo:",omitempty"`
	Uniq             int64   `dynamo:",omitempty"`
	SenderId         int64   `dynamo:",omitempty"`
	Message          string  `dynamo:",omitempty"`
	NoticeType       string  `dynamo:",omitempty"`
	Timestamp        string  `dynamo:",omitempty"`
	ActionPath       string  `dynamo:",omitempty"`
	ObjectId         int64   `dynamo:",omitempty"`
	Archived         bool    `dynamo:",omitempty"`
	Read             bool    `dynamo:",omitempty"`
	ContentThumbnail string  `dynamo:",omitempty"`
	ContentTypeId    int64   `dynamo:",omitempty"`
	Visible          bool    `dynamo:",omitempty"`
	Ttl              float32 `dynamo:",omitempty"`
	PushTimestamp    string  `dynamo:",omitempty"`
	Seen             bool    `dynamo:",omitempty"`
}

func AlertFromProto(alert *pb.Alert) *Alert {
	return &Alert{
		RecipientId:      alert.RecipientId,
		Uniq:             alert.Uniq,
		SenderId:         alert.SenderId,
		Message:          alert.Message,
		NoticeType:       alert.NoticeType,
		Timestamp:        alert.Timestamp,
		ActionPath:       alert.ActionPath,
		ObjectId:         alert.ObjectId,
		Archived:         alert.Archived,
		Read:             alert.Read,
		ContentThumbnail: alert.ContentThumbnail,
		ContentTypeId:    alert.ContentTypeId,
		Visible:          alert.Visible,
		Ttl:              alert.Ttl,
		PushTimestamp:    alert.PushTimestamp,
		Seen:             alert.Seen,
	}
}

type DB struct {
	db    *dynamo.DB
	table dynamo.Table
}

func NewDB() DB {
	db := getDB()
	table := getTable(db)
	return DB{db: db, table: table}
}

func (db DB) SaveAlert(alert *Alert) error {
	return saveAlertToTable(alert, db.table)
}

func (db DB) GetAlertsForRecipient(recipientId int64) ([]Alert, error){
	return getUserAlerts(recipientId, db.table)
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

func getTable(db *dynamo.DB) dynamo.Table {
	table := db.Table("alerts")
	return table
}

func getDB() *dynamo.DB {
	conf := &aws.Config{
		Endpoint: aws.String("http://localstack:4569"),
		Region: aws.String("us-east-1"),
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
	ha.Write([]byte(fmt.Sprintf("notice_type:%s", alert.NoticeType)))
	if alert.ContentTypeId != 0 {
		ha.Write([]byte(fmt.Sprintf("content_type_id:%d", alert.ContentTypeId)))
	}
	if alert.ObjectId != 0 {
		ha.Write([]byte(fmt.Sprintf("object_id:%d", alert.ObjectId)))
	}
	if alert.ActionPath != "" {
		ha.Write([]byte(fmt.Sprintf("action_path:%d", alert.ActionPath)))
	}
	return ha
}

// GetUniq ...
func getUniq(alert Alert) int64 {
	ha := getAlertHash(alert)
	uniq := computeKey(ha)
	return uniq
}
