package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sirena2/internal/sirenaplay"
	"strings"
	"syscall"
	"time"
)

func main() {
	var alarm bool
	key := flag.String("key", "", "api key")
	obl := flag.Int("oblast", 31, "Oblast")
	trevoga := flag.String("trevoga", "Sub.mp3", "File name alarm on")
	vidbiy := flag.String("vidbiy", "Sub.mp3", "File name alarm off")
	test := flag.String("test", "", "Test audio file")
	flag.Parse()
	if *test != "" {
		err := sirenaplay.SirenaPlay(*test)
		if err != nil {
			log.Println(err)
		}
		return
	}
	if *key == "" {
		log.Fatalln("No API key")
	}
	_, err := os.Stat(*trevoga)
	if err != nil {
		log.Fatalln(err)
	}
	_, err = os.Stat(*vidbiy)
	if err != nil {
		log.Fatalln(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	signal_chan := make(chan os.Signal, 1)
	signal.Notify(
		signal_chan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGABRT,
	)
	go func(ctx context.Context) {
		tk := time.NewTicker(time.Millisecond * 7500)
		for {
			select {
			case <-ctx.Done():
				return
			case <-tk.C:
				reqURL := fmt.Sprintf("https://api.alerts.in.ua/v1/iot/active_air_raid_alerts/%d.json", *obl) + "?token=" + *key
				req, err := http.NewRequest(http.MethodGet, reqURL, nil)
				if err != nil {
					log.Println(err)
					continue
				}
				req.Header.Add("Authorization", *key)
				resp, err := http.DefaultClient.Do(req)
				if err != nil {
					log.Println(err)
					continue
				}
				resBody, err := io.ReadAll(resp.Body)
				if err != nil {
					log.Println(err)
					continue
				}
				resp.Body.Close()
				stat := strings.Trim(string(resBody), `"`)
				log.Println(stat)
				switch stat {
				case "A", "P":
					if !alarm {
						log.Println("Trevoga")
						alarm = true
						err := sirenaplay.SirenaPlay(*trevoga)
						if err != nil {
							log.Println(err)
						}
					}
				case "N":
					if alarm {
						log.Println("Otboy")
						alarm = false
						err := sirenaplay.SirenaPlay(*vidbiy)
						if err != nil {
							log.Println(err)
						}
					}
				default:
					log.Println("unknow answer", stat)
				}
			}
		}
	}(ctx)
	s := <-signal_chan
	log.Println(s.String())
	cancel()
}
