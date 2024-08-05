package cache

import (
	"bytes"
	"compress/gzip"

	"github.com/vmihailenco/msgpack"
)

// Serialize and its inverse, Deserialize, provide functions suitable for storing most objects in
// caches.
func Serialize(v interface{}) ([]byte, error) {
	buf := &bytes.Buffer{}
	w := gzip.NewWriter(buf)
	if err := msgpack.NewEncoder(w).Encode(v); err != nil {
		return nil, err
	}
	w.Close()
	return buf.Bytes(), nil
}

// Deserialize and its inverse, Serialize, provide functions suitable for storing most objects in
// caches.
func Deserialize(b []byte, dest interface{}) error {
	r, err := gzip.NewReader(bytes.NewReader(b))
	if err != nil {
		return err
	}
	defer r.Close()
	return msgpack.NewDecoder(r).Decode(dest)
}
