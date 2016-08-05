package mailfull

import (
	"bufio"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

// Errors for parameter.
var (
	ErrNotEnoughAliasUserTargets = errors.New("AliasUser: targets not enough")
)

// AliasUser represents a AliasUser.
type AliasUser struct {
	name    string
	targets []string
}

// AliasUserSlice attaches the methods of sort.Interface to []*AliasUser.
type AliasUserSlice []*AliasUser

func (p AliasUserSlice) Len() int           { return len(p) }
func (p AliasUserSlice) Less(i, j int) bool { return p[i].Name() < p[j].Name() }
func (p AliasUserSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// NewAliasUser creates a new AliasUser instance.
func NewAliasUser(name string, targets []string) (*AliasUser, error) {
	au := &AliasUser{}

	if err := au.setName(name); err != nil {
		return nil, err
	}

	if err := au.SetTargets(targets); err != nil {
		return nil, err
	}

	return au, nil
}

// setName sets the name.
func (au *AliasUser) setName(name string) error {
	if !validAliasUserName(name) {
		return ErrInvalidAliasUserName
	}

	au.name = name

	return nil
}

// Name returns name.
func (au *AliasUser) Name() string {
	return au.name
}

// SetTargets sets targets.
func (au *AliasUser) SetTargets(targets []string) error {
	if len(targets) < 1 {
		return ErrNotEnoughAliasUserTargets
	}

	for _, target := range targets {
		if !validAliasUserTarget(target) {
			return ErrInvalidAliasUserTarget
		}
	}

	au.targets = targets

	return nil
}

// Targets returns targets.
func (au *AliasUser) Targets() []string {
	return au.targets
}

// AliasUsers returns a AliasUser slice.
func (r *Repository) AliasUsers(domainName string) ([]*AliasUser, error) {
	domain, err := r.Domain(domainName)
	if err != nil {
		return nil, err
	}
	if domain == nil {
		return nil, ErrDomainNotExist
	}

	file, err := os.Open(filepath.Join(r.DirMailDataPath, domainName, FileNameAliasUsers))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	aliasUsers := make([]*AliasUser, 0, 50)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		words := strings.Split(scanner.Text(), ":")
		if len(words) != 2 {
			return nil, ErrInvalidFormatAliasUsers
		}

		name := words[0]
		targets := strings.Split(words[1], ",")

		aliasUser, err := NewAliasUser(name, targets)
		if err != nil {
			return nil, err
		}

		aliasUsers = append(aliasUsers, aliasUser)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return aliasUsers, nil
}

// AliasUser returns a AliasUser of the input name.
func (r *Repository) AliasUser(domainName, aliasUserName string) (*AliasUser, error) {
	aliasUsers, err := r.AliasUsers(domainName)
	if err != nil {
		return nil, err
	}

	for _, aliasUser := range aliasUsers {
		if aliasUser.Name() == aliasUserName {
			return aliasUser, nil
		}
	}

	return nil, nil
}
