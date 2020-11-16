package main

import (
	"context"
	"encoding/json"
	"flag"
	_ "github.com/go-sql-driver/mysql"
	"github.com/snt999kk/colabnote/Logic"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"time"
)

var (
	/*upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024}*/
	configPath string
)

func init() {
	flag.StringVar(&configPath, "-configpath", "../../configs/colabnote.json", "path config file")
}
func main() {
	flag.Parse()
	conf := Logic.Config{}
	jsonFile, err := os.Open(configPath)
	defer jsonFile.Close()
	if err != nil {
		panic(err)
	}
	confjson, _ := ioutil.ReadAll(jsonFile)
	err = json.Unmarshal(confjson, &conf)
	if err != nil {
		panic(err)
	}
	myServer, err := Logic.NewServer(conf)
	if err != nil {
		log.Fatalln("could not initialize the server")
	}
	go func() {
		err = myServer.Run()
		if err != nil {
			log.Println("could not launch the server")
		}
	}()
	ctx, cancel := context.WithCancel(context.Background())
	go func(cancel context.CancelFunc) {
		sign := make(chan os.Signal)
		signal.Notify(sign, os.Interrupt)
		for sn := range sign {
			if sn == os.Interrupt {
				cancel()
			}
		}
	}(cancel)
	<-ctx.Done()
	ctxShutDown, cancel1 := context.WithTimeout(context.Background(), 5*time.Second)
	myServer.Shutdown(ctxShutDown)
	cancel1()
}
