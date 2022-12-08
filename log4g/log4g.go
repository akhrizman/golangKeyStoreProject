package log4g

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	logsDir        = "C:/Users/Alex.Khrizman/go_logs/"
	storeLogName   = "datasource.log"
	requestLogName = "htaccess.log"
)

var (
	Request        *log.Logger
	Warning        *log.Logger
	Info           *log.Logger
	Error          *log.Logger
	RequestLogFile *os.File
	ServerLogFile  *os.File
)

func SetupLogFiles() {
	SetupRequestLog()
	SetupServerLog()
}

func CloseLogFiles() {
	RequestLogFile.Close()
	ServerLogFile.Close()
}

type RequestLogEntry struct {
	SourceIP   string    `json:"source IP"`
	HttpMethod string    `json:"HTTP method"`
	Url        string    `json:"URL"`
	Time       time.Time `json:"time of request"`
}

func (entry RequestLogEntry) String() string {
	reqDetails, _ := json.Marshal(entry)
	return string(reqDetails)
}

func NewRequestLogEntry(request *http.Request) RequestLogEntry {
	return RequestLogEntry{
		SourceIP:   request.RemoteAddr,
		HttpMethod: request.Method,
		Url:        request.URL.String(),
		Time:       time.Now()}
}

// SetupRequestLog Create and setup request log
func SetupRequestLog() {
	requestLogFile, fileErr := os.OpenFile(fmt.Sprintf("%s%s", logsDir, requestLogName), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if fileErr != nil {
		Request.Fatal(fileErr)
	}

	Request = log.New(requestLogFile, "", log.Ldate|log.Ltime)
}

// SetupServerLog Create and setup server log
func SetupServerLog() {
	serverLogFile, fileErr := os.OpenFile(fmt.Sprintf("%s%s", logsDir, storeLogName), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if fileErr != nil {
		Error.Fatal(fileErr)
	}

	Info = log.New(serverLogFile, "INFO:", log.Ldate|log.Ltime|log.Lshortfile)
	Warning = log.New(serverLogFile, "WARNING:", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(serverLogFile, "ERROR:", log.Ldate|log.Ltime|log.Lshortfile)
}
