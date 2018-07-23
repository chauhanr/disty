package main

import "strconv"

type NodeInfo struct{
	NodeId int `json:"nodeId"`
	NodeIpAddr string `json:"nodeIpAddr"`
	Port string `json:"port"`
}

type ClusterHBMessage struct {
	Source NodeInfo `json:"source"`
	Dest NodeInfo `json:"dest"`
	Message string 	`json:"message"`
}

func (node NodeInfo) String() string{
	return "NodeInfo:{ nodeId:" + strconv.Itoa(node.NodeId) + ", nodeIpAddr:" + node.NodeIpAddr + ", port:" + node.Port + " }"
}

func (req ClusterHBMessage) String() string{
	return "AddToClusterMessage:{\n  source:" + req.Source.String() + ",\n  dest: " + req.Dest.String() + ",\n  message:" + req.Message + " }"
}

