package runnlib

import (
	"bytes"
	"testing"
)

func TestArchiDir(t *testing.T) {

	reader, err := ArchieveDir("..", "test", []byte("Hello"))
	if err != nil {
		t.Fatal(err)
	}

	//file, _ := os.Create("test.zip")
	//io.Copy(file, reader)

	reader2 := reader.(*bytes.Buffer)
	err = UnarchiveToDir("test-zip", bytes.NewReader(reader2.Bytes()), int64(reader2.Len()), []byte("Hello"))
	if err != nil {
		t.Fatal(err)
	}
}
