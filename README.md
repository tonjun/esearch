# esearch

esearch is an elasicsearch client library that supports AWS elasticsearch service

![Build Status](https://api.travis-ci.org/tonjun/esearch.svg?branch=master)

### Installation

`go get github.com/tonjun/esearch`

### Usage

```go
package main

import (
	"fmt"
	"time"

	"github.com/tonjun/esearch"
)

func main() {

	// connect to the elasticsearch server
	es := esearch.NewESearch(&esearch.Options{
		URL: "PUT_ELASTICSEARCH_SERVER_ENDPOINT_HERE",
		AWSAccessKeyID: "AWS_KEY_HERE",
		AWSSecretAccessKey: "ACCESS_KEY_HERE",
	})

	// put some documents
	err := es.Put("twitter", "tweet", "1", esearch.M{
		"user":     "kimchy",
		"postDate": "2009-11-15T13:12:00",
		"message":  "Trying out Elasticsearch, so far so good?",
	})
	if err != nil {
		panic(err)
	}

	err = es.Put("twitter", "tweet", "2", esearch.M{
		"user":     "kimchy",
		"postDate": "2009-11-15T14:12:12",
		"message":  "Another tweet, will it be indexed?",
	})
	if err != nil {
		panic(err)
	}

	// sleep for a while
	time.Sleep(1 * time.Second)

	// search
	res, err := es.Search("twitter", "tweet", esearch.M{
		"query": esearch.M{
			"match": esearch.M{
				"user": "kimchy",
			},
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("Total tweets:", res.Hits.Total)
	for _, hit := range res.Hits.Hits {
		fmt.Println("tweet: ", string(*hit.Source))
	}

	// delete the index
	err = es.DeleteIndex("twitter")
	if err != nil {
		panic(err)
	}
}
```

### License

MIT (https://opensource.org/licenses/mit-license.php)

