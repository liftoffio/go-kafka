package proto

import (
	"bytes"
	"io"
	"reflect"
	"testing"
	"time"
)

type Request interface {
	Bytes() ([]byte, error)
	WriteTo(io.Writer) (int64, error)
}

var _ Request = &MetadataReq{}
var _ Request = &ProduceReq{}
var _ Request = &FetchReq{}
var _ Request = &ConsumerMetadataReq{}
var _ Request = &OffsetReq{}
var _ Request = &OffsetCommitReq{}
var _ Request = &OffsetFetchReq{}

func testRequestSerialization(t *testing.T, r Request) {
	var buf bytes.Buffer
	if n, err := r.WriteTo(&buf); err != nil {
		t.Fatalf("could not write request to buffer: %s", err)
	} else if n != int64(buf.Len()) {
		t.Fatalf("writer returned invalid number of bytes written %d != %d", n, buf.Len())
	}
	b, err := r.Bytes()
	if err != nil {
		t.Fatalf("could not convert request to bytes: %s", err)
	}
	if !bytes.Equal(b, buf.Bytes()) {
		t.Fatal("Bytes() and WriteTo() serialized request is of different form")
	}
}

func TestMetadataRequest(t *testing.T) {
	req1 := &MetadataReq{
		CorrelationID: 123,
		ClientID:      "testcli",
		Topics:        nil,
	}
	testRequestSerialization(t, req1)
	b, _ := req1.Bytes()
	expected := []byte{0x0, 0x0, 0x0, 0x15, 0x0, 0x3, 0x0, 0x0, 0x0, 0x0, 0x0, 0x7b, 0x0, 0x7, 0x74, 0x65, 0x73, 0x74, 0x63, 0x6c, 0x69, 0x0, 0x0, 0x0, 0x0}

	if !bytes.Equal(b, expected) {
		t.Fatalf("expected different bytes representation: %v", b)
	}

	req2 := &MetadataReq{
		CorrelationID: 123,
		ClientID:      "testcli",
		Topics:        []string{"foo", "bar"},
	}
	testRequestSerialization(t, req2)
	b, _ = req2.Bytes()
	expected = []byte{0x0, 0x0, 0x0, 0x1f, 0x0, 0x3, 0x0, 0x0, 0x0, 0x0, 0x0, 0x7b, 0x0, 0x7, 0x74, 0x65, 0x73, 0x74, 0x63, 0x6c, 0x69, 0x0, 0x0, 0x0, 0x2, 0x0, 0x3, 0x66, 0x6f, 0x6f, 0x0, 0x3, 0x62, 0x61, 0x72}

	if !bytes.Equal(b, expected) {
		t.Fatalf("expected different bytes representation: %v", b)
	}

	r, _ := ReadMetadataReq(bytes.NewBuffer(expected))
	if !reflect.DeepEqual(r, req2) {
		t.Fatalf("malformed request: %#v", r)
	}
}

func TestMetadataResponse(t *testing.T) {
	msgb := []byte{0x0, 0x0, 0x1, 0xc7, 0x0, 0x0, 0x0, 0x7b, 0x0, 0x0, 0x0, 0x4, 0x0, 0x0, 0xc0, 0x10, 0x0, 0xb, 0x31, 0x37, 0x32, 0x2e, 0x31, 0x37, 0x2e, 0x34, 0x32, 0x2e, 0x31, 0x0, 0x0, 0xc0, 0x10, 0x0, 0x0, 0xc0, 0x12, 0x0, 0xb, 0x31, 0x37, 0x32, 0x2e, 0x31, 0x37, 0x2e, 0x34, 0x32, 0x2e, 0x31, 0x0, 0x0, 0xc0, 0x12, 0x0, 0x0, 0xc0, 0x11, 0x0, 0xb, 0x31, 0x37, 0x32, 0x2e, 0x31, 0x37, 0x2e, 0x34, 0x32, 0x2e, 0x31, 0x0, 0x0, 0xc0, 0x11, 0x0, 0x0, 0xc0, 0x13, 0x0, 0xb, 0x31, 0x37, 0x32, 0x2e, 0x31, 0x37, 0x2e, 0x34, 0x32, 0x2e, 0x31, 0x0, 0x0, 0xc0, 0x13, 0x0, 0x0, 0x0, 0x2, 0x0, 0x0, 0x0, 0x3, 0x66, 0x6f, 0x6f, 0x0, 0x0, 0x0, 0x6, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2, 0x0, 0x0, 0xc0, 0x13, 0x0, 0x0, 0x0, 0x3, 0x0, 0x0, 0xc0, 0x13, 0x0, 0x0, 0xc0, 0x10, 0x0, 0x0, 0xc0, 0x11, 0x0, 0x0, 0x0, 0x3, 0x0, 0x0, 0xc0, 0x13, 0x0, 0x0, 0xc0, 0x10, 0x0, 0x0, 0xc0, 0x11, 0x0, 0x0, 0x0, 0x0, 0x0, 0x5, 0x0, 0x0, 0xc0, 0x12, 0x0, 0x0, 0x0, 0x3, 0x0, 0x0, 0xc0, 0x12, 0x0, 0x0, 0xc0, 0x10, 0x0, 0x0, 0xc0, 0x11, 0x0, 0x0, 0x0, 0x3, 0x0, 0x0, 0xc0, 0x12, 0x0, 0x0, 0xc0, 0x10, 0x0, 0x0, 0xc0, 0x11, 0x0, 0x0, 0x0, 0x0, 0x0, 0x4, 0x0, 0x0, 0xc0, 0x11, 0x0, 0x0, 0x0, 0x3, 0x0, 0x0, 0xc0, 0x11, 0x0, 0x0, 0xc0, 0x13, 0x0, 0x0, 0xc0, 0x10, 0x0, 0x0, 0x0, 0x3, 0x0, 0x0, 0xc0, 0x11, 0x0, 0x0, 0xc0, 0x13, 0x0, 0x0, 0xc0, 0x10, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0xc0, 0x12, 0x0, 0x0, 0x0, 0x3, 0x0, 0x0, 0xc0, 0x12, 0x0, 0x0, 0xc0, 0x13, 0x0, 0x0, 0xc0, 0x10, 0x0, 0x0, 0x0, 0x3, 0x0, 0x0, 0xc0, 0x12, 0x0, 0x0, 0xc0, 0x13, 0x0, 0x0, 0xc0, 0x10, 0x0, 0x0, 0x0, 0x0, 0x0, 0x3, 0x0, 0x0, 0xc0, 0x10, 0x0, 0x0, 0x0, 0x3, 0x0, 0x0, 0xc0, 0x10, 0x0, 0x0, 0xc0, 0x11, 0x0, 0x0, 0xc0, 0x12, 0x0, 0x0, 0x0, 0x3, 0x0, 0x0, 0xc0, 0x10, 0x0, 0x0, 0xc0, 0x11, 0x0, 0x0, 0xc0, 0x12, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xc0, 0x11, 0x0, 0x0, 0x0, 0x3, 0x0, 0x0, 0xc0, 0x11, 0x0, 0x0, 0xc0, 0x12, 0x0, 0x0, 0xc0, 0x13, 0x0, 0x0, 0x0, 0x3, 0x0, 0x0, 0xc0, 0x11, 0x0, 0x0, 0xc0, 0x12, 0x0, 0x0, 0xc0, 0x13, 0x0, 0x0, 0x0, 0x4, 0x74, 0x65, 0x73, 0x74, 0x0, 0x0, 0x0, 0x2, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0xc0, 0x11, 0x0, 0x0, 0x0, 0x3, 0x0, 0x0, 0xc0, 0x11, 0x0, 0x0, 0xc0, 0x12, 0x0, 0x0, 0xc0, 0x13, 0x0, 0x0, 0x0, 0x3, 0x0, 0x0, 0xc0, 0x11, 0x0, 0x0, 0xc0, 0x12, 0x0, 0x0, 0xc0, 0x13, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xc0, 0x10, 0x0, 0x0, 0x0, 0x3, 0x0, 0x0, 0xc0, 0x10, 0x0, 0x0, 0xc0, 0x11, 0x0, 0x0, 0xc0, 0x12, 0x0, 0x0, 0x0, 0x3, 0x0, 0x0, 0xc0, 0x10, 0x0, 0x0, 0xc0, 0x11, 0x0, 0x0, 0xc0, 0x12}
	resp, err := ReadMetadataResp(bytes.NewBuffer(msgb))
	if err != nil {
		t.Fatalf("could not read metadata response: %s", err)
	}
	expected := &MetadataResp{
		CorrelationID: 123,
		Brokers: []MetadataRespBroker{
			MetadataRespBroker{NodeID: 49168, Host: "172.17.42.1", Port: 49168},
			MetadataRespBroker{NodeID: 49170, Host: "172.17.42.1", Port: 49170},
			MetadataRespBroker{NodeID: 49169, Host: "172.17.42.1", Port: 49169},
			MetadataRespBroker{NodeID: 49171, Host: "172.17.42.1", Port: 49171},
		},
		Topics: []MetadataRespTopic{
			MetadataRespTopic{
				Name: "foo",
				Err:  error(nil),
				Partitions: []MetadataRespPartition{
					MetadataRespPartition{Err: error(nil), ID: 2, Leader: 49171, Replicas: []int32{49171, 49168, 49169}, Isrs: []int32{49171, 49168, 49169}},
					MetadataRespPartition{Err: error(nil), ID: 5, Leader: 49170, Replicas: []int32{49170, 49168, 49169}, Isrs: []int32{49170, 49168, 49169}},
					MetadataRespPartition{Err: error(nil), ID: 4, Leader: 49169, Replicas: []int32{49169, 49171, 49168}, Isrs: []int32{49169, 49171, 49168}},
					MetadataRespPartition{Err: error(nil), ID: 1, Leader: 49170, Replicas: []int32{49170, 49171, 49168}, Isrs: []int32{49170, 49171, 49168}},
					MetadataRespPartition{Err: error(nil), ID: 3, Leader: 49168, Replicas: []int32{49168, 49169, 49170}, Isrs: []int32{49168, 49169, 49170}},
					MetadataRespPartition{Err: error(nil), ID: 0, Leader: 49169, Replicas: []int32{49169, 49170, 49171}, Isrs: []int32{49169, 49170, 49171}},
				},
			},
			MetadataRespTopic{
				Name: "test",
				Err:  error(nil),
				Partitions: []MetadataRespPartition{
					MetadataRespPartition{Err: error(nil), ID: 1, Leader: 49169, Replicas: []int32{49169, 49170, 49171}, Isrs: []int32{49169, 49170, 49171}},
					MetadataRespPartition{Err: error(nil), ID: 0, Leader: 49168, Replicas: []int32{49168, 49169, 49170}, Isrs: []int32{49168, 49169, 49170}},
				},
			},
		},
	}

	if !reflect.DeepEqual(resp, expected) {
		t.Fatalf("expected different message: %#v", resp)
	}

	if b, err := resp.Bytes(); err != nil {
		t.Fatalf("cannot serialize response: %s", err)
	} else {
		if !bytes.Equal(b, msgb) {
			t.Fatalf("serialized representation different from expected: %#v", b)
		}
	}
}

func TestProduceRequest(t *testing.T) {
	req := &ProduceReq{
		CorrelationID: 241,
		ClientID:      "test",
		RequiredAcks:  RequiredAcksAll,
		Timeout:       time.Second,
		Topics: []ProduceReqTopic{
			ProduceReqTopic{
				Name: "foo",
				Partitions: []ProduceReqPartition{
					ProduceReqPartition{
						ID: 0,
						Messages: []*Message{
							&Message{
								Offset: 0,
								Crc:    3099221847,
								Key:    []byte("foo"),
								Value:  []byte("bar"),
							},
						},
					},
				},
			},
		},
	}
	testRequestSerialization(t, req)
	b, _ := req.Bytes()
	expected := []byte{0x0, 0x0, 0x0, 0x49, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xf1, 0x0, 0x4, 0x74, 0x65, 0x73, 0x74, 0xff, 0xff, 0x0, 0x0, 0x3, 0xe8, 0x0, 0x0, 0x0, 0x1, 0x0, 0x3, 0x66, 0x6f, 0x6f, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x20, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x14, 0xb8, 0xba, 0x5f, 0x57, 0x0, 0x0, 0x0, 0x0, 0x0, 0x3, 0x66, 0x6f, 0x6f, 0x0, 0x0, 0x0, 0x3, 0x62, 0x61, 0x72}

	if !bytes.Equal(b, expected) {
		t.Fatalf("expected different bytes representation: %#v", b)
	}

	r, _ := ReadProduceReq(bytes.NewBuffer(expected))
	if !reflect.DeepEqual(r, req) {
		t.Fatalf("malformed request: %#v", r)
	}
}

func TestProduceResponse(t *testing.T) {
	msgb1 := []byte{0x0, 0x0, 0x0, 0x22, 0x0, 0x0, 0x0, 0xf1, 0x0, 0x0, 0x0, 0x1, 0x0, 0x6, 0x66, 0x72, 0x75, 0x69, 0x74, 0x73, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x5d, 0x0, 0x3, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
	resp1, err := ReadProduceResp(bytes.NewBuffer(msgb1))
	if err != nil {
		t.Fatalf("could not read metadata response: %s", err)
	}
	expected1 := &ProduceResp{
		CorrelationID: 241,
		Topics: []ProduceRespTopic{
			ProduceRespTopic{
				Name: "fruits",
				Partitions: []ProduceRespPartition{
					ProduceRespPartition{
						ID:     93,
						Err:    ErrUnknownTopicOrPartition,
						Offset: -1,
					},
				},
			},
		},
	}
	if !reflect.DeepEqual(resp1, expected1) {
		t.Fatalf("expected different message: %#v", resp1)
	}

	if b, err := resp1.Bytes(); err != nil {
		t.Fatalf("cannot serialize response: %s", err)
	} else {
		if !bytes.Equal(b, msgb1) {
			t.Fatalf("serialized representation different from expected: %#v", b)
		}
	}

	msgb2 := []byte{0x0, 0x0, 0x0, 0x1f, 0x0, 0x0, 0x0, 0xf1, 0x0, 0x0, 0x0, 0x1, 0x0, 0x3, 0x66, 0x6f, 0x6f, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1}
	resp2, err := ReadProduceResp(bytes.NewBuffer(msgb2))
	if err != nil {
		t.Fatalf("could not read metadata response: %s", err)
	}
	expected2 := &ProduceResp{
		CorrelationID: 241,
		Topics: []ProduceRespTopic{
			ProduceRespTopic{
				Name: "foo",
				Partitions: []ProduceRespPartition{
					ProduceRespPartition{
						ID:     0,
						Err:    error(nil),
						Offset: 1,
					},
				},
			},
		},
	}
	if !reflect.DeepEqual(resp2, expected2) {
		t.Fatalf("expected different message: %#v", resp2)
	}
	if b, err := resp2.Bytes(); err != nil {
		t.Fatalf("cannot serialize response: %s", err)
	} else {
		if !bytes.Equal(b, msgb2) {
			t.Fatalf("serialized representation different from expected: %#v", b)
		}
	}
}

func TestFetchRequest(t *testing.T) {
	req := &FetchReq{
		CorrelationID: 241,
		ClientID:      "test",
		MaxWaitTime:   time.Second * 2,
		MinBytes:      12454,
		Topics: []FetchReqTopic{
			FetchReqTopic{
				Name: "foo",
				Partitions: []FetchReqPartition{
					FetchReqPartition{ID: 421, FetchOffset: 529, MaxBytes: 4921},
					FetchReqPartition{ID: 0, FetchOffset: 11, MaxBytes: 92},
				},
			},
		},
	}
	testRequestSerialization(t, req)
	b, _ := req.Bytes()
	expected := []byte{0x0, 0x0, 0x0, 0x47, 0x0, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0xf1, 0x0, 0x4, 0x74, 0x65, 0x73, 0x74, 0xff, 0xff, 0xff, 0xff, 0x0, 0x0, 0x7, 0xd0, 0x0, 0x0, 0x30, 0xa6, 0x0, 0x0, 0x0, 0x1, 0x0, 0x3, 0x66, 0x6f, 0x6f, 0x0, 0x0, 0x0, 0x2, 0x0, 0x0, 0x1, 0xa5, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2, 0x11, 0x0, 0x0, 0x13, 0x39, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xb, 0x0, 0x0, 0x0, 0x5c}

	if !bytes.Equal(b, expected) {
		t.Fatalf("expected different bytes representation: %#v", b)
	}

	r, _ := ReadFetchReq(bytes.NewBuffer(expected))
	if !reflect.DeepEqual(r, req) {
		t.Fatalf("malformed request: %#v", r)
	}
}

func TestFetchResponse(t *testing.T) {
	msgb := []byte{0x0, 0x0, 0x0, 0x75, 0x0, 0x0, 0x0, 0xf1, 0x0, 0x0, 0x0, 0x1, 0x0, 0x3, 0x66, 0x6f, 0x6f, 0x0, 0x0, 0x0, 0x2, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x4, 0x0, 0x0, 0x0, 0x40, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2, 0x0, 0x0, 0x0, 0x14, 0xb8, 0xba, 0x5f, 0x57, 0x0, 0x0, 0x0, 0x0, 0x0, 0x3, 0x66, 0x6f, 0x6f, 0x0, 0x0, 0x0, 0x3, 0x62, 0x61, 0x72, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x3, 0x0, 0x0, 0x0, 0x14, 0xb8, 0xba, 0x5f, 0x57, 0x0, 0x0, 0x0, 0x0, 0x0, 0x3, 0x66, 0x6f, 0x6f, 0x0, 0x0, 0x0, 0x3, 0x62, 0x61, 0x72, 0x0, 0x0, 0x0, 0x1, 0x0, 0x3, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x0, 0x0, 0x0, 0x0}
	resp, err := ReadFetchResp(bytes.NewBuffer(msgb))
	if err != nil {
		t.Fatalf("could not read metadata response: %s", err)
	}
	expected := &FetchResp{
		CorrelationID: 241,
		Topics: []FetchRespTopic{
			FetchRespTopic{
				Name: "foo",
				Partitions: []FetchRespPartition{
					FetchRespPartition{
						ID:        0,
						Err:       error(nil),
						TipOffset: 4,
						Messages: []*Message{
							&Message{Offset: 2, Crc: 0xb8ba5f57, Key: []byte("foo"), Value: []byte("bar"), Topic: "foo", Partition: 0},
							&Message{Offset: 3, Crc: 0xb8ba5f57, Key: []byte("foo"), Value: []byte("bar"), Topic: "foo", Partition: 0},
						},
					},
					FetchRespPartition{
						ID:        1,
						Err:       ErrUnknownTopicOrPartition,
						TipOffset: -1,
						Messages:  []*Message{},
					},
				},
			},
		},
	}
	if !reflect.DeepEqual(resp, expected) {
		t.Fatalf("expected different message: %#v", resp)
	}
	b, err := resp.Bytes()
	if err != nil {
		t.Fatalf("cannot serialize response: %s", err)
	}
	if !bytes.Equal(b, msgb) {
		t.Fatalf("serialized representation different from expected: %#v", b)
	}
}

func TestFetchResponse2(t *testing.T) {
	msgb := []byte{0x0, 0x0, 0x0, 0x48, 0x0, 0x0, 0x0, 0xf1, 0x0, 0x0, 0x0, 0x1, 0x0, 0x4, 0x74, 0x65, 0x73, 0x74, 0x0, 0x0, 0x0, 0x3, 0x0, 0x0, 0x0, 0x0, 0x0, 0x3, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x3, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x8, 0x0, 0x3, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x0, 0x0, 0x0, 0x0}
	resp, err := ReadFetchResp(bytes.NewBuffer(msgb))
	if err != nil {
		t.Fatalf("could not read fetch response: %s", err)
	}
	expected := &FetchResp{
		CorrelationID: 241,
		Topics: []FetchRespTopic{
			FetchRespTopic{
				Name: "test",
				Partitions: []FetchRespPartition{
					FetchRespPartition{
						ID:        0,
						Err:       ErrUnknownTopicOrPartition,
						TipOffset: -1,
						Messages:  []*Message{},
					},
					FetchRespPartition{
						ID:        1,
						Err:       ErrUnknownTopicOrPartition,
						TipOffset: -1,
						Messages:  []*Message{},
					},
					FetchRespPartition{
						ID:        8,
						Err:       ErrUnknownTopicOrPartition,
						TipOffset: -1,
						Messages:  []*Message{},
					},
				},
			},
		},
	}
	if !reflect.DeepEqual(resp, expected) {
		t.Fatalf("expected different structure: %#v", resp)
	}

	b, err := resp.Bytes()
	if err != nil {
		t.Fatalf("cannot serialize response: %s", err)
	}
	if !bytes.Equal(b, msgb) {
		t.Fatalf("serialization failure: %#v", b)
	}
}

func TestSerializeEmptyMessageSet(t *testing.T) {
	var buf bytes.Buffer
	messages := []*Message{}
	if err := writeMessageSet(&buf, messages); err != nil {
		t.Fatalf("cannot serialize messages: %s", err)
	}
	expected := []byte{0, 0, 0, 0} // zero size, int32 type
	if !bytes.Equal(buf.Bytes(), expected) {
		t.Fatalf("expected different byte representation: %#v", buf.Bytes())
	}
}

func TestReadIncomleteMessage(t *testing.T) {
	var buf bytes.Buffer
	err := writeMessageSet(&buf, []*Message{
		&Message{Value: []byte("111111111111111")},
		&Message{Value: []byte("222222222222222")},
		&Message{Value: []byte("333333333333333")},
	})
	if err != nil {
		t.Fatalf("cannot serialize messages: %s", err)
	}

	b := buf.Bytes()
	// cut off the last bytes as kafka can do
	b = b[:len(b)-4]
	messages, err := readMessageSet(bytes.NewBuffer(b))
	if err != nil {
		t.Fatalf("cannot deserialize messages: %s", err)
	}
	if len(messages) != 2 {
		t.Fatalf("expected 2 messages, got %d", len(messages))
	}
	if messages[0].Value[0] != '1' || messages[1].Value[0] != '2' {
		t.Fatal("expected different messages content")
	}
}

func BenchmarkProduceRequestMarshal(b *testing.B) {
	messages := make([]*Message, 1000)
	for i := range messages {
		messages[i] = &Message{
			Offset: int64(i),
			Crc:    uint32(i),
			Key:    nil,
			Value:  []byte(`Lorem ipsum dolor sit amet, consectetur adipiscing elit. Donec a diam lectus. Sed sit amet ipsum mauris. Maecenas congue ligula ac quam viverra nec consectetur ante hendrerit. Donec et mollis dolor. Praesent et diam eget libero egestas mattis sit amet vitae augue. Nam tincidunt congue enim, ut porta lorem lacinia consectetur.`),
		}

	}
	req := &ProduceReq{
		CorrelationID: 241,
		ClientID:      "test",
		RequiredAcks:  RequiredAcksAll,
		Timeout:       time.Second,
		Topics: []ProduceReqTopic{
			ProduceReqTopic{
				Name: "foo",
				Partitions: []ProduceReqPartition{
					ProduceReqPartition{
						ID:       0,
						Messages: messages,
					},
				},
			},
		},
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if _, err := req.Bytes(); err != nil {
			b.Fatalf("could not serialize messages: %s", err)
		}
	}
}

func BenchmarkProduceResponseUnmarshal(b *testing.B) {
	resp := &ProduceResp{
		CorrelationID: 241,
		Topics: []ProduceRespTopic{
			ProduceRespTopic{
				Name: "foo",
				Partitions: []ProduceRespPartition{
					ProduceRespPartition{
						ID:     0,
						Err:    error(nil),
						Offset: 1,
					},
				},
			},
		},
	}
	raw, err := resp.Bytes()
	if err != nil {
		b.Fatalf("cannot serialize response: %s", err)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if _, err := ReadProduceResp(bytes.NewBuffer(raw)); err != nil {
			b.Fatalf("could not deserialize messages: %s", err)
		}
	}
}

func BenchmarkFetchRequestMarshal(b *testing.B) {
	req := &FetchReq{
		CorrelationID: 241,
		ClientID:      "test",
		MaxWaitTime:   time.Second * 2,
		MinBytes:      12454,
		Topics: []FetchReqTopic{
			FetchReqTopic{
				Name: "foo",
				Partitions: []FetchReqPartition{
					FetchReqPartition{ID: 421, FetchOffset: 529, MaxBytes: 4921},
					FetchReqPartition{ID: 0, FetchOffset: 11, MaxBytes: 92},
				},
			},
		},
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if _, err := req.Bytes(); err != nil {
			b.Fatalf("could not serialize messages: %s", err)
		}
	}
}

func BenchmarkFetchResponseUnmarshal(b *testing.B) {
	messages := make([]*Message, 100)
	for i := range messages {
		messages[i] = &Message{
			Offset: int64(i),
			Key:    nil,
			Value:  []byte(`Lorem ipsum dolor sit amet, consectetur adipiscing elit. Donec a diam lectus. Sed sit amet ipsum mauris. Maecenas congue ligula ac quam viverra nec consectetur ante hendrerit. Donec et mollis dolor. Praesent et diam eget libero egestas mattis sit amet vitae augue. Nam tincidunt congue enim, ut porta lorem lacinia consectetur.`),
		}

	}
	resp := &FetchResp{
		CorrelationID: 241,
		Topics: []FetchRespTopic{
			FetchRespTopic{
				Name: "foo",
				Partitions: []FetchRespPartition{
					FetchRespPartition{
						ID:        0,
						TipOffset: 444,
						Messages:  messages,
					},
					FetchRespPartition{
						ID:        123,
						Err:       ErrBrokerNotAvailable,
						TipOffset: -1,
						Messages:  []*Message{},
					},
				},
			},
		},
	}
	raw, err := resp.Bytes()
	if err != nil {
		b.Fatalf("cannot serialize response: %s", err)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if _, err := ReadFetchResp(bytes.NewBuffer(raw)); err != nil {
			b.Fatalf("could not deserialize messages: %s", err)
		}
	}
}

// vim has problem with coloring byte arrays in this file
// vim: set syntax=off:
