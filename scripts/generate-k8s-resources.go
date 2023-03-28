package main
import (
	"fmt"
	"os"
	"strings"
	"encoding/json"
)

func main() {
	svcData, err := os.ReadFile(os.Getenv("MICROSERVICES_JSON_FILE_PATH"))
	if err != nil { panic(err) }
	var svcDefs []SvcDefinition
	if err := json.Unmarshal([]byte(svcData), &svcDefs); err != nil {
		panic(err)
	}

	predefinedTemplate := strings.TrimSpace(os.Getenv("PREDEFINED_TEMPLATE"))
	if len(predefinedTemplate) > 0 {
		defs := strings.Split(predefinedTemplate, "|")
		frontend_def := strings.Split(defs[0], ":")
		backend_def := strings.Split(defs[1], ":")
		svcDefs[0].Type = frontend_def[0]
		svcDefs[0].Version = frontend_def[1]
		svcDefs[1].Type = backend_def[0]
		svcDefs[1].Version = backend_def[1]
	}

	fmt.Println("App definitions:", len(svcDefs))
	for _, svcDef := range svcDefs {
		fmt.Println("App definition:", string(svcJson))
	}
}
