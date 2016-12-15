package test

import (
	"encoding/json"
	"fmt"
	"testing"

	"gorecord"

	"changit.cn/contra/server/net"

	"gopkg.in/redis.v4"
)

var obj = net.RoomPlayer{
	Id:             "9_9999999",
	Name:           "test player data test player data test player data test player data test player data test player data test player data test player data ",
	Lv:             99999999,
	AreaId:         11111111111111,
	RoomId:         2222222222222222,
	Camp:           1,
	Icon:           1,
	RankScore:      2,
	Prepare:        true,
	IsJoinInTeam:   true,
	DiffRankScore:  99999999999,
	HeroId:         1111111111,
	GunId1:         222222222222,
	GunId2:         33333333333333,
	FightForce:     444444444444444444,
	TeamId:         22222222222222222,
	IsMatch:        true,
	StartFightTime: 11111111111111111,
	RoomSvrId:      222222222222222222,
}

var mobj = map[string]interface{}{
	"Id":             "9_9999999",
	"Name":           "test player data test player data test player data test player data test player data test player data test player data test player data ",
	"Lv":             99999999,
	"AreaId":         11111111111111,
	"RoomId":         2222222222222222,
	"Camp":           1,
	"Icon":           1,
	"RankScore":      23,
	"Prepare":        true,
	"IsJoinInTeam":   true,
	"DiffRankScore":  99999999999,
	"HeroId":         1111111111,
	"GunId1":         222222222222,
	"GunId2":         33333333333333,
	"FightForce":     444444444444444444,
	"TeamId":         22222222222222222,
	"IsMatch":        true,
	"StartFightTime": 11111111111111111,
	"RoomSvrId":      222222222222222222,
}

var Rds = redis.NewClient(&redis.Options{
	Addr:     "192.168.199.156:6380",
	Password: "",
	DB:       0,
})

var keyname = "testewfrewgfrg"

var rd = gorecord.NewResRecod(0, keyname, Rds)

// func Test_new(t *testing.T) {
// 	for i := 1; i < 1000; i++ {
// 		v := gorecord.NewResRecod(0, keyname, Rdshp)
// 		_ = v
// 	}
// }

func setflg(rb *gorecord.RdsRecord) {
	for k, _ := range mobj {
		rb.Setflog(k)
		//rb.Setflog("RoomSvrId")
		//rb.Setflog("Lv")
	}
}

func Test_update(t *testing.T) {
	rd.LoadFromMap(mobj)
	for i := 1; i < 10000; i++ {
		setflg(rd)
		rd.Update()
	}
}

// func Test_load(t *testing.T) {
// 	for i := 1; i < 1000; i++ {
// 		rd.LoadFromRds(keyname, reflect.ValueOf(net.RoomPlayer{}))
// 	}
// }

// func Test_Delete(t *testing.T) {
// 	for i := 1; i < 1000; i++ {
// 		rd.Delete()
// 	}
// }

// func Test_All(t *testing.T) {
// 	for i := 1; i < 1000; i++ {
// 		rd := gorecord.NewResRecod(0, keyname, Rdshp)
// 		rd.Update()
// 		rd.LoadFromRds(keyname, reflect.ValueOf(net.RoomPlayer{}))
// 		//rd.Delete()
// 	}
// }

func Test_Old(t *testing.T) {
	for i := 1; i < 10000; i++ {
		v, err := json.Marshal(obj)
		if err != nil {
			fmt.Println("errro r;;;;; ", err)
			return
		}

		r, err := Rds.HSet("99999999999999999s", "roomplayer", string(v)).Result()
		_, _ = r, err
	}
}
