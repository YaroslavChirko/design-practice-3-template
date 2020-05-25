package main
import "testing"
import ."gopkg.in/check.v1"

func Test(t *testing.T) { TestingT(t) }

type MySuite struct{}

var _ = Suite(&MySuite{})

func (s *MySuite) TestBalancer(c *C) {
	traffic := [3]int{10,1000,100}
	health := [3]bool{true,true,true}
	i,_ := getServerMoc(traffic,health)
	c.Assert(0,Equals,i)
	
	health[0]=false
	i1,_ := getServerMoc(traffic,health)
	c.Assert(2,Equals,i1)
	
	health[1]=false
	health[2]=false
	_,err := getServerMoc(traffic,health)
	c.Assert("No healthy servers",Equals,err.Error())
	
}
