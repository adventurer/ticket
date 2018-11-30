package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"sync/atomic"

	"github.com/satori/go.uuid"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

var TicketCount = make(chan int, 15)
var SetCnt = 5000

var BuyChan1 = make(chan Ticket, 0)
var BuyChan2 = make(chan Ticket, 0)
var BuyChan3 = make(chan Ticket, 0)
var BuyChan4 = make(chan Ticket, 0)
var BuyChan5 = make(chan Ticket, 0)
var BuyChanBak = make(chan Ticket, 0)

var SoldTicket sync.Map

type Ticket struct {
	No        string
	TicketSet []int
}

func main() {

	go initTicketPool()
	var wg sync.WaitGroup
	wg.Add(5)
	go sellPool(wg, BuyChan1, 1)
	go sellPool(wg, BuyChan2, 2)
	go sellPool(wg, BuyChan3, 3)
	go sellPool(wg, BuyChan4, 4)
	go sellPool(wg, BuyChan5, 5)

	errorChain := alice.New(loggerHandler, recoverHandler)
	r := mux.NewRouter()
	r.HandleFunc("/ticket", takeTicket)
	r.HandleFunc("/check", check)
	http.Handle("/", errorChain.Then(r))
	http.ListenAndServe(":8088", nil)

}

func checkTick() {

}

func sellPool(wg sync.WaitGroup, buyChan chan Ticket, chanCnt int) {
	wg.Done()
	emptyTicket := Ticket{}
	emptyTicketSet := []int{}
	emptyTicket.No = "0"
	emptyTicket.TicketSet = emptyTicketSet

	var ok = int32(1)
Begin:

	for {
		if ok == 1 {
			ticket := Ticket{}
			ticketSet := []int{}

			for index := 0; index < chanCnt; index++ {

				setNo := <-TicketCount

				ticketSet = append(ticketSet, setNo)
				if setNo == 0 {
					go func() {
						bakTicketCnt := len(ticketSet) - 1
						ticketSet = ticketSet[:len(ticketSet)-1]
						switch bakTicketCnt {
						case 1:
							u1 := uuid.Must(uuid.NewV4())

							ticket.No = fmt.Sprintf("%s", u1)
							ticket.TicketSet = ticketSet
							BuyChan1 <- ticket
							break
						case 2:
							u1 := uuid.Must(uuid.NewV4())

							ticket.No = fmt.Sprintf("%s", u1)
							ticket.TicketSet = ticketSet
							BuyChan2 <- ticket
							break
						case 3:
							u1 := uuid.Must(uuid.NewV4())

							ticket.No = fmt.Sprintf("%s", u1)
							ticket.TicketSet = ticketSet
							BuyChan3 <- ticket
							break
						case 4:
							u1 := uuid.Must(uuid.NewV4())

							ticket.No = fmt.Sprintf("%s", u1)
							ticket.TicketSet = ticketSet
							BuyChan4 <- ticket
							break
						case 5:
							u1 := uuid.Must(uuid.NewV4())

							ticket.No = fmt.Sprintf("%s", u1)
							ticket.TicketSet = ticketSet
							BuyChan5 <- ticket
							break
						}
					}()

					buyChan <- emptyTicket
					atomic.AddInt32(&ok, 1)
					goto Begin
				}
			}

			u1 := uuid.Must(uuid.NewV4())

			ticket.No = fmt.Sprintf("%s", u1)
			ticket.TicketSet = ticketSet
			log.Println(ticket)
			buyChan <- ticket
			continue
		}
		log.Println("on put ticket empty end")
		buyChan <- emptyTicket

	}

}

func initTicketPool() {
	for index := 1; ; index++ {
		if index <= SetCnt {
			TicketCount <- index
		} else {
			TicketCount <- 0
		}
	}
}

func check(w http.ResponseWriter, r *http.Request) {

	f := func(k, v interface{}) bool {
		key, _ := json.Marshal(k)
		data, _ := json.Marshal(v)
		w.Write(key)
		w.Write([]byte(":"))
		w.Write(data)
		w.Write([]byte("\n"))
		return true
	}
	SoldTicket.Range(f)

}

func takeTicket(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	c := r.Form.Get("c")
	i, _ := strconv.Atoi(c)

	switch i {
	case 1:
		ticket := <-BuyChan1
		log.Println("get ticket:", ticket)
		if ticket.No == "0" {
			w.Write([]byte("sold out"))
			return
		}
		SoldTicket.Store(ticket.No, ticket.TicketSet)
		// log.Println("take ticket:", ticket)
		w.Write([]byte(ticket.No))

	case 2:
		ticket := <-BuyChan2
		log.Println("get ticket:", ticket)

		if ticket.No == "0" {
			w.Write([]byte("sold out"))
			return
		}
		SoldTicket.Store(ticket.No, ticket.TicketSet)
		// log.Println("take ticket:", ticket)
		w.Write([]byte(ticket.No))

	case 3:
		ticket := <-BuyChan3
		log.Println("get ticket:", ticket)

		if ticket.No == "0" {
			w.Write([]byte("sold out"))
			return
		}
		SoldTicket.Store(ticket.No, ticket.TicketSet)
		// log.Println("take ticket:", ticket)
		w.Write([]byte(ticket.No))

	case 4:
		ticket := <-BuyChan4
		log.Println("get ticket:", ticket)

		if ticket.No == "0" {
			w.Write([]byte("sold out"))
			return
		}
		SoldTicket.Store(ticket.No, ticket.TicketSet)
		// log.Println("take ticket:", ticket)
		w.Write([]byte(ticket.No))

	case 5:
		ticket := <-BuyChan5
		log.Println("get ticket:", ticket)

		if ticket.No == "0" {
			w.Write([]byte("sold out"))
			return
		}
		SoldTicket.Store(ticket.No, ticket.TicketSet)
		// log.Println("take ticket:", ticket)
		w.Write([]byte(ticket.No))

	default:
		w.Write([]byte("other"))

	}

}

func recoverHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic: %#v", err)
				http.Error(w, http.StatusText(500), 500)
			}
		}()
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func loggerHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// start := time.Now()
		h.ServeHTTP(w, r)
		// log.Printf("<< %s %s %v", r.Method, r.URL.Path, time.Since(start))
	})
}
