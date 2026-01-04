package main

import (
	"strings"
	"testing"

	"github.com/micheldebree/retrospex/internal/dithering"
	"github.com/micheldebree/retrospex/internal/indexedimage"
	"github.com/micheldebree/retrospex/internal/pixels"
)

func TestKeysToString(t *testing.T) {
	// Test with map[string]struct{}
	m1 := map[string]struct{}{
		"a": {},
		"b": {},
		"c": {},
	}
	result1 := KeysToString(m1)
	// Check that all expected keys are present (order doesn't matter)
	expectedKeys := []string{"a", "b", "c"}
	resultKeys := strings.Split(result1, ", ")
	if len(resultKeys) != len(expectedKeys) {
		t.Errorf("Expected %d keys, got %d", len(expectedKeys), len(resultKeys))
	}

	keySet := make(map[string]bool)
	for _, key := range expectedKeys {
		keySet[key] = false
	}
	for _, key := range resultKeys {
		if _, exists := keySet[key]; !exists {
			t.Errorf("Unexpected key %q in result", key)
		}
		keySet[key] = true
	}
	for key, found := range keySet {
		if !found {
			t.Errorf("Missing key %q in result", key)
		}
	}

	// Test with map[string]int
	m2 := map[string]int{
		"x": 1,
		"y": 2,
		"z": 3,
	}
	result2 := KeysToString(m2)
	// Check if all keys are present
	keys := map[string]bool{
		"x": false,
		"y": false,
		"z": false,
	}
	for _, key := range strings.Split(result2, ", ") {
		keys[key] = true
	}
	for key, found := range keys {
		if !found {
			t.Errorf("Key %q not found in result %q", key, result2)
		}
	}

	// Test with RetrospecName map
	m3 := indexedimage.BitsPerPixel
	result3 := KeysToString(m3)
	if len(result3) == 0 {
		t.Error("Expected non-empty result for RetrospecName map")
	}

	// Test with real palettes
	result4 := KeysToString(pixels.C64Palettes)
	resultKeys4 := strings.Split(result4, ", ")
	expectedPaletteCount := len(pixels.C64Palettes)
	if len(resultKeys4) != expectedPaletteCount {
		t.Errorf("Expected %d palettes, got %d", expectedPaletteCount, len(resultKeys4))
	}

	// Test with real dither matrices
	result5 := KeysToString(dithering.DitherMatrices)
	resultKeys5 := strings.Split(result5, ", ")
	expectedMatrixCount := len(dithering.DitherMatrices)
	if len(resultKeys5) != expectedMatrixCount {
		t.Errorf("Expected %d dither matrices, got %d", expectedMatrixCount, len(resultKeys5))
	}
}
