package main

import (
	"github.com/vs-uulm/go-subjectivelogic/pkg/subjectivelogic"
	"fmt"

)

func main() {

	op1, _ := subjectivelogic.NewOpinion(0.2,0.3,0.5,0.5)
	op2, _ := subjectivelogic.NewOpinion(0.1,0.2,0.7,0.5)
	fmt.Println(op1, op2)
	
	op3, _ := subjectivelogic.Addition(&op1, &op2)
	fmt.Println(op3)

}

