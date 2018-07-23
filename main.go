package main

import (
	"net"
	"time"
	"fmt"
	"encoding/json"
	"flag"
	"math/rand"
	"strings"
)

func getCluserHBMessage(source NodeInfo, dest NodeInfo, message string) (ClusterHBMessage){
	return ClusterHBMessage{
		Source: NodeInfo{
			NodeId: source.NodeId,
			NodeIpAddr: source.NodeIpAddr,
			Port: source.Port,
		},
		Dest: NodeInfo{
			NodeId: dest.NodeId,
			NodeIpAddr: dest.NodeIpAddr,
			Port: dest.Port,
		},
		Message: message,
	}
}

func connectToCluster(me NodeInfo, dest NodeInfo) (bool) {
	connOut, err := net.DialTimeout("tcp", dest.NodeIpAddr+":"+dest.Port, time.Duration(10)*time.Second)
	if err != nil{
		if _, ok := err.(net.Error); ok {
			fmt.Println("Couldn't connect to the cluster, ", me.NodeId)
			return false
		}
	}else{
		fmt.Println("Connected to cluster. Sending message to node.")
		text := "Hi.., Add me to your cluster.."
		requestMessage := getCluserHBMessage(me, dest, text)
		json.NewEncoder(connOut).Encode(&requestMessage)

		decoder := json.NewDecoder(connOut)
		var responseMessage ClusterHBMessage
		decoder.Decode(&responseMessage)
		fmt.Println("Go reponse : \n"+responseMessage.String())
		return true
	}
	return false
}


func listenToPort(me NodeInfo) {
	/** Listen to the port for incoming message*/
	ln,_ := net.Listen("tcp", fmt.Sprintf(":"+me.Port))
	// accept connection
	for {
		connIn, err := ln.Accept()
		if err != nil {
			if _, ok := err.(net.Error); ok{
				fmt.Println("Error when listening to the service.", me.NodeId)
			}
		}else{
			var requestMessage ClusterHBMessage
			json.NewDecoder(connIn).Decode(&requestMessage)
			fmt.Println("Got Request: \n "+ requestMessage.String())
			text := "Sure buddy. this is easy..."
			respMessage := getCluserHBMessage(me, requestMessage.Source, text)
			fmt.Println("Sending Response: \n", respMessage.String())
			json.NewEncoder(connIn).Encode(&respMessage)
			connIn.Close()
		}
	}
}


func main(){
    /* parse the command line parameters */
    makeMasteronError := flag.Bool("makeMasterOnError", false, "make this node as the master")
    clusterIp := flag.String("clusterip", "127.0.0.1:8001", "ip address of any node to connect.")
    myPort := flag.String("myport", "8001", "ip address to run this node on")
    flag.Parse()

    rand.Seed(time.Now().UTC().UnixNano())
    myid := rand.Intn(99999999)

    myIp,_ := net.InterfaceAddrs()
    me := NodeInfo{NodeId: myid, NodeIpAddr: myIp[0].String(), Port: *myPort}
    dest := NodeInfo{NodeId: -1, NodeIpAddr: strings.Split(*clusterIp,":")[0], Port: strings.Split(*clusterIp,":")[1]}
    fmt.Println("Mydetails: ", me.String())

    // listen to incoming requests and connect the hb
	ableToConn := connectToCluster(me, dest)
	if ableToConn || (!ableToConn && *makeMasteronError){
		if *makeMasteronError {
			fmt.Println("Will start this node as master")
		}
		listenToPort(me)
	}else {
		fmt.Println("Quiting system. Set make master on error flag to make a not master.")
	}

}