package whalelid

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

type commands struct {
	command   string
	operation string
	options   []option
	target    string
}

type option struct {
	key   string
	value string
}

type Network struct {
	bridge     string
	subnet     string
	gateway    string
	containers []container
}

type container struct {
	name      string
	IPAddress string
}

type Err struct {
	command string
	message string
	status  int
}

func (e *Err) Error() string {
	return fmt.Sprintf("[ERROR] *** command : %v, Exit Status : %v, ErrorMessage : %v", e.command, e.status, e.message)
}

func (n *Network) networkInfo(b []byte) error {
	var i interface{}
	if err := json.Unmarshal(b, &i); err != nil {
		return err
	}
	net := network(i)
	n.bridge = net.bridge
	n.subnet = net.subnet
	n.gateway = net.gateway
	n.containers = net.containers
	return nil
}

func network(i interface{}) Network {
	info, _ := i.([]interface{})[0].(map[string]interface{})

	bridge, _ := info["Options"].(map[string]interface{})["com.docker.network.bridge.name"].(string)

	o, _ := info["IPAM"].(map[string]interface{})["Config"].([]interface{})[0].(map[string]interface{})

	subnet, _ := o["Subnet"].(string)
	gateway, _ := o["Gateway"].(string)

	c, _ := info["Containers"].(map[string]interface{})

	key := []string{}
	for k := range c {
		key = append(key, k)
	}

	containers := []container{}
	for _, val := range key {
		co, _ := c[val].(map[string]interface{})
		n, _ := co["Name"].(string)
		ip, _ := co["IPv4Address"].(string)
		con := container{
			name:      n,
			IPAddress: ip,
		}
		containers = append(containers, con)
	}

	return Network{
		bridge:     bridge,
		subnet:     subnet,
		gateway:    gateway,
		containers: containers,
	}

}

func createOptions(m map[string]string) []option {
	o := []option{}
	for k, v := range m {
		_o := option{
			key:   k,
			value: v,
		}
		o = append(o, _o)
	}
	return o
}

func (c *commands) output(b *[]byte) error {

	cs := []string{c.operation}
	for _, v := range c.options {
		cs = append(cs, v.key, v.value)
	}

	o, err := exec.Command(c.command, cs...).Output()
	if err != nil {
		return err
	}
	b = &o
	return nil
}

func (c *commands) run() error {

	cs := []string{c.operation}
	for _, v := range c.options {
		cs = append(cs, v.key, v.value)
	}

	cs = append(cs, c.target)

	err := exec.Command(c.command, cs...).Run()
	if err != nil {
		return err
	}

	return nil
}
