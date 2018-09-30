package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
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
	TestMode bool `yaml:"test_mode"`
	RabbitMQHost string `yaml:"rabbitmq_host"`
	RabbitMQPort uint `yaml:"rabbitmq_port"`
	RabbitMQUser string `yaml:"rabbitmq_user"`
	RabbitMQPass string `yaml:"rabbitmq_pass"`
}

var config Config
var conn *amqp.Connection
var ch amqp.Channel

func main() {

	b, err := ioutil.ReadFile("/config/configuration.yaml")
	FailOnError("failed to open configuration file", err)

	err = yaml.Unmarshal(b, &config)
	FailOnError("failed to parse configuration file", err)

	conn, err = amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%d/", config.RabbitMQUser, config.RabbitMQPass, config.RabbitMQHost, config.RabbitMQPort))
	FailOnError("failed to connect to RabbitMQ server", err)
	defer conn.Close()

	ch, err := conn.Channel()
	FailOnError("failed to open mq channel", err)

	err = ch.ExchangeDeclare("commands", "topic", true, false, false, false, nil)
	FailOnError("failed to declare exchange for commands", err)

	err = ch.ExchangeDeclare("team_join", "topic", true, false, false, false, nil)
	FailOnError("failed to declare exchange for commands", err)

	mux := http.NewServeMux()

	mux.HandleFunc("/", HandleEventRequest)
	mux.HandleFunc("/command/", HandleSlashCommand)

	srv := &http.Server{
		Addr:         ":" + strconv.Itoa(8080),
		Handler:      mux,
		ReadTimeout:  4 * time.Second,
		WriteTimeout: 4 * time.Second,
	}

	err = srv.ListenAndServe()
	FailOnError("failed to start http service", err)

	log.Println("Service stopped cleanly")
}

func FailOnError(msg string, err error) {
	if err != nil {
		log.Printf("[Error] %s: %s\n", msg, err.Error())
		panic(err)
	}
}

func HandleEventRequest(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if !config.TestMode && !IsValidSignature(r, b) {
		log.Println("[Error] Request not from Slack")
		http.Error(w, err.Error(), 500)
		return
	}

	// Unmarshal
	var msg messages.Event
	err = json.Unmarshal(b, &msg)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	EvaluateEvent(&msg)

	w.WriteHeader(http.StatusOK)
	return
}

func IsValidSignature(r *http.Request, body []byte) (bool) {
	signature := r.Header.Get("X-Slack-Signature")
	slack_timestamp := r.Header.Get("X-Slack-Request-Timestamp")
	var buffer bytes.Buffer
	buffer.WriteString("v0")
	buffer.WriteString(":")
	buffer.WriteString(slack_timestamp)
	buffer.WriteString(":")
	buffer.Write(body)
	mac := hmac.New(sha256.New, []byte(config.SignedSecret))
	mac.Write(buffer.Bytes())
	expected := mac.Sum(nil)
	return hmac.Equal(expected, []byte(signature))
}

func EvaluateEvent (e *messages.Event){
	var event_type messages.Type
	err := json.Unmarshal(e.EventInfo, &event_type)
	if err != nil {
		log.Printf("[Error] failed to parse event type to forward request: %s", err.Error())
		return
	}

	switch event_type.Type {
	case "team_join":
		log.Println("[Info] Received team_join event")
	default:
		log.Printf("[Info] Received %s event. Unable to forward", event_type.Type)
	}

	return
}

func HandleSlashCommand (w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if !config.TestMode && !IsValidSignature(r, b) {
		log.Println("[Error] Request not from Slack")
		http.Error(w, err.Error(), 500)
		return
	}

	r.ParseForm()

	switch r.FormValue("command") {
	case "/anon":
		log.Println("[Info] Received anon command")
		parseAndPublishAnon(w, r)
		break
	default:
		log.Printf("[Error] Received unknown command: %s", r.FormValue("command"))
		http.Error(w, err.Error(), 500)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}

func parseAndPublishAnon(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var anon messages.Anon

	anon.ChannelId = r.FormValue("channel_id")
	anon.ResponseUrl = r.FormValue("response_url")
	anon.Text = r.FormValue("text")
	anon.TriggerId = r.FormValue("trigger_id")

	b, err := json.Marshal(anon)
	if err != nil {
		log.Printf("[Error] Failed to encode anon command for publish: %s", err.Error())
		http.Error(w, err.Error(), 500)
		return
	}

	ch.Publish("commands",
		"command.anon",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body: b,
		})

	return
}