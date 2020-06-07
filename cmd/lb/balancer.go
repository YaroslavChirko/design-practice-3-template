package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
	"math"
	"errors"
	"github.com/YaroslavChirko/design-practice-3-template/httptools"
	"github.com/YaroslavChirko/design-practice-3-template/signal"
)

var (
	port = flag.Int("port", 8090, "load balancer port")
	timeoutSec = flag.Int("timeout-sec", 3, "request timeout time in seconds")
	https = flag.Bool("https", false, "whether backends support HTTPs")

	traceEnabled = flag.Bool("trace", false, "whether to include tracing information into responses")
)

var (
	timeout = time.Duration(*timeoutSec) * time.Second
	serversPool = []string{
		"server1:8080",
		"server2:8080",
		"server3:8080",
	}
	
	serversT = []int64{
		0,
		0,
		0,
	}
)


func scheme() string {
	if *https {
		return "https"
	}
	return "http"
}

func health(dst string) bool {
	ctx, _ := context.WithTimeout(context.Background(), timeout)
	req, _ := http.NewRequestWithContext(ctx, "GET",
		fmt.Sprintf("%s://%s/health", scheme(), dst), nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false
	}
	if resp.StatusCode != http.StatusOK {
		return false
	}
	return true
}


func forward(frm int,dst string, rw http.ResponseWriter, r *http.Request) error {
	ctx, _ := context.WithTimeout(r.Context(), timeout)
	fwdRequest := r.Clone(ctx)
	fwdRequest.RequestURI = ""
	fwdRequest.URL.Host = dst
	fwdRequest.URL.Scheme = scheme()
	fwdRequest.Host = dst

	resp, err := http.DefaultClient.Do(fwdRequest)
	if err == nil {
		for k, values := range resp.Header {
			for _, value := range values {
				rw.Header().Add(k, value)
			}
		}
		if *traceEnabled {
			rw.Header().Set("lb-from", dst)
		}
		log.Println("fwd", resp.StatusCode, resp.Request.URL)
		rw.WriteHeader(resp.StatusCode)
		defer resp.Body.Close()
		bs, err := io.Copy(rw, resp.Body)
		serversT[frm] += bs
		if err != nil {
			log.Printf("Failed to write response: %s", err)
		}
		return nil
	} else {
		log.Printf("Failed to get response from %s: %s", dst, err)
		rw.WriteHeader(http.StatusServiceUnavailable)
		return err
	}
}

func getServer (index int) (int, error){
var hiTraffic int64 =math.MaxInt64
	 
	for i := 0;i<3;i++ {
		if(health(serversPool[i])&&serversT[i]<=hiTraffic){
			index = i
			hiTraffic=serversT[i]
		}
	}
	if(index==-1){
		return -1,errors.New("No healthy servers")
	}
	return index,nil;
}

//dummy for test purpouses to verify getServer work
func getServerMoc (t [3] int, h [3] bool) (int, error){
index:=-1
var hiTraffic int =math.MaxInt32
	 
	for i := 0;i<3;i++ {
		if(h[i]&&t[i]<=hiTraffic){
			index = i
			hiTraffic=t[i]
		}
	}
	if(index==-1){
		return -1,errors.New("No healthy servers")
	}
	return index,nil;
}

func main() {
	var err error = nil
	index := -1;
	flag.Parse()

	// TODO: Використовуйте дані про стан сервреа, щоб підтримувати список тих серверів, яким можна відправляти ззапит.
	for _, server := range serversPool {
		server := server
		go func() {
			for range time.Tick(10 * time.Second) {
				log.Println(server, health(server))
			}
		}()
	}
		go func() {
			for range time.Tick(time.Hour) {
				serversT[0] =0
				serversT[1] =0
				serversT[2] =0
			}
		}()
	
		
	frontend := httptools.CreateServer(*port, http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		index,err = getServer(index)
		if(err!=nil){
			log.Printf("Encountered errors: %s", err.Error())
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(err.Error()))
		}else{
		forward(index,serversPool[index], rw, r)
		}
		
		index = -1
		err = nil
	}))

	log.Println("Starting load balancer...")
	log.Printf("Tracing support enabled: %t", *traceEnabled)
	frontend.Start()
	signal.WaitForTerminationSignal()
}
