package sqlcraft

import "errors"

var ErrDuplicatedOption = errors.New("cannot have duplicated options")

type optionKey string

func (o optionKey) IsInList(list optionKeys) bool {
	for _, item := range list {
		if item == o {
			return true
		}
	}

	return false
}

const (
	returning optionKey = "RETURNING"
	where     optionKey = "WHERE"
	order     optionKey = "ORDER"
	limit     optionKey = "LIMIT"
	offset    optionKey = "OFFSET"
)

type optionKeys []optionKey

type options struct {
	sql   string
	args  []any
	key   optionKey
	order uint
}

type Option func(option *options) error
