package esearch

import (
	"log"
	//"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const AWSElasticsearch = "PUT_ES_SERVER_URL_HERE"

// NOTE: to test AWS elasticsearch credentials must be in environment variable
// AWS_ACCESS_KEY_ID
// AWS_SECRET_ACCESS_KEY

func TestESearch(t *testing.T) {

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile | log.Lmicroseconds)

	es := NewESearch(&Options{URL: "http://localhost:9200"})
	assert.NotNil(t, es)

	err := es.DeleteIndex("tj-test-index")
	assert.Nil(t, err)

	data := M{"owner": "u1", "message": "hello"}

	err = es.Put("tj-test-index", "test", "doc1", data)
	assert.Nil(t, err)

	// allow to save
	time.Sleep(1 * time.Second)

	query := M{
		"query": M{
			"term": M{"owner": "u1"},
		},
	}
	res, err := es.Search("tj-test-index", "test", query)
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, int64(1), res.Hits.Total)

	log.Printf("total: %d", res.Hits.Total)
	if res.Hits.Total > 0 {
		log.Printf("source: %s", string(*res.Hits.Hits[0].Source))
	}

	err = es.DeleteQuery("tj-test-index", "test", query)
	assert.Nil(t, err)

	time.Sleep(1 * time.Second)

	res, err = es.Search("tj-test-index", "test", query)
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, int64(0), res.Hits.Total)

}

/*
func TestAWSEsearch(t *testing.T) {

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile | log.Lmicroseconds)

	es := NewESearch(&Options{
		URL:                AWSElasticsearch,
		AWSAccessKeyID:     os.Getenv("AWS_ACCESS_KEY_ID"),
		AWSSecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
	})
	assert.NotNil(t, es)

	var index = "tj-test-aws-index"
	var typ = "test"

	err := es.DeleteIndex(index)
	//assert.Nil(t, err)

	data := M{"owner": "u2", "message": "hello"}

	err = es.Put(index, typ, "docA", data)
	assert.Nil(t, err)

	// allow to save
	time.Sleep(1 * time.Second)

	query := M{
		"query": M{
			"term": M{"owner": "u2"},
		},
	}
	res, err := es.Search(index, typ, query)
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, int64(1), res.Hits.Total)

	log.Printf("total: %d", res.Hits.Total)
	if res.Hits.Total > 0 {
		log.Printf("source: %s", string(*res.Hits.Hits[0].Source))
	}

	err = es.DeleteQuery(index, typ, query)
	assert.Nil(t, err)

	time.Sleep(1 * time.Second)

	res, err = es.Search(index, typ, query)
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, int64(0), res.Hits.Total)

}
*/
