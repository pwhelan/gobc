package main

import (
	"bytes"
	"blockchain"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

var hashes int64
var blocks int64

var host = "http://127.0.0.1:8080"

type status struct {
	sync.Mutex
	Height uint64
	Difficulty uint64
	LastBlock *blockchain.Block
}

type condwakeup sync.Cond

// NewCondWakeup creates a cond wakeup thingie
func newCondWakeup(locker sync.Locker) *condwakeup {
	return &condwakeup{L: locker}
}

type condwakeupthread struct {
	C chan bool
	S chan bool
}

func (cw *condwakeup) Thread() *condwakeupthread {
	c := make(chan bool)
	s := make(chan bool)

	go func (c chan<-bool, s chan bool) {
		for {
			select {
			default:
				//cw.Wait()
				c <- true
			case <-s:
				close(c)
				close(s)
				return
			}
		}
	}(c, s)

	return &condwakeupthread{
		C: c,
		S: s,
	}
}

func getstatus() (*status, error) {
	var s status

	resp, err := http.Get(fmt.Sprintf("%s/chain/status", host))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&s); err != nil {
		return nil, err
	}

	return &s, nil
}

func mine(cstatus <-chan *status, cmined chan bool) {
	status := <-cstatus
	myhashes := int64(0)
	var nonce blockchain.Hash

	rand.Seed(time.Now().Unix())
	rand.Read(nonce[:])

	fmt.Printf("\r[-] Start Mining: %d (diff: %d)\n",
		status.Height, status.Difficulty)

	for {
		select {
		case status = <- cstatus:
			fmt.Printf("\r[-] New Height: %d (diff: %d)\n",
				status.Height, status.Difficulty)
		default:
			mined, err := status.LastBlock.Generate(
				uint64(time.Now().Unix()), status.Difficulty, nonce)
			if err != nil {
				panic(err)
			}
			if myhashes >= 5000 {
				atomic.AddInt64(&hashes, myhashes)
				myhashes = 0
			}
			myhashes++
			if mined.HasDifficulty(status.Difficulty) {
				fmt.Printf("\r[-] Block mined at diff: %d (%d)\n",
					status.Difficulty, mined.Nonce.CountDiff())
				b := new(bytes.Buffer)
				json.NewEncoder(b).Encode(&mined)

				resp, err := http.Post(
					fmt.Sprintf("%s/chain/block", host),
					"application/json; charset=utf-8", b)
				if err != nil {
					panic(err)
				}
				resp.Body.Close()
				if resp.StatusCode == 201 || resp.StatusCode == 200 {
					fmt.Println("\r[!] BLOCK MINED!")
				} else {
					fmt.Printf("\r[%d] ERROR! ERROR!!!\n", resp.StatusCode)
					nonce.Inc()
					continue
				}
				//<-time.NewTimer(10 * time.Second).C
				rand.Read(nonce[:])
				cmined <- true
				fmt.Print("\rbegin again...")
				<-cmined
				fmt.Print("\rbegin again!!!")
			}
			nonce.Inc()
		}
	}
}

func drawarrow(arrow int) string {
	switch arrow % 4 {
	case 0:
		return "[>===]"
	case 1:
		return "[=>==]"
	case 2:
		return "[==>=]"
	case 3:
		return "[===>]"
	}
	return "[====]"
}

func main() {
	lasthashes := int64(0)
	if len(os.Args) >= 2 {
		host = os.Args[1]
	}

	curstatus, err := getstatus()
	if err != nil {
		panic(err)
	}

	cstatus := make(chan *status)
	cmined := make(chan bool)
	go mine(cstatus, cmined)

	cstatus <- curstatus

	t := time.NewTicker(5 * time.Second)
	d := time.NewTicker(1 * time.Second)
	lasttick := time.Now().UnixNano()
	arrow := 0
	for {
		select {
		case <-d.C:
			arrow++
			thistick := time.Now().UnixNano()
			curhashes := atomic.LoadInt64(&hashes)
			ticks := float64(thistick-lasttick)/(1000.0*1000.0*1000.0)
			fmt.Printf("\r%s %.02f h/s", drawarrow(arrow) ,
				float64(curhashes-lasthashes)/ticks)
			lasthashes = curhashes
			lasttick = thistick
		case <-cmined:
			nstatus, err := getstatus()
			for {
				if err == nil && nstatus.Height > curstatus.Height {
					curstatus = nstatus
					cmined <- true
					cstatus <- curstatus
					break
				} else {
					<-time.NewTimer(100 * time.Millisecond).C
				}
			}
		case <-t.C:
			//<-time.NewTimer(10 * time.Second).C
			nstatus, err := getstatus()
			if err == nil && nstatus.Height > curstatus.Height {
				curstatus = nstatus
				cstatus <- curstatus
			}
		}
	}
}
