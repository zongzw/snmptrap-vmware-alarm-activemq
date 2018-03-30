package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"strings"

	stomp "github.com/go-stomp/stomp"
	g "github.com/soniah/gosnmp"
)

var MsgChan = make(chan *Alarm, 32)

var testtimes int
var trapdsite string

var mqsite string
var mquser string
var mqpass string
var queue string

type Alarm struct {
	OldStatus string `json:"oldstatus"`
	NewStatus string `json:"newstatus"`
	Object    string `json:"object"`
	Detail    string `json:"detail"`
	Source    string `json:"reported_source"`
}

func ConnectMQ() *stomp.Conn {
	options := []func(*stomp.Conn) error{
		stomp.ConnOpt.Login(mquser, mqpass),
		stomp.ConnOpt.Host("/"),
	}

	conn, err := stomp.Dial("tcp", mqsite, options...)
	if err != nil {
		fmt.Printf("failed to connect to mq: %s\n", err.Error())
		return nil
	}

	return conn
}

func Alarm2MQ() bool {
	a := <-MsgChan

	fmt.Printf("Channel ====> MQ\n")
	bytes, err := json.Marshal(*a)
	if err != nil {
		fmt.Printf("error for marshal data: %s\n", err.Error())
		return false
	}

	fmt.Printf("new alert: %s\n", string(bytes))

	mqconn := ConnectMQ()
	if mqconn == nil {
		fmt.Printf("failed to connect, abort.\n")
		return false
	}

	defer mqconn.Disconnect()

	fmt.Printf("connect to mq successfully.\n")
	err = mqconn.Send(queue, "text/plain", bytes)
	if err != nil {
		fmt.Printf("failed to send data to mq: %s\n", err.Error())
		return false
	}
	fmt.Printf("succeed to send message to mq\n")

	return true
}

func TestReadMQ() {
	fmt.Printf("start to test reading from mq. \n")
	mqconn := ConnectMQ()
	if mqconn == nil {
		fmt.Printf("failed to connect, abort.\n")
		return
	}

	defer mqconn.Disconnect()

	sub, err := mqconn.Subscribe(queue, stomp.AckAuto)
	if err != nil {
		fmt.Printf("Failed to sub from mq: %s\n", err.Error())
		return
	}
	for i := 0; i < testtimes; i++ {
		msg := <-sub.C
		fmt.Printf("MQ ====> Stdout\n")
		fmt.Printf("%d: Read from MQ: %v\n", i, msg)
	}
	err = sub.Unsubscribe()
	if err != nil {
		fmt.Printf("failed to unsub %s\n", err.Error())
	}
}

func myTrapHandler(packet *g.SnmpPacket, addr *net.UDPAddr) {

	var alarm Alarm
	alarm.Source = fmt.Sprintf("%s", addr.IP)

	for _, v := range packet.Variables {
		if strings.Contains(v.Name, "6876.4.3.304") {
			alarm.OldStatus = string(v.Value.([]byte))
		}
		if strings.Contains(v.Name, "6876.4.3.305") {
			alarm.NewStatus = string(v.Value.([]byte))
		}
		if strings.Contains(v.Name, "6876.4.3.306") {
			alarm.Detail = string(v.Value.([]byte))
		}
		if strings.Contains(v.Name, "6876.4.3.307") {
			alarm.Object = string(v.Value.([]byte))
		}
	}

	fmt.Printf("Trap ====> Channel\n")
	MsgChan <- &alarm
}

func StartTrapd() {
	tl := g.NewTrapListener()

	tl.OnNewTrap = myTrapHandler
	tl.Params = g.Default
	tl.Params.Logger = nil

	err := tl.Listen(trapdsite)
	if err != nil {
		log.Panicf("error in listen: %s", err)
	}
}

func main() {

	flag.StringVar(&trapdsite, "trapd", "0.0.0.0:162", "snmp trapd [ip:port]")
	flag.StringVar(&mqsite, "mqsite", "0.0.0.0:61613", "activemq site [ip:port]")
	flag.StringVar(&queue, "queue", "/queue/myq", "queue name, i.e. /queue/myq")
	flag.StringVar(&mquser, "mquser", "", "activemq user")
	flag.StringVar(&mqpass, "mqpass", "", "activemq pass")
	flag.IntVar(&testtimes, "test", 2, "times of testing reading from queue")

	flag.Parse()

	go StartTrapd()
	go TestReadMQ()
	for {
		if ok := Alarm2MQ(); !ok {
			fmt.Printf("failed to post data to mq.\n")
		}
	}
}
