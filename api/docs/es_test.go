package docs

import (
	"encoding/json"
	"strings"
	"testing"
)

type specOperation struct {
	Summary     string   `json:"summary"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
}

type spec struct {
	Info struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	} `json:"info"`
	Paths map[string]map[string]specOperation `json:"paths"`
}

func parseSpec(t *testing.T, doc string) spec {
	t.Helper()
	var s spec
	if err := json.Unmarshal([]byte(doc), &s); err != nil {
		t.Fatalf("cannot parse spec: %v", err)
	}
	return s
}

// TestSpanishDocCoversEveryOperation garantiza que cada operación del spec en
// inglés tenga traducción al español y que no queden traducciones huérfanas.
func TestSpanishDocCoversEveryOperation(t *testing.T) {
	english := parseSpec(t, SwaggerInfo.ReadDoc())
	spanish := parseSpec(t, SpanishDoc())

	if spanish.Info.Description != spanishInfoDescription || spanish.Info.Title != spanishInfoTitle {
		t.Errorf("info was not translated: %+v", spanish.Info)
	}

	operations := map[string]bool{}
	for path, methods := range english.Paths {
		for method, op := range methods {
			key := strings.ToUpper(method) + " " + path
			operations[key] = true

			tr, ok := spanishOperations[key]
			if !ok || strings.TrimSpace(tr.Summary) == "" || strings.TrimSpace(tr.Description) == "" {
				t.Errorf("missing Spanish summary/description for %q", key)
				continue
			}

			es := spanish.Paths[path][method]
			if es.Summary != tr.Summary || es.Description != tr.Description {
				t.Errorf("Spanish overlay not applied to %q", key)
			}
			for _, tag := range op.Tags {
				if _, ok := spanishTags[tag]; !ok {
					t.Errorf("missing Spanish tag for %q (used by %q)", tag, key)
				}
			}
		}
	}

	for key := range spanishOperations {
		if !operations[key] {
			t.Errorf("Spanish translation for %q does not match any documented operation", key)
		}
	}

	t.Logf("%d of %d operations translated to Spanish", len(spanishOperations), len(operations))
}
