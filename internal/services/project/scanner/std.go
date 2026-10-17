package scanner

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"golang.org/x/tools/go/packages"
)

var defaultStdPackages = map[string]struct{}{
	"archive/tar":          {},
	"archive/zip":          {},
	"bufio":                {},
	"bytes":                {},
	"cmp":                  {},
	"compress/bzip2":       {},
	"compress/flate":       {},
	"compress/gzip":        {},
	"compress/lzw":         {},
	"compress/zlib":        {},
	"container/heap":       {},
	"container/list":       {},
	"container/ring":       {},
	"context":              {},
	"crypto":               {},
	"crypto/aes":           {},
	"crypto/cipher":        {},
	"crypto/des":           {},
	"crypto/dsa":           {},
	"crypto/ecdh":          {},
	"crypto/ecdsa":         {},
	"crypto/ed25519":       {},
	"crypto/elliptic":      {},
	"crypto/hmac":          {},
	"crypto/md5":           {},
	"crypto/rand":          {},
	"crypto/rc4":           {},
	"crypto/rsa":           {},
	"crypto/sha1":          {},
	"crypto/sha256":        {},
	"crypto/sha512":        {},
	"crypto/subtle":        {},
	"crypto/tls":           {},
	"crypto/x509":          {},
	"crypto/x509/pkix":     {},
	"database/sql":         {},
	"database/sql/driver":  {},
	"debug/buildinfo":      {},
	"debug/dwarf":          {},
	"debug/elf":            {},
	"debug/gosym":          {},
	"debug/macho":          {},
	"debug/pe":             {},
	"debug/plan9obj":       {},
	"embed":                {},
	"encoding":             {},
	"encoding/ascii85":     {},
	"encoding/asn1":        {},
	"encoding/base32":      {},
	"encoding/base64":      {},
	"encoding/binary":      {},
	"encoding/csv":         {},
	"encoding/gob":         {},
	"encoding/hex":         {},
	"encoding/json":        {},
	"encoding/pem":         {},
	"encoding/xml":         {},
	"errors":               {},
	"expvar":               {},
	"flag":                 {},
	"fmt":                  {},
	"go/ast":               {},
	"go/build":             {},
	"go/build/constraint":  {},
	"go/constant":          {},
	"go/doc":               {},
	"go/doc/comment":       {},
	"go/format":            {},
	"go/importer":          {},
	"go/parser":            {},
	"go/printer":           {},
	"go/scanner":           {},
	"go/token":             {},
	"go/types":             {},
	"go/version":           {},
	"hash":                 {},
	"hash/adler32":         {},
	"hash/crc32":           {},
	"hash/crc64":           {},
	"hash/fnv":             {},
	"hash/maphash":         {},
	"html":                 {},
	"html/template":        {},
	"image":                {},
	"image/color":          {},
	"image/color/palette":  {},
	"image/draw":           {},
	"image/gif":            {},
	"image/jpeg":           {},
	"image/png":            {},
	"index/suffixarray":    {},
	"io":                   {},
	"io/fs":                {},
	"io/ioutil":            {},
	"iter":                 {},
	"log":                  {},
	"log/slog":             {},
	"log/syslog":           {},
	"maps":                 {},
	"math":                 {},
	"math/big":             {},
	"math/bits":            {},
	"math/cmplx":           {},
	"math/rand":            {},
	"math/rand/v2":         {},
	"mime":                 {},
	"mime/multipart":       {},
	"mime/quotedprintable": {},
	"net":                  {},
	"net/http":             {},
	"net/http/cgi":         {},
	"net/http/cookiejar":   {},
	"net/http/fcgi":        {},
	"net/http/httptest":    {},
	"net/http/httptrace":   {},
	"net/http/httputil":    {},
	"net/http/pprof":       {},
	"net/mail":             {},
	"net/netip":            {},
	"net/rpc":              {},
	"net/rpc/jsonrpc":      {},
	"net/smtp":             {},
	"net/textproto":        {},
	"net/url":              {},
	"os":                   {},
	"os/exec":              {},
	"os/signal":            {},
	"os/user":              {},
	"path":                 {},
	"path/filepath":        {},
	"plugin":               {},
	"reflect":              {},
	"regexp":               {},
	"regexp/syntax":        {},
	"runtime":              {},
	"runtime/cgo":          {},
	"runtime/coverage":     {},
	"runtime/debug":        {},
	"runtime/metrics":      {},
	"runtime/pprof":        {},
	"runtime/trace":        {},
	"slices":               {},
	"sort":                 {},
	"strconv":              {},
	"strings":              {},
	"structs":              {},
	"sync":                 {},
	"sync/atomic":          {},
	"syscall":              {},
	"testing":              {},
	"testing/fstest":       {},
	"testing/iotest":       {},
	"testing/quick":        {},
	"testing/slogtest":     {},
	"text/scanner":         {},
	"text/tabwriter":       {},
	"text/template":        {},
	"text/template/parse":  {},
	"time":                 {},
	"time/tzdata":          {},
	"unicode":              {},
	"unicode/utf16":        {},
	"unicode/utf8":         {},
	"unique":               {},
	"unsafe":               {},
	"weak":                 {},
}

func loadStdPackages() map[string]struct{} {
	stdMap := make(map[string]struct{}, len(defaultStdPackages)+64)
	for pkg := range defaultStdPackages {
		stdMap[pkg] = struct{}{}
	}

	goroot := os.Getenv("GOROOT")
	if goroot == "" {
		if out, err := exec.Command("go", "env", "GOROOT").Output(); err == nil {
			goroot = strings.TrimSpace(string(out))
		}
	}
	if goroot == "" {
		//nolint:staticcheck // SA1019: runtime.GOROOT is deprecated in Go 1.24+, used as fallback when go binary is unavailable
		goroot = runtime.GOROOT()
	}

	if goroot != "" {
		loadFromGoroot(goroot, stdMap)
	}

	loadFromPackages(goroot, stdMap)

	return stdMap
}

func loadFromGoroot(goroot string, stdMap map[string]struct{}) {
	cleanGoroot := filepath.Clean(goroot)
	if !filepath.IsAbs(cleanGoroot) {
		if abs, err := filepath.Abs(cleanGoroot); err == nil {
			cleanGoroot = abs
		}
	}

	srcDir := filepath.Clean(filepath.Join(cleanGoroot, "src"))
	rel, err := filepath.Rel(cleanGoroot, srcDir)
	if err != nil || rel != "src" {
		return
	}

	//nolint:gosec // G703: path is sanitized with filepath.Clean and verified to stay within goroot
	fi, err := os.Stat(srcDir)
	if err != nil || !fi.IsDir() {
		return
	}

	//nolint:gosec // G703: path is sanitized with filepath.Clean and verified to stay within goroot
	_ = filepath.WalkDir(srcDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if !d.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(srcDir, path)
		if err != nil || rel == "." || strings.HasPrefix(rel, "..") {
			return nil
		}
		pkg := filepath.ToSlash(rel)
		if pkg == "cmd" || strings.HasPrefix(pkg, "cmd/") ||
			pkg == "vendor" || strings.HasPrefix(pkg, "vendor/") ||
			pkg == "builtin" || strings.HasPrefix(pkg, "builtin/") ||
			strings.HasPrefix(d.Name(), ".") || strings.HasPrefix(d.Name(), "_") {
			return filepath.SkipDir
		}

		cleanPath := filepath.Clean(path)
		//nolint:gosec // G703: path is verified to be within srcDir
		entries, err := os.ReadDir(cleanPath)
		if err != nil {
			return nil
		}
		for _, e := range entries {
			if !e.IsDir() && strings.HasSuffix(e.Name(), ".go") && !strings.HasSuffix(e.Name(), "_test.go") {
				stdMap[pkg] = struct{}{}
				break
			}
		}
		return nil
	})
}

func loadFromPackages(goroot string, stdMap map[string]struct{}) {
	cfg := &packages.Config{
		Mode: packages.NeedName,
	}
	if goroot != "" {
		cleanGoroot := filepath.Clean(goroot)
		cfg.Dir = cleanGoroot
		cfg.Env = append(os.Environ(), "GO111MODULE=off", fmt.Sprintf("GOROOT=%s", cleanGoroot))
	}

	stdPackages, err := packages.Load(cfg, "std")
	if err != nil {
		return
	}

	for _, stdPackage := range stdPackages {
		if stdPackage.ID != "" &&
			!strings.HasPrefix(stdPackage.ID, "cmd/") &&
			!strings.HasPrefix(stdPackage.ID, "vendor/") {
			stdMap[stdPackage.ID] = struct{}{}
		}
	}
}

func isStandardImport(importPath, moduleName string) bool {
	if importPath == "" || importPath == "C" {
		return true
	}
	if importPath == moduleName || strings.HasPrefix(importPath, moduleName+"/") {
		return false
	}
	parts := strings.Split(importPath, "/")
	if len(parts) > 0 && !strings.Contains(parts[0], ".") && !strings.HasPrefix(parts[0], "_") {
		return true
	}
	return false
}
