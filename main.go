package main

import (
  "fmt"
  "log"
  "time"
  "strconv"
  "net/http"
  "html/template"
  "github.com/gorilla/mux"
  "github.com/proyecto-chaucha/chinchilla-go/types"
  "github.com/proyecto-chaucha/chinchilla-go/functions"
)

type blockPage struct {
  Container BlocksContainer
  Page int64
  PageNext int64
  PagePrev int64
}

var apiURL string = "http://localhost:21662/rest"
var maxBlocks int = 12

type BlocksContainer struct {
  Block []types.Block
}

func (container *BlocksContainer) AddItem(item types.Block) []types.Block {
    container.Block = append(container.Block, item)
    return container.Block
}

func getblocks(offset int) BlocksContainer {
  blocks := BlocksContainer{[]types.Block{}}
  target := types.Chain{}
  functions.GetJSON("http://localhost:21662/rest/chaininfo.json", &target)

  var off int
  if offset > maxBlocks { off = int(offset) - maxBlocks } else { off = 0 }

  if target.Height - off - maxBlocks + 1 >= 0 {
    for i := target.Height - off; i >= target.Height - off - maxBlocks + 1; i-- {
      hash := functions.GetHash(apiURL + "/getblockhash/" + strconv.Itoa(i) + ".json")

      blockTarget := types.Block{}
      functions.GetJSON(apiURL + "/block/" + hash + ".json", &blockTarget)

      blockTime := time.Unix(blockTarget.Time, 0)
      blockTarget.Datetime = blockTime.Format("02.01.2006 15:04:05")
      blockTarget.Txcount = len(blockTarget.Tx)

      blocks.AddItem(blockTarget)
    }
  }
  return blocks
}

func index(w http.ResponseWriter, r *http.Request)  {
  fmt.Println("> GET /")
  var page int64 = 1
  blocks := getblocks(int(page))

  tmpl := template.Must(template.ParseFiles("templates/home.html"))
  blockpage := blockPage{Container: blocks,
                         Page: page,
                         PageNext: page + 1,
                         PagePrev: page - 1}
  tmpl.Execute(w, blockpage)
}

func indexOffset(w http.ResponseWriter, r *http.Request)  {
  vars := mux.Vars(r)
  fmt.Println("> GET /" + vars["id"])
  page, err := strconv.ParseInt(vars["id"], 10, 64)

  blocks := getblocks(int(page)*maxBlocks)

  if err != nil || len(blocks.Block) == 0 {
    fmt.Fprintln(w, "Page not found")

  } else {
    tmpl := template.Must(template.ParseFiles("templates/home.html"))
    blockpage := blockPage{Container: blocks,
                           Page: page,
                           PageNext: page + 1,
                           PagePrev: page - 1}
    tmpl.Execute(w, blockpage)
  }
}

func blockHandler(w http.ResponseWriter, r *http.Request)  {
  vars := mux.Vars(r)
  hash := vars["id"]

  if len(hash) == 64 {
    fmt.Println("> GET block:", hash)

    var blockTarget = types.Block{}
    functions.GetJSON(apiURL + "/block/" + hash + ".json", &blockTarget)

    tmpl := template.Must(template.ParseFiles("templates/block.html"))
    tmpl.Execute(w, blockTarget)

  } else {
    fmt.Fprintln(w, "Page not found")
 }
}

func txHandler(w http.ResponseWriter, r *http.Request)  {
  vars := mux.Vars(r)
  hash := vars["id"]

  if len(hash) == 64 {
    var txTarget = types.Tx{}
    functions.GetJSON(apiURL + "/tx/" + hash + ".json", &txTarget)

    tmpl := template.Must(template.ParseFiles("templates/tx.html"))
    tmpl.Execute(w, txTarget)

  } else {
    fmt.Fprintln(w, "Page not found")
  }
}

func main() {
  fmt.Println("SERVER STARTED :D")

  r := mux.NewRouter()
  r.HandleFunc("/", index)
  r.HandleFunc("/{id}", indexOffset)
  r.HandleFunc("/block/{id}", blockHandler)
  r.HandleFunc("/tx/{id}", txHandler)

  log.Fatal(http.ListenAndServe(":8080", r))
}
