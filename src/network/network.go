package network

import (
	"fmt"
	"net"
	"constUtil"
)

type Network struct {
	list net.Listener
	player_one net.Conn
	player_two net.Conn
	clientList []net.Conn
	buf []byte
}


func New() (listner *Network){
	return &Network{
		list : nil,
	}
}

func (list *Network) Init(port string, buf_size int) int {
	var err error
	list.list, err = net.Listen(constUtil.Network_mode_tcp, port)
	
	if err != nil {
		return -1
	}
	
	list.buf = make([]byte, buf_size)
	list.clientList = make([]net.Conn, constUtil.List_size)
	
	return 0
}

func (list *Network) StartServer(c chan int) int{
	for {
		conn, err := list.list.Accept()
		if err != nil{
			return -2
		}
		
		list.clientList = append(list.clientList, conn)
		go list.waitingModeSelect(conn, c)
	}
}

func (list *Network) waitingModeSelect(conn net.Conn, c chan int) int{
	n, err := conn.Read(list.buf)
	if err != nil {
		fmt.Println(err)
		return -1
	}
	
	mode := string(list.buf[:n])
	switch(mode) {
	case constUtil.Mode_player_one:
		if list.player_one == nil {
			list.player_one = conn
			conn.Write([]byte(constUtil.Mode_player_one))
			fmt.Println("Player 1 enterd")
			break
		}
		fallthrough
	case constUtil.Mode_player_two:
		if list.player_two == nil {
			list.player_two = conn
			conn.Write([]byte(constUtil.Mode_player_two))
			fmt.Println("Player 2 enterd")
			break
		}else if list.player_one == nil {
			list.player_one = conn
			conn.Write([]byte(constUtil.Mode_player_one))
			fmt.Println("Player 1 enterd")
			break
		}
		fallthrough
	default:
			conn.Write([]byte(constUtil.Mode_observer))
			fmt.Println("Observer enterd")
	}
	
	if list.player_one != nil && list.player_two != nil {
		c <- 1
	}
	return 0
}

func (list *Network) InputPosition(player string) (int, int){
	var conn net.Conn
	var n int
	var err error
	
	if player == constUtil.Mode_player_one {
		n, err = list.player_one.Write([]byte(constUtil.Request_position_input))
		conn = list.player_one
	}else if player == constUtil.Mode_player_two {
		n, err = list.player_two.Write([]byte(constUtil.Request_position_input))
		conn = list.player_two
	}else {
		return -1, -1
	}
	if err != nil {
		fmt.Println(err)
	}
	
	n, err = conn.Read(list.buf)
	
	if err != nil {
		fmt.Println(err)
		return -2, -2
	}else{
		if n != 2 {
			return -2, -2
		}else {
			return int(list.buf[0]), int(list.buf[1])
		}
	}
}

func (list *Network) SendToClients(data []byte) {
	for _, value := range list.clientList {
		n, err := value.Write(data)
		if err != nil {
			fmt.Println(value.RemoteAddr, " : ", err)
		}else if n != len(data) {
			fmt.Println(value.RemoteAddr, " : Data length is different.")
		}
	}
}

func (list *Network) Finish(winner int) {
	defer list.clearNetword()
	data := make([]byte, 2)
	if winner == 1 {
		data = []byte(constUtil.Result_code_input_player1_win)
	} else if winner == 2 {
		data = []byte(constUtil.Result_code_input_player2_win)
	} else {
		fmt.Println("Error")
		return
	}
	list.SendToClients(data)
}

func (list *Network) SendResultOfInput(player int, result string) {
	data := []byte(result)
	if player == 1 {
		list.player_one.Write(data)
	}else if player == 2{
		list.player_two.Write(data)
	}
}

func (list *Network) clearNetword() {
	list.list.Close()
	list.player_one.Close()
	list.player_two.Close()
	
	if list.clientList != nil && len(list.clientList) > 0{
		for _, value := range list.clientList {
			value.Close()
		}
	}
}