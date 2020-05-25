package main
import "bytes"
import "net/http"
import "testing"
import "log"
import ."gopkg.in/check.v1"

func Test(t *testing.T) { TestingT(t) }

type MySuite struct{}

var _ = Suite(&MySuite{})

func (s *MySuite) TestBalancer(c *C) {
	req, _ := http.NewRequest("POST","http://server1:8080/traffic_set",bytes.NewBuffer([]byte("10")))
	http.DefaultClient.Do(req)
	req, _ = http.NewRequest("POST","http://server2:8080/traffic_set",bytes.NewBuffer([]byte("1000")))
	http.DefaultClient.Do(req)
	req, _ = http.NewRequest("POST","http://server3:8080/traffic_set",bytes.NewBuffer([]byte("100")))
	http.DefaultClient.Do(req)
	/*i:=-1
	log.Printf("Index is : %d now",i)
	i,_=getServer(i)
	log.Printf("Index is : %d now",i)*/
	//c.Assert(0,Equals,i)
	/*req, _ := http.NewRequest("GET","http://server1:8080/ill_set",nil)
	http.DefaultClient.Do(req)
	c.Assert(0,Equals,0)
	
	req, _ = http.NewRequest("GET","http://server1:8080/ill_set",nil)
	http.DefaultClient.Do(req)
	i :=1//,_=getServer(-1)
	c.Assert(1,Equals,i)*/
	req, _ = http.NewRequest("GET","http://server1:8080/ill_set",nil)
	http.DefaultClient.Do(req)
	req, _ = http.NewRequest("GET","http://server2:8080/ill_set",nil)
	http.DefaultClient.Do(req)
	req, _ = http.NewRequest("GET","http://server3:8080/ill_set",nil)
	http.DefaultClient.Do(req)
	
	_,err:=getServer(-1)
	c.Assert("No healthy servers",Equals,err.Error())
}
