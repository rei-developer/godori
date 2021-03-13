package game

type Score struct {
	Kill               int
	KillForWardrobe    int
	KillCombo          int
	Death              int
	DeathForWardrobe   int
	Assist             int
	Blast              int
	Rescue             int
	RescueCombo        int
	Escape             int
	Survive            int
	SurviveForWardrobe int
	FoundKey           int
	Sum                int
}

func NewScore() *Score {
	return &Score{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
}

func (s *Score) Send(u *User) {
	u.UserData.Kill += s.Kill + s.KillForWardrobe
	u.UserData.Death += s.Death + s.DeathForWardrobe
	u.UserData.Assist += s.Assist
	u.UserData.Blast += s.Blast
	u.UserData.Rescue += s.Rescue
	// 1 combo
	u.UserData.Escape += s.Escape
	u.UserData.Survive += s.Survive + s.SurviveForWardrobe
	s.Clear()
}

func (s *Score) Clear() {
	s.Kill = 0
	s.KillForWardrobe = 0
	s.KillCombo = 0
	s.Death = 0
	s.DeathForWardrobe = 0
	s.Assist = 0
	s.Blast = 0
	s.Rescue = 0
	s.RescueCombo = 0
	s.Escape = 0
	s.Survive = 0
	s.SurviveForWardrobe = 0
	s.FoundKey = 0
	s.Sum = 0
}
