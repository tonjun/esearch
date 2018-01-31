package esearch

import (
	"encoding/json"
	"log"
	//"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const AWSElasticsearch = "PUT_ES_SERVER_URL_HERE"

// NOTE: to test AWS elasticsearch credentials must be in environment variable
// AWS_ACCESS_KEY_ID
// AWS_SECRET_ACCESS_KEY

func TestESearch(t *testing.T) {

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile | log.Lmicroseconds)

	idx := "tj-test-index"

	es := NewESearch(&Options{URL: "http://localhost:9200"})
	assert.NotNil(t, es)

	es.DeleteIndex(idx)
	//assert.Nil(t, err)

	err := es.CreateIndex(idx)
	assert.Nil(t, err)
	err = es.PutMapping(idx, "test", M{
		"properties": M{
			"owner": M{
				"type": "keyword",
			},
		},
	})
	assert.Nil(t, err)

	data := M{"owner": "User-1", "message": "hello", "seen": false}

	err = es.Put(idx, "test", "doc1", data)
	assert.Nil(t, err)

	// allow to save
	es.RefreshIndex(idx)

	query := M{
		"query": M{
			"term": M{"owner": "User-1"},
		},
	}
	res, err := es.Search(idx, "test", query)
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, int64(1), res.Hits.Total)

	log.Printf("total: %d", res.Hits.Total)
	if res.Hits.Total > 0 {
		log.Printf("source: %s", string(*res.Hits.Hits[0].Source))
	}

	err = es.DeleteQuery(idx, "test", query)
	assert.Nil(t, err)

	es.RefreshIndex(idx)

	res, err = es.Search(idx, "test", query)
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, int64(0), res.Hits.Total)

	err = es.DeleteIndex(idx)
	assert.Nil(t, err)
}

func TestUpdate(t *testing.T) {

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile | log.Lmicroseconds)

	idx := "myindex"

	es := NewESearch(&Options{URL: "http://localhost:9200"})

	es.DeleteIndex(idx)

	data := M{"owner": "u2", "message": "hello there", "seen": false}

	err := es.Put(idx, "test", "doc1", data)
	assert.Nil(t, err)

	es.RefreshIndex(idx)

	query := M{
		"query": M{
			"term": M{"owner": "u2"},
		},
	}
	res, err := es.Search(idx, "test", query)
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, int64(1), res.Hits.Total)

	log.Printf("total: %d", res.Hits.Total)
	if res.Hits.Total > 0 {
		log.Printf("source: %s", string(*res.Hits.Hits[0].Source))
		r := make(map[string]interface{})
		err := json.Unmarshal(*res.Hits.Hits[0].Source, &r)
		assert.Nil(t, err)
		assert.Equal(t, r["owner"], "u2")
		assert.Equal(t, r["message"], "hello there")
		assert.Equal(t, r["seen"], false)
	}

	q := M{
		"doc": M{
			"seen": true,
		},
	}
	es.Update(idx, "test", "doc1", q)

	es.RefreshIndex(idx)

	res, err = es.Search(idx, "test", query)
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, int64(1), res.Hits.Total)
	log.Printf("total: %d", res.Hits.Total)
	if res.Hits.Total > 0 {
		log.Printf("source: %s", string(*res.Hits.Hits[0].Source))
		r := make(map[string]interface{})
		err := json.Unmarshal(*res.Hits.Hits[0].Source, &r)
		assert.Nil(t, err)
		assert.Equal(t, r["owner"], "u2")
		assert.Equal(t, r["message"], "hello there")
		assert.Equal(t, r["seen"], true)
	}

	err = es.DeleteIndex(idx)
	assert.Nil(t, err)
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
