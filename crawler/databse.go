package crawler

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log/slog"
	"time"
)

type Client struct {
	*mongo.Client
}

func ConnectDB() *Client {
	return &Client{
		MustGetClient(),
	}
}

func MustGetClient() *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	databaseURL := fmt.Sprintf("mongodb://%s:%s@%s:%s", App.Config.Database.Username, App.Config.Database.Password, App.Config.Database.Host, App.Config.Database.Port)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(databaseURL))
	if err != nil {
		panic(err)
	}

	return client
}
func (client *Client) GetCollection(collectionName string) *mongo.Collection {
	collection := client.Database(App.Config.Site.Name).Collection(collectionName)
	ensureUniqueIndex(collection)
	return collection
}

func ensureUniqueIndex(collection *mongo.Collection) {
	indexModel := mongo.IndexModel{
		Keys:    bson.M{"url": 1},
		Options: options.Index().SetUnique(true),
	}

	_, err := collection.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		slog.Error(fmt.Sprintf("Could not create index: %v", err))
	}
}
func (client *Client) Insert(urlCollections []UrlCollection, parent string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var documents []interface{}
	for _, urlCollection := range urlCollections {
		urlCollection := UrlCollection{
			Url:       urlCollection.Url,
			Parent:    parent,
			Status:    false,
			Error:     false,
			MetaData:  urlCollection.MetaData,
			Attempts:  0,
			CreatedAt: time.Now(),
			UpdatedAt: nil,
		}
		documents = append(documents, urlCollection)
	}

	opts := options.InsertMany().SetOrdered(false)

	collection := client.GetCollection(App.GetCollection())
	_, _ = collection.InsertMany(ctx, documents, opts)

}
func (client *Client) NewSite() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var documents interface{}

	documents = SiteCollection{
		Url:       App.Config.Site.Url,
		BaseUrl:   App.Config.Site.BaseUrl,
		Status:    false,
		Attempts:  0,
		StartedAt: time.Now(),
		EndedAt:   nil,
	}

	collection := client.GetCollection(App.GetCollection())
	_, _ = collection.InsertOne(ctx, documents)

}
func (client *Client) SaveProductDetail(productDetail *ProductDetail) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	collection := client.GetCollection(App.GetCollection())
	collection.ReplaceOne(ctx, bson.D{{Key: "url", Value: productDetail.Url}}, productDetail, options.Replace().SetUpsert(true))
	defer cancel()

}

func (client *Client) GetUrlsFromCollection(collection string) []string {
	filterCondition := bson.D{
		{Key: "status", Value: false},
		{Key: "attempts", Value: bson.D{{Key: "$lt", Value: 5}}},
	}
	return extractUrls(filterData(filterCondition, client.GetCollection(collection)))
}

func (client *Client) GetUrlCollections(collection string) []UrlCollection {
	filterCondition := bson.D{
		{Key: "status", Value: false},
		{Key: "attempts", Value: bson.D{{Key: "$lt", Value: 5}}},
	}
	return filterUrlData(filterCondition, client.GetCollection(collection))
}

func filterUrlData(filterCondition bson.D, mongoCollection *mongo.Collection) []UrlCollection {
	findOptions := options.Find().SetLimit(1000) // TODO: need to refactor

	cursor, err := mongoCollection.Find(context.TODO(), filterCondition, findOptions)
	if err != nil {
		slog.Error(err.Error())
	}

	var results []UrlCollection
	if err = cursor.All(context.TODO(), &results); err != nil {
		slog.Error(err.Error())
	}

	return results
}
func filterData(filterCondition bson.D, mongoCollection *mongo.Collection) []bson.M {
	findOptions := options.Find().SetLimit(1000) // TODO: need to refactor

	cursor, err := mongoCollection.Find(context.TODO(), filterCondition, findOptions)
	if err != nil {
		slog.Error(err.Error())
	}

	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		slog.Error(err.Error())
	}

	return results
}
func extractUrls(results []bson.M) []string {
	var urls []string
	for _, result := range results {
		if url, ok := result["url"].(string); ok {
			urls = append(urls, url)
		}
	}
	return urls
}
