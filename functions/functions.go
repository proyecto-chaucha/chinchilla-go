package functions

import (
  "log"
  "time"
  "strconv"
  "strings"
  "net/http"
  "io/ioutil"
  "encoding/json"
  "github.com/proyecto-chaucha/chinchilla-go/types"
)

var Maxblocks int = 39
var apiURL string = "http://localhost:21662/rest"
var myClient = &http.Client{Timeout: 10 * time.Second}

func GetBlocks(offset int) types.BlocksContainer {
  // Block container
  blocks := types.BlocksContainer{[]types.Block{}}

  // Get blockcount
  target := GetWeb("/blockcount.json")
  height64, err := strconv.ParseInt(strings.TrimSpace(target), 10, 0)
  if err != nil { log.Fatal(err) }
  height := int(height64)

  // Page offset
  var off int
  if offset > Maxblocks { off = int(offset) - Maxblocks } else { off = 0 }

  if height - off - Maxblocks + 1 >= 0 {
    for i := height - off; i >= height - off - Maxblocks + 1; i-- {
      STRi := strconv.Itoa(i)
      hashTarget := GetWeb("/getblockhash/" + STRi + ".json")
      hash := hashTarget

      blockTarget := types.Block{}
      GetJSON("/block/" + hash + ".json", &blockTarget)

      blockTime := time.Unix(blockTarget.Time, 0)
      blockTarget.Datetime = blockTime.Format("02.01.2006 15:04:05")
      blockTarget.Txcount = len(blockTarget.Tx)

      blocks.AddItem(blockTarget)
    }
  }
  return blocks
}

func GetJSON(url string, target interface{}) error {
  resp, err := myClient.Get(apiURL + url)
  if err != nil { log.Fatal(err) }
  defer resp.Body.Close()
  return json.NewDecoder(resp.Body).Decode(target)
}

func GetWeb(url string) string  {
  resp, err := myClient.Get(apiURL + url)
  if err != nil { log.Fatal(err) }
  defer resp.Body.Close()
  html, err := ioutil.ReadAll(resp.Body)
  if err != nil { log.Fatal(err) }
  return string(html)
}
