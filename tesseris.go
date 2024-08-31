package tesseris

import "github.com/tesseris-go/tesseris/router"

type Tesseris struct {
	Routes []router.InertiaRoute
}

func Init(r []router.InertiaRoute) *Tesseris {
	return &Tesseris{r}
}
