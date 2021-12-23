package nfs

import (
	"fmt"
	"log"
	"strings"

	"github.com/DblK/tinshop/config"
	collection "github.com/DblK/tinshop/gamescollection"
	"github.com/DblK/tinshop/repository"
	"github.com/DblK/tinshop/utils"
	"github.com/vmware/go-nfs-client/nfs"
	"github.com/vmware/go-nfs-client/nfs/rpc"
	"github.com/vmware/go-nfs-client/nfs/util"
)

func loadGamesNfs(share string) {
	if config.GetConfig().DebugNfs() {
		util.DefaultLogger.SetDebug(true)
	}

	shareInfos := strings.Split(share, ":")
	if len(shareInfos) != 2 {
		log.Printf("Error parsing the nfs share configuration (%s)\n", share)
		return
	}

	host := shareInfos[0]
	target := shareInfos[1]

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
				var newPath string
				if path == "." {
					newPath = "/" + dir.FileName
				} else {
					newPath = path + "/" + dir.FileName
				}
				subDirGameFiles := lookIntoNfsDirectory(v, share, newPath)
				newGameFiles = append(newGameFiles, subDirGameFiles...)
			} else {
				nfsRootPath := share
				if path != "." {
					nfsRootPath += path
				}

				newFile := repository.FileDesc{Size: dir.Size(), Path: nfsRootPath + "/" + dir.FileName}
				names := utils.ExtractGameID(dir.FileName)

				if names.ShortID() != "" {
					newFile.GameID = names.ShortID()
					newFile.GameInfo = names.FullID()
					newFile.HostType = repository.NFSShare
					newGameFiles = append(newGameFiles, newFile)
				} else {
					log.Println("Ignoring file because parsing failed", dir.FileName)
				}
			}
		}
	}

	return newGameFiles
}
