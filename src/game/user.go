package game

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"time"
	"unicode/utf8"

	"godori.com/db"
	"godori.com/getty"
	toClient "godori.com/packet/toClient"
	modelType "godori.com/util/constant/modelType"
	roomType "godori.com/util/constant/roomType"
	teamType "godori.com/util/constant/teamType"
	cFilter "godori.com/util/filter"
	cMath "godori.com/util/math"
	pix "godori.com/util/pix"
)

type UserData struct {
	Id           int
	Uid          string
	Uuid         string
	Name         string
	ClanName     string
	Rank         int
	Sex          int
	Level        int
	Exp          int
	MaxExp       int
	Coin         int
	Cash         int
	Point        int
	Win          int
	Lose         int
	Kill         int
	Death        int
	Assist       int
	Blast        int
	Rescue       int
	RescueCombo  int
	Survive      int
	Escape       int
	Graphics     string
	RedGraphics  string
	BlueGraphics string
	Memo         string
	LastChat     int
	Admin        int
}

type User struct {
	Model  int
	Index  int
	Client *getty.Client
	Room   *Room
	Place  int
	Alert  int
	Score  *Score
	Reward *Reward
	Character
	UserData *UserData
	GameData map[string]interface{}
}

var nextUserIndex int = 0
var Users map[*getty.Client]*User = make(map[*getty.Client]*User)

func NewUser(c *getty.Client, uid string, loginType int) (*User, bool) {
	if result, ok := db.GetUserByOAuth(uid, loginType); ok {
		nextUserIndex++
		user := &User{
			Model:  modelType.USER,
			Index:  nextUserIndex,
			Client: c,
			Place:  0,
			Alert:  0,
			Score:  NewScore(),
			Reward: NewReward(),
			UserData: &UserData{
				int(result.Id.Int32),
				result.Uid.String,
				result.Uuid.String,
				result.Name.String,
				"",
				1,
				int(result.Sex.Int32),
				int(result.Level.Int32),
				int(result.Exp.Int32),
				0,
				int(result.Coin.Int32),
				int(result.Cash.Int32),
				int(result.Point.Int32),
				int(result.Win.Int32),
				int(result.Lose.Int32),
				0, //int(result.Kill.Int32),
				int(result.Death.Int32),
				int(result.Assist.Int32),
				int(result.Blast.Int32),
				int(result.Rescue.Int32),
				int(result.RescueCombo.Int32),
				int(result.Survive.Int32),
				int(result.Escape.Int32),
				"",
				result.RedGraphics.String,
				result.BlueGraphics.String,
				result.Memo.String,
				int(result.LastChat.Int32),
				int(result.Admin.Int32),
			},
		}
		Users[c] = user
		user.Setting(user.Model, user.UserData.BlueGraphics)
		user.UserData.MaxExp = cMath.GetMaxExp(user.UserData.Level)
		return user, true
	}
	return nil, false
}

func (u *User) Remove() bool {
	_, ok := Users[u.Client]
	if ok {
		delete(Users, u.Client)
	}
	return ok
}

func (u *User) GetUserdata() (int, int, string, string, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, string, string, string, string, int) {
	return u.Index, u.UserData.Id, u.UserData.Name, u.UserData.ClanName,
		u.UserData.Rank, u.UserData.Sex, u.UserData.Level, u.UserData.Exp, u.UserData.MaxExp,
		u.UserData.Coin, u.UserData.Cash, u.UserData.Point,
		u.UserData.Win, u.UserData.Lose, u.UserData.Kill, u.UserData.Death, u.UserData.Assist,
		u.UserData.Blast, u.UserData.Rescue, u.UserData.Survive, u.UserData.Escape,
		u.UserData.Graphics, u.UserData.RedGraphics, u.UserData.BlueGraphics, u.UserData.Memo, u.UserData.Admin
}

func (u *User) GetCreateGameObject(hide bool) (model int, index int, name string, clanName string, team int, level int, image string, x int, y int, dirX int, dirY int, collider bool) {
	model = 1
	index = u.Index
	name = u.UserData.Name
	clanName = ""
	team = 0
	if t, ok := u.GameData["team"]; ok {
		team = t.(int)
	}
	level = u.UserData.Level
	image = u.Image
	x = u.X
	y = u.Y
	dirX = u.DirX
	dirY = u.DirY
	collider = false
	if hide {
		name = ""
		clanName = ""
		level = 0
	}
	return
}

func (u *User) SetUpLevel(v int) {
	u.UserData.Level += v
}

func (u *User) SetUpExp(v int) {
	if u.UserData.Level > 200 {
		return
	}
	u.UserData.Exp = cMath.Max(u.UserData.Exp+v, 0)
	for u.UserData.Exp >= u.UserData.MaxExp {
		u.UserData.Exp -= u.UserData.MaxExp
		u.SetUpLevel(1)
	}
}

func (u *User) SetUpCash(v int) {
	u.UserData.Cash += v
	// TODO : 개발중
}

func (u *User) SetGraphics(image string) {
	u.Image = image
	u.PublishMap(toClient.SetGraphics(1, u.Index, image))
}

func (u *User) ChangeName(name string) {
	if u.UserData.Cash < 500 {
		u.Send(toClient.MessageLobby("NOT_ENOUGH_CASH"))
		return
	}
	nameLen := utf8.RuneCountInString(name)
	if nameLen < 1 || nameLen > 6 {
		u.Send(toClient.MessageLobby("AN_IMPOSSIBLE_LENGTH"))
		return
	}
	if match, _ := regexp.MatchString("[^가-힣]", name); match {
		u.Send(toClient.MessageLobby("AN_IMPOSSIBLE_WORD"))
		return
	}
	// TODO : filtering
	// TODO : find user by name
	// TODO : update user name
	// TODO : rank change
	u.UserData.Name = name
	u.SetUpCash(-500)
	u.Send(toClient.MessageLobby("CHANGE_USERNAME_SUCCESS"))
	u.Send(toClient.UserData(u.GetUserdata()))
}

func (u *User) CreateClan(name string) {
	// TODO
}

func (u *User) InviteClan(name string) {
	// TODO
}

func (u *User) JoinClan(id int) {
	// TODO
}

func (u *User) CancelClan(id int) {
	// TODO
}

func (u *User) KickClan(id int) {
	// TODO
}

func (u *User) GetClan() {

}

func (u *User) LeaveClan() {
	// TODO
}

func (u *User) SetOptionClan(d []byte) {
	// TODO
}

func (u *User) PayClan(coin int) {
	// TODO
}

func (u *User) DonateClan(cash int) {
	// TODO
}

func (u *User) WithdrawClan(coin int) {
	// TODO
}

func (u *User) LevelUpClan() {
	// TODO
}

func (u *User) SetUpMemberLevelClan(id int) {
	// TODO
}

func (u *User) SetDownMemberLevelClan(id int) {
	// TODO
}

func (u *User) ChangeMasterClan(id int) {
	// TODO
}

func (u *User) GetBilling() {
	// TODO
}

func (u *User) GetPayInfoItem(id int) {
	// TODO
}

func (u *User) UseBilling(id int) {
	// TODO
}

func (u *User) RefundBilling(id int) {
	// TODO
}

func (u *User) GetShop(page int) {
	// TODO
}

func (u *User) GetInfoItem(id int) {
	// TODO
}

func (u *User) BuyItem(d []byte) {
	// TODO
}

func (u *User) AddItem(id int, num int, expiry int) {
	// TODO
}

func (u *User) CheckSkinExpiry() {
	// TODO
}

func (u *User) GetSkinList() {
	// TODO
}

func (u *User) GetRank(page int) {
	// TODO
}

func (u *User) GetUserInfoRank(id int) {
	// TODO
}

func (u *User) GetUserInfoRankByUserName(name string) {
	// TODO
}

func (u *User) GetNoticeMessageCount() {
	// TODO
}

func (u *User) GetNoticeMessage(deleted int) {
	// TODO
}

func (u *User) GetInfoNoticeMessage(id int) {
	// TODO
}

func (u *User) WithdrawNoticeMessage(id int) {
	// TODO
}

func (u *User) DeleteNoticeMessage(id int) {
	// TODO
}

func (u *User) RestoreNoticeMessage(id int) {
	// TODO
}

func (u *User) ClearNoticeMessage() {
	// TODO
}

func (u *User) AddNoticeMessage(d []byte) {
	// TODO
}

func (u *User) AddUserReport(d []byte) {
	// TODO
}

func (u *User) SetSkin(id int) {
	// TODO
}

func (u *User) SetState(state int) {
	// TODO
}

func (u *User) Turn(dirX int, dirY int) {
	if u.Room == nil {
		return
	}
	u.Character.Turn(dirX, dirY)
}

func (u *User) Move(x int, y int) {
	if u.Room == nil {
		return
	}
	u.Character.Turn(x, y)
	dir := u.GetDirection(x, y)
	r := u.Room
	if r.Passable(u.Place, u.X, u.Y, dir, false) && r.Passable(u.Place, u.X+x, u.Y-y, 10-dir, true) {
		u.Character.Move(x, -y)
		r.Portal(u)
	} else {
		u.Teleport(u.Place, u.X, u.Y)
	}
}

func (u *User) Chat(text string) {
	if u.Room == nil {
		return
	}
	r := u.Room
	size := len(text)
	if size > 128 {
		size = 128
	}
	text = text[:size]
	if u.UserData.LastChat > int(time.Now().Unix()) {
		u.Send(toClient.SystemMessage("<color=red>운영진에 의해 채팅이 금지되었습니다.</color> (" + time.Unix(int64(u.UserData.LastChat), 0).String() + ")"))
		return
	}
	if cFilter.Check(text) {
		u.Alert++
		if u.Alert >= 3 {
			u.Send(toClient.QuitGame())
		} else {
			u.Send(toClient.Vibrate())
			u.Send(toClient.SystemMessage("<color=red>금칙어를 언급하여 경고 " + strconv.Itoa(u.Alert) + "회를 받았습니다. 3회 이상시 자동 추방됩니다.</color>"))
		}
		return
	}
	if u.Command(text) || u.Alert >= 3 {
		return
	}
	fmt.Println(string(u.UserData.Name) + "(#" + strconv.Itoa(r.Index) + "@" + strconv.Itoa(u.Place) + "): " + text)
	switch r.RoomType {
	case roomType.GAME:
		if team, ok := u.GameData["team"]; ok {
			if team.(int) == teamType.RED {
				u.ChatToRedTeam(text)
			} else {
				u.ChatToBlueTeam(text)
			}
			break
		}
		fallthrough
	default:
		u.Publish(toClient.ChatMessage(u.Model, u.Index, u.UserData.Name, text))
	}
}

func (u *User) ChatToRedTeam(text string) {
	fmt.Println(u.Room)
	u.Publish(toClient.ChatMessage(u.Model, u.Index, "<color=#FF0000>"+u.UserData.Name+"</color>", text))
}

func (u *User) ChatToBlueTeam(text string) {
	fmt.Println(u.Room)
	if caught, ok := u.GameData["caught"]; ok {
		if caught.(bool) {
			u.PublishMap(toClient.ChatMessage(u.Model, u.Index, "<color=#808080>"+u.UserData.Name+"</color>", text))
			return
		}
	}
	u.Publish(toClient.ChatMessage(u.Model, u.Index, "<color=#00A2E8>"+u.UserData.Name+"</color>", text))
}

func (u *User) Command(text string) bool {
	if len(text) < 1 {
		return false
	}
	if text[:1] == "#" {
		if u.UserData.Admin < 1 {
			if u.UserData.Cash < 20 {
				u.Send(toClient.SystemMessage("<color=red>보석이 부족합니다. 보석 20개가 필요합니다.</color>"))
			} else {
				u.Client.Broadcast(toClient.SystemMessage("<color=#1DDB16>" + u.UserData.Name + "#" + strconv.Itoa(u.Room.Index) + ": " + text[1:] + "</color>")) // TODO : 확인 필요
				u.SetUpCash(-20)
			}
		} else {
			u.Client.Broadcast(toClient.SystemMessage("<color=#EFE4B0>@[운영진] " + u.UserData.Name + ": " + text[1:] + "</color>"))
		}
		return true
	}
	if len(text) < 3 || u.UserData.Admin < 1 {
		return false
	}
	//slice := strings.Split(text, ",")
	switch text[:3] {
	default:
		return false
	}
	return true
}

func (u *User) Ban(target *User, name string, description string, days int) {
	if target == nil {
		var findUser string
		if findUser == "" {
			return
		}
		// TODO : INSERT DB
	} else {
		// TODO : INSERT DB
		target.Send(toClient.QuitGame())
	}
	text := pix.Maker(name, "를", "을") + " " + string(days) + "일 동안 접속을 차단함. (" + description + ")"
	u.Publish(toClient.SystemMessage("<color=red>" + text + "</color>"))
	log.Println(text)
}

func (u *User) Entry(rType int) {
	if u.Room != nil {
		return
	}
	// TODO : set state, send
	r := AvailableRoom(rType)
	r.Join(u)
}

func (u *User) Leave() {
	if u.Room == nil {
		return
	}
	u.Room.Leave(u)
}

func (u *User) InputArrow(d *getty.Data) {
	var err error
	var pos [2]int8
	var dir [1]uint8
	buf := bytes.NewBuffer(d.Buffers)
	err = binary.Read(buf, binary.BigEndian, &pos)
	CheckError(err)
	err = binary.Read(buf, binary.BigEndian, &dir)
	CheckError(err)
	x, y := int(pos[0]), int(pos[1])
	if dir[0] == 0 {
		u.Turn(x, y)
	} else {
		u.Move(x, y)
	}
}

func (u *User) Hit() {
	if u.Room == nil {
		return
	}
	u.Room.Hit(u)
}

func (u *User) UseItem() {
	if u.Room == nil {
		return
	}
	u.Room.UseItem(u)
}

func (u *User) Portal(place int, x int, y int, dirX int, dirY int) {
	u.BroadcastMap(toClient.RemoveGameObject(u.Model, u.Index))
	u.Place = place
	u.SetPosition(x, y)
	if !(dirX == dirY && dirX == 0) {
		u.Character.Turn(dirX, dirY)
	}
	u.Send(toClient.Portal(place, x, y, u.DirX, u.DirY))
}

func (u *User) Teleport(place int, x int, y int) {
	if u.Room == nil {
		return
	}
	u.Room.Teleport(u, place, x, y, 0, -1)
}

func (u *User) Result(ad int) {
	if u.GameData == nil {
		return
	}
	if result, ok := u.GameData["result"]; ok {
		if result.(bool) {
			switch ad {
			case 1:
				u.Entry(roomType.GAME)
			case 2:
				u.Reward.Cash += 10
			}
		}
	}
	u.Score.Send(u)
	u.Reward.Send(u)
	u.GameData["result"] = false
}

func (u *User) UpdateUser() {
	go db.UpdateUser(
		u.UserData.Id,
		u.UserData.Uuid,
		u.UserData.Sex,
		u.UserData.Level,
		u.UserData.Exp,
		u.UserData.Coin,
		u.UserData.Cash,
		u.UserData.Point,
		u.UserData.Kill,
		u.UserData.Death,
		u.UserData.Assist,
		u.UserData.Blast,
		u.UserData.Rescue,
		u.UserData.RescueCombo,
		u.UserData.Survive,
		u.UserData.Escape,
		u.UserData.RedGraphics,
		u.UserData.BlueGraphics,
		u.UserData.Memo,
	)
}

func (u *User) Disconnect() {
	u.UpdateUser()
	u.Leave()
	u.Remove()
}

func (u *User) Send(d []byte) {
	u.Client.Send(d)
}

func (u *User) Publish(d []byte) {
	if u.Room == nil {
		return
	}
	u.Room.Publish(d)
}

func (u *User) PublishMap(d []byte) {
	if u.Room == nil {
		return
	}
	u.Room.PublishMap(u.Place, d)
}

func (u *User) Broadcast(d []byte) {
	if u.Room == nil {
		return
	}
	u.Room.Broadcast(u, d)
}

func (u *User) BroadcastMap(d []byte) {
	if u.Room == nil {
		return
	}
	u.Room.BroadcastMap(u, d)
}
