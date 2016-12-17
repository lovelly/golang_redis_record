package record

import (
	"fmt"
	"testing"

	"gopkg.in/redis.v4"
)

type RoomPlayer struct {
	Id   string //玩家唯一ID(服务器ID_玩家ID)
	Name string //名字
	Lv   int    //等级

	//不导出字段
	record *RdsRecord
}

func RdsCachePlayerKey(playerId string) string {
	return fmt.Sprintf("{player:%v}", playerId)
}
func NewPlayer(db *redis.Client, id string) *RoomPlayer {
	keyname := RdsCachePlayerKey(id)
	player := &RoomPlayer{record: NewResRecod(0, keyname, db)}
	player.Id = id
	player.SetFlg(0)
	return player
}

func (r *RoomPlayer) GetId() string {
	return r.Id
}

func (r *RoomPlayer) SetName(v string) {
	r.SetFlg(1)
	r.Name = v
}
func (r *RoomPlayer) GetName() string {
	return r.Name
}

func (r *RoomPlayer) SetLv(v int) {
	r.SetFlg(2)
	r.Lv = v
}
func (r *RoomPlayer) GetLv() int {
	return r.Lv
}

func (r *RoomPlayer) SetFlg(index uint64) {
	r.record.SetFlg(index)
}

func (r *RoomPlayer) ClearFlg() {
	r.record.ClearFlg()
}

func (r *RoomPlayer) Update() {
	r.record.Update(r)
}

func (r *RoomPlayer) LoadFromRds() bool {
	return r.record.LoadFromRds(r)
}

func (r *RoomPlayer) Delete() {
	r.record.Delete()
}

func (r *RoomPlayer) SetPlayerAttr() {
	r.SetName("tom............")
	r.SetLv(99)
}

var Rds = redis.NewClient(&redis.Options{
	Addr:     "192.168.199.156:6380",
	Password: "",
	DB:       0,
})

var player = NewPlayer(Rds, "99999999999")

func Test_Update(t *testing.T) {
	player.SetPlayerAttr()
	player.Update()
}

func Test_Load(t *testing.T) {
	player := NewPlayer(Rds, "99999999999")
	ok := player.LoadFromRds()
	_ = ok
}

// func Test_Delete(t *testing.T) {
// 	player.Delete()
// }
