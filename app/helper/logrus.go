package helper

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

type Log struct {
	Event        string
	ResponseTime time.Duration
	Response     interface{}
	Key          string
}

type Loge2e struct {
	Event        string
	StatusCode   int
	ResponseTime time.Duration
	Method       string
	Request      interface{}
	URL          string
	Message      string
	Tag          string
	Key          string
}

var log = logrus.New()

func HttpLog(level string, method string, order_id string, status int, message string, responseTime time.Duration, url string, response string) {
	log.Out = os.Stdout
	log_dir := "."

	if os.Getenv("ENV") == "production" || os.Getenv("ENV") == "preproduction" {
		log_dir = os.Getenv("DIR_LOGS")
	}

	file, err := os.OpenFile(log_dir+"/mylog.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	if err == nil {
		log.Out = file
	} else {
		log.Info("Failed to log to file, using default stderr (http log)")
	}

	payload := logrus.Fields{
		"method":        method,
		"order_id":      order_id,
		"status_code":   status,
		"resposne_time": responseTime,
		"url":           url,
		"response":      response,
	}

	if level == "info" {
		log.WithFields(payload).Info(message)
	} else if level == "warn" {
		log.WithFields(payload).Warn(message)
	} else if level == "error" {
		log.WithFields(payload).Error(message)
	}

}

func StringLog(level string, message string) {
	log.Out = os.Stdout
	log_dir := "../logs"

	if os.Getenv("ENV") == "production" || os.Getenv("ENV") == "preproduction" {
		log_dir = os.Getenv("DIR_LOGS")
	}

	file, err := os.OpenFile(log_dir+"/mylog.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	if err == nil {
		log.Out = file
	} else {
		log.Info("Failed to log to file, using default stderr (string log)")
	}

	if level == "info" {
		log.Info(message)
	} else if level == "warn" {
		log.Warn(message)
	} else if level == "error" {
		log.Error(message)
	}
}

func CreateLog(data *Log) error {

	var log = logrus.New()

	log_dir := "../logs"

	if os.Getenv("ENV") == "production" || os.Getenv("ENV") == "preproduction" {
		log_dir = os.Getenv("DIR_LOGS")
	}

	// You could set this to any `io.Writer` such as a file
	file, err := os.OpenFile(log_dir+"/mylog.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)
	if err == nil {
		log.Out = file
	} else {
		log.Info("Failed to log to file, using default stderr")
	}

	log.WithFields(logrus.Fields{
		"event":         data.Event,
		"response_time": data.ResponseTime,
		"response":      data.Response,
		"key":           data.Key,
	}).Info("LOG")

	// The API for setting attributes is a little different than the package level
	// exported logger. See Godoc.
	log.Out = os.Stdout

	return nil
}

func LogE2E(data *Loge2e, types string, order string) error {

	// You could set this to any `io.Writer` such as a file
	log_dir := "../logs"

	if os.Getenv("ENV") == "preproduction" || os.Getenv("ENV") == "production" {
		log_dir = os.Getenv("DIR_LOGS")

	}

	file, err := os.OpenFile(log_dir+"/mylog.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)
	if err == nil {
		log.Out = file
	} else {
		log.Info("Failed to log to file, using default stderr")
	}

	if types == "info" {
		log.WithFields(logrus.Fields{
			"event":         data.Event,
			"tag":           order,
			"status_code":   data.StatusCode,
			"response_time": data.ResponseTime,
			"method":        data.Method,
			"request":       data.Request,
			"url_payload":   data.URL,
			"key":           data.Key,
		}).Info(data.Message)
	}

	log.Out = os.Stdout

	return nil
}
