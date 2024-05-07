package hw10programoptimization

import (
	"bufio"
	"encoding/json"
	"io"
	"strings"
)

type User struct {
	Email string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	scanner := bufio.NewScanner(r)
	result := make(DomainStat)

	var user User
	var err error
	var countDomain string

	for scanner.Scan() {
		if err = json.Unmarshal(scanner.Bytes(), &user); err != nil {
			return nil, err
		}

		if strings.HasSuffix(user.Email, domain) && strings.Contains(user.Email, "@") {
			countDomain = strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])
			result[countDomain]++
		}
	}
	return result, nil
}
