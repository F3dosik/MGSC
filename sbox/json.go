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

// SaveToFile сохраняет S-блок в JSON-файл в формате, совместимом с LoadFromFile:
// объект с полями "name" и "sbox" (матрица 16×16).
func SaveToFile(sb *SBox, path, name string) error {
	t := sb.Table()
	rows := make([][]uint8, 16)
	for i := 0; i < 16; i++ {
		row := make([]uint8, 16)
		copy(row, t[i*16:(i+1)*16])
		rows[i] = row
	}

	data, err := json.MarshalIndent(sboxJSON{Name: name, SBox: rows}, "", "  ")
	if err != nil {
		return fmt.Errorf("sbox: ошибка сериализации JSON: %w", err)
	}

	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("sbox: не удалось записать %s: %w", path, err)
	}
	return nil
}
