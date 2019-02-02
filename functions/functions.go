package functions

import (
  "log"
  "time"
  "net/http"
  "io/ioutil"
  "encoding/json"
)

var myClient = &http.Client{Timeout: 10 * time.Second}

func GetJSON(url string, target interface{}) error {
  resp, err := myClient.Get(url)
  if err != nil { log.Fatal(err) }
  defer resp.Body.Close()
  return json.NewDecoder(resp.Body).Decode(target)
}

func GetHash(url string) string  {
  resp, err := myClient.Get(url)
  if err != nil { log.Fatal(err) }
  defer resp.Body.Close()
  html, err := ioutil.ReadAll(resp.Body)
  if err != nil { log.Fatal(err) }
  return string(html)
}
