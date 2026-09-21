package rom

import "testing"

func TestHelloName(t *testing.T) {
	tests := []struct {
		name      string
		input     byte
		isFastRom bool
		mapValue  byte
		isValid   bool
	}{
		{
			name:      "slow LoROM",
			input:     0b00100000,
			isFastRom: false,
			mapValue:  0b0000,
			isValid:   true,
		},
		{
			name:      "fast HiROM",
			input:     0b00110001,
			isFastRom: true,
			mapValue:  0b0001,
			isValid:   true,
		},
		{
			name:      "invalid fixed bits",
			input:     0b00000000,
			isFastRom: false,
			mapValue:  0b0000,
			isValid:   false,
		},
		{
			name:      "fast ExHiROM with S-DD1",
			input:     0b00110010,
			isFastRom: true,
			mapValue:  0b0010,
			isValid:   true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := NewHeaderMapInfo(test.input)

			if got.Raw != test.input {
				t.Errorf("Raw = 0x%02X, want 0x%02X", got.Raw, test.input)
			}
			if got.MapMode != test.mapValue {
				t.Errorf("MapMode = 0x%02X, want 0x%02X", got.MapMode, test.mapValue)
			}
			if got.MapValue != test.mapValue {
				t.Errorf("MapValue = 0x%02X, want 0x%02X", got.MapValue, test.mapValue)
			}
			if got.IsFastRom != test.isFastRom {
				t.Errorf("IsFastRom = %t, want %t", got.IsFastRom, test.isFastRom)
			}
			if got.IsValid != test.isValid {
				t.Errorf("IsValid = %t, want %t", got.IsValid, test.isValid)
			}
		})
	}
}

func TestDeveloperName(t *testing.T) {
	tests := []struct {
		name string
		id   byte
		want string
	}{
		{name: "Capcom", id: 0x08, want: "0x08 Capcom"},
		{name: "extended license code", id: 0x33, want: "0x33 Extended license code"},
		{name: "unknown", id: 0xFF, want: "0xFF Unknown developer"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			header := SNESHeader{DeveloperId: test.id}
			if got := header.DeveloperName(); got != test.want {
				t.Errorf("DeveloperName() = %q, want %q", got, test.want)
			}
		})
	}
}
