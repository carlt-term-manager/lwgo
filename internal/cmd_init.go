package internal

import (
	"fmt"
	"os"
	"path"
)

func readInitArgs(p *Entry) {
	fmt.Print(`Project Name: `)
	fmt.Scanln(&p.Name)

	fmt.Print(`Project Version(default=1.0.0):`)
	fmt.Scanln(&p.Version)
	if p.Version == "" {
		p.Version = "1.0.0"
	}
}

// init project and create dir, vendor, package.json
func CmdInit() {
	p := newEntry()
	readInitArgs(p)

	if _, err := os.Stat(p.Name); err != nil {
		if os.IsNotExist(err) {
			os.Mkdir(p.Name, 0755)
		}
	}

	handleError(checkVendorDir(path.Join(p.Name, Vendor)))

	p.Save(path.Join(p.Name, PackageFile))
}
