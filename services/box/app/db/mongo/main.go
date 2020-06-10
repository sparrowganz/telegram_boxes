package mongo

import (
	"time"
	"gopkg.in/mgo.v2"
)

type Mongo interface {
	Getter
	Connector
	Models
	createCollections()
}

type Getter interface {
	GetMainSession() *mgo.Session
	GetDatabaseName() string
	Models() Models
}

type DB struct {
	models       Models
	databaseName string
	session      *mgo.Session
}

func CreateMongoDB(host, port, username, password, database, mechanism string) (*DB, error) {
	m := &DB{
		models: createModels(database),
	}
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

	db.session, err = mgo.DialWithInfo(info)
	if err != nil {
		db.session, err = mgo.Dial(host)
		if err != nil {
			return err
		}
	}

	db.session.SetMode(mgo.Monotonic, false)
	err = db.session.Ping()
	if err != nil {
		return err
	}

	err = db.session.Login(&mgo.Credential{
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

func (db *DB) Models() Models {
	return db.models
}

//Close connect to MongoDB
func (db *DB) Close() {
	db.session.Close()
}

func (db *DB) createCollections() {
	_ = db.models.Users().queryUsers(db.session).Create(&mgo.CollectionInfo{})
}

func (db *DB) GetMainSession() *mgo.Session {
	return db.session
}

func (db *DB) GetDatabaseName() string {
	return db.databaseName
}
