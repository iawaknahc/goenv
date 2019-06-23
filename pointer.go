package goenv

import (
	"sort"
	"strconv"
	"strings"
)

type pointerPart struct {
	Name       string
	FieldIndex int
	SliceIndex int
}

func (p pointerPart) String() string {
	if p.Name != "" {
		return p.Name
	}
	return strconv.Itoa(p.SliceIndex)
}

type pointer []pointerPart

func (p pointer) StructField(name string, fieldIndex int) pointer {
	output := make([]pointerPart, len(p))
	copy(output, p)
	output = append(output, pointerPart{
		Name:       name,
		FieldIndex: fieldIndex,
	})
	return output
}

func (p pointer) SliceIndex(sliceIndex int) pointer {
	output := make([]pointerPart, len(p))
	copy(output, p)
	output = append(output, pointerPart{
		SliceIndex: sliceIndex,
	})
	return output

}

func (p pointer) String() string {
	strParts := make([]string, len(p))
	for i, part := range p {
		strParts[i] = part.String()
	}
	return strings.Join(strParts, "_")
}

func (p pointer) specialize(name string) (pointer, bool) {
	var output pointer
	l := len(p)
	for i, part := range p {
		if part.Name == "" {
			trimmed := strings.TrimLeft(name, "0123456789")
			if trimmed == name {
				return nil, false
			}
			idxStr := name[:len(name)-len(trimmed)]
			idx, _ := strconv.Atoi(idxStr)
			output = append(output, pointerPart{
				SliceIndex: idx,
			})
			name = trimmed
		} else {
			if !strings.HasPrefix(name, part.Name) {
				return nil, false
			}
			output = append(output, part)
			name = strings.TrimPrefix(name, part.Name)
		}
		if i != l-1 {
			// Skip _
			name = name[1:]
		}
	}
	if name != "" {
		return nil, false
	}
	return output, true
}

type ascPointers []pointer

func (p ascPointers) Len() int {
	return len(p)
}

func (p ascPointers) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p ascPointers) Less(i, j int) bool {
	ptr1 := p[i]
	ptr2 := p[j]
	i = 0
	j = 0
	for i < len(ptr1) && j < len(ptr2) {
		part1 := ptr1[i]
		part2 := ptr2[j]
		switch {
		case part1.Name != "" && part2.Name != "":
			switch {
			case part1.Name < part2.Name:
				return true
			case part1.Name > part2.Name:
				return false
			}
		case part1.Name == "" && part2.Name == "":
			switch {
			case part1.SliceIndex < part2.SliceIndex:
				return true
			case part1.SliceIndex > part2.SliceIndex:
				return false
			}
		}
		i++
		j++
	}

	if i == len(ptr1) && j < len(ptr2) {
		return true
	}

	return false
}

func sortPointers(ptrs []pointer) {
	sort.Stable(sort.Reverse(ascPointers(ptrs)))
}
