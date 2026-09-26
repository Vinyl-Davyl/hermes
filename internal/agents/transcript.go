package agents

import (
	"bufio"
	"encoding/json"
	"os"
	"strings"
)

func readJSONLAsMarkdown(path, agent string) (Transcript, error) {
	f, err := os.Open(path)
	if err != nil {
		return Transcript{}, err
	}
	defer f.Close()

	var b strings.Builder
	sc := bufio.NewScanner(f)
	buf := make([]byte, 0, 1024*64)
	sc.Buffer(buf, 8*1024*1024)

	msgs := 0
	title := ""
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		role, text := parseJSONLLine(line)
		if text == "" {
			continue
		}
		if title == "" && role == "user" {
			title = truncate(text, 80)
		}
		msgs++
		if msgs <= 80 {
			b.WriteString("### ")
			b.WriteString(strings.ToUpper(role))
			b.WriteString("\n\n")
			b.WriteString(text)
			b.WriteString("\n\n")
		}
	}
	if err := sc.Err(); err != nil {
		return Transcript{}, err
	}
	if msgs > 80 {
		b.WriteString("_…truncated after 80 messages. Full session stays on disk at the source path._\n")
	}
	if title == "" {
		title = firstLineHint(path)
	}
	return Transcript{
		Agent:    agent,
		Title:    title,
		Path:     path,
		Markdown: strings.TrimSpace(b.String()),
		Messages: msgs,
	}, nil
}

func parseJSONLLine(line string) (role, text string) {
	var raw map[string]any
	if err := json.Unmarshal([]byte(line), &raw); err != nil {
		return "", ""
	}
	role = stringField(raw, "role", "type", "kind")
	switch strings.ToLower(role) {
	case "human", "user_message", "user":
		role = "user"
	case "assistant_message", "ai", "model", "assistant":
		role = "assistant"
	case "system", "tool", "tool_result", "tool_use":
		return role, ""
	default:
		if role == "" {
			role = "message"
		}
	}
	text = stringField(raw, "text", "content", "message", "prompt")
	if text == "" {
		if msg, ok := raw["message"].(map[string]any); ok {
			role = firstNonEmpty(stringField(msg, "role"), role)
			text = stringField(msg, "content", "text")
			if text == "" {
				text = anyText(msg["content"])
			}
		}
	}
	if text == "" {
		text = anyText(raw["content"])
	}
	if strings.HasPrefix(text, "[") || strings.HasPrefix(text, "{") {
		text = flattenContent(text)
	}
	return role, strings.TrimSpace(text)
}

func flattenContent(s string) string {
	var parts []any
	if err := json.Unmarshal([]byte(s), &parts); err == nil {
		var b strings.Builder
		for _, p := range parts {
			switch v := p.(type) {
			case string:
				b.WriteString(v)
			case map[string]any:
				if t := stringField(v, "text", "content"); t != "" {
					b.WriteString(t)
					b.WriteByte('\n')
				}
			}
		}
		return strings.TrimSpace(b.String())
	}
	var obj map[string]any
	if err := json.Unmarshal([]byte(s), &obj); err == nil {
		return stringField(obj, "text", "content")
	}
	return s
}

func stringField(m map[string]any, keys ...string) string {
	for _, k := range keys {
		switch v := m[k].(type) {
		case string:
			if strings.TrimSpace(v) != "" {
				return v
			}
		}
	}
	return ""
}

func anyText(v any) string {
	switch x := v.(type) {
	case string:
		return x
	case []any:
		var b strings.Builder
		for _, p := range x {
			if m, ok := p.(map[string]any); ok {
				if t := stringField(m, "text", "content"); t != "" {
					b.WriteString(t)
					b.WriteByte('\n')
				}
			}
			if s, ok := p.(string); ok {
				b.WriteString(s)
			}
		}
		return strings.TrimSpace(b.String())
	case map[string]any:
		return stringField(x, "text", "content")
	default:
		return ""
	}
}

func firstNonEmpty(xs ...string) string {
	for _, x := range xs {
		if strings.TrimSpace(x) != "" {
			return x
		}
	}
	return ""
}
