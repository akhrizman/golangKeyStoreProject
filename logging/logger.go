package logging

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

var logsDir = "C:/Users/Alex.Khrizman/go_logs/"
var storeLogName = "datasource.log"
var requestLogName = "htaccess.log"
var (
	RequestLogger  *log.Logger
	WarningLogger  *log.Logger
	InfoLogger     *log.Logger
	ErrorLogger    *log.Logger
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
		RequestLogger.Fatal(fileErr)
	}

	RequestLogger = log.New(requestLogFile, "", log.Ldate|log.Ltime)
}

// SetupServerLog Create and setup server log
func SetupServerLog() {
	serverLogFile, fileErr := os.OpenFile(fmt.Sprintf("%s%s", logsDir, storeLogName), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if fileErr != nil {
		ErrorLogger.Fatal(fileErr)
	}

	InfoLogger = log.New(serverLogFile, "INFO:", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLogger = log.New(serverLogFile, "WARNING:", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(serverLogFile, "ERROR:", log.Ldate|log.Ltime|log.Lshortfile)
}

// LogAppStart Display available endpoints
func LogAppStart(port int) {
	// TODO Maybe put endpoints in an slice and loop through them, or find better way to display endpoints automatically
	host := fmt.Sprintf("http://localhost:%d", port)

	// TODO Remove these extra print statements after building application
	fmt.Printf("Server available, see -")
	fmt.Printf("\n      %s", host)
	fmt.Printf("\n      %s%s", host, "/ping")
	fmt.Printf("\n      %s%s", host, "/store/{key}")
	fmt.Println()

	InfoLogger.Println("Server available, see - %s", host)
}
