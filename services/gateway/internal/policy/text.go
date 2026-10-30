package policy

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

var urlRe = regexp.MustCompile(`https?://[^\s]+`)

type Contact struct {
	Number       string
	Lname        string
	Fname        string
	Sname        string
	Param1       string
	Param2       string
	Param3       string
	OperatorName string
}

func Personalize(text string, c Contact) string {
	repl := strings.NewReplacer(
		"{Фамилия}", c.Lname,
		"{Имя}", c.Fname,
		"{Отчество}", c.Sname,
		"{Параметр1}", c.Param1,
		"{Параметр2}", c.Param2,
		"{Параметр3}", c.Param3,
		"{lname}", c.Lname,
		"{fname}", c.Fname,
		"{sname}", c.Sname,
	)
	return repl.Replace(text)
}

func Transliterate(text string) string {
	var b strings.Builder
	for _, r := range text {
		if lat, ok := ruToLat[r]; ok {
			b.WriteString(lat)
			continue
		}
		b.WriteRune(r)
	}
	return b.String()
}

func Shorten(text, publicBase string, links map[string]string) string {
	return urlRe.ReplaceAllStringFunc(text, func(found string) string {
		code, ok := links[found]
		if !ok {
			return found
		}
		base := strings.TrimRight(publicBase, "/")
		return base + "/s/" + code
	})
}

func ExtractURLs(text string) []string {
	return urlRe.FindAllString(text, -1)
}

func Analyze(text string) (parts int, unicode bool) {
	if text == "" {
		return 0, false
	}
	unicode = !isGSM(text)
	n := utf8.RuneCountInString(text)
	single, multi := 160, 153
	if unicode {
		single, multi = 70, 67
	}
	if n <= single {
		return 1, unicode
	}
	parts = (n + multi - 1) / multi
	return parts, unicode
}

func isGSM(text string) bool {
	for _, r := range text {
		if !gsm[r] {
			return false
		}
	}
	return true
}

var gsm = func() map[rune]bool {
	const basic = "@£$¥èéùìòÇ\nØø\rÅåΔ_ΦΓΛΩΠΨΣΘΞÆæßÉ !\"#¤%&'()*+,-./0123456789:;<=>?¡ABCDEFGHIJKLMNOPQRSTUVWXYZÄÖÑÜ§¿abcdefghijklmnopqrstuvwxyzäöñüà"
	set := map[rune]bool{}
	for _, r := range basic {
		set[r] = true
	}
	for _, r := range "^{}\\[~]|€" {
		set[r] = true
	}
	return set
}()

var ruToLat = map[rune]string{
	'а': "a", 'б': "b", 'в': "v", 'г': "g", 'д': "d", 'е': "e", 'ё': "e",
	'ж': "zh", 'з': "z", 'и': "i", 'й': "y", 'к': "k", 'л': "l", 'м': "m",
	'н': "n", 'о': "o", 'п': "p", 'р': "r", 'с': "s", 'т': "t", 'у': "u",
	'ф': "f", 'х': "h", 'ц': "ts", 'ч': "ch", 'ш': "sh", 'щ': "sch",
	'ъ': "", 'ы': "y", 'ь': "", 'э': "e", 'ю': "yu", 'я': "ya",
	'А': "A", 'Б': "B", 'В': "V", 'Г': "G", 'Д': "D", 'Е': "E", 'Ё': "E",
	'Ж': "Zh", 'З': "Z", 'И': "I", 'Й': "Y", 'К': "K", 'Л': "L", 'М': "M",
	'Н': "N", 'О': "O", 'П': "P", 'Р': "R", 'С': "S", 'Т': "T", 'У': "U",
	'Ф': "F", 'Х': "H", 'Ц': "Ts", 'Ч': "Ch", 'Ш': "Sh", 'Щ': "Sch",
	'Ъ': "", 'Ы': "Y", 'Ь': "", 'Э': "E", 'Ю': "Yu", 'Я': "Ya",
}
