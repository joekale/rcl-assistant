package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"gopkg.in/yaml.v2"
	"strconv"
	"time"
	"rcl-assistant/assistant/messages"
)

type Config struct {
	Token string `yaml:"token"`
	AppId string `yaml:"app_id"`
	SignedSecret string `yaml:"signed_secret"`
	ClientId string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
}

var config Config

func main() {

	b, err := ioutil.ReadFile("/config/configuration.yaml")
	if err != nil {
		log.Printf("[Error] failed to start http service %s\n", err.Error())
		panic(err)
	}

	err = yaml.Unmarshal(b, &config)

	mux := http.NewServeMux()

	mux.HandleFunc("/", HandleEventRequest)

	srv := &http.Server{
		Addr:         ":" + strconv.Itoa(8080),
		Handler:      mux,
		ReadTimeout:  4 * time.Second,
		WriteTimeout: 4 * time.Second,
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Printf("[Error] failed to start http service %s\n", err.Error())
		panic(err)
	}

	log.Println("Service stopped cleanly")
}

func HandleEventRequest(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	signature := r.Header.Get("X-Slack-Signature")
	slack_timestamp := r.Header.Get("X-Slack-Request-Timestamp")
	var buffer bytes.Buffer
	buffer.WriteString("v0")
	buffer.WriteString(":")
	buffer.WriteString(slack_timestamp)
	buffer.WriteString(":")
	buffer.Write(b)

	mac := hmac.New(sha256.New, []byte(config.SignedSecret))
	mac.Write(buffer.Bytes())
	expected := mac.Sum(nil)

	if !hmac.Equal(expected, []byte(signature)) {
		log.Println("[Error] Request not from Slack")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// Unmarshal
	var msg messages.Event
	err = json.Unmarshal(b, &msg)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	go Respond(&msg)

	w.WriteHeader(http.StatusOK)
	return
}

func Respond (e *messages.Event){

}