package directory

import (
	"fmt"
	"os"

	collection "github.com/DblK/tinshop/gamescollection"
	"github.com/DblK/tinshop/nsp"
	"github.com/DblK/tinshop/repository"
)

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

	valid, err := nsp.IsTicketValid(f, key)
	if err != nil {
		fmt.Println("Error while opening NSP", err)
	}
	if !valid {
		fmt.Println("Your file", file.Path, "is not valid!")
	}
}
