package main

import (
	"encoding/json"
	"flag"
	"net/http"
	"os"
	"strconv"
	"time"
	"log"
	"bytes"
	"github.com/YaroslavChirko/design-practice-3-template/httptools"
	"github.com/YaroslavChirko/design-practice-3-template/signal"
)

var port = flag.Int("port", 8080, "server port")
var Traffic int = 0;
var ill bool = false
const confResponseDelaySec = "CONF_RESPONSE_DELAY_SEC"
const confHealthFailure = "CONF_HEALTH_FAILURE"

func main() {
	h := new(http.ServeMux)

	go func() {
			for range time.Tick(time.Hour) {
				Traffic=0
			}
		}()	
	
	h.HandleFunc("/health", func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("content-type", "text/plain")
		if failConfig := os.Getenv(confHealthFailure); failConfig == "true" {
			rw.WriteHeader(http.StatusInternalServerError)
			_, _ = rw.Write([]byte("FAILURE"))
		} else if(ill){
			rw.WriteHeader(http.StatusInternalServerError)
			_, _ = rw.Write([]byte("FAILURE"))
		}else{
		var n int
			rw.WriteHeader(http.StatusOK)
			n, _ = rw.Write([]byte("OK"))
			Traffic += n
			
		}
	})
	
	//for testing purposes
	h.HandleFunc("/traffic_set", func(rw http.ResponseWriter,r *http.Request) {
		reqBuf  := new(bytes.Buffer)
		reqBuf.ReadFrom(r.Body)
		strReq := string(reqBuf.Bytes())
		count,_ := strconv.Atoi(strReq)
		Traffic = count;
		log.Printf("My traffic is %d",Traffic)
	})
	
	//for testing purposes
	h.HandleFunc("/ill_set", func(rw http.ResponseWriter,r *http.Request) {
		ill = true;
	})
	
	h.HandleFunc("/traffic", func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("content-type", "text/plain")
			str := strconv.Itoa(Traffic)
			n, _ := rw.Write([]byte(str))
			Traffic += n
			
	})

	report := make(Report)

	h.HandleFunc("/api/v1/some-data", func(rw http.ResponseWriter, r *http.Request) {
		respDelayString := os.Getenv(confResponseDelaySec)
		if delaySec, parseErr := strconv.Atoi(respDelayString); parseErr == nil && delaySec > 0 && delaySec < 300 {
			time.Sleep(time.Duration(delaySec) * time.Second)
		}

		report.Process(r)

		rw.Header().Set("content-type", "application/json")
		rw.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(rw).Encode([]string{
			"1", "2",
		})
		Traffic += len([]string{"1", "2",})
	})

	h.Handle("/report", report)

	server := httptools.CreateServer(*port, h)
	server.Start()
	signal.WaitForTerminationSignal()
}
