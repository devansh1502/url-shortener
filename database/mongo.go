package database

import (
	"context"
	"fmt"
	"log"
	"url-shortener/interfaces"
	"url-shortener/models"
	"url-shortener/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDB
type MongoDB struct {
	context           context.Context
	db                *mongo.Database
	client            *mongo.Client
	urlCollection     *mongo.Collection
	metricsCollection *mongo.Collection
}

func mongoDBConn() *mongo.Client {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Printf("Error while connecting to db, %v", err)
		return nil
	}

	return client
}

func NewMongo(ctx context.Context) interfaces.Store {
	client := mongoDBConn()
	db := client.Database("url-shortner")
	mg := &MongoDB{
		context:           ctx,
		db:                db,
		client:            client,
		urlCollection:     db.Collection("url"),
		metricsCollection: db.Collection("metrics"),
	}
	return mg
}

func (mg *MongoDB) Create(url, shortURL string) bool {
	domain := utils.GetDomain(url)
	session, err := mg.client.StartSession()
	if err != nil {
		log.Printf("Failed to start sesssion. %v", err)
		return false
	}

	err = mongo.WithSession(context.Background(), session, func(sessionContext mongo.SessionContext) error {
		if err = session.StartTransaction(); err != nil {
			return err
		}
		_, err := mg.urlCollection.InsertOne(mg.context, models.UrlCollection{URL: url, ShortURL: shortURL, Domain: domain})
		if err != nil {
			log.Printf("Error while inserting the value for %v. %v", url, err)
			return err
		}

		dm := &models.DomainMetricsCollection{}
		searchFilter := bson.M{"domain": domain}
		err = mg.metricsCollection.FindOne(mg.context, searchFilter).Decode(dm)
		if dm.Counter == 0 {
			_, err := mg.metricsCollection.InsertOne(mg.context, models.DomainMetricsCollection{Domain: domain, Counter: 1})
			if err != nil {

				log.Printf("Error while updating the counter for %v in the db. %v", domain, err)
				return err
			}
			return nil
		}
		if err != nil {
			log.Printf("Error while finding the value for %v in the db. %v", domain, err)
			return err
		}

		value := dm.Counter + 1
		update := bson.D{{"$set", bson.D{{"counter", value}}}}
		_, err = mg.metricsCollection.UpdateOne(mg.context, searchFilter, update)
		if err != nil {
			log.Printf("Error while updating the counter for %v in the db. %v", domain, err)
			return err
		}

		if err = session.CommitTransaction(sessionContext); err != nil {
			log.Printf("Failed to commit the DB transaction. %v", err)
			return err
		}
		return nil
	})

	if err != nil {
		if abortErr := session.AbortTransaction(context.Background()); abortErr != nil {
			panic(abortErr)
		}
		log.Printf("Failed to create mongo txn session. %v", err)
		return false
	}

	defer session.EndSession(mg.context)
	return true
}

func (mg *MongoDB) GetByURL(url string) string {
	urlColl := &models.UrlCollection{}
	searchFilter := bson.M{"url": url}
	err := mg.urlCollection.FindOne(mg.context, searchFilter).Decode(urlColl)
	if err != nil {
		log.Printf("Could not find %v in the database", url)
		return ""
	}

	fmt.Println("result", urlColl.ShortURL)
	return urlColl.ShortURL
}

func (mg *MongoDB) GetByShortURL(shortUrl string) string {
	urlColl := &models.UrlCollection{}
	searchFilter := bson.M{"short_url": shortUrl}
	err := mg.urlCollection.FindOne(mg.context, searchFilter).Decode(urlColl)
	if err != nil {
		log.Printf("Error while finding %v in the database", shortUrl)
		return ""
	}
	return urlColl.URL
}

func (mg *MongoDB) GetTopThreeDomains() []models.DomainMetricsCollection {
	metrics := &[]models.DomainMetricsCollection{}
	result := []models.DomainMetricsCollection{}

	opts := options.Find().SetSort(bson.D{{"counter", -1}})
	cur, err := mg.metricsCollection.Find(mg.context, bson.D{}, opts)
	if err != nil {
		log.Printf("Error getting details from the database. %v", err)
		return *metrics
	}
	err = cur.All(mg.context, metrics)
	if err != nil {
		log.Printf("Error getting all records from the database. %v", err)
		return *metrics
	}

	fmt.Println("metrics", metrics)
	i := 0
	for _, metric := range *metrics {
		if i == 3 {
			return result
		}
		result = append(result, metric)
		i++
	}

	return nil
}
