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

func indexHandler(w http.ResponseWriter, r *http.Request)  {
  pageSTR := r.URL.Query().Get("page")
  if pageSTR == "" {
    pageSTR = "1"
  }
  fmt.Println("> GET /" + pageSTR)

  page, err := strconv.ParseInt(pageSTR, 10, 64)
  blocks := functions.GetBlocks(int(page)*functions.Maxblocks)

  if err != nil || len(blocks.Block) == 0 || page <= 0 {
    fmt.Fprintln(w, "Page not found")

  } else {
    tmpl := template.Must(template.ParseFiles("templates/home.html"))
    blockpage := types.BlockPage{Container: blocks,
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
    functions.GetJSON("/block/" + hash + ".json", &blockTarget)

    blockTime := time.Unix(blockTarget.Time, 0)
    blockTarget.Datetime = blockTime.Format("02.01.2006 15:04:05")

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
    fmt.Println("> GET tx:", hash)
    var txTarget = types.Tx{}
    functions.GetJSON("/tx/" + hash + ".json", &txTarget)

    tmpl := template.Must(template.ParseFiles("templates/tx.html"))
    tmpl.Execute(w, txTarget)

  } else {
    fmt.Fprintln(w, "Page not found")
  }
}

func main() {
  fmt.Println("SERVER STARTED :D")

  r := mux.NewRouter()
  r.HandleFunc("/", indexHandler)
  r.HandleFunc("/block/{id}", blockHandler)
  r.HandleFunc("/tx/{id}", txHandler)

  log.Fatal(http.ListenAndServe(":8080", r))
}
