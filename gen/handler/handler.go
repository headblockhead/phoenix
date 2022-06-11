package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image/jpeg"
	"log"
	"net/http"
	"time"

	gen "github.com/headblockhead/phoenix"
)

type Response struct {
	Frames    [][]byte
	Objection []byte
}

func Handle(w http.ResponseWriter, r *http.Request) {
	var resp Response
	seed := int(time.Now().Unix())
	frames, objection, err := gen.Generate(seed)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to generate: %v", err), http.StatusInternalServerError)
		return
	}
	if objection != nil {
		o := new(bytes.Buffer)
		err = jpeg.Encode(o, objection, &jpeg.Options{Quality: 100})
		if err != nil {
			log.Printf("failed to encode the objection: %v\n", err)
		}
		resp.Objection = o.Bytes()
	}
	resp.Frames = make([][]byte, len(frames))
	for i := 0; i < len(frames); i++ {
		b := new(bytes.Buffer)
		err = jpeg.Encode(b, frames[i], &jpeg.Options{Quality: 100})
		if err != nil {
			log.Printf("failed to encode frame %d: %v\n", i, err)
		}
		resp.Frames[i] = b.Bytes()
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", " ")
	enc.Encode(resp)
}