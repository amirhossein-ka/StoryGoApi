package redis

import (
	"fmt"
	"time"
)

// GenGuestRedisKey  => token:<token>:limit:<limit>:offset:<offset>:from:<time>:to:<time>:sort_by:<sort>
func GenGuestRedisKey(token, sort string, limit, offset int, from, to time.Time) string {
	var (
		fromStr        = "0" // zero is default for zero time (when no time set in query options)
		toStr          = "0"
		sortStr string = "def" // default, i guess
	)

	if !from.IsZero() {
		fromStr = from.Format(time.RFC3339)
	}
	if !to.IsZero() {
		toStr = to.Format(time.RFC3339)
	}

	if sort != "" {
		sortStr = sort
	}
	return fmt.Sprintf("token:%s:limit:%d:offset:%d:from:%s:to:%s:sort_by:%s", token, limit, offset, fromStr, toStr, sortStr)
}

// GenGuestListKey  => token:<token>
func GenGuestListKey(token string) string {
	return fmt.Sprintf("token:%s", token)
}

// GenUserRedisKey => user:id:limit:<LIMIT>:offset:<OFFSET>:from:<time>:to:<time>:sort_by:<sort>
func GenUserRedisKey(id, limit, offset int, from, to time.Time, sort string) string {
	var (
		fromStr        = "0"
		toStr          = "0"
		sortStr string = "def" // default, i guess
	)

	if !from.IsZero() {
		fromStr = from.Format(time.RFC3339)
	}
	if !to.IsZero() {
		toStr = to.Format(time.RFC3339)
	}
	if sort != "" {
		sortStr = sort
	}
	return fmt.Sprintf("user:%d:limit:%d:offset:%d:from:%s:to:%s:sort_by:%s", id, limit, offset, fromStr, toStr, sortStr)
}

// GenUserListKey => user:id
func GenUserListKey(id int) string {
	return fmt.Sprintf("user:%d", id)
}

// GenStoryKey => story:id
func GenStoryKey(id int) string {
	return fmt.Sprintf("story:%d", id)
}
