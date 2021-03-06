package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"runtime"
	"strconv"
	"sync"
	"time"
)

func main() {
	start := time.Now() // 获取当前时间

	runtime.GOMAXPROCS(1)//自定义核数
	//支持参数
	var (
		count int    // 起始数
		total int    // 截至数
		index string // index
		id    string // id
		title string // title
	)
	flag.IntVar(&count, "c", 1, "起始数")
	flag.IntVar(&total, "e", 1, "截至数")
	flag.StringVar(&index, "i", "", "index")
	flag.StringVar(&id, "d", "", "id")
	flag.StringVar(&title, "t", "", "title")
	// 解析参数
	flag.Parse()
	if index == "" {
		index = "demo"
	}
	if id == "" {
		id = "id_1"
	}
	if title == "" {
		title = "世界"
	}
	fmt.Println("count：", count)
	fmt.Println("total：", total)
	fmt.Println("index：", index)
	fmt.Println("id：", id)
	fmt.Println("title：", title)

	addresses := []string{"http://127.0.0.1:9200", "http://127.0.0.1:9201"}
	config := elasticsearch.Config{
		Addresses: addresses,
		Username:  "",
		Password:  "",
		CloudID:   "",
		APIKey:    "",
	}
	// new client
	es, err := elasticsearch.NewClient(config)
	if err != nil {
		fmt.Println(err, "Error creating the client")
	}

	//Get(*es, index, id)
	//Update(*es, index, id)
	//Get(*es, index, id)
	create(*es, index, count, total)
	//Search(*es, index, title)

	elapsed := time.Since(start)
	fmt.Println("该函数执行完成耗时：", elapsed)

}

func create(es elasticsearch.Client, index string, count int, total int) bool {
	var wg sync.WaitGroup
	// Create creates a new document in the index.
	// Returns a 409 response when a document with a same ID already exists in the index.
	for i := count; i < total; i++ {
		wg.Add(1)
		k := strconv.Itoa(i)
		var buf bytes.Buffer
		doc := map[string]interface{}{
			"title":   "你看到外面的世界是什么样的？" + k,
			"content": "外面的世界真的很精彩-" + k,
			"time":    time.Now().Unix(),
			"date":    time.Now(),
		}
		if err := json.NewEncoder(&buf).Encode(doc); err != nil {
			fmt.Println(err, "Error encoding doc")
			return false
		}
		go func() {
			time.Sleep(1 * time.Millisecond)
			res, err := es.Create(index, "idx_"+k, &buf)
			if err != nil {
				fmt.Println(err, "Error create response")
			}
			wg.Done()
			defer res.Body.Close()
			fmt.Println(res.String())
		}()
	}
	wg.Wait()
	return true
}

func Search(es elasticsearch.Client, index string, title string) {
	// info
	res, err := es.Info()
	if err != nil {
		fmt.Println(err, "Error getting response")
	}
	fmt.Println(res.String())
	// search - highlight
	var buf bytes.Buffer
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"title": title,
			},
		},
		"highlight": map[string]interface{}{
			"pre_tags":  []string{"<font color='red'>"},
			"post_tags": []string{"</font>"},
			"fields": map[string]interface{}{
				"title": map[string]interface{}{},
			},
		},
	}
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		fmt.Println(err, "Error encoding query")
	}
	// Perform the search request.
	res, err = es.Search(
		es.Search.WithContext(context.Background()),
		es.Search.WithIndex(index),
		es.Search.WithBody(&buf),
		es.Search.WithTrackTotalHits(true),
		es.Search.WithFrom(0),
		es.Search.WithSize(10),
		es.Search.WithSort("time:desc"),
		es.Search.WithPretty(),
	)
	if err != nil {
		fmt.Println(err, "Error getting response")
	}
	defer res.Body.Close()
	fmt.Println(res.String())
}

func DeleteByQuery(es elasticsearch.Client) {
	// DeleteByQuery deletes documents matching the provided query
	var buf bytes.Buffer
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"title": "外面",
			},
		},
	}
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		fmt.Println(err, "Error encoding query")
	}
	index := []string{"demo"}
	res, err := es.DeleteByQuery(index, &buf)
	if err != nil {
		fmt.Println(err, "Error delete by query response")
	}
	defer res.Body.Close()
	fmt.Println(res.String())
}

func Delete(es elasticsearch.Client, index string, id string) {
	// Delete removes a document from the index
	res, err := es.Delete(index, id)
	if err != nil {
		fmt.Println(err, "Error delete by id response")
	}
	defer res.Body.Close()
	fmt.Println(res.String())
}

func Get(es elasticsearch.Client, index string, id string) {
	res, err := es.Get(index, id)
	if err != nil {
		fmt.Println(err, "Error get response")
	}
	defer res.Body.Close()
	fmt.Println(res.String())
}

func Update(es elasticsearch.Client, index string, id string) {
	// Update updates a document with a script or partial document.
	var buf bytes.Buffer
	doc := map[string]interface{}{
		"doc": map[string]interface{}{
			"title":       "更新你看到外面的世界是什么样的？",
			"content":     "更新外面的世界真的很精彩",
			"update_time": time.Now().Unix(),
			"update_date": time.Now(),
		},
	}
	if err := json.NewEncoder(&buf).Encode(doc); err != nil {
		fmt.Println(err, "Error encoding doc")
	}
	res, err := es.Update(index, id, &buf)
	if err != nil {
		fmt.Println(err, "Error Update response")
	}
	defer res.Body.Close()
	fmt.Println(res.String())
}

func UpdateByQuery(es elasticsearch.Client, indexs string) {
	// UpdateByQuery performs an update on every document in the index without changing the source,
	// for example to pick up a mapping change.
	index := []string{indexs}
	var buf bytes.Buffer
	doc := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"title": "外面",
			},
		},
		// 根据搜索条件更新title
		/*
		   "script": map[string]interface{}{
		       "source": "ctx._source['title']='更新你看到外面的世界是什么样的？'",
		   },
		*/
		// 根据搜索条件更新title、content
		/*
		   "script": map[string]interface{}{
		       "source": "ctx._source=params",
		       "params": map[string]interface{}{
		           "title": "外面的世界真的很精彩",
		           "content": "你看到外面的世界是什么样的？",
		       },
		       "lang": "painless",
		   },
		*/
		// 根据搜索条件更新title、content
		"script": map[string]interface{}{
			"source": "ctx._source.title=params.title;ctx._source.content=params.content;",
			"params": map[string]interface{}{
				"title":   "看看外面的世界真的很精彩",
				"content": "他们和你看到外面的世界是什么样的？",
			},
			"lang": "painless",
		},
	}
	if err := json.NewEncoder(&buf).Encode(doc); err != nil {
		fmt.Println(err, "Error encoding doc")
	}
	res, err := es.UpdateByQuery(
		index,
		es.UpdateByQuery.WithBody(&buf),
		es.UpdateByQuery.WithContext(context.Background()),
		es.UpdateByQuery.WithPretty(),
	)
	if err != nil {
		fmt.Println(err, "Error Update response")
	}
	defer res.Body.Close()
	fmt.Println(res.String())
}
