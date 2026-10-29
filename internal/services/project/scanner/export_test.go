package scanner

type PackagesLoader = packagesLoader

func NewScannerWithLoader(loader PackagesLoader) *Scanner {
	return newScannerWithLoader(loader)
}

func (r *Scanner) HasStdPackage(name string) bool {
	_, ok := r.stdPackages[name]
	return ok
}

func (r *Scanner) StdPackagesCount() int {
	return len(r.stdPackages)
}
