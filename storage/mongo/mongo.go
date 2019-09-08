package mongo

import (
	"context"
	"fully-testable-go-project/storage"
	"log"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"gopkg.in/mgo.v2/bson"
)

// Mongo - MongoHost
type Mongo struct {
	HostName        string        `json:"hostName"`
	Server          string        `json:"server"`
	Port            int           `json:"port"`
	Username        string        `json:"username"`
	Password        string        `json:"password"`
	Database        string        `json:"database"`
	IsDefault       bool          `json:"isDefault"`
	MaxIdleConns    int           `json:"maxIdleConns" `
	MaxOpenConns    int           `json:"maxOpenConns"`
	ConnMaxLifetime time.Duration `json:"connMaxLifetime" `
	IsDisabled      bool          `json:"isDisabled" `
}

type MongoClient struct {
	Client      *mongo.Client
	HostDetails Mongo
}

func InitMongoClient(hostDetails Mongo) {
	clientOption := options.Client()
	clientOption.SetHosts([]string{bindMongoServerWithPort(hostDetails.Server, hostDetails.Port)}).
		SetConnectTimeout(hostDetails.ConnMaxLifetime).
		SetMaxPoolSize(uint64(hostDetails.MaxOpenConns)).
		SetReadPreference(readpref.Primary()).
		SetDirect(true) // important if in cluster, connect to primary only.
	if hostDetails.Username != "" {
		cred := options.Credential{}
		cred.Username = hostDetails.Username
		cred.Password = hostDetails.Password
		cred.AuthSource = hostDetails.Database
		clientOption.SetAuth(cred)
	}
	client, err := mongo.NewClient(clientOption)
	if err != nil {
		log.Fatalln(err)
	}
	err = client.Connect(context.Background())
	if err != nil {
		log.Fatalln(err)
	}
	cl := MongoClient{}
	cl.Client = client
	cl.HostDetails = hostDetails
	storage.Init(cl)
}

func bindMongoServerWithPort(server string, port int) string {
	// if port is empty then used default port 27017 & bind to server ip
	var serverURI string
	if port <= 0 || strings.TrimSpace(strconv.Itoa(port)) == "" {
		serverURI = server + ":27017"
	} else {
		serverURI = server + ":" + strconv.Itoa(port)
	}
	return serverURI
}

func (mg MongoClient) SaveData(collectionName string, data interface{}) (string, error) {

	collection := mg.Client.Database(mg.HostDetails.Database).Collection(collectionName)
	opts, insertError := collection.InsertOne(context.Background(), data)
	if insertError != nil {
		return "", insertError
	}
	return opts.InsertedID.(primitive.ObjectID).Hex(), nil
}

// UpdateAll update all
func (mg MongoClient) UpdateAll(collectionName string, selector map[string]interface{}, data interface{}) error {
	session := mg.Client
	collection := session.Database(mg.HostDetails.Database).Collection(collectionName)

	_, err := collection.UpdateMany(context.Background(), selector, bson.M{"$set": data})
	if err != nil {
		return nil
	}
	return nil
}

// Update will update single entry
func (mg MongoClient) Update(collectionName string, selector map[string]interface{}, data interface{}) error {
	session := mg.Client
	collection := session.Database(mg.HostDetails.Database).Collection(collectionName)
	_, err := collection.UpdateOne(context.Background(), selector, bson.M{"$set": data})
	if err != nil {
		return nil
	}
	return nil
}

// GetData will return query for selector
func (mg MongoClient) GetData(collectionName string, selector map[string]interface{}, projector map[string]interface{}, val interface{}) error {
	session := mg.Client
	collection := session.Database(mg.HostDetails.Database).Collection(collectionName)
	ops := &options.FindOptions{}
	if projector != nil {
		ops.Projection = projector
	}
	cur, err := collection.Find(context.Background(), selector, ops)
	if err != nil {
		log.Println(err)
		return err
	}
	defer cur.Close(context.Background())
	err = cur.Decode(&val)
	if err != nil {
		return err
	}
	return nil
}

// DeleteData will delete data given for selector
func (mg MongoClient) DeleteData(collectionName string, selector map[string]interface{}) error {
	session := mg.Client
	collection := session.Database(mg.HostDetails.Database).Collection(collectionName)
	_, err := collection.DeleteOne(context.Background(), selector)
	if err != nil {
		return err
	}
	return nil
}

// DeleteAll will delete all the matching data given for selector
func (mg MongoClient) DeleteAll(collectionName string, selector map[string]interface{}) error {
	session := mg.Client
	collection := session.Database(mg.HostDetails.Database).Collection(collectionName)
	_, err := collection.DeleteMany(context.Background(), selector)
	if err != nil {
		return err
	}
	return err
}
