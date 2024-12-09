package redis

import (
	"bytes"
	"compress/gzip"
	"encoding/gob"
)

// Marshal 将对象序列化为字节数组
func (s *Store) Marshal(obj interface{}) ([]byte, error) {
	var buf bytes.Buffer

	var writer interface {
		Write(p []byte) (n int, err error)
	} = &buf

	// 如果启用压缩，使用 gzip
	var gzWriter *gzip.Writer
	if s.compress {
		gzWriter = gzip.NewWriter(&buf)
		writer = gzWriter
	}

	enc := gob.NewEncoder(writer)
	if err := enc.Encode(obj); err != nil {
		return nil, err
	}

	if s.compress && gzWriter != nil {
		if err := gzWriter.Close(); err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

// Unmarshal 从字节数组反序列化为对象
func (s *Store) Unmarshal(data []byte, obj interface{}) error {
	reader := bytes.NewReader(data)

	var decoder *gob.Decoder
	if s.compress {
		gzReader, err := gzip.NewReader(reader)
		if err != nil {
			return err
		}
		defer func(gzReader *gzip.Reader) {
			_ = gzReader.Close()
		}(gzReader)
		decoder = gob.NewDecoder(gzReader)
	} else {
		decoder = gob.NewDecoder(reader)
	}

	return decoder.Decode(obj)
}
