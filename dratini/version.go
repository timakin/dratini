package dratini

import (
	"fmt"
	"runtime"
)

func PrintVersion() {
	fmt.Printf(`Dratini %s
Compiler: %s %s
`,
		Version,
		runtime.Compiler,
		runtime.Version())

}

func serverHeader() string {
	return fmt.Sprintf("Dratini %s", Version)
}
