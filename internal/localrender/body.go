package localrender

import (
	"encoding/json"
	"regexp"
	"strings"

	"github.com/geekjourneyx/md2wechat-skill/internal/layoutcatalog"
)

type field struct{ key, value string }

// fieldSet keeps body fields in document order; repeated keys are allowed.
type fieldSet []field

var genericFieldKey = regexp.MustCompile(`^[A-Za-z][A-Za-z0-9_.-]*$`)

// get returns the last non-empty value for key.
func (f fieldSet) get(keys ...string) string {
	for _, key := range keys {
		for i := len(f) - 1; i >= 0; i-- {
			if f[i].key == key && f[i].value != "" {
				return f[i].value
			}
		}
	}
	return ""
}

func (f fieldSet) all(key string) []string {
	var values []string
	for _, item := range f {
		if item.key == key && item.value != "" {
			values = append(values, item.value)
		}
	}
	return values
}

// list splits a "a | b | c" value into its non-empty parts.
func (f fieldSet) list(keys ...string) []string {
	return splitList(f.get(keys...), "|")
}

func splitList(value, sep string) []string {
	var parts []string
	for _, part := range strings.Split(value, sep) {
		if part = strings.TrimSpace(part); part != "" {
			parts = append(parts, part)
		}
	}
	return parts
}

// splitCells splits a row and keeps empty cells so column positions stay put.
func splitCells(line, sep string) []string {
	cells := strings.Split(line, sep)
	for i := range cells {
		cells[i] = strings.TrimSpace(cells[i])
	}
	return cells
}

func cell(cells []string, i int) string {
	if i < len(cells) {
		return cells[i]
	}
	return ""
}

func declaredFields(block *layoutcatalog.Block) map[string]bool {
	if block.Spec == nil {
		return nil
	}
	return block.Spec.DeclaredFieldNames()
}

// parseFieldLine accepts "key: value" when key is declared (or looks like a
// field name for modules the catalog does not know).
func parseFieldLine(line string, declared map[string]bool) (string, string, bool) {
	idx := strings.Index(line, ":")
	if idx <= 0 {
		return "", "", false
	}
	key := strings.TrimSpace(line[:idx])
	if declared != nil {
		if !declared[key] {
			return "", "", false
		}
	} else if !genericFieldKey.MatchString(key) {
		return "", "", false
	}
	value := strings.TrimSpace(line[idx+1:])
	value = strings.Trim(value, `"`)
	return key, value, true
}

// parseFields reads key: value lines; anything else is returned as rest.
func parseFields(body []string, declared map[string]bool) (fieldSet, []string) {
	var fields fieldSet
	var rest []string
	for _, raw := range body {
		line := strings.TrimSpace(strings.TrimRight(raw, "\r"))
		if line == "" {
			continue
		}
		if key, value, ok := parseFieldLine(line, declared); ok {
			fields = append(fields, field{key, value})
			continue
		}
		rest = append(rest, line)
	}
	return fields, rest
}

// parseRows splits data rows by delimiter; declared field lines are pulled out.
func parseRows(block *layoutcatalog.Block) ([][]string, fieldSet) {
	delimiter := "|"
	if block.Spec != nil && block.Spec.Rows != nil && block.Spec.Rows.Delimiter != "" {
		delimiter = block.Spec.Rows.Delimiter
	}
	declared := map[string]bool{}
	if block.Spec != nil && block.Spec.Fields != nil {
		declared = block.Spec.DeclaredFieldNames()
	}
	var rows [][]string
	var fields fieldSet
	for _, raw := range block.Body {
		line := strings.TrimSpace(strings.TrimRight(raw, "\r"))
		if line == "" {
			continue
		}
		if key, value, ok := parseFieldLine(line, declared); ok && len(declared) > 0 {
			fields = append(fields, field{key, value})
			continue
		}
		rows = append(rows, splitCells(line, delimiter))
	}
	return rows, fields
}

func parseJSON(body []string) any {
	var value any
	if err := json.Unmarshal([]byte(strings.Join(body, "\n")), &value); err != nil {
		return nil
	}
	return value
}

func jsonString(obj map[string]any, keys ...string) string {
	for _, key := range keys {
		if v, ok := obj[key]; ok {
			switch t := v.(type) {
			case string:
				if t != "" {
					return t
				}
			case float64, bool:
				b, _ := json.Marshal(t)
				return string(b)
			}
		}
	}
	return ""
}

func jsonStrings(obj map[string]any, key string) []string {
	items, _ := obj[key].([]any)
	var out []string
	for _, item := range items {
		if s, ok := item.(string); ok && strings.TrimSpace(s) != "" {
			out = append(out, s)
		}
	}
	return out
}

func jsonObjects(value any) []map[string]any {
	items, _ := value.([]any)
	var out []map[string]any
	for _, item := range items {
		if obj, ok := item.(map[string]any); ok {
			out = append(out, obj)
		}
	}
	return out
}

type turn struct {
	speaker string
	text    string
	group   int // blank lines start a new group
}

func parseDialogue(body []string) []turn {
	var turns []turn
	group := 0
	for _, raw := range body {
		line := strings.TrimSpace(strings.TrimRight(raw, "\r"))
		if line == "" {
			if len(turns) > 0 {
				group++
			}
			continue
		}
		idx := strings.Index(line, ":")
		width := 1
		if full := strings.Index(line, "："); idx < 0 || (full >= 0 && full < idx) {
			idx, width = full, len("：")
		}
		if idx <= 0 {
			if len(turns) > 0 {
				turns[len(turns)-1].text += "\n" + line
			}
			continue
		}
		turns = append(turns, turn{
			speaker: strings.TrimSpace(line[:idx]),
			text:    strings.TrimSpace(line[idx+width:]),
			group:   group,
		})
	}
	return turns
}

type imageItem struct {
	src   string
	alt   string
	parts []string // text after the image, split by |
}

var imageLinePattern = regexp.MustCompile(`!\[([^\]]*)\]\(<?([^)\s>]+)>?(?:\s+"[^"]*")?\)`)

func parseImageLines(lines []string) []imageItem {
	var items []imageItem
	for _, raw := range lines {
		line := strings.TrimSpace(strings.TrimRight(raw, "\r"))
		m := imageLinePattern.FindStringSubmatchIndex(line)
		if m == nil {
			continue
		}
		item := imageItem{alt: line[m[2]:m[3]], src: line[m[4]:m[5]]}
		rest := strings.TrimSpace(line[m[1]:])
		rest = strings.TrimSpace(strings.TrimPrefix(rest, "|"))
		if rest != "" {
			item.parts = splitCells(rest, "|")
		}
		items = append(items, item)
	}
	return items
}

// splitBody splits lines at the first standalone separator.
func splitBody(body []string, separator string) ([]string, []string, bool) {
	for i, raw := range body {
		if strings.TrimSpace(raw) == separator {
			return body[:i], body[i+1:], true
		}
	}
	return body, nil, false
}
