package main

import (
	"github.com/vs-uulm/go-subjectivelogic/pkg/subjectivelogic"
	"fmt"

)

func main() {

	op1, _ := subjectivelogic.NewOpinion(0.2,0.3,0.5,0.5)
	fmt.Println(op1)
	

}

