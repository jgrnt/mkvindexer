package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/jgrnt/mkvindexer/mkvextract"
	"gopkg.in/olivere/elastic.v2"
)

const elasticHost = "http://127.0.0.1:9200"

var client *elastic.Client

func init() {
	var err error
	client, err = elastic.NewClient(elastic.SetURL(elasticHost), elastic.SetTraceLog(log.New(os.Stdout, "", log.LstdFlags)))
	if err != nil {
		panic(err)
	}
}

func ReAddMapping() {

}

func AddMovie(movie mkvextract.MkvInfo) (Result, error) {
	put, err := client.Index().
		Index("movie").
		Id(filepath.Base(movie.FileName)).
		Type("movie").
		BodyJson(movie).
		Do()
	if err != nil {
		return Error, err

	}
	if put.Created {
		return Indexed, nil
	} else {
		return AlreadyIndexed, nil
	}
}
