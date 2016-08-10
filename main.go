// dto-skeleton-broker
package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/cloudfoundry-community/go-cfenv"

	"github.com/AusDTO/dto-s3-broker/internal/broker"
	"github.com/AusDTO/dto-s3-broker/internal/s3"
	"github.com/cloudfoundry-community/types-cf"
)

func mustEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("%q must be set", key)
	}
	return val
}

func envOr(key, def string) string {
	val := os.Getenv(key)
	if val == "" {
		return def
	}
	return val
}

func main() {
	port := flag.String("p", envOr("PORT", "3000"), "port to listen")
	flag.Parse()

	addr := ":" + *port

	appEnv, err := cfenv.Current()
	if err != nil {
		log.Fatal(err)
	}

	c := cf.Catalog{
		Services: []*cf.Service{{
			ID:          "d89d04f4-ece8-4bff-a0b7-e99a5b952da2",
			Name:        "dto-s3-broker",
			Description: "dto-s3-broker",
			Bindable:    true,
			Tags:        []string{"s3", "storage"},
			Plans: []*cf.Plan{{
				ID:          "e3c61d15-f74b-4589-9945-31d96f9fbba5",
				Name:        "basic",
				Description: "S3 bucket",
				Free:        true,
			}},
		}},
	}

	b := s3.S3Broker{
		Config: aws.Config{
			Region: aws.String(mustEnv("AWS_REGION")),
			Credentials: credentials.NewStaticCredentials(
				mustEnv("AWS_ACCESS_KEY"),
				mustEnv("AWS_SECRET_KEY"),
				""),
		},
	}
	api := broker.NewAPI(appEnv, &b, os.Getenv("AUTH_USER"), os.Getenv("AUTH_PASS"), &c)

	log.Println(os.Args[0], "listening on", addr)
	if err := http.ListenAndServe(addr, api); err != nil {
		log.Fatal(err)
	}
}
