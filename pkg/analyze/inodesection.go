package analyze

import (
	"bufio"
	"fmt"
	"io"
	"os"

	pb "fsanalyze/protofiles/src/hadoop_hdfs_fsimage"

	"github.com/VictoriaMetrics/fastcache"
)

var CacheFile = "/data01/acceldata/cache"

func parseInodeSection(info *pb.FileSummary_Section, imageFile *os.File, cache *fastcache.Cache) (map[InodeId]INode, EntityCount, error) {
	var (
		nameIdMap           = make(map[InodeId]INode)
		files        uint32 = 0
		dirs         uint32 = 0
		symlinks     uint32 = 0
		inode               = &pb.INodeSection_INode{}
		inodeSection        = &pb.INodeSection{}
		entityCount         = EntityCount{}
	)

	pos, err := imageFile.Seek(int64(info.GetOffset()), io.SeekStart)
	if err != nil {
		return nameIdMap, entityCount, err
	}
	fmt.Println("Position of seek: ", pos)
	bReader := bufio.NewReader(imageFile)
	_, _, err = readDelimited(bReader, inodeSection)
	if err != nil {
		return nameIdMap, entityCount, err
	}
	for a := uint64(0); a < inodeSection.GetNumInodes(); a++ {
		_, _, err := readDelimited(bReader, inode)
		if err != nil {
			return nameIdMap, entityCount, err
		}
		switch inode.GetType() {
		case 1:
			//nameIdMap[InodeId(inode.GetId())] = INode{inode.GetName(), InodeId(inode.GetId()), FILE}
			files++
		case 2:
			//nameIdMap[InodeId(inode.GetId())] = INode{inode.GetName(), InodeId(inode.GetId()), DIRECTORY}
			dirs++
		case 3:
			//nameIdMap[InodeId(inode.GetId())] = INode{inode.GetName(), InodeId(inode.GetId()), SYMLINK}
			symlinks++
		}
	}
	entityCount = EntityCount{Files: files, Directories: dirs, Symlinks: symlinks}
	return nameIdMap, entityCount, nil
}

func parseInodeDirectorySection(info *pb.FileSummary_Section, imageFile *os.File) (map[ParentId][]uint64, error) {
	var (
		parChildrenMap = make(map[ParentId][]uint64)
		dirEntry       = &pb.INodeDirectorySection_DirEntry{}
	)
	length := info.GetLength()
	_, err := imageFile.Seek(int64(info.GetOffset()), io.SeekStart)
	if err != nil {
		return parChildrenMap, err
	}
	bReader := bufio.NewReader(imageFile)

	for a := length; a > 0; {
		_, newPos, err := readDelimited(bReader, dirEntry)
		if err != nil {
			return parChildrenMap, err
		}
		//parChildrenMap[ParentId(dirEntry.GetParent())] = dirEntry.GetChildren()
		a -= newPos
	}
	return parChildrenMap, nil
}
