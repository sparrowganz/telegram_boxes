package mongo

import (
	"gopkg.in/mgo.v2"
	"time"
)

type Mongo interface {
	Getter
	Models
	Connector
	createCollections()
}

type Getter interface {
	GetMainSession() *mgo.Session
	GetDatabaseName() string
}

type Models interface {
	Bots
	Tasks
	Payments
	Broadcasts
}

type DB struct {
	DatabaseName string
	Session      *mgo.Session
}

func CreateMongoDB(host, port, username, password, database, mechanism string) (*DB, error) {
	m := &DB{}
	return m, m.Connect(host, port, username, password, database, mechanism)
}

type Connector interface {
	Connect(host, port, username, password, database, mechanism string) (err error)
	Close()
}

//Connect to MongoDB
func (db *DB) Connect(host, port, username, password, database, mechanism string) (err error) {

	info := &mgo.DialInfo{
		Addrs:    []string{host + ":" + port},
		Username: username,
		Password: password,
		Timeout:  60 * time.Second,
		Database: database,
	}

	db.Session, err = mgo.DialWithInfo(info)
	if err != nil {
		db.Session, err = mgo.Dial(host)
		if err != nil {
			return err
		}
	}

	db.Session.SetMode(mgo.Monotonic, false)
	err = db.Session.Ping()
	if err != nil {
		return err
	}

	err = db.Session.Login(&mgo.Credential{
		Username:  username,
		Password:  password,
		Source:    database,
		Mechanism: mechanism,
	})
	if err != nil {
		return err
	}

	return nil
}

//Close connect to MongoDB
func (db *DB) Close() {
	db.Session.Close()
}

func (db *DB) createCollections() {
	_ = db.queryBot(db.Session).Create(&mgo.CollectionInfo{})
	_ = db.queryPayments(db.Session).Create(&mgo.CollectionInfo{})
	_ = db.queryTasks(db.Session).Create(&mgo.CollectionInfo{})

}

func (db *DB) GetMainSession() *mgo.Session {
	return db.Session
}

func (db *DB) GetDatabaseName() string {
	return db.DatabaseName
}
