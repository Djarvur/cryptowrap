package cryptowrap_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/Djarvur/cryptowrap"
)

type TestData struct {
	Field1 string
	Field2 string
	Field3 string
}

func TestWrapper128(t *testing.T) {
	testWrapper(t, 16)
}

func TestWrapper256(t *testing.T) {
	testWrapper(t, 32)
}

func testWrapper(t *testing.T, keyLen int) {
	keys := [][]byte{
		randBytes(keyLen),
		randBytes(keyLen),
	}

	orig := TestData{
		Field1: "Field1",
		Field2: "Field2",
		Field3: "Field3",
	}

	src := cryptowrap.Wrapper{
		Keys:    keys[1:],
		Payload: &orig,
	}

	data, err := json.Marshal(&src)
	if err != nil {
		t.Fatal(err)
	}

	dst := cryptowrap.Wrapper{
		Keys:    keys[1:],
		Payload: &TestData{},
	}

	err = json.Unmarshal(data, &dst)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(&orig, dst.Payload) {
		t.Fatal("decrypted is not equal to original")
	}

	dst = cryptowrap.Wrapper{
		Keys:    keys,
		Payload: &TestData{},
	}

	err = json.Unmarshal(data, &dst)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(&orig, dst.Payload) {
		t.Fatal("decrypted is not equal to original")
	}

	dst = cryptowrap.Wrapper{
		Keys:    [][]byte{keys[1], keys[0]},
		Payload: &TestData{},
	}

	err = json.Unmarshal(data, &dst)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(&orig, dst.Payload) {
		t.Fatal("decrypted is not equal to original")
	}

	err = json.Unmarshal(
		data,
		&cryptowrap.Wrapper{
			Keys: keys[:1],
		},
	)
	if err == nil {
		t.Fatal("decrypted undecryptable")
	}

	err = json.Unmarshal(
		data,
		&cryptowrap.Wrapper{
			Payload: &TestData{},
		},
	)
	if err == nil {
		t.Fatal("decrypted undecryptable")
	}

	badKey := [][]byte{randBytes(15)}

	err = json.Unmarshal(data, &cryptowrap.Wrapper{Keys: badKey})
	if err == nil {
		t.Fatal("decrypted undecryptable")
	}

	_, err = json.Marshal(&cryptowrap.Wrapper{Keys: badKey, Payload: &orig})
	if err == nil {
		t.Fatal("encrypted unencryptable")
	}

	_, err = json.Marshal(&cryptowrap.Wrapper{Payload: &orig})
	if err == nil {
		t.Fatal("encrypted unencryptable")
	}
}

func randBytes(keyLen int) []byte {
	data, err := cryptowrap.RandBytes(keyLen)
	if err != nil {
		panic(err)
	}
	return data
}
