package policy

type op struct {
	code int
	name string
}

func GuessOperator(number string) (code int, name string) {
	if len(number) >= 4 && number[0] == '7' {
		def := number[1:4]
		if found, ok := ruDEF[def]; ok {
			return found.code, found.name
		}
	}
	return 0, "OTHER"
}

var ruDEF = map[string]op{
	"900": {6, "t2"}, "901": {6, "t2"}, "902": {6, "t2"}, "904": {6, "t2"},
	"908": {6, "t2"}, "950": {6, "t2"}, "951": {6, "t2"}, "952": {6, "t2"},
	"953": {6, "t2"}, "958": {6, "t2"}, "977": {6, "t2"}, "991": {6, "t2"},
	"992": {6, "t2"}, "993": {6, "t2"}, "994": {6, "t2"}, "995": {6, "t2"},
	"996": {6, "t2"}, "999": {6, "t2"},
	"910": {4, "MTS"}, "911": {4, "MTS"}, "912": {4, "MTS"}, "913": {4, "MTS"},
	"914": {4, "MTS"}, "915": {4, "MTS"}, "916": {4, "MTS"}, "917": {4, "MTS"},
	"918": {4, "MTS"}, "919": {4, "MTS"}, "980": {4, "MTS"}, "981": {4, "MTS"},
	"982": {4, "MTS"}, "983": {4, "MTS"}, "984": {4, "MTS"}, "985": {4, "MTS"},
	"986": {4, "MTS"}, "987": {4, "MTS"}, "988": {4, "MTS"}, "989": {4, "MTS"},
	"903": {5, "BEELINE"}, "905": {5, "BEELINE"}, "906": {5, "BEELINE"},
	"909": {5, "BEELINE"}, "960": {5, "BEELINE"}, "961": {5, "BEELINE"},
	"962": {5, "BEELINE"}, "963": {5, "BEELINE"}, "964": {5, "BEELINE"},
	"965": {5, "BEELINE"}, "966": {5, "BEELINE"}, "967": {5, "BEELINE"},
	"968": {5, "BEELINE"}, "969": {5, "BEELINE"},
	"920": {1, "MEGAFON"}, "921": {1, "MEGAFON"}, "922": {1, "MEGAFON"},
	"923": {1, "MEGAFON"}, "924": {1, "MEGAFON"}, "925": {1, "MEGAFON"},
	"926": {1, "MEGAFON"}, "927": {1, "MEGAFON"}, "928": {1, "MEGAFON"},
	"929": {1, "MEGAFON"}, "930": {1, "MEGAFON"}, "931": {1, "MEGAFON"},
	"932": {1, "MEGAFON"}, "933": {1, "MEGAFON"}, "934": {1, "MEGAFON"},
	"936": {1, "MEGAFON"}, "937": {1, "MEGAFON"}, "938": {1, "MEGAFON"},
	"939": {1, "MEGAFON"},
}
