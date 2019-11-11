//todo: tem que ter um log avisando quando foi a ultima requisição
//todo: microsservico deveria ter um id inicial que permita seguir uma requisição em caso de erro
package main

import (
	"github.com/helmutkemper/iotmaker.capibaribe.module"

	"flag"

	"sync"
)

func main() {
	var wg sync.WaitGroup
	var capibaribe iotmaker_capibaribe_module.MainConfig

	wg.Add(1)

	filePath := flag.String("f", "./capibaribe-config.yml", "./your-capibaribe-config-file.yml")
	flag.Parse()

	capibaribe.LoadConfAndStart(*filePath)

	wg.Wait()
}
