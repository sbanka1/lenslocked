package dbx

import (
	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox"
	dbxfiles "github.com/dropbox/dropbox-sdk-go-unofficial/dropbox/files"
)

type Folder struct {
	Name string
	Path string
}

type File struct {
	Name string
	Path string
}

func List(accessToken, path string) ([]Folder, []File, error) {
	config := dropbox.Config{
		Token: accessToken,
	}
	client := dbxfiles.New(config)
	res, err := client.ListFolder(&dbxfiles.ListFolderArg{
		Path: path,
	})
	if err != nil {
		return nil, nil, err
	}

	var folders []Folder
	var files []File

	for _, entry := range res.Entries {
		switch meta := entry.(type) {
		case *dbxfiles.FolderMetadata:
			folders = append(folders, Folder{
				Name: meta.Name,
				Path: meta.PathLower,
			})
		case *dbxfiles.FileMetadata:
			files = append(files, File{
				Name: meta.Name,
				Path: meta.PathLower,
			})
		}
	}

	return folders, files, nil
}
