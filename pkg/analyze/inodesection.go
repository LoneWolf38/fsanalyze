package analyze

import (
	"bufio"
	"io"
	"os"

	pb "fsanalyze/protofiles/src/hadoop_hdfs_fsimage"
)

func parseInodeSection(info *pb.FileSummary_Section, imageFile *os.File) (map[InodeId]INode, EntityCount, error) {
	var (
		nameIdMap           = make(map[InodeId]INode)
		files        uint32 = 0
		dirs         uint32 = 0
		symlinks     uint32 = 0
		inode               = &pb.INodeSection_INode{}
		inodeSection        = &pb.INodeSection{}
		entityCount         = EntityCount{}
	)
	_, err := imageFile.Seek(int64(info.GetOffset()), io.SeekCurrent)
	if err != nil {
		return nameIdMap, entityCount, err
	}
	bReader := bufio.NewReader(imageFile)
	_, err = readDelimited(bReader, inodeSection)
	if err != nil {
		return nameIdMap, entityCount, err
	}
	for a := uint64(0); a < inodeSection.GetNumInodes(); a++ {
		_, err := readDelimited(bReader, inode)
		if err != nil {
			return nameIdMap, entityCount, err
		}
		switch inode.GetType() {
		case 1:
			nameIdMap[InodeId(inode.GetId())] = INode{inode.GetName(), InodeId(inode.GetId()), FILE}
			files++
		case 2:
			nameIdMap[InodeId(inode.GetId())] = INode{inode.GetName(), InodeId(inode.GetId()), DIRECTORY}
			dirs++
		case 3:
			nameIdMap[InodeId(inode.GetId())] = INode{inode.GetName(), InodeId(inode.GetId()), SYMLINK}
			symlinks++
		}
	}
	entityCount = EntityCount{Files: files, Directories: dirs, Symlinks: symlinks}
	return nameIdMap, entityCount, nil
}
