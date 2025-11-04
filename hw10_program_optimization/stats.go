package hw10programoptimization

import (
	"bufio"
	"errors"
	"io"
	"strings"

	"github.com/buger/jsonparser" //nolint:depguard
)

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	return countDomains(r, domain)
}

func countDomains(r io.Reader, domain string) (DomainStat, error) {
	result := make(DomainStat)

	scanner := bufio.NewScanner(r)
	for i := 0; scanner.Scan(); i++ {
		email, err := jsonparser.GetString(scanner.Bytes(), "Email")
		if err != nil {
			return nil, err
		}

		if !strings.Contains(email, "@") {
			return nil, errors.New("incorrect email")
		}

		builder := strings.Builder{}
		builder.WriteString(".")
		builder.WriteString(domain)
		if strings.HasSuffix(email, builder.String()) {
			domain := strings.ToLower(strings.SplitN(email, "@", 2)[1])
			result[domain]++
		}
	}

	return result, nil
}
