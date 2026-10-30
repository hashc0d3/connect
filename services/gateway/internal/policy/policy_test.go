package policy

import "testing"

func TestNormalizePhone(t *testing.T) {
	got, err := NormalizePhone("358403536665")
	if err != nil || got != "358403536665" {
		t.Fatalf("valid: %q %v", got, err)
	}
	if _, err := NormalizePhone("+79990000000"); err == nil {
		t.Fatal("plus must be rejected")
	}
	if _, err := NormalizePhone("12345"); err == nil {
		t.Fatal("short must be rejected")
	}
	got, err = NormalizePhone("7 (999) 000-00-00")
	if err != nil || got != "79990000000" {
		t.Fatalf("separators: %q %v", got, err)
	}
}

func TestAnalyze(t *testing.T) {
	parts, unicode := Analyze("Test Sms")
	if parts != 1 || unicode {
		t.Fatalf("gsm: %d %v", parts, unicode)
	}
	parts, unicode = Analyze("Привет")
	if parts != 1 || !unicode {
		t.Fatalf("ucs: %d %v", parts, unicode)
	}
}

func TestTranslitAndPersonalize(t *testing.T) {
	text := Personalize("{Имя}, ждём", Contact{Fname: "Анна"})
	if text != "Анна, ждём" {
		t.Fatal(text)
	}
	if got := Transliterate(text); got != "Anna, zhdem" {
		t.Fatal(got)
	}
}

func TestCallback(t *testing.T) {
	if err := ValidateCallback("https://client.example/hook"); err != nil {
		t.Fatal(err)
	}
	if err := ValidateCallback("http://127.0.0.1/hook"); err == nil {
		t.Fatal("private ip must fail")
	}
	if err := ValidateCallback(""); err != nil {
		t.Fatal(err)
	}
}
