// String Helper Extensions

package stringer

import (
	"strings"
	"unicode"
)

// these are the special characters which will use to split the
// sentences based on these characters
type specialCharacters []rune

func (scs *specialCharacters) contains(c rune) bool {
	for _, e := range *scs {
		if e == c {
			return true
			break
		}
	}
	return false
}

//var stringSpliter = rune("_"[0])
//var SpecialCharacters = specialCharacters{32, 45, 58, 95}                                         //=> [' ', '-', ':', '_']
var SpecialCharacters = specialCharacters{rune("_"[0]), rune(" "[0]), rune(":"[0]), rune("-"[0])} //=> [' ', '-', ':', '_']

// convert the given string into camecase string
// Camelize("abc def htk dd_ppp") => "AbcDefHtkDdPpp"
// Camelize("a b c_d") => "ABCD"

func Camelize(s string) (r string) {
	for _, e := range splitSentence(s) {
		r += strings.Title(e)
	}
	return
}

// convert the given string into plurals (by adding suffix 's')
// Pluralize("abc def htk dd_ppp") => "abc_def_htk_dd_ppps"
// Pluralize("a b c_d") => "a_b_c_ds"

func Pluralize(s string) string {
	if len(s) > 0 && string(s[len(s)-1]) != "s" {
		return s + "s"
	}
	return s
}

// convert the given string into underscored string
// Underscore("abc def htk dd_ppp") => "abc_def_htk_dd_ppp"
// Underscore("a b c_d") => "a_b_c_d"

func Underscore(s string) string {
	return strings.Join(splitSentence(s), "_")
}

// spliting the given sentence string into an slice of string ,
// assume the separators are SpecialCharacters
// splitSentence("abc def htk dd_ppp") => ["abc","def","htk","dd","ppp"]

func splitSentence(s string) (words []string) {
	word := ""
	sc := false
	for _, c := range s {
		sc = SpecialCharacters.contains(c)
		if sc && len(word) > 0 {
			words = append(words, word)
			word = ""
		} else if unicode.IsUpper(c) {
			if len(word) > 0 {
				words = append(words, word)
			}
			word = string(unicode.ToLower(c))
		} else if !sc {
			word += string(c)
		}
	}
	words = append(words, word)
	return
}

// check wether the string slice 'a' has 'b' in it list ?
func Contains(a []string, b string) bool {
	for _, e := range a {
		if e == b {
			return true
		}
	}
	return false
}

// remove fist occurance if 'b' in the string slice
func Remove(a []string, b string) []string {
	for i, e := range a {
		if e == b {
			a[i] = a[len(a)-1]
			a = a[0 : len(a)-1]
			break
		}
	}
	return a
}
