package version

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

var (
	// TODO(minkezhang): Use
	//
	// 1. gitver https://git.rootprojects.org/root/go-gitver, or
	// 1. ldflags https://www.digitalocean.com/community/tutorials/using-ldflags-to-set-version-information-for-go-applications
	version string
)

func hash() string {
	e, err := os.Executable()
	if err != nil {
		return ""
	}

	fp, err := filepath.EvalSymlinks(e)
	if err != nil {
		return ""
	}

	data, err := os.ReadFile(fp)
	if err != nil {
		return ""
	}

	buf := sha256.New()
	buf.Write([]byte(data))

	// Only need
	return hex.EncodeToString(buf.Sum(nil)[:4])
}

// Sem returns a valid semver.org string of the build.
func Sem() string {
	if version != "" {
		return version
	}

	v := hash()
	if v != "" {
		v = fmt.Sprintf(
			"v0.0.0+%s.sha.%s",
			time.Now().Format("200601021504"),
			v,
		)
	}
	return v
}
