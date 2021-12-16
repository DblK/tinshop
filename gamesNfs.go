package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/vmware/go-nfs-client/nfs"
	"github.com/vmware/go-nfs-client/nfs/rpc"
	"github.com/vmware/go-nfs-client/nfs/util"
)

var nfsShares []string
var debugNfs bool

func loadGamesNfsShares() {
	for _, share := range nfsShares {
		loadGamesNfs(share)
	}
}

func loadGamesNfs(share string) {
	if debugNfs {
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
		AddNewGames(nfsGames)
	}
}

func nfsConnect(host string, target string) (*nfs.Mount, *nfs.Target) {
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

func lookIntoNfsDirectory(v *nfs.Target, share string, path string) []FileDesc {
	// Retrieve all directories
	log.Printf("Retrieving all files in directory ('%s')...\n", path)

	dirs, err := v.ReadDirPlus(path)
	if err != nil {
		_ = fmt.Errorf("readdir error: %s", err.Error())
		return nil
	}

	var newGameFiles []FileDesc

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
					nfsRootPath = nfsRootPath + path
				}

				newFile := FileDesc{size: dir.Size(), path: nfsRootPath + "/" + dir.FileName}
				names := ExtractGameId(dir.FileName)

				if names.ShortId != "" {
					newFile.gameId = names.ShortId
					newFile.gameInfo = names.FullId
					newFile.hostType = NFSShare
					newGameFiles = append(newGameFiles, newFile)
				} else {
					log.Println("Ignoring file because parsing failed", dir.FileName)
				}
			}
		}
	}

	return newGameFiles
}

func downloadNfsFile(w http.ResponseWriter, r *http.Request, share string) {
	if debugNfs {
		util.DefaultLogger.SetDebug(true)
	}

	shareInfos := strings.Split(share, ":")
	if len(shareInfos) != 2 {
		log.Printf("Error parsing the nfs share configuration (%s)\n", share)
		return
	}

	// Cut the share string
	host := shareInfos[0]
	path := shareInfos[1]
	name := path[strings.LastIndex(path, "/")+1:]
	target := path[:strings.LastIndex(path, "/")]

	// Connection
	mount, v := nfsConnect(host, target)
	defer mount.Close()
	defer v.Close()

	// Open file
	rdr, err := v.Open(name)
	if err != nil {
		util.Errorf("read error: %v", err)
		return
	}
	// Stats file
	fsInfo, _, err := v.Lookup(name)
	if err != nil {
		log.Fatalf("lookup error: %s", err.Error())
	}

	byteRange := strings.Split(strings.Replace(strings.Join(r.Header["Range"], ""), "bytes=", "", -1), "-")
	start, _ := strconv.Atoi(byteRange[0])
	end, _ := strconv.Atoi(byteRange[1])

	if start > int(fsInfo.Size()) || end > int(fsInfo.Size()) {
		w.WriteHeader(http.StatusRequestedRangeNotSatisfiable)
		return
	}

	// Check if partial content
	if end-start+1 == int(fsInfo.Size()) {
		// Full Content
		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Length", fmt.Sprint(fsInfo.Size()))
		_, _ = io.Copy(w, rdr)
	} else {
		// Partial Content
		rng := make([]byte, end-start+1)
		if start != 0 {
			_, _ = rdr.Seek(int64(start), 0)
		}
		_, err = rdr.Read(rng)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("Error while trying to read file from nfs", err)
			return
		}
		w.WriteHeader(http.StatusPartialContent)
		w.Header().Add("Content-Range", "bytes "+fmt.Sprint(start)+"-"+fmt.Sprint(end)+"/"+fmt.Sprint(fsInfo.Size()))
		w.Header().Add("Accept-Ranges", "bytes")
		w.Header().Add("Content-Length", fmt.Sprint(end-start+1))
		_, _ = w.Write(rng)
	}

	if err = mount.Unmount(); err != nil {
		log.Fatalf("unable to umount target: %v", err)
	}
}
