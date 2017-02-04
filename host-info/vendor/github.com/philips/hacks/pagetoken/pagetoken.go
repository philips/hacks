package main

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"fmt"
)

type Token struct {
	Limit uint8
	Page uint16
}

func main() {
	// Encode
	buf := new(bytes.Buffer)
	tok := Token{100, 4}
	binary.Write(buf, binary.LittleEndian, tok)
	enc := base64.URLEncoding.EncodeToString(buf.Bytes())
	println(enc)

	it := &Token{}
	dec, _ := base64.URLEncoding.DecodeString(enc)
	db := bytes.NewBuffer(dec)
	binary.Read(db, binary.LittleEndian, it)
	fmt.Printf("%v", it)
}
