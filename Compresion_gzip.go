// You can edit this code!
// Click here and start typing.
package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"time"
)

type Embeddings struct {
	Version    int       `json:"v"`
	Embeddings []float32 `json:"e"`
}

const (
	desiredStdDev = 0.5
)

func GenerateEmbeddings(size int) (*[]float64, error) {
	if size < 0 {
		return nil, fmt.Errorf("size can't be negative")
	}
	data := make([]float64, size)
	for i := range data {
		// NormFloat64 will generate a random normalized float number with mean = 0 and stddev = 1
		// To produce a different normal distribution we can use
		// NormFloat64() * desiredStdDev + desiredMean
		data[i] = rand.NormFloat64() * desiredStdDev
	}
	return &data, nil
}

func CutPrecision(data *[]float64) (*[]float32, error) {
	newData := make([]float32, len(*data))
	for i, v := range *data {
		newData[i] = float32(v)
	}
	return &newData, nil
}

func main() {
	var start time.Time
	start = time.Now()
	randomEmbbedingsData, _ := GenerateEmbeddings(50)
	data, _ := CutPrecision(randomEmbbedingsData)
	fmt.Printf("[GENERATION] Elapsed: %+v\n", time.Since(start))
	fmt.Printf("Data: %+v\n\n", *data)
	start = time.Now()
	msg := Embeddings{Version: 20220505, Embeddings: *data}
	value, err := json.Marshal(msg)
	if err != nil {
		fmt.Printf("error: %v", err)
	}
	fmt.Printf("[JSON] Elapsed: %+v\n", time.Since(start))
	fmt.Printf("Msg: %+v\n\n", msg)
	fmt.Printf("Value: %+v\n", value)
	fmt.Printf("Value: %T - Original size: %d\n\n", value, len(value))

	// bytes
	start = time.Now()
	var gzipValue bytes.Buffer
	gz, err := gzip.NewWriterLevel(&gzipValue, gzip.BestCompression)
	if err != nil {
		fmt.Printf("error: %v", err)
	}
	gz.Write(value)
	gz.Close()
	fmt.Printf("[GZIP 1] Elapsed: %+v\n", time.Since(start))
	gzipBytes := gzipValue.Bytes()
	fmt.Printf("[GZIP 1] buffer: %+v\n", gzipBytes)
	fmt.Printf("[GZIP 1] buffer: %T - Compressed size: %d\n\n", gzipBytes, len(gzipBytes))

	// read
	start = time.Now()
	zr, err := gzip.NewReader(&gzipValue)
	if err != nil {
		fmt.Printf("error: %v", err)
	}
	var unGZIP bytes.Buffer
	_, err = unGZIP.ReadFrom(zr)
	if err != nil {
		fmt.Printf("error: %v", err)
	}
	defer zr.Close()
	fmt.Printf("[READ] Elapsed: %+v\n", time.Since(start))
	unGZIPBytes := unGZIP.Bytes()
	fmt.Printf("[READ] Uncompressed Size: %+v\n", len(unGZIPBytes))
	mismatch := 0
	for i := 0; i < len(unGZIPBytes); i++ {
		if value[i] != unGZIPBytes[i] {
			mismatch++
		}
	}
	if mismatch > 0 {
		fmt.Printf("[READ] Slices doesn't match: %+v\n", unGZIPBytes)
	} else {
		fmt.Println("[READ] Slices match!!")
	}
	fmt.Println()

	// Another method of compression
	start = time.Now()
	var inlineGZIP bytes.Buffer
	ngz, err := gzip.NewWriterLevel(&inlineGZIP, gzip.BestCompression)
	if err != nil {
		fmt.Printf("error: %v", err)
	}
	// encode adds 0xOA char (new line) to the stream
	if err := json.NewEncoder(ngz).Encode(msg); err != nil {
		fmt.Printf("error: %v", err)
	}
	ngz.Close()
	fmt.Printf("[GZIP 2] Elapsed: %+v\n", time.Since(start))
	inlineGZIPBytes := inlineGZIP.Bytes()
	inlineGZIP2 := inlineGZIP
	fmt.Printf("[GZIP 2] ngz: %+v\n", inlineGZIPBytes)
	fmt.Printf("[GZIP 2] ngz: %T - Compressed size: %d\n\n", inlineGZIPBytes, len(inlineGZIPBytes))

	// read
	start = time.Now()
	nzr, err := gzip.NewReader(&inlineGZIP)
	if err != nil {
		fmt.Printf("error: %v", err)
	}
	var nunGZIP bytes.Buffer
	_, err = nunGZIP.ReadFrom(nzr)
	if err != nil {
		fmt.Printf("error: %v", err)
	}
	defer nzr.Close()
	fmt.Printf("[READ] Elapsed: %+v\n", time.Since(start))
	nunGZIPBytes := nunGZIP.Bytes()
	fmt.Printf("[READ] Uncompressed Size: %+v\n", len(nunGZIPBytes))
	mismatch = 0
	for i := 0; i < len(nunGZIPBytes)-1; i++ {
		if value[i] != nunGZIPBytes[i] {
			mismatch++
		}
	}
	if mismatch > 0 {
		fmt.Printf("[READ] Slices doesn't match: %+v\n", nunGZIPBytes)
	} else {
		fmt.Println("[READ] Slices match!!")
	}

	fmt.Println()

	// another read example
	start = time.Now()
	n2zr, err := gzip.NewReader(&inlineGZIP2)
	if err != nil {
		fmt.Printf("error: %v", err)
	}
	uncompressedData, err := ioutil.ReadAll(n2zr)
	if err != nil {
		fmt.Printf("error: %v", err)
	}
	defer nzr.Close()
	fmt.Printf("[READ] Elapsed: %+v\n", time.Since(start))
	nunGZIPBytes = uncompressedData
	fmt.Printf("[READ] Uncompressed Size: %+v\n", len(nunGZIPBytes))
	mismatch = 0
	for i := 0; i < len(nunGZIPBytes)-1; i++ {
		if value[i] != nunGZIPBytes[i] {
			mismatch++
		}
	}
	if mismatch > 0 {
		fmt.Printf("[READ] Slices doesn't match: %+v\n", nunGZIPBytes)
	} else {
		fmt.Println("[READ] Slices match!!")
	}
	fmt.Println()
	// unmarshall
	var resultado Embeddings
	json.Unmarshal(nunGZIPBytes, &resultado)
	fmt.Printf("Resultado: %+v\n", resultado)
	fmt.Println()
}
