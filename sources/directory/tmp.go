package directory

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"os"

	collection "github.com/DblK/tinshop/gamescollection"
	"github.com/DblK/tinshop/repository"
)

type pfs0_header struct {
	// magic uint32
	Magic          [4]byte
	File_cnt       uint32
	Str_table_size uint32
	Reserved       uint32
}

func nspCheck(file repository.FileDesc) {
	fmt.Println("GameID:", file.GameID)
	key := collection.GetKey(file.GameID)
	fmt.Println("Key:", key)
	fmt.Println()

	f, err := os.Open(file.Path)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()

	// Working
	// nspHeader := &pfs0_header{}
	// binary.Read(f, binary.LittleEndian, &nspHeader.magic)
	// binary.Read(f, binary.LittleEndian, &nspHeader.file_cnt)
	// binary.Read(f, binary.LittleEndian, &nspHeader.str_table_size)
	// binary.Read(f, binary.LittleEndian, &nspHeader.reserved)

	// head := make([]byte, 4)
	// f.Read(head)

	nspHeader := pfs0_header{}
	// data := make([]byte, 4)
	data := make([]byte, 4*4)
	_, _ = f.Read(data)
	fmt.Println(len(data))
	buffer := bytes.NewBuffer(data)
	fmt.Println(buffer)
	fmt.Printf("%s", hex.Dump(data))
	err = binary.Read(buffer, binary.LittleEndian, &nspHeader)

	// fmt.Println(string(head))
	fmt.Println(nspHeader)
}
