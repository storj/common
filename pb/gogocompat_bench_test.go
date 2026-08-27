// Copyright (C) 2026 Storj Labs, Inc.
// See LICENSE for copying information.

// The benchmarks here measure the picobuf encoding against the gogo/protobuf
// one in internal/gogopb, on the messages the upload and download paths spend
// their time in. BenchmarkMarshalReuse is the shape drpc uses on the wire:
// encoding into a buffer that is kept between messages.
package pb_test

import (
	"reflect"
	"testing"
	"time"

	gogo "github.com/gogo/protobuf/proto"

	"storj.io/common/pb"
	"storj.io/common/pb/internal/gogopb"
	"storj.io/common/storj"
	"storj.io/picobuf"
)

func benchBytes(n int, seed byte) []byte {
	b := make([]byte, n)
	for i := range b {
		b[i] = byte(i) ^ seed
	}
	return b
}

func benchNodeID(seed byte) (id storj.NodeID) {
	copy(id[:], benchBytes(len(id), seed))
	return id
}

func benchPieceID(seed byte) (id storj.PieceID) {
	copy(id[:], benchBytes(len(id), seed))
	return id
}

func benchSerial(seed byte) (s storj.SerialNumber) {
	copy(s[:], benchBytes(len(s), seed))
	return s
}

func benchPubKey() storj.PiecePublicKey {
	pub, _, err := storj.NewPieceKey()
	if err != nil {
		panic(err)
	}
	return pub
}

var (
	benchNow    = time.Date(2026, 9, 8, 12, 0, 0, 0, time.UTC)
	benchPub    = benchPubKey()
	benchChunk  = benchBytes(64*1024, 0x5a)
	benchSig    = benchBytes(64, 0x11)
	benchAddr   = "storagenode-1234.example.com:28967"
	benchMetaKI = benchBytes(16, 0x22)
	benchMeta   = benchBytes(64, 0x33)
)

func picoOrderLimit() *pb.OrderLimit {
	return &pb.OrderLimit{
		SerialNumber:           benchSerial(1),
		SatelliteId:            benchNodeID(2),
		UplinkPublicKey:        benchPub,
		StorageNodeId:          benchNodeID(3),
		PieceId:                benchPieceID(4),
		Limit:                  1 << 22,
		Action:                 pb.PieceAction_PUT,
		PieceExpiration:        benchNow.Add(720 * time.Hour),
		OrderExpiration:        benchNow.Add(48 * time.Hour),
		OrderCreation:          benchNow,
		EncryptedMetadataKeyId: benchMetaKI,
		EncryptedMetadata:      benchMeta,
		SatelliteSignature:     benchSig,
	}
}

func gogoOrderLimit() *gogopb.OrderLimit {
	return &gogopb.OrderLimit{
		SerialNumber:           benchSerial(1),
		SatelliteId:            benchNodeID(2),
		UplinkPublicKey:        benchPub,
		StorageNodeId:          benchNodeID(3),
		PieceId:                benchPieceID(4),
		Limit:                  1 << 22,
		Action:                 gogopb.PieceAction_PUT,
		PieceExpiration:        benchNow.Add(720 * time.Hour),
		OrderExpiration:        benchNow.Add(48 * time.Hour),
		OrderCreation:          benchNow,
		EncryptedMetadataKeyId: benchMetaKI,
		EncryptedMetadata:      benchMeta,
		SatelliteSignature:     benchSig,
	}
}

func picoOrder() *pb.Order {
	return &pb.Order{SerialNumber: benchSerial(1), Amount: 1 << 20, UplinkSignature: benchSig}
}

func gogoOrder() *gogopb.Order {
	return &gogopb.Order{SerialNumber: benchSerial(1), Amount: 1 << 20, UplinkSignature: benchSig}
}

func picoPieceHash() *pb.PieceHash {
	return &pb.PieceHash{
		PieceId:       benchPieceID(4),
		Hash:          benchBytes(32, 0x44),
		PieceSize:     1 << 20,
		Timestamp:     benchNow,
		Signature:     benchSig,
		HashAlgorithm: pb.PieceHashAlgorithm_BLAKE3,
	}
}

func gogoPieceHash() *gogopb.PieceHash {
	return &gogopb.PieceHash{
		PieceId:       benchPieceID(4),
		Hash:          benchBytes(32, 0x44),
		PieceSize:     1 << 20,
		Timestamp:     benchNow,
		Signature:     benchSig,
		HashAlgorithm: gogopb.PieceHashAlgorithm_BLAKE3,
	}
}

func picoAddressedLimits(n int) []*pb.AddressedOrderLimit {
	limits := make([]*pb.AddressedOrderLimit, n)
	for i := range limits {
		limit := picoOrderLimit()
		limit.StorageNodeId = benchNodeID(byte(i))
		limits[i] = &pb.AddressedOrderLimit{
			Limit:              limit,
			StorageNodeAddress: &pb.NodeAddress{Address: benchAddr, DebounceLimit: 3, Features: 1},
		}
	}
	return limits
}

func gogoAddressedLimits(n int) []*gogopb.AddressedOrderLimit {
	limits := make([]*gogopb.AddressedOrderLimit, n)
	for i := range limits {
		limit := gogoOrderLimit()
		limit.StorageNodeId = benchNodeID(byte(i))
		limits[i] = &gogopb.AddressedOrderLimit{
			Limit:              limit,
			StorageNodeAddress: &gogopb.NodeAddress{Address: benchAddr, DebounceLimit: 3, Features: 1},
		}
	}
	return limits
}

func picoDownloadSegmentResponse() *pb.DownloadSegmentResponse {
	return &pb.DownloadSegmentResponse{
		SegmentId:       storj.SegmentID(benchBytes(64, 0x55)),
		AddressedLimits: picoAddressedLimits(39),
		PrivateKey:      storj.PiecePrivateKey{},
		SegmentSize:     64 * 1024 * 1024,
		EncryptedKey:    benchBytes(48, 0x66),
		RedundancyScheme: &pb.RedundancyScheme{
			Type: pb.RedundancyScheme_RS, MinReq: 29, Total: 110,
			RepairThreshold: 35, SuccessThreshold: 80, ErasureShareSize: 256,
		},
	}
}

func gogoDownloadSegmentResponse() *gogopb.DownloadSegmentResponse {
	return &gogopb.DownloadSegmentResponse{
		SegmentId:       storj.SegmentID(benchBytes(64, 0x55)),
		AddressedLimits: gogoAddressedLimits(39),
		PrivateKey:      storj.PiecePrivateKey{},
		SegmentSize:     64 * 1024 * 1024,
		EncryptedKey:    benchBytes(48, 0x66),
		RedundancyScheme: &gogopb.RedundancyScheme{
			Type: gogopb.RedundancyScheme_RS, MinReq: 29, Total: 110,
			RepairThreshold: 35, SuccessThreshold: 80, ErasureShareSize: 256,
		},
	}
}

type benchCase struct {
	name string
	pico func() pb.Message
	gogo func() gogo.Message
}

var benchCases = []benchCase{
	{"OrderLimit", func() pb.Message { return picoOrderLimit() }, func() gogo.Message { return gogoOrderLimit() }},
	{"Order", func() pb.Message { return picoOrder() }, func() gogo.Message { return gogoOrder() }},
	{"PieceHash", func() pb.Message { return picoPieceHash() }, func() gogo.Message { return gogoPieceHash() }},
	{"UploadRequestFirst", func() pb.Message {
		return &pb.PieceUploadRequest{Limit: picoOrderLimit(), HashAlgorithm: pb.PieceHashAlgorithm_BLAKE3,
			Chunk: &pb.PieceUploadRequest_Chunk{Offset: 0, Data: benchChunk}}
	}, func() gogo.Message {
		return &gogopb.PieceUploadRequest{Limit: gogoOrderLimit(), HashAlgorithm: gogopb.PieceHashAlgorithm_BLAKE3,
			Chunk: &gogopb.PieceUploadRequest_Chunk{Offset: 0, Data: benchChunk}}
	}},
	{"UploadRequestChunk", func() pb.Message {
		return &pb.PieceUploadRequest{Order: picoOrder(), Chunk: &pb.PieceUploadRequest_Chunk{Offset: 1 << 20, Data: benchChunk}}
	}, func() gogo.Message {
		return &gogopb.PieceUploadRequest{Order: gogoOrder(), Chunk: &gogopb.PieceUploadRequest_Chunk{Offset: 1 << 20, Data: benchChunk}}
	}},
	{"UploadRequestDone", func() pb.Message {
		return &pb.PieceUploadRequest{Done: picoPieceHash()}
	}, func() gogo.Message {
		return &gogopb.PieceUploadRequest{Done: gogoPieceHash()}
	}},
	{"UploadResponse", func() pb.Message {
		return &pb.PieceUploadResponse{Done: picoPieceHash()}
	}, func() gogo.Message {
		return &gogopb.PieceUploadResponse{Done: gogoPieceHash()}
	}},
	{"DownloadRequestFirst", func() pb.Message {
		return &pb.PieceDownloadRequest{Limit: picoOrderLimit(),
			Chunk: &pb.PieceDownloadRequest_Chunk{Offset: 0, ChunkSize: 1 << 20}, MaximumChunkSize: 64 * 1024}
	}, func() gogo.Message {
		return &gogopb.PieceDownloadRequest{Limit: gogoOrderLimit(),
			Chunk: &gogopb.PieceDownloadRequest_Chunk{Offset: 0, ChunkSize: 1 << 20}, MaximumChunkSize: 64 * 1024}
	}},
	{"DownloadRequestOrder", func() pb.Message {
		return &pb.PieceDownloadRequest{Order: picoOrder()}
	}, func() gogo.Message {
		return &gogopb.PieceDownloadRequest{Order: gogoOrder()}
	}},
	{"DownloadResponseChunk", func() pb.Message {
		return &pb.PieceDownloadResponse{Chunk: &pb.PieceDownloadResponse_Chunk{Offset: 1 << 20, Data: benchChunk}}
	}, func() gogo.Message {
		return &gogopb.PieceDownloadResponse{Chunk: &gogopb.PieceDownloadResponse_Chunk{Offset: 1 << 20, Data: benchChunk}}
	}},
	{"DownloadSegmentResponse", func() pb.Message { return picoDownloadSegmentResponse() },
		func() gogo.Message { return gogoDownloadSegmentResponse() }},
}

func BenchmarkMarshal(b *testing.B) {
	for _, c := range benchCases {
		picoMsg, gogoMsg := c.pico(), c.gogo()
		data, err := pb.Marshal(picoMsg)
		if err != nil {
			b.Fatal(err)
		}
		gogoData, err := gogo.Marshal(gogoMsg)
		if err != nil {
			b.Fatal(err)
		}
		if len(data) != len(gogoData) {
			b.Fatalf("%s: encoding size differs pico=%d gogo=%d", c.name, len(data), len(gogoData))
		}

		b.Run(c.name+"/impl=pico", func(b *testing.B) {
			b.SetBytes(int64(len(data)))
			for b.Loop() {
				if _, err := pb.Marshal(picoMsg); err != nil {
					b.Fatal(err)
				}
			}
		})
		b.Run(c.name+"/impl=gogo", func(b *testing.B) {
			b.SetBytes(int64(len(data)))
			for b.Loop() {
				if _, err := gogo.Marshal(gogoMsg); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkUnmarshal(b *testing.B) {
	for _, c := range benchCases {
		data, err := pb.Marshal(c.pico())
		if err != nil {
			b.Fatal(err)
		}

		picoType := reflect.TypeOf(c.pico()).Elem()
		gogoType := reflect.TypeOf(c.gogo()).Elem()

		b.Run(c.name+"/impl=pico", func(b *testing.B) {
			b.SetBytes(int64(len(data)))
			for b.Loop() {
				msg := reflect.New(picoType).Interface().(pb.Message)
				if err := pb.Unmarshal(data, msg); err != nil {
					b.Fatal(err)
				}
			}
		})
		b.Run(c.name+"/impl=pico-reuse", func(b *testing.B) {
			b.SetBytes(int64(len(data)))
			var dec picobuf.Decoder
			for b.Loop() {
				msg := reflect.New(picoType).Interface().(pb.Message)
				if err := dec.Unmarshal(data, msg, picobuf.UnmarshalOptions{}); err != nil {
					b.Fatal(err)
				}
			}
		})
		b.Run(c.name+"/impl=gogo", func(b *testing.B) {
			b.SetBytes(int64(len(data)))
			for b.Loop() {
				msg := reflect.New(gogoType).Interface().(gogo.Message)
				if err := gogo.Unmarshal(data, msg); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

// BenchmarkMarshalReuse encodes into a buffer that is kept between messages.
func BenchmarkMarshalReuse(b *testing.B) {
	for _, c := range benchCases {
		picoMsg, gogoMsg := c.pico(), c.gogo()
		data, err := pb.Marshal(picoMsg)
		if err != nil {
			b.Fatal(err)
		}

		b.Run(c.name+"/impl=pico", func(b *testing.B) {
			b.SetBytes(int64(len(data)))
			buf := []byte{}
			for b.Loop() {
				buf, err = picobuf.MarshalAppend(buf[:0], picoMsg)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
		b.Run(c.name+"/impl=gogo", func(b *testing.B) {
			b.SetBytes(int64(len(data)))
			buf := gogo.NewBuffer(nil)
			for b.Loop() {
				buf.SetBuf(buf.Bytes()[:0])
				if err := buf.Marshal(gogoMsg); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
