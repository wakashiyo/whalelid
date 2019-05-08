package whalelid

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

type Commands struct {
	Command   string
	Operation string
	Options   []string
	Target    string
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

func NetworkInfo(c Commands, n *Network) error {
	bytes := []byte{}
	if err := c.output(&bytes); err != nil {
		return err
	}
	if err := n.networkInfo(bytes); err != nil {
		return err
	}
	return nil
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

func (c *Commands) output(b *[]byte) error {

	cs := append([]string{c.Operation}, c.Options...)

	o, err := exec.Command(c.Command, cs...).Output()
	if err != nil {
		return err
	}
	b = &o
	return nil
}

func (c *Commands) run() error {

	cs := append([]string{c.Operation}, c.Options...)

	cs = append(cs, c.Target)

	err := exec.Command(c.Command, cs...).Run()
	if err != nil {
		return err
	}

	return nil
}

////////////////////////////////////////////////////////////////////////////////////////////////

type ExecuteCommand interface {
	Run() error
	Output() ([]byte, error)
}

func output(ec ExecuteCommand, b *[]byte) error {
	bytes, err := ec.Output()
	if err != nil {
		return err
	}
	*b = bytes
	return nil
}

func run(ec ExecuteCommand) error {
	if err := ec.Run(); err != nil {
		return err
	}
	return nil
}
