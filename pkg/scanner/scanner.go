// Package scanner recursively (if specified) scans for directory content
// and builds up a representation tree that exactly represents the directory
// structure in given path. Example: Output of `tree` command
package scanner

import (
	"os"
	"path"
)

type File struct {
	FullPath    string
	Contents    []*File // Empty for a file
	IsDirectory bool
	Size        int64
	Errored     bool   // Specifies if we encountered error while reading this file
	Remarks     string // Remarks on file
}

func (f *File) readDir(channel chan<- *File, recursiveScan bool) {
	if !f.IsDirectory {
		// Neat and tricky
		channel <- f
		return
	}

	// Read current directory contents
	dirContents, err := os.ReadDir(f.FullPath)
	if err != nil {
		f.Errored = true
		f.Remarks = err.Error()

		channel <- f

		return
	}

	// Make list of specific size to accomomdate all
	// dirContents as `File` into f.Contents
	f.Contents = make([]*File, len(dirContents))

	for index, file := range dirContents {
		// File path will be parent path + filename
		filePath := path.Join(f.FullPath, file.Name())

		// Required to get .Size() method later
		info, err := file.Info()
		if err != nil {
			f.Errored = true
			f.Remarks = err.Error()
		}

		newFile := File{
			FullPath:    filePath,
			Contents:    nil,
			IsDirectory: file.IsDir(),
			Size:        info.Size(),
		}

		// Push new file into the f.Contents
		f.Contents[index] = &newFile
	}

	// Make them hear we explored one directory
	channel <- f

	// Recursively go for others
	if recursiveScan {
		for _, file := range f.Contents {
			file.readDir(channel, true)
		}
	}
}

func Scan(path string, channel chan *File) File {
	// Create instance of first root directory
	root := File{
		FullPath:    path,
		IsDirectory: true, // Should be a directory right?
	}

	root.readDir(channel, true)

	close(channel)

	return root
}
