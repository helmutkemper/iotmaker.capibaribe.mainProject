//todo: tem que ter um log avisando quando foi a ultima requisição
//todo: microsservico deveria ter um id inicial que permita seguir uma requisição em caso de erro
package main

import (
	capib "./capibaribe"

	"flag"

	"sync"
)

func main() {
	var wg sync.WaitGroup
	var capibaribe capib.MainConfig

	wg.Add(1)

	//mutexTerminateServer = make( map[string]*sync.WaitGroup )

	filePath := flag.String("f", "./capibaribe-config.yml", "./your-capibaribe-config-file.yml")
	flag.Parse()

	capibaribe.LoadConfAndStart(*filePath)

	wg.Wait()
}
