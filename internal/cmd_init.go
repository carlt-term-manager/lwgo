package internal

import (
	"fmt"
	"os"
	"path"
)

func readInitArgs(p *mod) {
	fmt.Print(`Project Name: `)
	fmt.Scanln(&p.Name)

	fmt.Print(`Project Version(default=1.00.00):`)
	fmt.Scanln(&p.Version)
	if p.Version == "" {
		p.Version = "1.00.00"
	}
}

// init project and create dir, vendor, package.json
func CmdInit() {
	p := newMod()
	readInitArgs(p)

	if _, err := os.Stat(p.Name); err != nil {
		if os.IsNotExist(err) {
			os.Mkdir(p.Name, 0755)
		}
	}

	handleError(checkVendorDir(path.Join(p.Name, Vendor)))

	p.Save(path.Join(p.Name, PackageFile))
}
