package main

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

func main()  {
	addrs, err := net.InterfaceAddrs()
	if err != nil{
		fmt.Println(err.Error())
		return
	}

	var ip4 net.IP
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ip4 = ipnet.IP.To4()
				break
			}

		}
	}
	if ip4 == nil{
		return
	}

	lAddr := &net.UDPAddr{
		IP:   ip4,
		Port: 9999,
	}
	conn, err := net.ListenUDP("udp", lAddr)
	if err != nil{
		fmt.Println(err.Error())
		return
	}

	for{
		data := make([]byte, 1024)
		n, rAddr, err := conn.ReadFromUDP(data)
		if err != nil{
			fmt.Println(err.Error())
			return
		}
		// checkSum
		checkSum(conn, rAddr, n, data)

		//
		str := string(data[:n])
		fmt.Println("rec msg: ",str)

		_, err = conn.WriteToUDP([]byte(strings.ToUpper(str)), rAddr)
		if err != nil{
			fmt.Println(err.Error())
			return
		}
	}
}

func checkSum(conn *net.UDPConn, rAddr *net.UDPAddr, n int, data []byte){
	list := []int{}

	// 来源地址
	list = append(list, int(rAddr.IP[0]) << 8 + int(rAddr.IP[1]))
	list = append(list, int(rAddr.IP[2]) << 8 + int(rAddr.IP[3]))
	fmt.Printf("来源地址：%04x %04x\n",list[len(list)-2],list[len(list)-1])

	// 目的地址
	lList := strings.Split(conn.LocalAddr().String(),":")
	lPort, _ := strconv.Atoi(lList[1])
	lIpList := strings.Split(lList[0],".")
	for i := 0; i < len(lIpList);  i+=2{
		j, _ := strconv.Atoi(lIpList[i])
		k, _ := strconv.Atoi(lIpList[i+1])
		list = append(list, j << 8 + k)
	}
	fmt.Printf("目的地址：%04x %04x\n",list[len(list)-2],list[len(list)-1])

	// 全零 协议名 UDP报文长度
	list = append(list, 0x0011, n+8)
	fmt.Printf("全零 协议名 UDP报文长度：%04x %04x\n",list[len(list)-2],list[len(list)-1])

	// 来源端口
	list = append(list, rAddr.Port)
	// 目标端口
	list = append(list, lPort)
	fmt.Printf("来源端口 目标端口：%04x %04x\n",list[len(list)-2],list[len(list)-1])

	// 报文长度 检验和
	list = append(list, n+8, 0)
	fmt.Printf("报文长度 检验和：%04x %04x\n",list[len(list)-2],list[len(list)-1])

	// 数据
	newData := data[:n]
	if n%2 >0 {
		newData = append(newData, 0)
	}
	fmt.Print("数据：")
	for i := 0; i < len(newData);  i+=2{
		list = append(list, int(newData[i]) << 8 + int(newData[i+1]))
		fmt.Printf("%04x ", list[len(list)-1])
	}
	fmt.Print("\n")

	// 校验
	sum := 0
	for _, num := range list {
		sum += num
		if ((sum >> 16) > 0){
			sum = (sum >> 16) + (sum & 0xffff)
		}
	}
	fmt.Printf("检验和：%04x\n", ^uint(sum))
}
