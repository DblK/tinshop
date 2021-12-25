package directory

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"unsafe"

	collection "github.com/DblK/tinshop/gamescollection"
	"github.com/DblK/tinshop/repository"
)

type nsp_file struct {
	Header    pfs0_header
	FileEntry []pfs0_file_entry
	FileName  []string
}

type pfs0_header struct {
	// magic uint32
	Magic          [4]byte
	File_cnt       uint32
	Str_table_size uint32
	Reserved       uint32
}

type pfs0_file_entry struct {
	File_offset     uint64
	File_size       uint64
	Filename_offset uint32
	Reserved        uint32
}

func nspCheck(file repository.FileDesc) {
	fmt.Println("GameID:", file.GameID)
	key := collection.GetKey(file.GameID)
	fmt.Println("Key:", key)
	fmt.Println()

	err := openNSP(file.Path)
	if err != nil {
		fmt.Println("Error while opening NSP", err)
	}
}

func openNSP(file string) error {
	f, err := os.Open(file)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()

	newNSP := &nsp_file{}

	// Read Header
	nspHeader := pfs0_header{}
	data := make([]byte, unsafe.Sizeof(nspHeader))
	_, _ = f.Read(data)
	buffer := bytes.NewBuffer(data)
	_ = binary.Read(buffer, binary.LittleEndian, &nspHeader)
	newNSP.Header = nspHeader

	if string(newNSP.Header.Magic[:]) != "PFS0" {
		return errors.New("header Magic is not present")
	}
	// fmt.Println(nspHeader)

	// Read file entry
	for i := 0; i < int(nspHeader.File_cnt); i++ {
		nspEntry := pfs0_file_entry{}
		data := make([]byte, unsafe.Sizeof(nspEntry))
		_, _ = f.Read(data)
		buffer := bytes.NewBuffer(data)
		_ = binary.Read(buffer, binary.LittleEndian, &nspEntry)
		newNSP.FileEntry = append(newNSP.FileEntry, nspEntry)
		// fmt.Println(nspEntry)
	}
	// fmt.Println(len(newNSP.FileEntry))
	// fmt.Println(newNSP.FileEntry)

	// Read nspStrTable + Display file_name
	nspStrTable := make([]byte, nspHeader.Str_table_size)
	_, _ = f.Read(nspStrTable)
	// fmt.Println(nspStrTable)
	// fmt.Println(nspHeader.Str_table_size)

	var tikOffset uint64
	var tikSize uint64

	for i := 0; i < int(nspHeader.File_cnt); i++ {
		start := newNSP.FileEntry[i].Filename_offset
		if i != int(nspHeader.File_cnt)-1 {
			end := newNSP.FileEntry[i+1].Filename_offset - 1
			// fmt.Println(string(nspStrTable[start:end]))
			newNSP.FileName = append(newNSP.FileName, string(nspStrTable[start:end]))
		} else {
			// fmt.Println(string(nspStrTable[start:]))
			newNSP.FileName = append(newNSP.FileName, string(nspStrTable[start:]))
		}

		// Compute Ticket information
		if newNSP.FileName[i][len(newNSP.FileName[i])-4:] == ".tik" {
			fmt.Println("Found ticket!")
			tikOffset = (uint64(unsafe.Sizeof(nspHeader)) + (uint64(unsafe.Sizeof(newNSP.FileEntry[i])) * uint64(len(newNSP.FileEntry))) + uint64(len(nspStrTable)) + newNSP.FileEntry[i].File_offset)
			tikSize = newNSP.FileEntry[i].File_size
		}
	}
	// Display Ticket
	fmt.Println(newNSP.FileName)
	fmt.Println(tikOffset, tikSize)

	return nil
}
