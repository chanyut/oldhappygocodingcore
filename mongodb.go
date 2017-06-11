package core

import (
	"fmt"

	"github.com/revel/revel"

	"gopkg.in/mgo.v2"
)

var (
	mainSession *mgo.Session
	mainDb      *mgo.Database
)

// MgoDb struct holds current mongo session, using Database(Db), and collection(Col)
//
type MgoDb struct {
	Session *mgo.Session
	Db      *mgo.Database
	Col     *mgo.Collection
	GridFS  *mgo.GridFS
}

// IsDup Helper function to verify that error is duplicated keys error or not
func IsDup(err error) bool {
	return mgo.IsDup(err)
}

// Init initialize new Mongo session
// host: database URI example: "mongodb://bnkvideoserver:bnkvideoserver@localhost:27017/"
// dbName: database name
func (mgoDb *MgoDb) Init(host string, dbName string) *mgo.Session {
	if mainSession == nil {
		var err error
		mongoDBHost := host + dbName
		fmt.Println("[mongodb] db host:", mongoDBHost)
		mainSession, err = mgo.Dial(mongoDBHost)
		if err != nil {
			fmt.Println("[mongodb] cannot connect to", mongoDBHost, "due to error:", err.Error())
			panic(err)
		}

		mainSession.SetMode(mgo.Monotonic, true)
		mainDb = mainSession.DB(dbName)
	}

	mgoDb.Session = mainSession.Copy()
	mgoDb.Db = mgoDb.Session.DB(dbName)
	return mgoDb.Session
}

// InitByRevelConfig initialize new Mongo session
// host: database URI example: "mongodb://bnkvideoserver:bnkvideoserver@localhost:27017/"
// dbName: database name
func (mgoDb *MgoDb) InitByRevelConfig() *mgo.Session {
	host := revel.Config.StringDefault("mongodb.host", "")
	dbName := revel.Config.StringDefault("mongodb.databasename", "")
	if host == "" || dbName == "" {
		panic("[MgoDb] failed to find database host and database name on revel config")
	}
	return mgoDb.Init(host, dbName)
}

// C implies to Collection
func (mgoDb *MgoDb) C(collection string) *mgo.Collection {
	mgoDb.Col = mgoDb.Db.C(collection)
	return mgoDb.Col
}

// GFS implies to gridFileSystem
func (mgoDb *MgoDb) GFS(fs string) *mgo.GridFS {
	mgoDb.GridFS = mgoDb.Db.GridFS(fs)
	return mgoDb.GridFS
}

// Close ...
func (mgoDb *MgoDb) Close() bool {
	defer mgoDb.Session.Close()
	return true
}

// DropDb ...
func (mgoDb *MgoDb) DropDb() {
	err := mgoDb.Db.DropDatabase()
	if err != nil {
		panic(err)
	}
}

// RemoveAll ...
func (mgoDb *MgoDb) RemoveAll(collection string) bool {
	mgoDb.Db.C(collection).RemoveAll(nil)
	databaseName := revel.Config.StringDefault("mongodb.databasename", "bnkvideoserver")
	mgoDb.Col = mgoDb.Session.DB(databaseName).C(collection)
	return true
}

// Index ...
func (mgoDb *MgoDb) Index(collection string, keys []string) bool {
	index := mgo.Index{
		Key:        keys,
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err := mgoDb.C(collection).EnsureIndex(index)
	if err != nil {
		panic(err)
	}
	return true
}