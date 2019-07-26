package framing

import (
	"bytes"
	"encoding/binary"
	"github.com/netifi/netifi-go/tags"
)

type Shard []byte

func (g Shard) Group() string {
	offset := HeaderOffset()

	groupLength := int(binary.BigEndian.Uint32(g[offset : offset+4]))
	offset += 4

	return string(g[offset : offset+groupLength])
}

func (g Shard) Metadata() []byte {
	offset := HeaderOffset()

	groupLength := int(binary.BigEndian.Uint32(g[offset : offset+4]))
	offset += 4 + groupLength

	metadataLength := int(binary.BigEndian.Uint32(g[offset : offset+4]))
	offset += 4

	return g[offset : offset+metadataLength]
}

func (g Shard) ShardKey() []byte {
	offset := HeaderOffset()

	groupLength := int(binary.BigEndian.Uint32(g[offset : offset+4]))
	offset += 4 + groupLength

	metadataLength := int(binary.BigEndian.Uint32(g[offset : offset+4]))
	offset += 4 + metadataLength

	shardKeyLength := int(binary.BigEndian.Uint32(g[offset : offset+4]))
	offset += 4

	return g[offset : offset+shardKeyLength]
}

func (g Shard) Tags() (t tags.Tags) {
	offset := HeaderOffset()

	groupLength := int(binary.BigEndian.Uint32(g[offset : offset+4]))
	offset += 4 + groupLength

	metadataLength := int(binary.BigEndian.Uint32(g[offset : offset+4]))
	offset += 4 + metadataLength

	shardKeyLength := int(binary.BigEndian.Uint32(g[offset : offset+4]))
	offset += 4 + shardKeyLength

	t = tags.Tags{}
	for {
		keyLength := int(binary.BigEndian.Uint32(g[offset : offset+4]))
		offset += 4

		key := string(g[offset+keyLength])
		offset += keyLength

		valueLength := int(binary.BigEndian.Uint32(g[offset : offset+4]))
		offset += 4

		value := string(g[offset+valueLength])
		offset += valueLength

		t = t.And(tags.New(key, value))

		if offset >= len(g) {
			break
		}
	}

	return
}

func EncodeShard(group string, metadata []byte, shardKey []byte, tags tags.Tags) (g Broadcast, err error) {
	w := &bytes.Buffer{}
	f, err := encodeFrameHeader(FrameTypeBroadcast)
	if err != nil {
		return
	}

	err = binary.Write(w, binary.BigEndian, f)
	if err != nil {
		return
	}

	err = binary.Write(w, binary.BigEndian, uint32(len(group)))
	if err != nil {
		return
	}

	err = binary.Write(w, binary.BigEndian, group)
	if err != nil {
		return
	}

	err = binary.Write(w, binary.BigEndian, uint32(len(metadata)))
	if err != nil {
		return
	}

	err = binary.Write(w, binary.BigEndian, metadata)
	if err != nil {
		return
	}

	err = binary.Write(w, binary.BigEndian, uint32(len(shardKey)))
	if err != nil {
		return
	}

	err = binary.Write(w, binary.BigEndian, shardKey)
	if err != nil {
		return
	}

	for _, tag := range tags {
		key := tag.Key()
		value := tag.Value()

		err = binary.Write(w, binary.BigEndian, uint32(len(key)))
		if err != nil {
			return
		}

		err = binary.Write(w, binary.BigEndian, key)
		if err != nil {
			return
		}

		err = binary.Write(w, binary.BigEndian, uint32(len(value)))
		if err != nil {
			return
		}

		err = binary.Write(w, binary.BigEndian, value)
		if err != nil {
			return
		}
	}
	g = w.Bytes()
	return
}