package integration

import (
	"fmt"
	"net/http"
	"testing"
	."gopkg.in/check.v1"
	"time"
)

const baseAddress = "http://balancer:8090"

var client = http.Client{
	Timeout: 3 * time.Second,
}


func Test(t *testing.T) { TestingT(t) }
type MySuite struct{}

var _ = Suite(&MySuite{})

func (s *MySuite) TestBalancer(c *C) {
	m :=make(map[string]int)
	m["server1:8080"]=0
	m["server2:8080"]=0
	m["server3:8080"]=0
	for i:=0;i<9;i++{
		resp, err := client.Get(fmt.Sprintf("%s/api/v1/some-data", baseAddress))
		if err != nil {
			c.Error(err)
		}
		fmt.Printf("response from [%s]", resp.Header.Get("lb-from"))
		c.Logf("response from [%s]", resp.Header.Get("lb-from"))
		m[resp.Header.Get("lb-from")] +=1
	}
	
	c.Assert(m["server1:8080"],Equals,3)
	c.Assert(m["server2:8080"],Equals,3)
	c.Assert(m["server3:8080"],Equals,3)
}

func BenchmarkBalancer(b *testing.B) {
	// TODO: Реалізуйте інтеграційний бенчмарк для балансувальникка.
}
