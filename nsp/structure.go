package nsp

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
