package sirenaplay

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"

	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto"
)

var bz sync.Mutex

func SirenaPlay(file string) {
	bz.Lock()
	defer bz.Unlock()
	f, err := os.Open(file)
	if err != nil {
		log.Println(err)
		return
	}
	defer f.Close()
	d, err := mp3.NewDecoder(f)
	if err != nil {
		log.Println(err)
		return
	}
	c, err := oto.NewContext(d.SampleRate(), 2, 2, 8192)
	if err != nil {
		log.Println(err)
		return
	}
	defer c.Close()
	p := c.NewPlayer()
	defer p.Close()
	fmt.Printf("Length: %d[bytes]\n", d.Length())
	if _, err := io.Copy(p, d); err != nil {
		log.Println(err)
		return
	}
}
