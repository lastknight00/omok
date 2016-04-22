package baduk

import (
	"fmt"
	"constUtil"
	"network"
)

type Baduk struct{
	ground [][]int
	net *network.Network
}


func New() (baduk *Baduk) {
    return &Baduk{
    	ground : nil,
    	net : network.New(),
    }
}

func (baduk *Baduk) InitGround() [][]int {
	size := -1;
	fmt.Print("Enter a ground size : ")
	fmt.Scanf("%d\n", &size)
	
	if size > 0{
		baduk.makeGround(size)
		return baduk.ground 
	}else{
		return nil;
	}
}

func (baduk *Baduk) makeGround(size int) {
	baduk.ground = make([][]int, size)
	
	for index, _ :=range baduk.ground {
		baduk.ground[index] = make([]int, size)
	}
}

func (baduk *Baduk) ShowGround() {
	fmt.Print("\t")
	for index := 0; index < len(baduk.ground); index++ {
		fmt.Printf("%d\t", index)
	}
	fmt.Println()
	for index, _ := range baduk.ground {
		fmt.Printf("%d\t", index)
		for _, value := range baduk.ground[index] {
			switch(value) {
			case 0:
				fmt.Printf("%s\t", "┼")
			case 1:
				fmt.Printf("%s\t", "○")
			case 2:
				fmt.Printf("%s\t", "●")
			}
		}
		fmt.Println()
	}
}

func (baduk *Baduk) StartOmok() int {
	player := 1
	flag := 1
	x, y := -1, -1
	result := 0
	fmt.Println("Starting server...")
	//net := network.New()
	baduk.net.Init(":8608", 4096)
	c := make(chan int)
	fmt.Println("Started server.")
	
	fmt.Println("Waiting to enter players...")
	go baduk.net.StartServer(c)
	<- c
	fmt.Println("Game start.")
	baduk.ShowGround()
	for result == 0 {
		//x, y = inputPosition()
		baduk.sendGroundToClients()
		x, y = baduk.inputPositionForNetwork(player)
		inputResult := plackBadukPiece(baduk.ground, player, x, y)
		
		failCount := 0
		for inputResult != constUtil.Result_code_input_success {
			failCount++
			if failCount >= constUtil.Max_fail_count {
				result := player + flag
				baduk.net.Finish(result)
				return result
			}
			baduk.net.SendResultOfInput(player, inputResult)
			x, y = baduk.inputPositionForNetwork(player)
			inputResult = plackBadukPiece(baduk.ground, player, x, y)
		}
		player = player + flag
		flag *= -1
		baduk.ShowGround()
		result = baduk.checkFinish(x, y)
	}
	baduk.sendGroundToClients()
	baduk.net.Finish(result)	
	return result
}

func inputPosition() (int, int) {
	var x, y int
	fmt.Printf("Enter the position : ")
	fmt.Scanf("%d %d\n", &x, &y)
	return x, y
}

func (baduk *Baduk) inputPositionForNetwork(player int) (int, int) {
	if player == 1 {
		return baduk.net.InputPosition(constUtil.Mode_player_one)
	}else if player == 2 {
		return baduk.net.InputPosition(constUtil.Mode_player_two)
	}else {
		return -1, -1
	}
}
func (baduk *Baduk) checkFinish(x int, y int) int {
	value := baduk.ground[x][y]
	count := 1
	for indexX := x - 1; (indexX >= 0) && baduk.ground[indexX][y] == value; indexX-- {
		count++
	}
	
	for indexX := x + 1; (indexX < len(baduk.ground)) && baduk.ground[indexX][y] == value; indexX++ {
		count++
	}
	
	if count >= 5 {
		return value
	}
	
	count = 1
	
	for indexY := y - 1; (indexY >= 0) && baduk.ground[x][indexY] == value; indexY-- {
		count++
	}
	
	for indexY := y + 1; (indexY < len(baduk.ground)) && baduk.ground[x][indexY] == value; indexY++ {
		count++
	}
	
	if count >= 5 {
		return value
	}
	
	count = 1
	
	for indexX, indexY := x - 1, y - 1; (indexX >= 0) && (indexY >= 0) && baduk.ground[indexX][indexY] == value; indexX, indexY = indexX - 1, indexY - 1 {
		count++
	}
	
	for indexX, indexY := x + 1, y + 1; (indexX < len(baduk.ground)) && (indexY < len(baduk.ground)) && baduk.ground[indexX][indexY] == value; indexX, indexY = indexX + 1, indexY + 1 {
		count++
	}
	
	if count >= 5 {
		return value
	}
	return 0
}

func plackBadukPiece(ground [][]int, player int, x int, y int) string {
	if x >= len(ground) || y >= len(ground) {
		return constUtil.Result_code_input_out_of_bound
	}else if ground[x][y] != 0 {
		return constUtil.Result_code_input_already_occupied
	}else{
		ground[x][y] = player
		return constUtil.Result_code_input_success
	}
}

func (baduk *Baduk) sendGroundToClients() {
	size := len(baduk.ground)
	data := make([]byte, size * size + 1)
	data[0] = byte(size)
	for index, _ := range baduk.ground {
		for index_, value_ := range baduk.ground[index] {
			data[index * size + index_ +1] = byte(value_)
		}
	}
	
	baduk.net.SendToClients(data)
}