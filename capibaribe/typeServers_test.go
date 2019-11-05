package capibaribe

import (
	"fmt"
)

func ExampleNewServerStruct() {

	server := NewServerStruct()
	server.AddExecutionTime(30)
	server.AddExecutionTime(40)
	server.AddExecutionTime(50)
	server.AddExecutionTime(60)
	server.AddExecutionTime(70)
	server.AddExecutionTime(80)
	server.AddExecutionTime(20)
	server.AddExecutionTime(90)
	server.AddExecutionTime(100)
	server.AddExecutionTime(110)
	server.AddExecutionTime(120)
	server.AddExecutionTime(130)
	server.AddExecutionTime(140)
	server.AddExecutionTime(150)

	fmt.Printf("Minimal execution time: %v\n", server.executionTimeMin)
	fmt.Printf("Maximal execution time: %v\n", server.executionTimeMax)
	fmt.Printf("List of execution time: %v\n", server.executionTimeList)

	// Output:
	// Minimal execution time: 20
	// Maximal execution time: 150
	// List of execution time: [70 80 20 90 100 110 120 130 140 150]

}
