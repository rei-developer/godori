package packet

import (
	"encoding/json"
	"log"
)

func UserData(index int, id int, name string, clanName string, rank int, sex int, level int, exp int, maxExp int, coin int, cash int, point int, win int, lose int, kill int, death int, assist int, blast int, rescue int, survive int, escape int, grphics string, redGraphics string, blueGraphics string, memo string, admin int) []byte {
	return PakcetWrapper(json.Marshal(struct {
		H            int
		Index        int
		Id           int
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
	}{USER_DATA, index, id, name, clanName, rank, sex, level, exp, maxExp, coin, cash, point, win, lose, kill, death, assist, blast, rescue, survive, escape, grphics, redGraphics, blueGraphics, memo, admin}))
}

func Vibrate() []byte {
	return PakcetWrapper(json.Marshal(struct {
		H int
	}{VIBRATE}))
}

func ConnectionCount(count int) []byte {
	return PakcetWrapper(json.Marshal(struct {
		H     int
		Count int
	}{CONNECTION_COUNT, count}))
}

func SystemMessage(text string) []byte {
	return PakcetWrapper(json.Marshal(struct {
		H    int
		Text string
	}{SYSTEM_MESSAGE, text}))
}

func InformMessage(text string) []byte {
	return PakcetWrapper(json.Marshal(struct {
		H    int
		Text string
	}{INFORM_MESSAGE, text}))
}

func NoticeMessage(text string) []byte {
	return PakcetWrapper(json.Marshal(struct {
		H    int
		Text string
	}{NOTICE_MESSAGE, text}))
}

func ComboMessage(text string) []byte {
	return PakcetWrapper(json.Marshal(struct {
		H    int
		Text string
	}{COMBO_MESSAGE, text}))
}

func ChatMessage(model int, index int, name string, text string) []byte {
	return PakcetWrapper(json.Marshal(struct {
		H     int
		Model int
		Index int
		Name  string
		Text  string
	}{CHAT_MESSAGE, model, index, name, text}))
}

func Portal(place int, x int, y int, dirX int, dirY int) []byte {
	return PakcetWrapper(json.Marshal(struct {
		H     int
		Place int
		X     int
		Y     int
		DirX  int
		DirY  int
	}{PORTAL, place, x, y, dirX, dirY}))
}

func RemoveGameObject(model int, index int) []byte {
	return PakcetWrapper(json.Marshal(struct {
		H     int
		Model int
		Index int
	}{REMOVE_GAME_OBJECT, model, index}))
}

func CreateGameObject(model int, index int, name string, clanName string, team int, level int, image string, x int, y int, dirX int, dirY int, collider bool) []byte {
	return PakcetWrapper(json.Marshal(struct {
		H        int
		Model    int
		Index    int
		Name     string
		ClanName string
		Team     int
		Level    int
		Image    string
		X        int
		Y        int
		DirX     int
		DirY     int
		Collider bool
	}{CREATE_GAME_OBJECT, model, index, name, clanName, team, level, image, x, y, dirX, dirY, collider}))
}

func SetGraphics(model int, index int, image string) []byte {
	return PakcetWrapper(json.Marshal(struct {
		H     int
		Model int
		Index int
		Image string
	}{SET_GRAPHICS, model, index, image}))
}

func PlaySound(name string) []byte {
	return PakcetWrapper(json.Marshal(struct {
		H    int
		Name string
	}{PLAY_SOUND, name}))
}

func UpdateRoomUserCount(count int) []byte {
	return PakcetWrapper(json.Marshal(struct {
		H     int
		Count int
	}{UPDATE_ROOM_USER_COUNT, count}))
}

func UpdateModeCount(count int) []byte {
	return PakcetWrapper(json.Marshal(struct {
		H     int
		Count int
	}{UPDATE_MODE_COUNT, count}))
}

func UpdateGameItem(name string, num int) []byte {
	return PakcetWrapper(json.Marshal(struct {
		H    int
		Name string
		Num  int
	}{UPDATE_GAME_ITEM, name, num}))
}

func RemoveGameItem() []byte {
	return PakcetWrapper(json.Marshal(struct {
		H int
	}{REMOVE_GAME_ITEM}))
}

func SetGameTeam(index int, team int) []byte {
	return PakcetWrapper(json.Marshal(struct {
		H     int
		Index int
		Team  int
	}{SET_GAME_TEAM, index, team}))
}

func ModeData(mode int) []byte {
	return PakcetWrapper(json.Marshal(struct {
		H    int
		Mode int
	}{MODE_DATA, mode}))
}

// TODO : get clan ~ invite clan

func MemberInfoClan(memberId int) []byte {
	return PakcetWrapper(json.Marshal(struct {
		H        int
		MemberId int
	}{MEMBER_INFO_CLAN, memberId}))
}

func UpdateClan(level int, coin int, cash int) []byte {
	return PakcetWrapper(json.Marshal(struct {
		H     int
		Level int
		Coin  int
		Cash  int
	}{UPDATE_CLAN, level, coin, cash}))
}

func MessageClan(state string) []byte {
	return PakcetWrapper(json.Marshal(struct {
		H     int
		State string
	}{MESSAGE_CLAN, state}))
}

func DeadAnimation() []byte {
	return PakcetWrapper(json.Marshal(struct {
		H int
	}{DEAD_ANIMATION}))
}

func ResultGame(winnder int, rank int, persons int, mission string, exp int, coin int) []byte {
	return PakcetWrapper(json.Marshal(struct {
		H       int
		Winner  int
		Rank    int
		Persons int
		Mission string
		Exp     int
		Coin    int
	}{RESULT_GAME, winnder, rank, persons, mission, exp, coin}))
}

func EnterWardrobe() []byte {
	return PakcetWrapper(json.Marshal(struct {
		H int
	}{ENTER_WARDROBE}))
}

func LeaveWardrobe() []byte {
	return PakcetWrapper(json.Marshal(struct {
		H int
	}{LEAVE_WARDROBE}))
}

func SwitchLight(flag bool) []byte {
	return PakcetWrapper(json.Marshal(struct {
		H    int
		Flag bool
	}{SWITCH_LIGHT, flag}))
}

func QuitGame() []byte {
	return PakcetWrapper(json.Marshal(struct {
		H int
	}{QUIT_GAME}))
}

// TODO : get billing

func UpdateBilling(id int, use bool, refund bool) []byte {
	return PakcetWrapper(json.Marshal(struct {
		H      int
		Id     int
		Use    bool
		Refund bool
	}{UPDATE_BILLING, id, use, refund}))
}

func GetPayInfoItem(id int, cash int, memo string, purchaseDate string, use bool, refund bool) []byte {
	return PakcetWrapper(json.Marshal(struct {
		H       int
		Id      int
		Cash    int
		Memo    string
		Regdate string
		Use     bool
		Refund  bool
	}{GET_PAY_INFO_ITEM, id, cash, memo, purchaseDate, use, refund}))
}

// TODO : get shop

func GetSkinItem(model int, id int, icon string, name string, creator string, desc string, cost int, pay bool, expiry string) []byte {
	return PakcetWrapper(json.Marshal(struct {
		H       int
		Model   int
		Id      int
		Icon    string
		Name    string
		Creator string
		Desc    string
		Cost    int
		Pay     bool
		Expiry  string
	}{GET_SKIN_ITEM, model, id, icon, name, creator, desc, cost, pay, expiry}))
}

func MessageShop(state string) []byte {
	return PakcetWrapper(json.Marshal(struct {
		H     int
		State string
	}{MESSAGE_SHOP, state}))
}

func MessageLobby(state string) []byte {
	return PakcetWrapper(json.Marshal(struct {
		H     int
		State string
	}{MESSAGE_LOBBY, state}))
}

// TODO : get skin list

func UpdateCoinAndCash(coin int, cash int) []byte {
	return PakcetWrapper(json.Marshal(struct {
		H    int
		Coin int
		Cash int
	}{UPDATE_COIN_AND_CASH, coin, cash}))
}

// TODO : get rank

func GetUserInfoRank(name string, clanName string, rank int, level int, exp int, maxExp int, kill int, death int, assist int, likes int, memo string, avatar string) []byte {
	return PakcetWrapper(json.Marshal(struct {
		H        int
		Name     string
		ClanName string
		Rank     int
		Level    int
		Exp      int
		MaxExp   int
		Kill     int
		Death    int
		Assist   int
		Likes    int
		Memo     string
		Avatar   string
	}{GET_USER_INFO_RANK, name, clanName, rank, level, exp, 100000, kill, death, assist, likes, memo, avatar}))
}

func MessageRank(state string) []byte {
	return PakcetWrapper(json.Marshal(struct {
		H     int
		State string
	}{MESSAGE_RANK, state}))
}

func GetNoticeMessageCount(count int) []byte {
	return PakcetWrapper(json.Marshal(struct {
		H     int
		Count int
	}{GET_NOTICE_MESSAGE_COUNT, count}))
}

// TODO : get notice message

func GetInfoNoticeMessage(id int, avatar string, author string, title string, content string, coin int, cash int, created string, deleted bool, rewarded int) []byte {
	return PakcetWrapper(json.Marshal(struct {
		H        int
		Id       int
		Avatar   string
		Author   string
		Title    string
		Content  string
		Coin     int
		Cash     int
		Created  string
		Deleted  bool
		Rewarded int
	}{GET_INFO_NOTICE_MESSAGE, id, avatar, author, title, content, coin, cash, created, deleted, rewarded}))
}

func DeleteNoticeMessage(id int) []byte {
	return PakcetWrapper(json.Marshal(struct {
		H  int
		Id int
	}{DELETE_NOTICE_MESSAGE, id}))
}

func MessageGame(state string) []byte {
	return PakcetWrapper(json.Marshal(struct {
		H     int
		State string
	}{MESSAGE_GAME, state}))
}

func SetAnimation(model int, index int, anim string, sound string) []byte {
	return PakcetWrapper(json.Marshal(struct {
		H     int
		Model int
		Index int
		Anim  string
		Sound string
	}{MESSAGE_GAME, model, index, anim, sound}))
}

func PakcetWrapper(d []byte, err error) []byte {
	CheckError(err)
	return d
}

func CheckError(err error) {
	if err != nil {
		log.Println(err)
	}
}

const (
	USER_DATA = iota
	VIBRATE
	CONNECTION_COUNT
	SYSTEM_MESSAGE
	INFORM_MESSAGE
	NOTICE_MESSAGE
	CHAT_MESSAGE
	PORTAL
	CREATE_GAME_OBJECT
	REMOVE_GAME_OBJECT
	SET_GRAPHICS
	PLAY_SOUND
	UPDATE_ROOM_USER_COUNT
	UPDATE_MODE_COUNT
	UPDATE_GAME_ITEM
	REMOVE_GAME_ITEM
	SET_GAME_TEAM
	MODE_DATA
	GET_CLAN
	INVITE_CLAN
	DEAD_ANIMATION
	RESULT_GAME
	ENTER_WARDROBE
	LEAVE_WARDROBE
	SWITCH_LIGHT
	QUIT_GAME
	UPDATE_MODE_SCORE
	SET_OPTION_CLAN
	MEMBER_INFO_CLAN
	UPDATE_CLAN
	MESSAGE_CLAN
	GET_BILLING
	UPDATE_BILLING
	GET_SHOP
	GET_SKIN_ITEM
	MESSAGE_SHOP
	MESSAGE_LOBBY
	GET_SKIN_LIST
	UPDATE_COIN_AND_CASH
	GET_PAY_INFO_ITEM
	GET_RANK
	GET_USER_INFO_RANK
	MESSAGE_RANK
	GET_NOTICE_MESSAGE_COUNT
	GET_NOTICE_MESSAGE
	GET_INFO_NOTICE_MESSAGE
	DELETE_NOTICE_MESSAGE
	MESSAGE_GAME
	SET_ANIMATION
	COMBO_MESSAGE
)
