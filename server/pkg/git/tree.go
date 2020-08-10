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
