package config

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/cluster"
	"github.com/arangodb/go-driver/http"
)

//arangoConfig struct
type arangoConfig struct {
	Username string
	Password string
	Database string
}

//GetArangoConnection func
func GetArangoConnection() (driver.Client, driver.Database, error) {
	arangoConfig := arangoConfig{
		Username: "root",
		Password: "root",
		Database: "examples_books",
	}

	var arangoCon driver.Client
	var arangoDB driver.Database

	//Get Connection
	conn, err := http.NewConnection(http.ConnectionConfig{
		Endpoints: []string{"http://localhost:8529"},
		ConnectionConfig: cluster.ConnectionConfig{
			DefaultTimeout: 60 * time.Second,
		},
		TLSConfig: &tls.Config{ /*...*/ },
	})

	if err != nil {
		// Handle error
		fmt.Println(err.Error())
		return arangoCon, arangoDB, err
	}

	//Get Client Authentication
	arangoCon, err = driver.NewClient(driver.ClientConfig{
		Connection:     conn,
		Authentication: driver.JWTAuthentication(arangoConfig.Username, arangoConfig.Password),
	})

	if err != nil {
		// Handle error
		fmt.Println(err.Error())
		return arangoCon, arangoDB, err
	}

	//Open Database
	ctx := context.Background()
	arangoDB, err = arangoCon.Database(ctx, arangoConfig.Database)

	if err != nil {
		// Handle error
		fmt.Println(err.Error())
		return arangoCon, arangoDB, err
	}

	return arangoCon, arangoDB, err
}
