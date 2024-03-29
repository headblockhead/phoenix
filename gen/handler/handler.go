package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image/jpeg"
	"log"
	"net/http"
	"strconv"
	"time"

	gen "github.com/headblockhead/phoenix"
	"github.com/nfnt/resize"
)

type Response struct {
	Frames            [][]byte
	Objection         []byte
	Seed              int
	ObjectionLocation int
}

func Handle(w http.ResponseWriter, r *http.Request) {
	var resp Response
	var queries = r.URL.Query()
	seed := int(time.Now().UTC().UnixNano())
	if queries.Get("seed") != "" {
		var err error
		seed, err = strconv.Atoi(queries.Get("seed"))
		if err != nil {
			log.Println(err)
			http.Error(w, "Invalid seed", http.StatusBadRequest)
			return
		}
	}
	if queries.Get("scale") == "" {
		http.Error(w, "Missing scale query", http.StatusBadRequest)
		return
	}
	resp.Seed = seed
	frames, objection, ObjectionLocation, err := gen.Generate(seed)
	resp.ObjectionLocation = ObjectionLocation
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to generate: %v", err), http.StatusInternalServerError)
		return
	}
	imageScaleInt, err := strconv.Atoi(queries.Get("scale"))
	if err != nil {
		http.Error(w, "Bad scale value", http.StatusBadRequest)
		return
	}
	if objection != nil {
		o := new(bytes.Buffer)
		if (imageScaleInt) != 100 {
			err = jpeg.Encode(o, objection, &jpeg.Options{Quality: 50})
			if err != nil {
				log.Printf("failed to encode the objection: %v\n", err)
			}
		} else {
			err = jpeg.Encode(o, objection, &jpeg.Options{Quality: 50})
			if err != nil {
				log.Printf("failed to encode the objection: %v\n", err)
			}
		}
		resp.Objection = o.Bytes()
	}
	imageScalePercent := float32(imageScaleInt) / 100
	for i := 0; i < len(frames); i++ {
		newImage := resize.Resize(uint(960*imageScalePercent), uint(640*imageScalePercent), frames[i], resize.Lanczos3)
		frames[i] = newImage
	}
	resp.Frames = make([][]byte, len(frames))
	for i := 0; i < len(frames); i++ {
		b := new(bytes.Buffer)
		if (imageScaleInt) != 100 {
			err = jpeg.Encode(b, frames[i], &jpeg.Options{Quality: 50})
			if err != nil {
				log.Printf("failed to encode frame %d: %v\n", i, err)
			}
		} else {
			err = jpeg.Encode(b, frames[i], &jpeg.Options{Quality: 100})
			if err != nil {
				log.Printf("failed to encode frame %d: %v\n", i, err)
			}
		}

		resp.Frames[i] = b.Bytes()
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", " ")
	enc.Encode(resp)
}
