package record

import (
	"fmt"
	"reflect"
	"strconv"

	"changit.cn/contra/server/zaplogger"
	"github.com/uber-go/zap"

	redis "gopkg.in/redis.v4"
)

//redis hset的封装成对象

type RdsRecord struct {
	dbIndex int //在哪个数据库
	keyname string
	db      *redis.Client //
	flg     uint64
}

func NewResRecod(index int, name string, db *redis.Client) *RdsRecord {
	v := &RdsRecord{dbIndex: index, keyname: name, db: db}
	return v
}

func (r *RdsRecord) SetFlg(index uint64) {
	index = 1 << index
	r.flg = r.flg | index
}

func (r *RdsRecord) ClearFlg() {
	r.flg = 0
}

func (r *RdsRecord) Update(obj interface{}) {
	upmap := make(map[string]string)
	fmt.Println("obj  === ", obj, r.flg)
	tmp := reflect.ValueOf(obj)
	if tmp.Type().Kind() != reflect.Ptr {
		zaplogger.Error("at record obj not ptr .. ")
		return
	}

	value := tmp.Elem()
	types := value.Type()
	index := 0

	for r.flg != 0 {
		if (r.flg & 1) != 0 {
			v := value.Field(index)
			switch v.Type().Kind() {
			case reflect.String:
				upmap[types.Field(index).Name] = v.String()
			case reflect.Int:
				upmap[types.Field(index).Name] = strconv.Itoa(int(v.Int()))
			case reflect.Int64:
				upmap[types.Field(index).Name] = strconv.FormatInt(v.Int(), 10)
			case reflect.Float64:
				upmap[types.Field(index).Name] = strconv.FormatFloat(v.Float(), 'f', -1, 64)
			case reflect.Bool:
				upmap[types.Field(index).Name] = strconv.FormatBool(v.Bool())
			}
		}
		index += 1
		r.flg = r.flg >> 1
	}

	fmt.Println("redis ceche data is ", upmap)
	_, err := r.db.HMSet(r.keyname, upmap).Result()
	if err != nil {
		zaplogger.Debug("at update roomplayer error", zap.String("err:=", err.Error()))
	}
}

func (r *RdsRecord) LoadFromRds(obj interface{}) bool {
	mapobj, err := r.db.HGetAll(r.keyname).Result()
	fmt.Println("redis get  data is ", mapobj)
	if err != nil {
		return false
	}

	tmp := reflect.ValueOf(obj)
	if tmp.Type().Kind() != reflect.Ptr {
		zaplogger.Error("at record obj not ptr .. ")
		return false
	}

	value := tmp.Elem()
	types := value.Type()

	numfiled := value.NumField()
	if numfiled < 1 {
		return false
	}
	fideCount := uint64(numfiled - 1) // -1 is RdsRecord
	r.flg = (1 << fideCount) - 1
	index := 0
	for r.flg != 0 {
		v := value.Field(index)
		strval, ok := mapobj[types.Field(index).Name]
		if !ok {
			zaplogger.Error("at LoadFromRds field can't find field", zap.String("obj", types.Field(index).Name))
			return false
		}
		if !v.CanSet() {
			zaplogger.Error("at LoadFromRds field can'tset")
			return false
		}

		switch v.Type().Kind() {
		case reflect.String:
			v.SetString(strval)
		case reflect.Int:
			i64v, err := strconv.ParseInt(strval, 10, 64)
			if err == nil {
				v.SetInt(i64v)
			}

		case reflect.Int64:
			i64v, err := strconv.ParseInt(strval, 10, 64)
			if err == nil {
				v.SetInt(i64v)
			}

		case reflect.Float64:
			fv, err := strconv.ParseFloat(strval, 64)
			if err == nil {
				v.SetFloat(fv)
			}

		case reflect.Bool:
			bv, err := strconv.ParseBool(strval)
			if err == nil {
				v.SetBool(bv)
			}
		}
		index += 1
		r.flg = r.flg >> 1
	}
	return true
}

func (r *RdsRecord) Delete() {
	r.db.HDel(r.keyname)
}
