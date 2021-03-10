package game

import (
	"fmt"
	"regexp"
	"unicode/utf8"

	"godori.com/db"
	"godori.com/getty"
	toClient "godori.com/packet/toClient"
	modelType "godori.com/util/constant/modelType"
	roomType "godori.com/util/constant/roomType"
	teamType "godori.com/util/constant/teamType"
	cMath "godori.com/util/math"
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
	Survive      int
	Escape       int
	Graphics     string
	RedGraphics  string
	BlueGraphics string
	Memo         string
	Admin        int
}

type User struct {
	Model  int
	Index  int
	Client *getty.Client
	Room   *Room
	Place  int
	Alert  int
	Character
	UserData *UserData
	GameData map[string]interface{}
}

var nextUserIndex int
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
				int(result.Survive.Int32),
				int(result.Escape.Int32),
				"",
				result.RedGraphics.String,
				result.BlueGraphics.String,
				result.Memo.String,
				int(result.Admin.Int32),
			},
		}
		Users[c] = user
		user.Setting(user.Model, user.UserData.RedGraphics, user.UserData.BlueGraphics)
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
	// TODO
}

func (u *User) LeaveClan() {
	// TODO
}

func (u *User) SetOptionClan(d []byte) {
	// TODO
}

func (u *User) PayClan(d []byte) {
	// TODO
}

func (u *User) DonateClan(cash int) {
	// TODO
}

func (u *User) WithdrawClan(d []byte) {
	// TODO
}

func (u *User) LevelUpClan() {
	// TODO
}

func (u *User) SetUpMemberLevelClan(d []byte) {
	// TODO
}

func (u *User) SetDownMemberLevelClan(d []byte) {
	// TODO
}

func (u *User) ChangeMasterClan(d []byte) {
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

func (u *User) GetNoticeMessage(deleted bool) {
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
	u.Turn(dirX, dirY)
}

func (u *User) Move(x int, y int) {
	if u.Room == nil {
		return
	}
	u.Turn(x, y)
	dir := u.GetDirection(x, y)
	r := u.Room
	if r.Passable(u.Place, u.X, u.Y, dir, false) && r.Passable(u.Place, u.X+x, u.Y-y, 10-dir, true) {
		u.Move(x, -y)
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
	text = text[:35]
	// TODO : 채팅 금지
	// TODO : filtering
	if u.Command(text) {
		return
	}
	fmt.Println(string(u.UserData.Name) + "(#" + string(r.Index) + "@" + string(u.Place) + "): " + text)
	switch r.RoomType {
	case roomType.PLAYGROUND:
		u.Publish(toClient.ChatMessage(u.Model, u.Index, u.UserData.Name, text))
	case roomType.GAME:
		if team, ok := u.GameData["team"]; ok {
			if team.(int) == teamType.RED {
				u.ChatToRedTeam(text)
			} else {
				u.ChatToBlueTeam(text)
			}
		}
	}
}

func (u *User) Command(text string) bool {
	return true
}

func (u *User) ChatToRedTeam(text string) {
	u.Publish(toClient.ChatMessage(u.Model, u.Index, "<color=#00A2E8>"+u.UserData.Name+"</color>", text))
}

func (u *User) ChatToBlueTeam(text string) {
	if caught, ok := u.GameData["caught"]; ok {
		if caught.(bool) {
			u.PublishMap(toClient.ChatMessage(u.Model, u.Index, "<color=#808080>"+u.UserData.Name+"</color>", text))
			return
		}
	}
	u.Publish(toClient.ChatMessage(u.Model, u.Index, "<color=#00A2E8>"+u.UserData.Name+"</color>", text))
}

func (u *User) Ban(target *User, name string, description string, days int) {
	if target == nil {
		// find
	} else {
		// send
	}
	// TODO
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
	// TODO : broadcast
	u.Place = place
	u.SetPosition(x, y)
	if !(dirX == dirY && dirX == 0) {
		u.Turn(dirX, dirY)
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
				// TODO
			}
		}
	}
	// TODO
	u.GameData = nil
}

func (u *User) Disconnect() {
	u.Leave()
	u.Remove()
	// TODO : db 저장
}

func (u *User) Send(d []byte) {
	u.Client.Send(d)
}

// TODO : notice는 clients의 broadcast 사용할 것.

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
