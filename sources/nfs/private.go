package nfs

import (
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/DblK/tinshop/config"
	collection "github.com/DblK/tinshop/gamescollection"
	"github.com/DblK/tinshop/nsp"
	"github.com/DblK/tinshop/repository"
	"github.com/DblK/tinshop/utils"
	"github.com/vmware/go-nfs-client/nfs"
	"github.com/vmware/go-nfs-client/nfs/rpc"
	"github.com/vmware/go-nfs-client/nfs/util"
)

func getHostTarget(share string) (string, string, error) {
	shareInfos := strings.Split(share, ":")

	if len(shareInfos) != 2 {
		return "", "", errors.New("Error parsing the nfs share configuration " + share)
	}
	return shareInfos[0], shareInfos[1], nil
}

func loadGamesNfs(share string) {
	if config.GetConfig().DebugNfs() {
		util.DefaultLogger.SetDebug(true)
	}

	host, target, err := getHostTarget(share)
	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("Loading games from nfs (host=%s target=%s)\n", host, target)

	// Connection
	mount, v := nfsConnect(host, target)
	defer mount.Close()
	defer v.Close()

	nfsGames := lookIntoNfsDirectory(v, share, ".")

	mount.Close()
	gameFiles = append(gameFiles, nfsGames...)

	// Add all files
	if len(nfsGames) > 0 {
		collection.AddNewGames(nfsGames)
	}
}

func nfsConnect(host, target string) (*nfs.Mount, *nfs.Target) {
	mount, err := nfs.DialMount(host)
	if err != nil {
		log.Fatalf("unable to dial MOUNT service: %v", err)
	}

	// Mount drive
	v, err := mount.Mount(target, rpc.AuthNull)
	if err != nil {
		log.Fatalf("unable to mount volume: %v", err)
	}

	return mount, v
}

func lookIntoNfsDirectory(v *nfs.Target, share, path string) []repository.FileDesc {
	// Retrieve all directories
	log.Printf("Retrieving all files in directory ('%s')...\n", path)

	dirs, err := v.ReadDirPlus(path)
	if err != nil {
		_ = fmt.Errorf("readdir error: %s", err.Error())
		return nil
	}

	var newGameFiles []repository.FileDesc

	for _, dir := range dirs {
		if dir.FileName != "." && dir.FileName != ".." {
			if dir.IsDir() {
				subDirGameFiles := lookIntoNfsDirectory(v, share, computePath(path, dir))
				newGameFiles = append(newGameFiles, subDirGameFiles...)
			} else {
				nfsRootPath := share
				if path != "." {
					nfsRootPath += path
				}
				extension := filepath.Ext(dir.FileName)
				if extension != ".nsp" && extension != ".nsz" {
					continue
				}

				newFile := repository.FileDesc{Size: dir.Size(), Path: nfsRootPath + "/" + dir.FileName}
				names := utils.ExtractGameID(dir.FileName)

				if names.ShortID() != "" {
					newFile.GameID = names.ShortID()
					newFile.GameInfo = names.FullID()
					newFile.HostType = repository.NFSShare

					if config.GetConfig().VerifyNSP() {
						valid, errTicket := nspCheck(newFile)
						if valid || (errTicket != nil && errTicket.Error() == "TitleDBKey for game "+newFile.GameID+" is not found") {
							newGameFiles = append(newGameFiles, newFile)
						} else {
							log.Println(errTicket)
						}
					} else {
						newGameFiles = append(newGameFiles, newFile)
					}
				} else {
					log.Println("Ignoring file because parsing failed", dir.FileName)
				}
			}
		}
	}

	return newGameFiles
}

func computePath(path string, dir *nfs.EntryPlus) string {
	var newPath string
	if path == "." {
		newPath = "/" + dir.FileName
	} else {
		newPath = path + "/" + dir.FileName
	}
	return newPath
}

func nspCheck(file repository.FileDesc) (bool, error) {
	key, err := collection.GetKey(file.GameID)
	if err != nil {
		return false, err
	}

	host, target, err := getHostTarget(file.Path)
	if err != nil {
		return false, err
	}

	mount, v := nfsConnect(host, filepath.Dir(target))
	defer mount.Close()
	defer v.Close()

	f, err := v.Open(filepath.Base(target))
	if err != nil {
		return false, err
	}
	defer f.Close()

	log.Println("Verifying Ticket:", file.Path)
	valid, err := nsp.IsTicketValid(f, key)
	if err != nil {
		return false, err
	}
	if !valid {
		return false, errors.New("Your file" + file.Path + "is not valid!")
	}

	return valid, err
}
