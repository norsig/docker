package fs

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type memoryGroup struct {
}

func (s *memoryGroup) Set(d *data) error {
	dir, err := d.join("memory")
	// only return an error for memory if it was not specified
	if err != nil && (d.c.Memory != 0 || d.c.MemorySwap != 0) {
		return err
	}
	defer func() {
		if err != nil {
			os.RemoveAll(dir)
		}
	}()

	if d.c.Memory != 0 || d.c.MemorySwap != 0 {
		if d.c.Memory != 0 {
			if err := writeFile(dir, "memory.limit_in_bytes", strconv.FormatInt(d.c.Memory, 10)); err != nil {
				return err
			}
			if err := writeFile(dir, "memory.soft_limit_in_bytes", strconv.FormatInt(d.c.Memory, 10)); err != nil {
				return err
			}
		}
		// By default, MemorySwap is set to twice the size of RAM.
		// If you want to omit MemorySwap, set it to `-1'.
		if d.c.MemorySwap != -1 {
			if err := writeFile(dir, "memory.memsw.limit_in_bytes", strconv.FormatInt(d.c.Memory*2, 10)); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *memoryGroup) Remove(d *data) error {
	return removePath(d.path("memory"))
}

func (s *memoryGroup) Stats(d *data) (map[string]float64, error) {
	paramData := make(map[string]float64)
	path, _ := d.path("memory")
	paramPath := filepath.Join(path, "memory.stat")
	f, err := os.Open(paramPath)
	if err != nil {
		return paramData, err
	}
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		fields := strings.Fields(sc.Text())
		t := fields[0]
		v, err := strconv.Atoi(fields[1])
		if err != nil {
			fmt.Errorf("Error parsing %s stats: %s", t, err)
			continue
		}
		paramData[t] = float64(v)
	}
	return paramData, nil
}
