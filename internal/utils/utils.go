package utils

import (
	"crypto/rand"
	"errors"
	"math/big"
	"strings"
)

func JoinStrings(sep string, str ...string) string {
	return strings.Join(str, sep)
}

var (
	rander = rand.Reader // random function

	// AlphaNum contains runes [abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789].
	AlphaNum = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	// Alpha contains runes [abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ].
	Alpha = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	// AlphaLowerNum contains runes [abcdefghijklmnopqrstuvwxyz0123456789].
	AlphaLowerNum = []rune("abcdefghijklmnopqrstuvwxyz0123456789")
	// AlphaUpperNum contains runes [ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789].
	AlphaUpperNum = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	// AlphaLower contains runes [abcdefghijklmnopqrstuvwxyz].
	AlphaLower = []rune("abcdefghijklmnopqrstuvwxyz")
	// AlphaUpper contains runes [ABCDEFGHIJKLMNOPQRSTUVWXYZ].
	AlphaUpper = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	// Numeric contains runes [0123456789].
	Numeric = []rune("0123456789")
)

// RuneSequence returns a random sequence using the defined allowed runes.
func RuneSequence(l int, allowedRunes []rune) (seq []rune, err error) {
	c := big.NewInt(int64(len(allowedRunes)))
	seq = make([]rune, l)

	for i := 0; i < l; i++ {
		r, err := rand.Int(rander, c)
		if err != nil {
			return seq, err
		}
		rn := allowedRunes[r.Uint64()]
		seq[i] = rn
	}

	return seq, nil
}

var (
	ErrNormilizeEmail = errors.New("failed normilize email")
)

func NormalizeEmail(email string) (string, error) {
	// all_lowercase
	normEmail := strings.ToLower(email)

	// validate email
	strs := strings.Split(normEmail, "@")
	if len(strs) < 2 {
		return "", ErrNormilizeEmail
	}

	user := strs[0]
	domain := strs[1]

	return user + "@" + domain, nil
}
