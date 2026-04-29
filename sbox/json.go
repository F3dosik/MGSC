package sbox

import (
	"encoding/json"
	"fmt"
	"os"
)

type sboxJSON struct {
	Name string    `json:"name"`
	SBox [][]uint8 `json:"sbox"`
}

// LoadFromFile загружает S-блок из JSON файла.
// Ожидает матрицу 16×16.
func LoadFromFile(path string) (*SBox, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("sbox: не удалось прочитать файл %s: %w", path, err)
	}

	var raw sboxJSON
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("sbox: ошибка парсинга JSON: %w", err)
	}

	if len(raw.SBox) != 16 {
		return nil, fmt.Errorf("sbox: ожидается 16 строк, получено %d", len(raw.SBox))
	}

	flat := make([]uint8, 0, 256)
	for i, row := range raw.SBox {
		if len(row) != 16 {
			return nil, fmt.Errorf("sbox: строка %d: ожидается 16 элементов, получено %d", i, len(row))
		}
		flat = append(flat, row...)
	}

	return FromSlice(flat)
}
