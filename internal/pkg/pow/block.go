package pow

//Block represents necessary information for pow.
type Block struct {
	Timestamp int    `json:"timestamp"`
	Data      string `json:"data"`
	Hash      string `json:"hash"`
	Nonce     int    `json:"nonce"`
}
