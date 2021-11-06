package link

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

type Link struct {
	ID          uint64
	UserID      uint64 `db:"user_id"`
	Url         string
	Description string
	Tags        string
	CreatedAt   time.Time `db:"created_at"`
}

func New(userID uint64, url string) *Link {
	return &Link{0, userID, url, "", "", time.Now()}
}

func (l *Link) String() string {
	return fmt.Sprintf(
		"ID: %v,\nUserID: %v,\nURL: %q,\nDescription: %q,\nTags: %q,\nDateCreated: %v",
		l.ID,
		l.UserID,
		l.Url,
		l.Description,
		l.Tags,
		l.CreatedAt.Format(time.RFC3339),
	)
}

func (l *Link) Equals(inputLink *Link) bool {
	if l.ID != inputLink.ID {
		return false
	}
	if l.UserID != inputLink.UserID {
		return false
	}
	if l.Url != inputLink.Url {
		return false
	}
	if l.Description != inputLink.Description {
		return false
	}
	tags := l.GetTagsAsSlice()
	inputTags := inputLink.GetTagsAsSlice()
	if len(tags) != len(inputTags) {
		return false
	}
	sort.Strings(tags)
	sort.Strings(inputTags)
	if strings.Join(tags, "#") != strings.Join(inputTags, "#") {
		return false
	}

	if !l.CreatedAt.Equal(inputLink.CreatedAt) {
		return false
	}

	return true
}

func (l *Link) AddTag(tag string) {
	if strings.Index(l.Tags, tag) == -1 {
		if len(l.Tags) > 0 {
			l.Tags += "#" + tag
		} else {
			l.Tags += tag
		}
	}
}

func (l *Link) RemoveTag(tag string) {
	l.Tags = strings.Replace(l.Tags, tag, "", 1)
	l.Tags = strings.ReplaceAll(l.Tags, "##", "#")
	l.Tags = strings.TrimPrefix(l.Tags, "#")
	l.Tags = strings.TrimSuffix(l.Tags, "#")
}

func (l *Link) GetTagsAsSlice() []string {
	return strings.Split(l.Tags, "#")
}

func (l *Link) SetTagsAsSlice(tags []string) {
	l.Tags = strings.Join(tags, "#")
}
