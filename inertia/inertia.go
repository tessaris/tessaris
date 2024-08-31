package inertia

import (
	"os"
	"os/exec"

	"github.com/petaki/inertia-go"
)

func CreateInertiaClient() *inertia.Inertia {
	url := "http://localhost"                // Application URL for redirect
	rootTemplate := "./resources/app.gohtml" // Root template, see the example below
	version := ""                            // Asset version

	i := inertia.New(url, rootTemplate, version)
	i.EnableSsrWithDefault()

	return i
}

type InertiaServer struct {
	p *os.Process
}

func RunInertiaServer(p chan<- *os.Process) {
	cmd := exec.Command("node", "../../bootstrap/ssr/ssr.js")
	cmd.Dir = "./resources/js"

	// cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	if err != nil {
		panic(err)
	}

	defer cmd.Process.Kill()

	p <- cmd.Process
}

func CreateInertiaServer() *InertiaServer {
	pChan := make(chan *os.Process)

	go RunInertiaServer(pChan)
	p := <-pChan

	return &InertiaServer{p}
}
