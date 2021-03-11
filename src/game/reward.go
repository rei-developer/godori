package game

type Reward struct {
	Exp   int
	Coin  int
	Cash  int
	Point int
}

func NewReward() *Reward {
	return &Reward{0, 0, 0, 0}
}

func (r *Reward) Send(u *User) {
	u.SetUpExp(r.Exp)
	u.UserData.Coin += r.Coin
	u.SetUpCash(r.Cash)
	u.UserData.Point += r.Point
	r.Clear()
}

func (r *Reward) Clear() {
	r.Exp = 0
	r.Coin = 0
	r.Cash = 0
	r.Point = 0
}
