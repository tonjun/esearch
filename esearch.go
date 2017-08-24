// Package esearch is an implementation of elasticsearch client that supports
// AWS authentication for AWS Elasticsearch Service
package esearch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	//"net/http/httputil"
	"time"

	awsauth "github.com/smartystreets/go-aws-auth"
)

// ESearch represents the interface to elasticsearch
type ESearch struct {
	opts        Options
	signRequest bool
}

// M is a convenient alias for map[string]interface{}
type M map[string]interface{}

// Options is the options for creating ESearch in NewESearch
type Options struct {

	// URL is the elasticsearch server URL. e.g. http://localhost:9200
	URL string

	// AWSAccessKeyID is the optional AWS authentication access key
	AWSAccessKeyID string

	// AWSSecretAccessKey is the optional AWS authentication secret key
	AWSSecretAccessKey string
}

// NewESearch returns a new instance of ESearch given the options
func NewESearch(opts *Options) *ESearch {
	signRequest := false
	if len(opts.AWSAccessKeyID) > 0 {
		log.Printf("Using AWS signed requests: URL: %s", opts.URL)
		signRequest = true
	}
	return &ESearch{
		opts:        *opts,
		signRequest: signRequest,
	}
}

const (
	httpTimeout = 30
)

// Put inserts a document to elasticsearch.
// idx is the index, typ is the type, id is the unique ID, and data is the JSON data to insert
func (es *ESearch) Put(idx, typ, id string, data M) error {
	if len(idx) == 0 || len(typ) == 0 || len(id) == 0 {
		return fmt.Errorf("Invalid Input: idx: \"%s\" typ: \"%s\" id: \"%s\"", idx, typ, id)
	}
	b, err := json.Marshal(data)
	if err != nil {
		log.Printf("Put: json.Marshal error: %s", err.Error())
		return err
	}
	uri := fmt.Sprintf("%s/%s/%s/%s", es.opts.URL, idx, typ, id)
	req, err := http.NewRequest("POST", uri, bytes.NewBuffer(b))
	if err != nil {
		log.Printf("Insert: http.NewRequest error: %s", err.Error())
		return err
	}

	if es.signRequest {
		awsauth.Sign4(req, awsauth.Credentials{
			AccessKeyID:     es.opts.AWSAccessKeyID,
			SecretAccessKey: es.opts.AWSSecretAccessKey,
		})
	}

	client := &http.Client{
		Timeout: (httpTimeout * time.Second),
	}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	b, err = ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("Put: ioutil.ReadAll: error: %s", err.Error())
		return err
	}
	//log.Printf("elasticsearch Put response: %s", string(b))
	return nil
}

// Search searches elasticsearch
func (es *ESearch) Search(idx, typ string, query M) (*Result, error) {
	uri := fmt.Sprintf("%s/%s/%s/_search", es.opts.URL, idx, typ)
	b, err := json.Marshal(query)
	if err != nil {
		log.Printf("Search: json.Marshal error: %s", err.Error())
		return nil, err
	}
	req, err := http.NewRequest("POST", uri, bytes.NewBuffer(b))
	if err != nil {
		log.Printf("Search: http.NewRequest error: %s", err.Error())
		return nil, err
	}

	if es.signRequest {
		awsauth.Sign4(req, awsauth.Credentials{
			AccessKeyID:     es.opts.AWSAccessKeyID,
			SecretAccessKey: es.opts.AWSSecretAccessKey,
		})
	}

	client := &http.Client{
		Timeout: (httpTimeout * time.Second),
	}
	res, err := client.Do(req)
	defer res.Body.Close()
	b, err = ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("Insert: ioutil.ReadAll: error: %s", err.Error())
		return nil, err
	}
	//log.Printf("elasticsearch Search response: %s", string(b))
	if res.StatusCode < 200 || res.StatusCode > 300 {
		log.Printf("Error in response: %d", res.StatusCode)
		return nil, fmt.Errorf("%s", string(b))
	}
	result := &Result{}
	err = json.Unmarshal(b, result)
	if err != nil {
		log.Printf("Search: Unmarshal error: %s", err.Error())
		return nil, err
	}
	return result, nil
}

// DeleteIndex deletes the given index
func (es *ESearch) DeleteIndex(idx string) error {
	uri := fmt.Sprintf("%s/%s", es.opts.URL, idx)
	req, err := http.NewRequest("DELETE", uri, nil)
	if err != nil {
		log.Printf("DeleteIndex: http.NewRequest error: %s", err.Error())
		return err
	}

	if es.signRequest {
		awsauth.Sign4(req, awsauth.Credentials{
			AccessKeyID:     es.opts.AWSAccessKeyID,
			SecretAccessKey: es.opts.AWSSecretAccessKey,
		})
	}

	//dump, err := httputil.DumpRequest(req, false)
	//if err != nil {
	//	log.Printf("DumpRequest error: %s", err.Error())
	//}
	////log.Printf("Request: %q", dump)
	//fmt.Printf("%q", dump)

	client := &http.Client{
		Timeout: (httpTimeout * time.Second),
	}
	res, err := client.Do(req)
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("DeleteIndex: ioutil.ReadAll: error: %s", err.Error())
		return err
	}
	//log.Printf("elasticsearch DeleteIndex response: %s", string(b))
	if res.StatusCode < 200 || res.StatusCode > 300 {
		log.Printf("DeleteIndex: Error in response: %d", res.StatusCode)
		return fmt.Errorf("%s", string(b))
	}
	return nil
}

// DeleteQuery deletes documents using a query
func (es *ESearch) DeleteQuery(idx, typ string, query M) error {
	res, err := es.Search(idx, typ, query)
	if err != nil {
		return err
	}
	if res != nil && res.Hits != nil {
		for _, hit := range res.Hits.Hits {
			es.DeleteDocument(hit.Index, hit.Type, hit.ID)
		}
	}
	return nil
}

// DeleteDocument deletes a single document given the ID
func (es *ESearch) DeleteDocument(idx, typ, id string) error {
	if len(idx) == 0 || len(typ) == 0 || len(id) == 0 {
		return fmt.Errorf("Invalid Input: idx: \"%s\" typ: \"%s\" id: \"%s\"", idx, typ, id)
	}
	uri := fmt.Sprintf("%s/%s/%s/%s", es.opts.URL, idx, typ, id)
	req, err := http.NewRequest("DELETE", uri, nil)
	if err != nil {
		log.Printf("DeleteDocument: http.NewRequest error: %s", err.Error())
		return err
	}

	if es.signRequest {
		awsauth.Sign4(req, awsauth.Credentials{
			AccessKeyID:     es.opts.AWSAccessKeyID,
			SecretAccessKey: es.opts.AWSSecretAccessKey,
		})
	}

	client := &http.Client{
		Timeout: (httpTimeout * time.Second),
	}
	res, err := client.Do(req)
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("DeleteDocument: ioutil.ReadAll: error: %s", err.Error())
		return err
	}
	//log.Printf("elasticsearch DeleteDocument response: %s", string(b))
	if res.StatusCode < 200 || res.StatusCode > 300 {
		log.Printf("Error in response: %d", res.StatusCode)
		return fmt.Errorf("%s", string(b))
	}

	return nil
}

// RefreshIndex calls elasticsearch's _refresh API on an index
func (es *ESearch) RefreshIndex(idx string) error {
	if len(idx) == 0 {
		return fmt.Errorf("Empty index")
	}
	uri := fmt.Sprintf("%s/%s/_refresh", es.opts.URL, idx)
	req, err := http.NewRequest("POST", uri, nil)
	if err != nil {
		log.Printf("RefreshIndex: http.NewRequest error: %s", err.Error())
		return err
	}
	if es.signRequest {
		awsauth.Sign4(req, awsauth.Credentials{
			AccessKeyID:     es.opts.AWSAccessKeyID,
			SecretAccessKey: es.opts.AWSSecretAccessKey,
		})
	}
	client := &http.Client{
		Timeout: (httpTimeout * time.Second),
	}
	res, err := client.Do(req)
	if err != nil {
		log.Printf("client.Do error: %s", err.Error())
		return err
	}
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("RefreshIndex: ioutil.ReadAll: error: %s", err.Error())
		return err
	}
	log.Printf("RefreshIndex response: \"%s\"", string(b))
	return nil
}

// Update does a partial update.
// https://www.elastic.co/guide/en/elasticsearch/guide/current/partial-updates.html
func (es *ESearch) Update(idx, typ, id string, data M) error {
	if len(idx) == 0 || len(typ) == 0 || len(id) == 0 {
		return fmt.Errorf("Invalid Input: idx: \"%s\" typ: \"%s\" id: \"%s\"", idx, typ, id)
	}
	b, err := json.Marshal(data)
	if err != nil {
		log.Printf("Update: json.Marshal error: %s", err.Error())
		return err
	}
	uri := fmt.Sprintf("%s/%s/%s/%s/_update", es.opts.URL, idx, typ, id)
	req, err := http.NewRequest("POST", uri, bytes.NewBuffer(b))
	if err != nil {
		log.Printf("Update: http.NewRequest error: %s", err.Error())
		return err
	}
	if es.signRequest {
		awsauth.Sign4(req, awsauth.Credentials{
			AccessKeyID:     es.opts.AWSAccessKeyID,
			SecretAccessKey: es.opts.AWSSecretAccessKey,
		})
	}
	client := &http.Client{
		Timeout: (httpTimeout * time.Second),
	}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	_, err = ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("Update: ioutil.ReadAll: error: %s", err.Error())
		return err
	}
	//log.Printf("Update response: %s", string(b))
	return nil
}
