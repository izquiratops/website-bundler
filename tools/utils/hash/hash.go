package hash

import (
	"crypto/md5"
	"fmt"
	"path/filepath"
	"strings"
)

func GenerateHash(sourcePath string, content []byte) string {
	hash := md5.Sum(content)
	fileName := strings.TrimSuffix(filepath.Base(sourcePath), filepath.Ext(sourcePath))

	return fmt.Sprintf("%s-%x%s", fileName, hash[:8], filepath.Ext(sourcePath))
}
