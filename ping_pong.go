//
// goroutines interactions
//
package main

import(
	"fmt"
	"time"
)

const (
	APP_TIMEOUT = 800 // ms
	RESPONSE_TIMEOUT = 8000 // us
)


//
// main driver
//
func main() {
	fmt.Println("[APP] started")
	var done chan struct{};

	done = make(chan struct{})

	// channels for control signals
	var pingSigCh chan struct{command string}
	var pongSigCh chan struct{command string}
	// channels for messaging
	var pingMsgCh chan struct{text string}
	var pongMsgCh chan struct{text string}
	// channel for workers sync
	var workersCh chan bool

	// init chans
	pingSigCh = make(chan struct{command string})
	pongSigCh = make(chan struct{command string})
	pingMsgCh = make(chan struct{text string})
	pongMsgCh = make(chan struct{text string})
	workersCh = make(chan bool)

	// ping Chi Sin worker
	go func() {
		// robust system)
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("[W] Ping FAULT")
			}
		}()
		// worker done signal
		defer func() {
			workersCh<- true
		}()
		fmt.Println("[PING] started")

		// start interaction
		fmt.Println("[PING] Ping message to Pong")
		pingSigCh<- struct{command string}{command: "start",}

		// reading messages and control sygnal from multiple channels
		for stop := false; stop != true; {
			// process multiple channels
			select {
			case signal := <-pongSigCh: // process signal
				switch signal.command {
				case "stop":
					fmt.Println("[PING] Slowing down...")
					fmt.Println("[PING] Signal Pong to stop too")
					pingSigCh<- struct{command string}{command: "stop"}
					stop = true;
				default:
					fmt.Println("[PING] Pong singnal ->", signal)
				}
			case message := <-pongMsgCh: // process messages
				switch message {
				default:
					fmt.Println("[PING] Pong message ->", message.text)
					// reply
					pingMsgCh<- struct{text string}{text: "Hello, Pong",}
				}
			case <-time.After(APP_TIMEOUT * time.Microsecond): // safety timeout
				fmt.Println("[PING] Pong response timeout")
			} // eof select
		} // eof for
		fmt.Println("[PING] trerminated")
	}()

	// pong Lu Eng worker
	go func() {
		// robust system)
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("[W] Pong FAULT")
			}
		}()
		// worker done signal
		defer func() {
			workersCh<- true
		}()
		fmt.Println("[PONG] started")
		// reading messages and control sygnal from multiple channels
		for stop := false; stop != true; {
			// process multiple channels
			select {
			case signal := <-pingSigCh: // process signal
				switch signal.command {
				case "stop":
					fmt.Println("[PONG] Slowing down...")
					stop = true;
				case "start":
					fmt.Println("[PONG] Start a converstaion")
					// send first message
					pongMsgCh<- struct{text	string}{text: "Hello, Ping"}
				default:
					fmt.Println("[PONG] Pong singnal ->", signal.command)
				}
			case message := <-pingMsgCh: // process message
				switch message {
				default:
					fmt.Println("[PONG] Pong message ->", message.text)
					// Reply
					//pingMsgCh<- struct{text	string}{text: "Hello, Ping"}
					pongSigCh<- struct{command string}{command: "stop"}
				}
			case <-time.After(RESPONSE_TIMEOUT * time.Microsecond): // safety timeout
				fmt.Println("[PONG] Pong response timeout")
			} // eof select
		} // eof for
		fmt.Println("[PONG] trerminated")
	}()


	// group lider
	go func() {
		// robust system)
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("[W] Supervisor FAULT")
			}
		}()
		// sync siglanl when everything is done
		defer close(done)
		fmt.Println("[SUP] started")
		workersDoneCount := 0
		for range workersCh {
			workersDoneCount++
			fmt.Println("[SUP] workers done:",workersDoneCount)
			if workersDoneCount >= 2 { // we have got just Ping and Pong workers
				fmt.Println("[SUP] all workers done!")
				// stop supervising
				close(workersCh)
			}
		}
		fmt.Println("[SUP] terminated")
	}()

	// controller-syncronizer
	select {
	case <-done:
		fmt.Println("[APP] done")
	case <-time.After(APP_TIMEOUT * time.Millisecond):
		fmt.Println("[APP] timeout")
	} // eof select
} // eof main

