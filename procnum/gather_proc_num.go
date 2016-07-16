package procnum

import (
	"io"
	"os"
	"strconv"
	"strings"
)

type ProcResult struct {
	Procnum *float64
}

func GatherProcInfo(statpath string) (*ProcResult, error) {
	p := new(ProcResult)
	procstat, err := os.Open(statpath)
	if err != nil {
		return nil, err
	}
	defer procstat.Close()
	buf := make([]byte, 4096)
	var strbuf string
	for {
		n, err := procstat.Read(buf)
		if n == 0 && err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		// get lines
		lines := strings.Split(string(buf), "\n")
		if strbuf != "" {
			lines[0] = strbuf + lines[0]
		}
		if len(lines) > 0 && (len(buf) > 0 && buf[len(buf)-1] == 0x0a) {
			strbuf = lines[len(lines)-1]
		}
		for _, l := range lines {
			if strings.HasPrefix(l, "processes") {
				k := strings.Split(l, " ")
				num, err := strconv.ParseFloat(k[1], 64)
				if err != nil {
					return nil, err
				}
				p.Procnum = &num
				return p, nil
			}
		}
	}
	return nil, nil
}
