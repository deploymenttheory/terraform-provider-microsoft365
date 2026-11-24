package main

import (
	"regexp"
	"strings"
	"unicode"
)

const nameOK = "ok"

// bannedNames are generic helper names the guide rules out.
var bannedNames = map[string]bool{"matches": true, "settings": true, "requestBody": true, "tryClaim": true, "readbackOptions": true, "writableApp": true}

// legacyStateMapper matches other spellings of the state mapper, such as MapRemoteResourceStateToTerraform or mapRemoteState.
// Entity mappers such as MapRemoteOwnersToTerraform don't match and are checked against the map…ToState form instead.
var legacyStateMapper = regexp.MustCompile(`(?i)^mapRemote(Resource)?State`)

// dataSourceStateMapper matches data source mapper spellings. In a resource package they are left alone, because the data source
// name may belong to a data source mapper kept next to the resource.
var dataSourceStateMapper = regexp.MustCompile(`(?i)ToDataSource$`)

// verbFiles maps each convention verb to the file its helpers belong in.
var verbFiles = map[string]string{"construct": "construct.go", "map": "state.go", "find": "state.go", "waitFor": "crud.go", "lock": "crud.go"}

// checkName applies the verb + Graph entity + qualifier convention and returns nameOK or the reason the name fails.
// mapper is the package's standard state mapper name; other spellings of it are flagged for renaming.
func checkName(name, file, mapper string) string {
	if bannedNames[name] {
		return "banned generic name"
	}
	if isLegacyStateMapper(name, mapper) {
		return "non-standard state mapper: rename to " + mapper
	}
	words := splitWords(name)
	verb, rest := "", words
	if len(words) > 0 {
		verb, rest = strings.ToLower(words[0]), words[1:]
	}
	if verb == "wait" && len(rest) > 0 && rest[0] == "For" {
		verb, rest = "waitFor", rest[1:]
	}
	wantFile, known := verbFiles[verb]
	if !known {
		return "verb not in convention"
	}
	switch verb {
	case "map":
		if len(rest) < 3 || rest[len(rest)-2] != "To" || rest[len(rest)-1] != "State" {
			return "map… must name an entity and end ToState"
		}
	case "find":
		by := indexOf(rest, "By")
		if by < 1 || by == len(rest)-1 {
			return "find… must name an entity then By…"
		}
	default:
		if len(rest) < 2 {
			return "missing entity or qualifier"
		}
	}
	if file != wantFile {
		return verb + "… belongs in " + wantFile
	}
	return nameOK
}

// splitWords splits a camel-case identifier into words, keeping acronyms such as ID or URL together.
func splitWords(name string) []string {
	runes := []rune(name)
	var words []string
	start := 0
	for i := 1; i < len(runes); i++ {
		prev, cur := runes[i-1], runes[i]
		nextLower := i+1 < len(runes) && unicode.IsLower(runes[i+1])
		boundary := unicode.IsUpper(cur) && (unicode.IsLower(prev) || unicode.IsDigit(prev) || (unicode.IsUpper(prev) && nextLower))
		if boundary {
			words = append(words, string(runes[start:i]))
			start = i
		}
	}
	if len(runes) > 0 {
		words = append(words, string(runes[start:]))
	}
	return words
}

// indexOf returns the position of word in words, or -1.
func indexOf(words []string, word string) int {
	for i, w := range words {
		if w == word {
			return i
		}
	}
	return -1
}

// isLegacyStateMapper reports whether name is another spelling of the package's standard state mapper.
func isLegacyStateMapper(name, mapper string) bool {
	if name == mapper || !legacyStateMapper.MatchString(name) {
		return false
	}
	return mapper == dataSourceStateMapperName || !dataSourceStateMapper.MatchString(name)
}

// checkStateMapperFile returns nameOK when the state mapper is in state.go, or the reason it is not.
func checkStateMapperFile(name, file string) string {
	if file != "state.go" {
		return name + " belongs in state.go"
	}
	return nameOK
}
