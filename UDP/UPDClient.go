package main

import (
	"fmt"
	"net"
)

func main(){
	rAddr := &net.UDPAddr{
		// TODO 修改为对应的ip和端口
		IP:   net.IPv4(192,168,16,228),
		Port: 9999,
	}

	conn, err := net.DialUDP("udp",nil, rAddr)
	if err != nil{
		fmt.Println(err.Error())
		return
	}
	defer conn.Close()

	msg := "hello word"
	_, err = conn.Write([]byte(msg))
	if err != nil{
		fmt.Println(err.Error())
		return
	}
	fmt.Println("send msg: ", msg)

	data := make([]byte, 1024)
	_, _, err = conn.ReadFromUDP(data)
	if err != nil{
		fmt.Println(err.Error())
		return
	}
	fmt.Println("rec msg: ",string(data))
}