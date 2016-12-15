package gorecord

import (
	"reflect"
	"strconv"

	redis "gopkg.in/redis.v4"
)

//redis hset的封装成对象

type RdsRecord struct {
	dbIndex int //在哪个数据库
	keyname string
	db      *redis.Client //
	data    map[string]interface{}
	flg     map[string]bool
}

func NewResRecod(index int, name string, db *redis.Client) *RdsRecord {
	v := &RdsRecord{dbIndex: index, keyname: name, db: db}
	v.data = make(map[string]interface{})
	v.flg = make(map[string]bool)
	return v
}

func (r *RdsRecord) Setflog(key string) {
	r.flg[key] = true
}

//获取设置float64
func (r *RdsRecord) getFloat64(name string) (float64, bool) {
	v, ok := r.data[name]
	if !ok {
		return 0, false
	}

	fv, ok1 := v.(float64)
	if !ok1 {
		return 0, false
	}
	return fv, true
}

func (r *RdsRecord) setFloat64(name string, v float64) {
	r.data[name] = v
	r.flg[name] = true
}

//获取设置int
func (r *RdsRecord) getInt(name string) (int, bool) {
	v, ok := r.data[name]
	if !ok {
		return 0, false
	}

	iv, ok1 := v.(int)
	if !ok1 {
		return 0, false
	}
	return iv, true
}

func (r *RdsRecord) setInt(name string, v int) {
	r.data[name] = v
	r.flg[name] = true
}

//获取设置int64
func (r *RdsRecord) getInt64(name string) (int64, bool) {
	v, ok := r.data[name]
	if !ok {
		return 0, false
	}

	iv, ok1 := v.(int64)
	if !ok1 {
		return 0, false
	}
	return iv, true
}

func (r *RdsRecord) setInt64(name string, v int64) {
	r.data[name] = v
	r.flg[name] = true
}

//获取设置strng
func (r *RdsRecord) getString(name string) (string, bool) {
	v, ok := r.data[name]
	if !ok {
		return "", false
	}

	sv, ok1 := v.(string)
	if !ok1 {
		return "", false
	}
	return sv, true
}

func (r *RdsRecord) setString(name string, v string) {
	r.data[name] = v
	r.flg[name] = true
}

func (r *RdsRecord) Update() {
	upmap := make(map[string]string)
	for k, flg := range r.flg {
		if flg {
			v := r.data[k]
			switch reflect.TypeOf(v).Kind() {
			case reflect.String:
				upmap[k] = v.(string)
			case reflect.Int:
				upmap[k] = strconv.Itoa(v.(int))
			case reflect.Int64:
				upmap[k] = strconv.FormatInt(v.(int64), 10)
			case reflect.Float64:
				upmap[k] = strconv.FormatFloat(v.(float64), 'f', -1, 64)
			case reflect.Bool:
				upmap[k] = strconv.FormatBool(v.(bool))
			}
		}
	}
	ok, err := r.db.HMSet(r.keyname, upmap).Result()
	_, _ = ok, err
}

func (r *RdsRecord) LoadFromRds(key string, tv reflect.Value) {
	mapobj, err := r.db.HGetAll(r.keyname).Result()
	if err != nil {
		return
	}
	for k, strval := range mapobj {
		field := tv.FieldByName(k)
		switch field.Type().Kind() {
		case reflect.String:
			r.data[k] = strval
		case reflect.Int:
			iv, err := strconv.Atoi(strval)
			if err != nil {
				continue
			}
			r.data[k] = iv
		case reflect.Int64:
			i64v, err := strconv.ParseInt(strval, 10, 64)
			if err != nil {
				continue
			}
			r.data[k] = i64v
		case reflect.Float64:
			fv, err := strconv.ParseFloat(strval, 64)
			if err != nil {
				continue
			}
			r.data[k] = fv
		case reflect.Bool:
			bv, err := strconv.ParseBool(strval)
			if err != nil {
				continue
			}
			r.data[k] = bv
		}
	}
}

func (r *RdsRecord) LoadFromMap(obj map[string]interface{}) {
	r.data = obj
}

func (r *RdsRecord) Delete() {
	r.db.HDel(r.keyname)
}
