package main

import (
	"github.com/gorilla/mux"
  "fmt"
  "log"
  "time"
  "strconv"
  "net/http"
  "io/ioutil"
  "encoding/json"
	"html/template"
)

type blocksContainer struct {
  Block []block
}

type chain struct {
  Height int `json:"height"`
  Total_amount float32 `json:"total_amount"`
}

type block struct {
  Hash string `json:"hash"`
  Height int `json:"height"`
  Time int64 `json:"time"`
  Datetime string
  Difficulty float32 `json:"difficulty"`
  Size int `json:"size"`
  Version int `json:"versionHex"`
  Previousblockhash string `json:"previousblockhash"`
  Merkleroot string `json:"merkleroot"`
  Bits string `json:"bits"`
  Nonce int64 `json:"nonce"`
  Txcount int
  Tx []tx `json:"tx"`
}

type tx struct {
  Txid string `json:"txid"`
  Vin []struct {
      Coinbase string `json:"coinbase"`
      Txid string `json:"txid"`
      Vout int `json:"vout"`
      } `json:"vin"`
  Vout []struct {
      Value float32 `json:"value"`
      N int `json:"n"`
      ScriptPubKey struct {
          Hex string `json:"hex"`
          Type string `json:"type"`
          Address [1]string `json:"addresses"`
          } `json:"scriptPubKey"`
      } `json:"vout"`
}

var myClient = &http.Client{Timeout: 10 * time.Second}
var apiURL string = "http://localhost:21662/rest"

func getJson(url string, target interface{}) error {
  resp, err := myClient.Get(url)
  if err != nil { log.Fatal(err) }
  defer resp.Body.Close()
  return json.NewDecoder(resp.Body).Decode(target)
}

func getHash(url string) string  {
  resp, err := myClient.Get(url)
  if err != nil { log.Fatal(err) }
  defer resp.Body.Close()
  html, err := ioutil.ReadAll(resp.Body)
  if err != nil { log.Fatal(err) }
  return string(html)
}

func (container *blocksContainer) AddItem(item block) []block {
    container.Block = append(container.Block, item)
    return container.Block
}

func index(w http.ResponseWriter, r *http.Request)  {
  fmt.Println("> GET /")
  blocks := blocksContainer{[]block{}}

  chainTarget := chain{}
  getJson("http://localhost:21662/rest/chaininfo.json", &chainTarget)

  for i := chainTarget.Height; i > chainTarget.Height - 20; i-- {
    // get blockhash
    hash := getHash(apiURL + "/getblockhash/" + strconv.Itoa(i) + ".json")

    // get each block
    blockTarget := block{}
    getJson(apiURL + "/block/" + hash + ".json", &blockTarget)

    // block manipulation
    blockTime := time.Unix(blockTarget.Time, 0)
    blockTarget.Datetime = blockTime.Format("02.01.2006 15:04:05")
    blockTarget.Txcount = len(blockTarget.Tx)

    // add block to container
    blocks.AddItem(blockTarget)
  }
  tmpl := template.Must(template.ParseFiles("templates/home.html"))
  tmpl.Execute(w, blocks)
}

func blockHandler(w http.ResponseWriter, r *http.Request)  {
  vars := mux.Vars(r)
  hash := vars["id"]
  if len(hash) == 64 {
    fmt.Println("> GET block:", hash)
    var blockTarget = block{}
    getJson(apiURL + "/block/" + hash + ".json", &blockTarget)
    tmpl := template.Must(template.ParseFiles("templates/block.html"))
    tmpl.Execute(w, blockTarget)
  } else {
    fmt.Fprintln(w, "Error :c")
 }
}

func txHandler(w http.ResponseWriter, r *http.Request)  {
  vars := mux.Vars(r)
  hash := vars["id"]
  if len(hash) == 64 {
    var txTarget = tx{}
    getJson(apiURL + "/tx/" + hash + ".json", &txTarget)
    tmpl := template.Must(template.ParseFiles("templates/tx.html"))
    tmpl.Execute(w, txTarget)
  } else {
    fmt.Fprintln(w, "Error :c")
  }
}

func main() {
  fmt.Println("SERVER STARTED :D")

  r := mux.NewRouter()
	r.HandleFunc("/", index)
	r.HandleFunc("/block/{id}", blockHandler)
	r.HandleFunc("/tx/{id}", txHandler)
  log.Fatal(http.ListenAndServe(":8080", r))
}
