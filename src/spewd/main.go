package main

import (
	"blockchain"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
)

// Replace the current chain if newBlocks wins
/*
func Replace(newBlocks []*Block) {
	if len(newBlocks) > len(Blockchain) {
		Blockchain = newBlocks
	}
}
*/


// Blockchain represents the current block chain
var Blockchain blockchain.Blockchain

func runhttpd() error {
	router := httprouter.New()

	router.POST("/chain", handleGetBlockchain)

	router.GET("/chain/status", handleGetChainStatus)
	router.GET("/chain/status/difficulty",
		handleGetDifficulty)

	router.GET("/chain/block", handleGetBlockchain)
	router.POST("/chain/block", handleWriteBlock)
	router.GET("/chain/block/:id", handleGetBlock)

	httpAddr := os.Getenv("HTTP_ADDR")
	log.Println("Listening on ", httpAddr)
	s := &http.Server{
		Addr:           ":" + httpAddr,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := s.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

func handleConn(conn net.Conn) {
	defer conn.Close()

}

func getRange(r *http.Request) (uint64, uint64, error) {
	var start, end uint64
	var err error

	if _, exists := r.Header["Range"]; exists {
		// Accept-Ranges: blocks
		// Range: blocks=123-532
		r := r.Header.Get("Range")
		if !strings.HasPrefix(r, "blocks=") {
			return 0, 0, fmt.Errorf("unsupported range type")
		}
		p := strings.Split(r[7:], "-")
		if len(p) != 2 {
			return 0, 0, fmt.Errorf("missing range part")
		}
		start, err = strconv.ParseUint(p[0], 10, 64)
		if err != nil {
			return 0, 0, err
		}
		end, err = strconv.ParseUint(p[1], 10, 64)
		if err != nil {
			return 0, 0, err
		}
	} else {
		start = 0
		end = Blockchain.Get(-1).Index
	}
	return start, end, nil
}

func handleGetBlockchain(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	start, end, err := getRange(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Accept-Ranges", "blocks")
	w.Header().Set("Content-Range", fmt.Sprintf("%d-%d", start, end))
	//w.Header().SetContentLength(-1)

	cursor := Blockchain.Cursor(int(start))
	io.WriteString(w, "[")
	first := true
	for block := cursor.Next(); block != nil && block.Index <= end; block = cursor.Next() {
		if !first {
			io.WriteString(w, ",")
		} else {
			first = false
		}
		bytes, _ := json.MarshalIndent(block, "", "\t")
		io.WriteString(w, string(bytes))
	}
	io.WriteString(w, "]")
}

func handleGetBlock(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.Atoi(ps.ByName("id"))
	block := Blockchain.Get(id)
	if block == nil {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
	bytes, err := json.MarshalIndent(block, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, string(bytes))
}

// https://developer.mozilla.org/en-US/docs/Web/HTTP/Range_requests
// Add HEAD handler to tell people we accept Range with blocks unit.
// Accept-Ranges: blocks
// Range: blocks=123-532
func handleGetBlocks(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	start, err := strconv.Atoi(ps.ByName("start"))
	if err != nil {
		respondWithJSON(w, r, http.StatusBadRequest, struct{
			Error string `json:"error"`
		}{
			Error: "bad start",
		})
		return
	}
	end, err := strconv.Atoi(ps.ByName("end"))
	if err != nil {
		respondWithJSON(w, r, http.StatusBadRequest, struct{
			Error string `json:"error"`
		}{
			Error: "bad end",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	cursor := Blockchain.Cursor(start)
	if cursor == nil {
		respondWithJSON(w, r, http.StatusNotFound, struct{
			Error string `json:"error"`
		}{
			Error: "no such block",
		})
		return
	}
	io.WriteString(w, "[")
	first := true
	for block := cursor.Next(); block != nil && block.Index <= uint64(end); block = cursor.Next() {
		if !first {
			io.WriteString(w, ",")
		} else {
			first = false
		}
		bytes, _ := json.MarshalIndent(block, "", "\t")
		io.WriteString(w, string(bytes))
	}
	io.WriteString(w, "]")
}

func handleGetDifficulty(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	respondWithJSON(w, r, http.StatusOK, Blockchain.GetDifficulty())
}

func handleGetChainStatus(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	block := Blockchain.Get(-1)
	status := struct{
		Height uint64
		Difficulty uint64
		LastBlock *blockchain.Block
		TotalDifficulty string
	}{
		Height: block.Index+1,
		Difficulty: Blockchain.GetDifficulty(),
		LastBlock: block,
		TotalDifficulty: Blockchain.CumulativeDifficulty.String(),
	}
	respondWithJSON(w, r, http.StatusOK, status)
}

func handleWriteBlock(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var block blockchain.Block

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&block); err != nil {
		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}
	defer r.Body.Close()

	if err := Blockchain.Add(&block); err != nil {
		respondWithJSON(w, r, http.StatusBadRequest, struct{
			Error string `json:"error"`
		}{
			Error: "invaid block",
		})
		return
	}

	fmt.Println("[!] Block Added")
	respondWithJSON(w, r, http.StatusCreated, block)
}

func respondWithJSON(w http.ResponseWriter, r *http.Request, code int, payload interface{}) {
	response, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("HTTP 500: Internal Server Error"))
		return
	}
	w.WriteHeader(code)
	w.Write(response)
}

func runnode() error {
	// start TCP and serve TCP server
	server, err := net.Listen("tcp", ":"+os.Getenv("NODE_ADDR"))
	if err != nil {
		log.Fatal(err)
	}
	defer server.Close()

	for {
		conn, err := server.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handleConn(conn)
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	genesisBlock := blockchain.Block{
		Index: 0,
		Timestamp: 1527949914,
	}

	genesisBlock.Hash = genesisBlock.GenerateHash()
	Blockchain.Genesis(&genesisBlock)
	go startPeer()

	log.Fatal(runhttpd())
}
