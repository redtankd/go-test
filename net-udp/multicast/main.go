// An example for UDP Multicast
// UDP Multicast 依赖于 IGMP(Internet Group Management Protocol)。
// Multicast != Broadcast
// TODO：网络设备是否都支持IGMP，是否可以在广域网环境下跨越网络边界使用IGMP

package main

import (
    "bufio"
    "bytes"
    "flag"
    "fmt"
    "net"
    "os"
)

// 从命令行读入信息
func fromUser() []byte {
    reader := bufio.NewReader(os.Stdin)
    fmt.Printf("multicast>")
    line, _, err := reader.ReadLine()
    check(err)
    return line
}

// 服务器线程
// 从命令行读入信息后，将其广播给所有客户端
// 从命令行读入"quit"，退出服务器
func multicastServer(conn *net.UDPConn, address *net.UDPAddr) {
    for str := fromUser(); !bytes.Equal(str, []byte("quit")); str = fromUser() {
        buffer := make([]byte, 512)
        copy(buffer, str) //TODO 为什么要copy到一个新的byte数组
        _, err := conn.WriteToUDP(buffer, address)
        check(err)
    }
}

// 客户端线程
func listen(conn *net.UDPConn) {
    for {
        b := make([]byte, 256)
        _, _, err := conn.ReadFromUDP(b)
        check(err)
        fmt.Println("server:", string(b))
    }
}

func check(err error) {
    if err != nil {
        panic(err)
    }
}

func main() {
    var ip = flag.String("ip", "224.0.1.60:1888", "The ip and port for multicast.")
    var server = flag.Bool("server", false, "Server or Client. The default is Client.")
    flag.Parse()

    // 服务器端使用的UDP连接，实际并不接收接入连接，仅仅用于发送组播消息
    // TODO 这个UDP连接，是否可以随意创建？是否可以同时用于一对一的通信？
    localAddress, err := net.ResolveUDPAddr("udp", ":0")
    check(err)
    localConn, err := net.ListenUDP("udp", localAddress)
    check(err)


    // 组播通信的地址，客户端侦听这个地址
    // 服务器直接向这个地址发送组播，并不需要建立一个到这个地址的连接
    // IGMP协议规定了可以用于组播的地址
    multicastAddress, err := net.ResolveUDPAddr("udp", *ip)
    check(err)
    multicastConn, err := net.ListenMulticastUDP("udp", nil, multicastAddress)
    check(err)

    if *server {
        fmt.Printf("Server Mode [%s]\n", *ip)
        multicastServer(localConn, multicastAddress)
    } else {
        fmt.Printf("Client Mode [%s]\n", *ip)
        listen(multicastConn)
    }
}