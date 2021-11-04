package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/MartinPfe/go-microservices-firststeps/handlers"
)

func main() {
	l := log.New(os.Stdout, "product-api", log.LstdFlags)

	//Armamos los 2 hadnlers 
	hh := handlers.NewHellow(l)
	gh := handlers.NewGoodbye(l)

	//Creamos el servidor, que va a contener nuestros handlers
	sm := http.NewServeMux()

	//Armamos las rutas para cada uno de los handlers
	sm.Handle("/", hh)
	sm.Handle("/goodbye", gh)

	//Config del servidor
	s:= &http.Server{
		Addr: 			"localhost:9090",
		Handler: 		sm,
		IdleTimeout:  	120 *time.Second,
		ReadTimeout: 	1 * time.Second,
		WriteTimeout: 	1 * time.Second,
	}

	//Lo iniciamos en una goroutine para despues poder terminarlos
	go func() {
		err := s.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}()

	//Seteamos que los channels van a esperar a que llegue una señal de kill
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	//Cuando llega la señal de kill, retomamso
	sig := <- sigChan
	l.Println("Senal de de fin recibida", sig)

	//Matamos el server
	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(tc)
}