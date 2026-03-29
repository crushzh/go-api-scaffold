package querybuilder

import (
	"fmt"
	"reflect"
	"strings"

	"gorm.io/gorm"
)

const tagName = "search"

type condition struct {
	Type   string
	Column string
	Table  string
}

// Apply 根据 struct 的 search tag 自动构建 GORM WHERE 条件
// 支持的 type: exact, contains, in, gt, gte, lt, lte
// 跳过零值字段
func Apply(db *gorm.DB, query any) *gorm.DB {
	v := reflect.ValueOf(query)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return db
	}

	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get(tagName)
		if tag == "" || tag == "-" {
			continue
		}
		fv := v.Field(i)
		if isZero(fv) {
			continue
		}

		cond := parseTag(tag)
		if cond.Column == "" {
			continue
		}
		col := cond.Column
		if cond.Table != "" {
			col = cond.Table + "." + col
		}
		val := getValue(fv)

		switch cond.Type {
		case "exact":
			db = db.Where(fmt.Sprintf("%s = ?", col), val)
		case "contains":
			db = db.Where(fmt.Sprintf("%s LIKE ?", col), "%"+fmt.Sprintf("%v", val)+"%")
		case "in":
			db = db.Where(fmt.Sprintf("%s IN ?", col), val)
		case "gt":
			db = db.Where(fmt.Sprintf("%s > ?", col), val)
		case "gte":
			db = db.Where(fmt.Sprintf("%s >= ?", col), val)
		case "lt":
			db = db.Where(fmt.Sprintf("%s < ?", col), val)
		case "lte":
			db = db.Where(fmt.Sprintf("%s <= ?", col), val)
		}
	}
	return db
}

func parseTag(tag string) condition {
	c := condition{Type: "exact"}
	for _, part := range strings.Split(tag, ";") {
		kv := strings.SplitN(strings.TrimSpace(part), ":", 2)
		if len(kv) != 2 {
			continue
		}
		switch strings.TrimSpace(kv[0]) {
		case "type":
			c.Type = strings.TrimSpace(kv[1])
		case "column":
			c.Column = strings.TrimSpace(kv[1])
		case "table":
			c.Table = strings.TrimSpace(kv[1])
		}
	}
	return c
}

func isZero(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Ptr, reflect.Interface:
		return v.IsNil()
	case reflect.String:
		return v.String() == ""
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Slice, reflect.Array:
		return v.Len() == 0
	default:
		return false
	}
}

func getValue(v reflect.Value) any {
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil
		}
		return v.Elem().Interface()
	}
	return v.Interface()
}
