package sirenaplay

import (
	"io"
	"os"
	"sync"

	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto"
)

var bz sync.Mutex

func SirenaPlay(file string) error {
	bz.Lock()
	defer bz.Unlock()
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()
	d, err := mp3.NewDecoder(f)
	if err != nil {
		return err
	}
	c, err := oto.NewContext(d.SampleRate(), 2, 2, 8192)
	if err != nil {
		return err
	}
	defer c.Close()
	p := c.NewPlayer()
	defer p.Close()
	if _, err := io.Copy(p, d); err != nil {
		return err
	}
	return nil
}
