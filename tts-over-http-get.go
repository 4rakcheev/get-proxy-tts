package main

import (
	"bytes"
	"crypto/md5"
	"errors"
	"flag"
	"fmt"
	"github.com/4rakcheev/golang-tts"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"log"
	"log/syslog"
	"net/http"
	"os"
)

const (
	appNameDefault          = "get-proxy-tts"
	appHTTPServePortDefault = "80"
	cacheDirectory          = "./cache"
)

var router *gin.Engine

type ErrorApi struct {
	Code    string `json:"code"`
	Reason  string `json:"reason"`
	Message string `json:"message"`
}

func InitHTTPHandlersSimplePollySpeech() {
	router.GET("/polly", func(c *gin.Context) {

		accessKey := c.Query("access_key")
		secretKey := c.Query("secret_key")
		voice := c.Query("voice")
		lang := c.Query("language")
		text := c.Query("text")

		if accessKey == "" || secretKey == "" {
			c.JSON(http.StatusBadRequest, ErrorApi{Code:"ERR001",Reason:http.StatusText(http.StatusBadRequest),Message:"`access_key` and `secret_key` query params must be set for AWS Polly service"})
			return
		}
		if text == "" || voice == "" || lang == "" {
			c.JSON(http.StatusBadRequest, ErrorApi{Code:"ERR001",Reason:http.StatusText(http.StatusBadRequest),Message:"`voice` and `language` and `text` query must be set for AWS Polly service"})
			return
		}

		queryPolly := golang_tts.New(accessKey, secretKey)
		queryPolly.Format(golang_tts.MP3)
		queryPolly.Voice(voice)
		queryPolly.Language(lang)

		tFilePath, err := generatePollyAudio(*queryPolly, text)
		if err != nil {
			Error("speech generation failed: %s", err)
			c.String(http.StatusBadRequest, fmt.Sprintf("speech generation failed: %s", err))
			return
		}

		// Return generated to to client
		c.File(tFilePath)
	})
}

func generatePollyAudio(q golang_tts.TTS, text string) (string, error) {

	h := md5.New()
	io.WriteString(h, fmt.Sprintf("%v%s", q, text))
	tmpFile := fmt.Sprintf("%s/%x.mp3", cacheDirectory, h.Sum(nil))

	if _, err := os.Stat(tmpFile); os.IsNotExist(err) {
		Debug("cache of speech `%s` does not exist. Generate new polly speech request", text)

		// Create the file
		retBytes, err := q.Speech(text)
		if err != nil {
			return "", err
		}
		if len(retBytes) == 0 {
			return "", errors.New("polly response has zero length")
		}

		err = ioutil.WriteFile(tmpFile, retBytes, 0644)
		if err != nil {
			return "", err
		}
	} else {
		Debug("file `%s` for text `%s` already generated. Return it", tmpFile, text)
	}

	return tmpFile, nil
}

func main() {

	// Read launch params
	sPort := flag.String("p", appHTTPServePortDefault, "Port for handle requests")
	lTag := flag.String("lt", appNameDefault, "Log tag and application name (used in start message and syslog tag)")
	ls := flag.Bool("ls", true, "Syslog logging enabled")
	flag.Parse()

	// Configure logger to write to the syslog. You could do this in init(), too.
	lw := io.Writer(os.Stdout)
	if *ls == true {
		slog, _ := syslog.New(syslog.LOG_INFO, *lTag)
		lw = io.MultiWriter(os.Stdout, slog)
	}
	log.SetOutput(lw)

	// Launch the server
	// Polly service controllers
	router = gin.Default()
	// Init http handlers
	InitHTTPHandlersSimplePollySpeech()

	if err := router.Run(fmt.Sprintf(":%s", *sPort)); err != nil {
		log.Fatal(err)
		return
	}
}

// todo refactor to trdLogger
func Debug(msg string, args ...interface{}) {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, msg, args...)
	msglog := buf.String()
	log.Println("Debug: " + msglog)
}

// todo refactor to trdLogger
func Error(msg string, args ...interface{}) {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, msg, args...)
	msglog := buf.String()
	log.Println("ERROR: " + msglog)
	//panic(msglog)
}

