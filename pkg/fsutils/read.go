package fsutils

import (
	"os"
)

var osReadFile = os.ReadFile

func ReadFileData(fullName string, max int) (data []byte, err error) {
	if max == 0 {
		data, err = osReadFile(fullName)
	} else {
		var file *os.File
		if file, err = os.Open(fullName); err == nil {
			defer func() {
				_ = file.Close()
			}()
			if max > 0 {
				data = make([]byte, max)
				var n int
				if n, err = file.Read(data); err == nil {
					data = data[:n]
				}
			} else {
				absMax := int64(-max)
				var fi os.FileInfo
				if fi, err = file.Stat(); err == nil {
					size := fi.Size()
					if size > absMax {
						if _, err = file.Seek(-absMax, 2); err == nil {
							data = make([]byte, absMax)
							var n int
							if n, err = file.Read(data); err == nil {
								data = data[:n]
							}
						}
					} else {
						data = make([]byte, size)
						var n int
						if n, err = file.Read(data); err == nil {
							data = data[:n]
						}
					}
				}
			}
		}
	}
	return
}
