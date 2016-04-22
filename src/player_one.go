package main

import (
	"net"
	"fmt"
)

const Mode_player_one = "1"
const Mode_player_two = "2"
const Mode_observer = "3"

const Network_mode_tcp = "tcp"

const Request_position_input = "99"

const Result_code_input_success = "00"
const Result_code_input_already_occupied = "11"
const Result_code_input_out_of_bound = "12"

const Result_code_input_player1_win = "01"
const Result_code_input_player2_win = "02"

var ground [][]int
func main() {
	client, err := net.Dial("tcp", "127.0.0.1:8608")
	
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Close()
	
	data := make([]byte, 4096)
	data = []byte("1")
	client.Write([]byte(data))
	
	_, err = client.Read(data)
	
	if err != nil {
		fmt.Println(err)
		return
	}
	
	player_num := string(data)
	switch(player_num) {
	case Mode_player_one:
		fmt.Println("You are player 1.")
	case Mode_player_two:
		fmt.Println("You are player 2.")
	default:
		fmt.Println("You are observer.")
	}	
	for{
		data = make([]byte, 4096)
		n, err := client.Read(data)
		if err != nil {
			fmt.Println(err)
		}
		switch(string(data[:n])) {
		case Request_position_input:
			data = make([]byte,2)
			x, y := inputPosition()
			data[0], data[1] = byte(x), byte(y)
			client.Write(data)
		case Result_code_input_already_occupied:
			client.Read(data)
			data = make([]byte,2)
			x, y := solveDupPosition()
			data[0], data[1] = byte(x), byte(y)
			client.Write(data)
		case Result_code_input_out_of_bound:
			client.Read(data)
			data = make([]byte,2)
			x, y := solveOutOfIndex()
			data[0], data[1] = byte(x), byte(y)
			client.Write(data)
		case Result_code_input_player1_win:
			if player_num == Mode_player_one {
				fmt.Println("You win!!")
			}else if player_num == Mode_player_two {
				fmt.Println("You lose..")
			}else {
				fmt.Println("Player 1 win!!")
			}
			return
		case Result_code_input_player2_win:
			if player_num == Mode_player_two {
				fmt.Println("You win!!")
			}else if player_num == Mode_player_one {
				fmt.Println("You lose..")
			}else {
				fmt.Println("Player 1 win!!")
			}
			return
		default:
			makeGround(data)
			showGround()
			size := int(data[0])
			ground := make([][]byte, size + 1)
			for index, _ := range ground {
				ground[index] = make([]byte, size)
				for index_, _ := range ground[index] {
					ground[index][index_] = data[index * size + index_ +1]
				}
			}
		}
	}
}

func inputPosition() (int, int) {
	x, y := -1, -1
	isNum := true
	fmt.Print("Enter positions(x y) : ")
	
	for x < 0 && y < 0 {
		if !isNum {
			fmt.Print("Not number, enter again(x y) : ")
			fmt.Scanf("%d %d\n", &x, &y)
			fmt.Scanf("%d %d\n", &x, &y)
		} 
		fmt.Scanf("%d %d\n", &x, &y)
		isNum = false
	}
	
	return x, y
}

func solveDupPosition() (int, int){
	fmt.Println("Already occupied. Try again")
	return inputPosition()
}

func solveOutOfIndex() (int, int){
	fmt.Println("Ouf of index, Try again")
	return inputPosition()
}

func makeGround(data []byte) {
	size := int(data[0])
	if ground == nil || len(ground) < 1 {
		ground = make([][]int, size)
		for index := 0; index < size; index++ {
			ground[index] = make([]int, size)
		}
	}
	
	for index := 0; index < size; index++ {
		for index_ := 0; index_ < size; index_++ {
			ground[index][index_] = int(data[index * size + index_ + 1])
		}
	}
}

func showGround() {
	for index, _ := range ground {
		for _, value := range ground[index] {
			fmt.Print(value, "\t")
		}
		fmt.Println()
	}
}