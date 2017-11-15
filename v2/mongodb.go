package core

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/revel/revel"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
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
		mongoDBHost := host
		log.Println("[mongodb::Init] db host:", mongoDBHost)

		useSSL := strings.Contains(host, "ssl=true")
		if useSSL {
			log.Println("[mongodb::Init] use SSL")

			// try to remove ssl option from host uri
			host = strings.Replace(host, "&ssl=true", "", 1)
			host = strings.Replace(host, "?ssl=true", "", 1)
			log.Println("[mongodb::Init] clean up host uri:", host)

			// set tls config...
			tlsConfig := &tls.Config{}
			tlsConfig.InsecureSkipVerify = true

			// dial...
			dialInfo, err := mgo.ParseURL(host)
			if err != nil {
				log.Println("[mongodb::Init] failed to parse host uri")
				panic(err)
			}
			dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
				conn, err := tls.Dial("tcp", addr.String(), tlsConfig)
				return conn, err
			}

			mainSession, err = mgo.DialWithInfo(dialInfo)
		} else {
			if mongoDBHost[len(mongoDBHost)-1:] != "/" {
				mongoDBHost = mongoDBHost + "/"
			}
			mainSession, err = mgo.Dial(mongoDBHost + dbName)
		}

		if err != nil {
			log.Println("[mongodb::Init] cannot connect to", mongoDBHost, "due to error:", err.Error())
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

func (mgoDb *MgoDb) UploadFile(fileCollectionName string, fileName string, data []byte) (*mgo.GridFile, int, error) {
	// save photo source file into database
	gfsFile, err := mgoDb.GFS(fileCollectionName).Create(fileName)
	if err != nil {
		gfsFile.Close()
		err = fmt.Errorf("[MongoDB::Uploadfile] failed to create file on database due to error: %v", err)
		return nil, 0, err
	}
	n, err := gfsFile.Write(data)
	if err != nil {
		gfsFile.Close()
		err := fmt.Errorf("[MongoDB::Uploadfile] failed to write filedue to error: %v", err)
		return nil, 0, err
	}
	println("[MongoDB::Uploadfile] finished uploading file[", gfsFile.Id().(bson.ObjectId).Hex(), "] to database...", n, "bytes")
	err = gfsFile.Close()
	if err != nil {
		return nil, 0, fmt.Errorf("[MongoDB::UploadFile] failed to close file")
	}
	return gfsFile, n, nil
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
