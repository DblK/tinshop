package directory

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"strings"
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

type rsa2048_sha256_ticket struct {
	Sig_type       uint32
	Signature      [0x100]uint8
	Padding        [0x3C]uint8
	Sig_issuer     [0x40]byte
	Titlekey_block [0x100]uint8
	Unk1           uint8
	Titlekey_type  uint8
	Unk2           [0x03]uint8
	Master_key_rev uint8
	Unk3           [0x0A]uint8
	Ticket_id      uint64
	Device_id      uint64
	Rights_id      [0x10]uint8
	Account_id     uint32
	Unk4           [0x0C]uint8
}

func nspCheck(file repository.FileDesc) {
	fmt.Println("GameID:", file.GameID)
	key := collection.GetKey(file.GameID)
	fmt.Println("Key:", key)
	fmt.Println()

	err := openNSP(file.Path, key)
	if err != nil {
		fmt.Println("Error while opening NSP", err)
	}
}

func openNSP(file string, titleDBKey string) error {
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

	ticket := &rsa2048_sha256_ticket{}
	fmt.Println(unsafe.Sizeof(ticket))
	f.Seek(int64(tikOffset), 0)

	data = make([]byte, tikSize)
	_, _ = f.Read(data)
	buffer = bytes.NewBuffer(data)
	err = binary.Read(buffer, binary.LittleEndian, ticket)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(ticket.Sig_issuer[:]))
	fmt.Println("Titlekey_block", ticket.Titlekey_block)
	// fmt.Println("Titlekey_block", strconv.Itoa(int(ticket.Titlekey_block[0])))
	var titleKey []byte
	for i := 0; i < 16; i++ {
		titleKey = append(titleKey, ticket.Titlekey_block[i])
	}
	fmt.Println(strings.ToUpper(hex.EncodeToString(titleKey)))
	fmt.Println(len(hex.EncodeToString(titleKey)))

	fmt.Println("Titlekey_type", ticket.Titlekey_type)
	fmt.Println("Master_key_rev", ticket.Master_key_rev)
	fmt.Println("Ticket_id", ticket.Ticket_id)
	fmt.Println("Device_id", ticket.Device_id)
	fmt.Println("Rights_id", ticket.Rights_id)
	fmt.Println("Account_id", ticket.Account_id)
	// fmt.Println(data)
	// fmt.Println(hex.Dump(data))

	if strings.ToUpper(hex.EncodeToString(titleKey)) == titleDBKey {
		fmt.Println("Good!!!")
	}

	return nil
}
