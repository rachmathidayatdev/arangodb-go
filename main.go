package main

import (
	"context"
	"fmt"

	"github.com/arangodb/config"
	"github.com/arangodb/go-driver"
)

// cara akses arango di terminal: arangod
// cara akses arango di terminal: arangosh
// cara create database: masuk ke arangosh trus ketik db._createDatabase("nama db");
// cara use database: masuk ke arangosh trus ketik db._useDatabase("nama db");
// cara create collection: masuk ke arangosh trus ketik db._create("nama collection");

func main() {
	arangoCon, arangoDB, err := config.GetArangoConnection()

	if err != nil {
		// Handle error
		fmt.Println(err.Error())
		return
	}

	// CreateDocument(arangoCon, arangoDB)
	ViewSingleDocument(arangoCon, arangoDB)
	ViewQueryDocument(arangoCon, arangoDB)
}

//Book struct
type Book struct {
	Title   string
	NoPages int
}

//CreateDocument func
func CreateDocument(arangoCon driver.Client, arangoDB driver.Database) {
	// Open Collection
	collection, err := GetCollection(arangoDB, "books")

	if err != nil {
		// Handle error
		fmt.Println(err.Error())
		return
	}

	book := Book{
		Title:   "ArangoDB Cookbook",
		NoPages: 260,
	}

	ctx := context.Background()
	meta, err := collection.CreateDocument(ctx, book)

	if err != nil {
		// Handle error
		fmt.Println(err.Error())
		return
	}

	fmt.Printf("Created document in collection '%s' in database '%s' in meta '%s'\n", collection.Name(), arangoDB.Name(), meta.Key)
}

//ViewSingleDocument func
func ViewSingleDocument(arangoCon driver.Client, arangoDB driver.Database) {
	var book Book

	collection, err := GetCollection(arangoDB, "books")

	if err != nil {
		// Handle error
		fmt.Println(err.Error())
		return
	}

	ctx := context.Background()
	meta, err := collection.ReadDocument(ctx, "2175", &book)

	if err != nil {
		// handle error
		fmt.Println(err)
		return
	}

	fmt.Println(meta)
	fmt.Println(book)
}

//ViewQueryDocument func
func ViewQueryDocument(arangoCon driver.Client, arangoDB driver.Database) {
	ctx := context.Background()
	query := "FOR d IN books FILTER d.NoPages == @nopages RETURN d"
	bindVars := map[string]interface{}{
		"nopages": 257,
	}

	cursor, err := arangoDB.Query(ctx, query, bindVars)

	if err != nil {
		// handle error
		fmt.Println(err.Error())
		return
	}

	defer cursor.Close()

	for {
		var book Book
		meta, err := cursor.ReadDocument(ctx, &book)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			// handle other errors
			fmt.Println(err.Error())
			return
		}
		fmt.Printf("Got doc with key '%s' from query\n", meta.Key)
		fmt.Println(book)
	}
}

//GetCollection func
func GetCollection(arangoDB driver.Database, collectionName string) (driver.Collection, error) {
	ctx := context.Background()
	found, err := arangoDB.CollectionExists(ctx, "books")

	if !found {
		options := &driver.CreateCollectionOptions{ /* ... */ }
		collection, err := arangoDB.CreateCollection(ctx, "books", options)

		if err != nil {
			// handle error
			fmt.Println(err.Error())
			return collection, err
		}
	}

	collection, err := arangoDB.Collection(ctx, "books")

	if err != nil {
		// Handle error
		fmt.Println(err.Error())
		return collection, err
	}

	return collection, err
}
