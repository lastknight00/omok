package main

import (
	"fmt"
	"baduk"
)

func main(){
	for {
		b := baduk.New()
		fmt.Println("Initiating baduk ground...")
		ground := b.InitGround()
		fmt.Println("Finished initiation.")
		if ground == nil {
			fmt.Println("Program is finished")
			return
		}else{
			winner := b.StartOmok()
			if winner > 0 {
				fmt.Printf("Winner is %d player. Congratulations!!\n", winner)
			}else{
				fmt.Printf("Error : %d", winner)//TODO
			}
		}
	}
}
