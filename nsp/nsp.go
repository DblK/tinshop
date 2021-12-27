package nsp

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"strings"
	"unsafe"
)

// CheckNCAKey ensure that the extracted key meet titledb
func CheckNCAKey() {
	fmt.Println("NSP.CheckNCAKey")
}

// IsTicketValid return if ticket is valid or not
func IsTicketValid(file io.ReadSeeker, titleDBKey string) (bool, error) {
	newNSP := &nspFile{}

	// Read Header
	nspHeader := pfs0Header{}
	data := make([]byte, unsafe.Sizeof(nspHeader))
	_, _ = file.Read(data)
	buffer := bytes.NewBuffer(data)
	_ = binary.Read(buffer, binary.LittleEndian, &nspHeader)
	newNSP.Header = nspHeader

	if string(newNSP.Header.Magic[:]) != "PFS0" {
		return false, errors.New("header Magic is not present")
	}
	// fmt.Println(nspHeader)

	// Read file entry
	for i := 0; i < int(nspHeader.FileCnt); i++ {
		nspEntry := pfs0FileEntry{}
		data := make([]byte, unsafe.Sizeof(nspEntry))
		_, _ = file.Read(data)
		buffer := bytes.NewBuffer(data)
		_ = binary.Read(buffer, binary.LittleEndian, &nspEntry)
		newNSP.FileEntry = append(newNSP.FileEntry, nspEntry)
		// fmt.Println(nspEntry)
	}
	// fmt.Println(len(newNSP.FileEntry))
	// fmt.Println(newNSP.FileEntry)

	// Read nspStrTable + Display file_name
	nspStrTable := make([]byte, nspHeader.StrTableSize)
	_, _ = file.Read(nspStrTable)
	// fmt.Println(nspStrTable)
	// fmt.Println(nspHeader.Str_table_size)

	var tikOffset uint64
	var tikSize uint64

	for i := 0; i < int(nspHeader.FileCnt); i++ {
		start := newNSP.FileEntry[i].FilenameOffset
		if i != int(nspHeader.FileCnt)-1 {
			end := newNSP.FileEntry[i+1].FilenameOffset - 1
			// fmt.Println(string(nspStrTable[start:end]))
			newNSP.FileName = append(newNSP.FileName, string(nspStrTable[start:end]))
		} else {
			// fmt.Println(string(nspStrTable[start:]))
			newNSP.FileName = append(newNSP.FileName, string(nspStrTable[start:]))
		}

		// Compute Ticket information
		if newNSP.FileName[i][len(newNSP.FileName[i])-4:] == ".tik" {
			fmt.Println("Found ticket!")
			tikOffset = (uint64(unsafe.Sizeof(nspHeader)) + (uint64(unsafe.Sizeof(newNSP.FileEntry[i])) * uint64(len(newNSP.FileEntry))) + uint64(len(nspStrTable)) + newNSP.FileEntry[i].FileOffset)
			tikSize = newNSP.FileEntry[i].FileSize
		}
	}
	// Display Ticket
	fmt.Println(newNSP.FileName)
	fmt.Println(tikOffset, tikSize)

	ticket := &rsa2048SHA256Ticket{}
	fmt.Println(unsafe.Sizeof(ticket))
	_, _ = file.Seek(int64(tikOffset), 0)

	data = make([]byte, tikSize)
	_, _ = file.Read(data)
	buffer = bytes.NewBuffer(data)
	err := binary.Read(buffer, binary.LittleEndian, ticket)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(ticket.SigIssuer[:]))
	fmt.Println("Titlekey_block", ticket.TitlekeyBlock)
	// fmt.Println("Titlekey_block", strconv.Itoa(int(ticket.TitlekeyBlock[0])))
	var titleKey []byte
	for i := 0; i < 16; i++ {
		titleKey = append(titleKey, ticket.TitlekeyBlock[i])
	}
	fmt.Println(strings.ToUpper(hex.EncodeToString(titleKey)))
	fmt.Println(len(hex.EncodeToString(titleKey)))

	fmt.Println("Titlekey_type", ticket.TitlekeyType)
	fmt.Println("Master_key_rev", ticket.MasterKeyRev)
	fmt.Println("Ticket_id", ticket.TicketID)
	fmt.Println("Device_id", ticket.DeviceID)
	fmt.Println("Rights_id", ticket.RightsID)
	fmt.Println("Account_id", ticket.AccountID)
	// fmt.Println(data)
	// fmt.Println(hex.Dump(data))

	if strings.ToUpper(hex.EncodeToString(titleKey)) == titleDBKey {
		fmt.Println("Good!!!")
		return true, nil
	}

	return false, nil
}
