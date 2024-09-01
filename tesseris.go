package tesseris

import (
	"github.com/tesseris-go/tesseris/cli"
	"github.com/tesseris-go/tesseris/config"
	"github.com/tesseris-go/tesseris/router"
)

type Tesseris struct {
	Config *config.Config
	cli    *cli.Cli
}

func Init(c *config.Config, r router.Routes) *Tesseris {
	c.Routes = r
	cli := cli.New(c)

	return &Tesseris{c, cli}
}

func (t *Tesseris) Run() {
	t.cli.Run()
}
