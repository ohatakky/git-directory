package git

// todo: support tree view
// fzf, err := c.Tree()
// if err != nil {
// 	return err
// }
// walker := object.NewTreeWalker(fzf, true, nil)
// defer walker.Close()
// for {
// 	name, entry, err := walker.Next()
// 	if err == io.EOF {
// 		break
// 	}
// 	if !entry.Mode.IsFile() {
// 		fmt.Println(entry.Name, name)
// 	}
// }

type Tree struct {
	Uuid     string
	Name     string
	Path     string
	IsFile   bool
	Children *Tree
}

// func (c *Client) Trees(commits []Commit) {
// 	defer close(c.Send)
// 	for i := len(commits) - 1; i >= 0; i-- {

// 		err := commits[i].Files.ForEach(func(f *object.File) error {
// 			t := &Tree{
// 				Uuid:   uuid.String(),
// 				Name:   filepath.Base(f.Name),
// 				Path:   filepath.Dir(f.Name),
// 				IsFile: f.Mode.IsFile(),
// 			}
// 			return nil
// 		})

// 	}
// }

//
