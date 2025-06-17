package location

import (
	"testing"
)

func TestConvertToPinyin(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"empty string", "", ""},
		{"single character", "北", "bei"},
		{"multiple characters", "北京市", "beijingshi"},
		{"with special characters", "重-庆市", "zhongqingshi"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ConvertToPinyin(tt.input); got != tt.want {
				t.Errorf("ConvertToPinyin() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalculateBestMatch(t *testing.T) {
	candidates := []string{"北京市", "上海市", "广州市"}
	tests := []struct {
		name      string
		input     string
		wantMatch string
		wantRatio float64
	}{
		{"exact match", "北京市", "北京市", 1.0},
		{"similar match", "北京", "北京市", 0.7},
		{"no match", "invalid", "", 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMatch, gotRatio := calculateBestMatch(tt.input, candidates)
			if gotMatch != tt.wantMatch {
				t.Errorf("calculateBestMatch() gotMatch = %v, want %v", gotMatch, tt.wantMatch)
			}
			if gotRatio < tt.wantRatio {
				t.Errorf("calculateBestMatch() gotRatio = %v, want at least %v", gotRatio, tt.wantRatio)
			}
		})
	}
}

func TestMatchProvince(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantMatch string
		wantRatio float64
	}{
		{"exact province match", "北京市", "北京市", 1.0},
		{"similar province match", "广东", "广东省", 0.6},
		{"no match", "invalid", "", 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMatch, gotRatio := MatchProvince(tt.input)
			if gotMatch != tt.wantMatch {
				t.Errorf("MatchProvince() gotMatch = %v, want %v", gotMatch, tt.wantMatch)
			}
			if gotRatio < tt.wantRatio {
				t.Errorf("MatchProvince() gotRatio = %v, want at least %v", gotRatio, tt.wantRatio)
			}
		})
	}
}

func TestMatchCity(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		province  string
		wantMatch string
		wantRatio float64
	}{
		{"exact city match", "平谷", "北京市", "平谷", 1.0},
		{"similar city match", "平", "北京市", "平谷", 0.6},
		{"no match", "invalid", "北京市", "", 0.0},
		{"invalid province", "平谷", "invalid", "", 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMatch, gotRatio := MatchCity(tt.input, tt.province)
			if gotMatch != tt.wantMatch {
				t.Errorf("MatchCity() gotMatch = %v, want %v", gotMatch, tt.wantMatch)
			}
			if gotRatio < tt.wantRatio {
				t.Errorf("MatchCity() gotRatio = %v, want at least %v", gotRatio, tt.wantRatio)
			}
		})
	}
}

func TestFuzzyMatch(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantMatch string
		wantRatio float64
	}{
		{"exact province match", "北京市", "北京市", 1.0},
		{"similar province match", "北京", "北京市 北京", 0.8},
		{"exact city match", "北京 平谷", "北京市", 0.5},
		{"similar city match", "伤害", "上海市", 0.7},
		{"no match", "invalid", "", 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMatch, gotRatio := FuzzyMatch(tt.input)
			if gotMatch != tt.wantMatch {
				t.Errorf("FuzzyMatch() gotMatch = %v, want %v", gotMatch, tt.wantMatch)
			}
			if gotRatio < tt.wantRatio {
				t.Errorf("FuzzyMatch() gotRatio = %v, want at least %v", gotRatio, tt.wantRatio)
			}
		})
	}
}
