package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestApi(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	tempDir := fmt.Sprintf("/tmp/alive-test-%s", randStringBytes(10))
	defer os.RemoveAll(tempDir)
	args := []string{"", "-b", tempDir, "--api-port", "12346", "--port", "12347"}
	config = getConfiguration(args)
	go main()
	time.Sleep(100 * time.Millisecond)

	apiURL := fmt.Sprintf("http://127.0.0.1:%s/api/v1/new", config.apiPort)

	jsonData := map[string]string{
		"id":          "testCreate",
		"name":        "testCreate",
		"size":        "dMedium",
		"status":      "grey",
		"lastMessage": "Box created",
	}

	jsonValue, _ := json.Marshal(jsonData)
	response, err := http.Post(apiURL, "application/json", bytes.NewBuffer(jsonValue))

	if err != nil {
		t.Error("Got error trying to create box through api" + err.Error())
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		var b Box
		_ = json.Unmarshal(data, &b)

		if b.ID != "testCreate" || b.Name != "testCreate" || b.Status != "grey" ||
			b.LastMessage != "Box created" || b.Size != "dMedium" || b.ExpireAfter != "" ||
			b.MaxTBU != "" {
			t.Error("Api didn't return the correct details")
		}
	}

	response, err = http.Post(apiURL, "application/json", bytes.NewBuffer(jsonValue))

	if err != nil {
		t.Error("Got error trying to create box through api" + err.Error())
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		var e ErrorMessage
		_ = json.Unmarshal(data, &e)
		if e.Error != "Cannot create box, the ID requested already exists." {
			t.Error("Didn't get error trying to create box which already exists" + string(data))
		}
	}

	jsonData = map[string]string{
		"name":        "testCreate2",
		"size":        "small",
		"status":      "red",
		"lastMessage": "Box2 created",
		"maxTBU":      "60",
	}

	jsonValue, _ = json.Marshal(jsonData)
	response, err = http.Post(apiURL, "application/json", bytes.NewBuffer(jsonValue))

	if err != nil {
		t.Error("Got error trying to create box without api through api" + err.Error())
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		var b Box
		_ = json.Unmarshal(data, &b)

		if b.ID == "" || b.Name != "testCreate2" || b.Status != "red" ||
			b.LastMessage != "Box2 created" || b.Size != "small" || b.ExpireAfter != "" ||
			b.MaxTBU != "60" {
			t.Error("Api didn't return the correct details")
		}
	}

}
