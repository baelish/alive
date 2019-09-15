package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

var apiURL string
var siteURL string

func TestMain(t *testing.T) {

  t.Run("default settings", testDefaults())
  t.Run("argument processing",testArgumentProcessing())

	// Create temp dir
	rand.Seed(time.Now().UnixNano())
	tempDir := fmt.Sprintf("/tmp/alive-test-%s", randStringBytes(10))

	if _, err := os.Stat(tempDir); os.IsNotExist(err) {
		err = os.MkdirAll(tempDir, 0755)

  	if err != nil {
			panic(err)
		}
	}

	defer os.RemoveAll(tempDir)

	// Copy test data file
	testDataFile := "testdata/test-data.json"
	destinationFile := tempDir + "/data.json"
	input, err := ioutil.ReadFile(testDataFile)

	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	err = ioutil.WriteFile(destinationFile, input, 0644)

	if err != nil {
		fmt.Println("Error creating", destinationFile)
		panic(err)
	}

	// Set arguments for running api and start it
	args := []string{"", "-b", tempDir, "--api-port", "12346", "--port", "12347", "--updater", "--default-static"}
	config = getConfiguration(args)
	apiURL = fmt.Sprintf("http://127.0.0.1:%s/api/v1/", config.apiPort)
	siteURL = fmt.Sprintf("http://127.0.0.1:%s", config.sitePort)

	go main()

	time.Sleep(2 * time.Millisecond)
	t.Run("get box", testGetBox())
	t.Run("create box", testCreateBox())
	t.Run("update box", testUpdateBox())
	t.Run("send event", testSendEvent())
	t.Run("delete box", testDeleteBox())
	t.Run("load front page", testLoadingFrontPage())
	t.Run("subscribe to events", testSubscribeEvent())
  time.Sleep(2 * time.Millisecond)
}


func testDefaults() func(t *testing.T) {
  return func(t *testing.T) {
  	args := []string{}
  	config := getConfiguration(args)

  	if config.apiPort != "8081" {
  		t.Error("Expected 8081, got", config.apiPort)
  	}

  	if config.sitePort != "8080" {
  		t.Error("Expected 8080, got", config.sitePort)
  	}

  	baseDir := fmt.Sprintf("%s/.alive", os.Getenv("HOME"))

  	if config.baseDir != baseDir {
  		t.Error("Expected ", baseDir, " got ", config.baseDir)
  	}

  	if config.updater != false {
  		t.Error("Expected false got ", config.updater)
  	}

  	if config.useDefaultStatic != false {
  		t.Error("Expected false got ", config.useDefaultStatic)
  	}
  }
}

func testArgumentProcessing() func (t *testing.T) {
  return func(t *testing.T) {
  	config := &Config{}
  	args := []string{"", "-b", "/data", "-p", "1234", "--foo"}
  	config.processArguments(args)

  	if config.baseDir != "/data" {
  		t.Error("Expected /data got", config.baseDir)
  	}

  	if config.sitePort != "1234" {
  		t.Error("Expected 1234 got", config.sitePort)
  	}

  	args = []string{"", "--api-port", "1233", "--port", "1235", "--base-dir", "/var/data", "--updater", "--default-static"}
  	config.processArguments(args)

  	if config.apiPort != "1233" {
  		t.Error("Expected 1233 got", config.apiPort)
  	}

  	if config.baseDir != "/var/data" {
  		t.Error("Expected /var/data got", config.baseDir)
  	}

  	if config.sitePort != "1235" {
  		t.Error("Expected 1235 got", config.sitePort)
  	}

  	if config.updater != true {
  		t.Error("Expected true got ", config.updater)
  	}

  	if config.useDefaultStatic != true {
  		t.Error("Expected true got ", config.useDefaultStatic)
  	}
  }
}


func testGetBox() func(t *testing.T) {
	return func(t *testing.T) {
		// Get all boxes
		response, err := http.Get(apiURL)

		if err != nil {
			t.Error("Got error trying to get all boxes %s" + err.Error())
		} else {
			data, _ := ioutil.ReadAll(response.Body)
			var b []Box
			_ = json.Unmarshal(data, &b)

			if len(b) != 35 {
				t.Error(fmt.Printf("Did not receive the expected number of boxes: %d", len(b)))
			}

			if b[0].ID != "status-bar" {
				t.Error(fmt.Printf("Box 0 is not status-bar: %s", b[0].ID))
			}

			if b[2].Name != "Baboon" || !(b[2].Name < b[3].Name && b[3].Name < b[4].Name && b[4].Name < b[5].Name) {
				t.Error(fmt.Printf("Unexepected results, is sorting working correctly? %s %s %s %s", b[2].Name, b[3].Name, b[4].Name, b[5].Name))
			}
		}

		// Get specific box
		response, err = http.Get(apiURL + "10")

		if err != nil {
			t.Error("Got error trying to get createBox box %s" + err.Error())
		} else {
			data, _ := ioutil.ReadAll(response.Body)
			var b Box
			_ = json.Unmarshal(data, &b)

			if b.Name != "Crow" || b.Size != "medium" {
				t.Error(fmt.Printf("Did not receive the expected data requesting single box: \n%s", string(data)))
			}
		}

		// Get box that doesn't exist
		response, err = http.Get(apiURL + "fakeBox")

		if err != nil {
			t.Error("Got error trying to get fakeBox box %s" + err.Error())
		} else {
			data, _ := ioutil.ReadAll(response.Body)

			if strings.TrimSpace(string(data)) != `{"error":"id not found"}` {
				t.Error("Api didn't return the correct details:" + string(data))
			}
		}
	}
}

func testCreateBox() func(t *testing.T) {
	return func(t *testing.T) {
		// Test box creation
		jsonData := map[string]string{
			"id":          "testCreate",
			"name":        "testCreate",
			"size":        "dMedium",
			"status":      "grey",
			"lastMessage": "Box created",
		}

		jsonValue, _ := json.Marshal(jsonData)
		response, err := http.Post(apiURL+"new", "application/json", bytes.NewBuffer(jsonValue))

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

		// Test response to creating box that already exists
		response, err = http.Post(apiURL+"new", "application/json", bytes.NewBuffer(jsonValue))

		if err != nil {
			t.Error("Got error trying to create box through api" + err.Error())
		} else {
			data, _ := ioutil.ReadAll(response.Body)

			if strings.TrimSpace(string(data)) != `{"error":"Cannot create box, the ID requested already exists."}` {
				t.Error(fmt.Println(strings.TrimSpace(string(data))))
			}
		}

		// Test box creation without supplying ID.
		jsonData = map[string]string{
			"name":        "testCreate2",
			"size":        "small",
			"status":      "red",
			"lastMessage": "Box2 created",
			"maxTBU":      "60",
		}

		jsonValue, _ = json.Marshal(jsonData)
		response, err = http.Post(apiURL+"new", "application/json", bytes.NewBuffer(jsonValue))

		if err != nil {
			t.Error("Got error trying to create box without ID through api" + err.Error())
		} else {
			data, _ := ioutil.ReadAll(response.Body)
			var b Box
			_ = json.Unmarshal(data, &b)

			if b.ID == "" || b.Name != "testCreate2" || b.Status != "red" ||
				b.LastMessage != "Box2 created" || b.Size != "small" || b.ExpireAfter != "" ||
				b.MaxTBU != "60" {
				stringData, _ := json.Marshal(b)
				t.Error(fmt.Sprintf("Api didn't return the correct details %s", stringData))
			}
		}
	}
}

func testUpdateBox() func(t *testing.T) {
	return func(t *testing.T) {
		// Test box creation through update
		jsonData := map[string]string{
			"id":          "testUpdate",
			"name":        "testUpdateCreate",
			"size":        "dLarge",
			"status":      "grey",
			"lastMessage": "Box created",
		}

		jsonValue, _ := json.Marshal(jsonData)
		response, err := http.Post(apiURL+"update", "application/json", bytes.NewBuffer(jsonValue))

		if err != nil {
			t.Error("Got error trying to create box through api (update)" + err.Error())
		} else {
			data, _ := ioutil.ReadAll(response.Body)
			var b Box
			_ = json.Unmarshal(data, &b)

			if b.ID != "testUpdate" || b.Name != "testUpdateCreate" || b.Status != "grey" ||
				b.LastMessage != "Box created" || b.Size != "dLarge" || b.ExpireAfter != "" ||
				b.MaxTBU != "" {
				t.Error("Api didn't return the correct details")
			}
		}

		// Test response to updating box that already exists
		jsonData = map[string]string{
			"id":          "testUpdate",
			"name":        "testUpdate",
			"size":        "dMedium",
			"status":      "green",
			"lastMessage": "Box updated",
		}

		jsonValue, _ = json.Marshal(jsonData)
		response, err = http.Post(apiURL+"update", "application/json", bytes.NewBuffer(jsonValue))

		if err != nil {
			t.Error("Got error trying to update box through api" + err.Error())
		} else {
			data, _ := ioutil.ReadAll(response.Body)
			var b Box
			_ = json.Unmarshal(data, &b)

			if b.ID != "testUpdate" || b.Name != "testUpdate" || b.Status != "green" ||
				b.LastMessage != "Box updated" || b.Size != "dMedium" || b.ExpireAfter != "" ||
				b.MaxTBU != "" {
				t.Error("Api didn't return the correct details")
			}
		}

		// Test box update without supplying ID.
		jsonData = map[string]string{
			"name":        "testCreate2",
			"size":        "small",
			"status":      "red",
			"lastMessage": "Box2 created",
			"maxTBU":      "60",
		}

		jsonValue, _ = json.Marshal(jsonData)
		response, err = http.Post(apiURL+"update", "application/json", bytes.NewBuffer(jsonValue))

		if err != nil {
			t.Error("Got error trying to create box without ID through api" + err.Error())
		} else {
			data, _ := ioutil.ReadAll(response.Body)
			if strings.TrimSpace(string(data)) != `{"error":"Cannot update box without an ID."}` {
				t.Error("Api didn't return the correct details:" + string(data))
			}
		}
	}
}

func testSendEvent() func(t *testing.T) {
	return func(t *testing.T) {
		// Send event
		jsonData := map[string]string{
			"id":          "testCreate",
			"status":      "red",
			"lastMessage": "Box updated",
			"expireAfter": "60",
			"maxTBU":      "5",
			"type":        "update",
		}

		jsonValue, _ := json.Marshal(jsonData)
		response, err := http.Post(apiURL+"events/testCreate", "application/json", bytes.NewBuffer(jsonValue))

		if err != nil {
			t.Error("Got error sending update event through api" + err.Error())
		} else {
			data, _ := ioutil.ReadAll(response.Body)
			var e Event
			_ = json.Unmarshal(data, &e)

			if e.ID != "testCreate" || e.Status != "red" ||
				e.Message != "Box updated" || e.ExpireAfter != "60" ||
				e.MaxTBU != "5" {
				t.Error("Api didn't return the correct details" + string(data))
			}
		}
	}
}

func testDeleteBox() func(t *testing.T) {
	return func(t *testing.T) {
		// Delete existing box
		req, err := http.NewRequest("DELETE", apiURL+"testCreate", nil)

		if err != nil {
			t.Error("Got error trying create delete request" + err.Error())
		}

		response, err := http.DefaultClient.Do(req)

		if err != nil {
			t.Error("Got error trying to delete box through api" + err.Error())
		} else {
			data, _ := ioutil.ReadAll(response.Body)
			if strings.TrimSpace(string(data)) != `{"info":"deleted box testCreate"}` {
				t.Error("Api didn't return the correct details:" + string(data))
			}
		}

		// Delete box that doesn't exist
		req, err = http.NewRequest("DELETE", apiURL+"testCreate", nil)
		if err != nil {
			t.Error("Got error trying create delete request" + err.Error())
		}

		response, err = http.DefaultClient.Do(req)

		if err != nil {
			t.Error("Got error trying to delete box through api" + err.Error())
		} else {
			data, _ := ioutil.ReadAll(response.Body)

			if strings.TrimSpace(string(data)) != `{"error":"box not found"}` {
				t.Error("Api didn't return the correct details:" + string(data))
			}
		}
	}
}

func testLoadingFrontPage() func(t *testing.T) {
	return func(t *testing.T) {
		// Get front page
		_, err := http.Get(siteURL)

		if err != nil {
			t.Error("Got error viewing front page." + err.Error())
		}
	}
}


func removeStringFromSlice(s []string, r string) []string {
  for i, v := range s {
    if v == r {
        return append(s[:i], s[i+1:]...)
    }
  }
  return s
}



func testSubscribeEvent() func(t *testing.T) {
	return func(t *testing.T) {
    if testing.Short() {
      t.Skip("Skipping test in short mode.")
    }

		var Client = &http.Client{}
		uri := siteURL + "/events/"
		req, err := http.NewRequest("GET", uri, nil)

		if err != nil {
			t.Error("Got error connecting to events stream" + err.Error())
		}

		req.Header.Set("Accept", "text/event-stream")
		res, err := Client.Do(req)

		if err != nil {
			t.Errorf("error performing request for %s: %v", uri, err)
		}

    var event *Event
    var jsonData string
    found := []string{ "status", "red" ,"amber", "green", "grey" }

    for len(found) > 0 {
		  br := bufio.NewReader(res.Body)
		  defer res.Body.Close()
		  bs, err := br.ReadBytes('\n')

  		if err != nil && err != io.EOF {
  			t.Error(err)
  		}

      jsonData = strings.TrimLeft(string(bs), "data: ")
      err = json.Unmarshal(json.RawMessage(jsonData), &event)

      if err != nil {
        t.Errorf("Couldn't get json info from event")
      }

      if event.Type == "keepalive" {
        found = removeStringFromSlice(found, "status")
      } else {
        switch event.Status {
        case "green":
          found = removeStringFromSlice(found, "green")
        case "amber":
          found = removeStringFromSlice(found, "amber")
        case "red":
          found = removeStringFromSlice(found, "red")
        case "grey":
          found = removeStringFromSlice(found, "grey")
        }
      }
  		t.Log(jsonData)
	  }
  t.Log(found)
  }
}
