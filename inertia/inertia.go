package inertia

import (
	"fmt"
	"os"
	"os/exec"

	inertiaGo "github.com/petaki/inertia-go"
	"github.com/tessaris/tessaris/config"
)

func CreateInertiaClient(cfg *config.Config) *inertiaGo.Inertia {
	url := "http://localhost"
	rootTemplate := fmt.Sprintf("./resources/views/%s.gohtml", cfg.InertiaView)
	version := ""

	i := inertiaGo.New(url, rootTemplate, version)

	return i
}

func RunInertiaServer(prod bool) {
	if prod {
		cmd := exec.Command("node", "../../bootstrap/ssr/ssr.js")
		cmd.Dir = "./resources/js"

		// cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err := cmd.Run()

		if err != nil {
			panic(err)
		}

		defer cmd.Process.Kill()
	} else {
		cmd := exec.Command("npm", "run", "dev")
		cmd.Dir = "."

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err := cmd.Run()

		if err != nil {
			panic(err)
		}

		defer cmd.Process.Kill()
	}
}
