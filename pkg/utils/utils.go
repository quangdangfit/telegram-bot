package utils

import (
	"fmt"
	"net/url"
	"time"
)

func URLEncode(queryObj map[string]string) (string, error) {
	values := url.Values{}
	for k, v := range queryObj {
		values.Add(k, v)
	}
	strQuery := values.Encode()
	return strQuery, nil
}

func InList(list []interface{}, element interface{}) int {
	for i, e := range list {
		if element == e {
			return i
		}
	}

	return -1
}

func GetArchivedCollection(collection string, createdTime time.Time) string {
	// return collection + "_" + createdTime.Format("0601")
	return MonthPartitionName(collection, createdTime)
}

func MonthPartitionName(s string, t time.Time) string {
	m := ""
	if t.Month() < 10 {
		m = fmt.Sprintf("0%d", t.Month())
	} else {
		m = fmt.Sprintf("%d", t.Month())
	}

	return fmt.Sprintf("%s_%s_%d", s, m, t.Year())
}
