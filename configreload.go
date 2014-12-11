package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"gopkg.in/fsnotify.v1"
)

type Config struct {
	Name string
	Old  int
}

var (
	config *Config
	path   string
)

func ExampleNewWatcher() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				//log.Println("event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					reloadConfig()
					log.Println("modified file:", event.Name)
				}
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(path)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}

func reloadConfig() {
	file, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Println("open config: ", err)
	}

	if err = json.Unmarshal(file, config); err != nil {
		log.Println("parse config: ", err)
	}
}

func sayhelloName(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello %s, %i", config.Name, config.Old) //这个写入到w的是输出到客户端的
}

func main() {
	path = "/Users/leonzhu/Project/company/go/src/github.com/LeonZhu1981/gorp"
	config = new(Config)

	reloadConfig()
	go ExampleNewWatcher()

	http.HandleFunc("/", sayhelloName) //设置访问的路由

	err := http.ListenAndServe(":9090", nil) //设置监听的端口

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
