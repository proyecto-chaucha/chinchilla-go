package types

type Block struct {
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
  Tx []Tx `json:"tx"`
}

type Tx struct {
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
