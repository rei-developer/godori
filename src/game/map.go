package game

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"godori.com/db"
	cMath "godori.com/util/math"
)

type GameMap struct {
	Name       string  `json: "name"`
	BGM        string  `json: "bgm"`
	Width      int     `json: "width"`
	Height     int     `json: "height"`
	Data       [][]int `json: "data"`
	Collisions []int   `json: "collisions"`
	Priorities []int   `json: "priorities"`
	Portals    []*Portal
}

type Portal struct {
	Place     int
	X         int
	Y         int
	NextPlace int
	NextX     int
	NextY     int
	NextDirX  int
	NextDirY  int
	Sound     string
}

const maxMapCnt int = 257

var GameMaps map[int]*GameMap = make(map[int]*GameMap)

func init() {
	fmt.Println("맵 로딩중...")
	for i := 1; i < maxMapCnt; i++ {
		jsonFile, err := os.Open("./lib/maps/" + strconv.Itoa(i) + ".json")
		if err != nil {
			fmt.Println(err)
		}
		defer jsonFile.Close()
		byteValue, _ := ioutil.ReadAll(jsonFile)
		var mapData GameMap
		json.Unmarshal(byteValue, &mapData)
		GameMaps[i] = &GameMap{
			Name:       mapData.Name,
			BGM:        mapData.BGM,
			Width:      mapData.Width,
			Height:     mapData.Height,
			Data:       mapData.Data,
			Collisions: mapData.Collisions,
			Priorities: mapData.Priorities,
		}
		portals := db.GetPortals(i)
		for _, p := range portals {
			GameMaps[i].Portals = append(GameMaps[i].Portals, &Portal{
				Place:     int(p.Place.Int32),
				X:         int(p.X.Int32),
				Y:         int(p.Y.Int32),
				NextPlace: int(p.NextPlace.Int32),
				NextX:     int(p.NextX.Int32),
				NextY:     int(p.NextY.Int32),
				NextDirX:  int(p.NextDirX.Int32),
				NextDirY:  int(p.NextDirY.Int32),
				Sound:     p.Sound.String,
			})
		}
	}
	fmt.Println("맵 로딩 완료")
}

func (m *GameMap) GetPortal(x int, y int) (*Portal, bool) {
	for _, p := range m.Portals {
		if p.X == x && p.Y == y {
			return p, true
		}
	}
	return nil, false
}

func (m *GameMap) RangePortal(x int, y int, rng int) bool {
	for _, p := range m.Portals {
		if cMath.Pow(cMath.Abs(p.X-x), 2)+cMath.Pow(cMath.Abs(p.Y-y), 2) <= rng*rng {
			return true
		}
	}
	return false
}

func (m *GameMap) Passable(x int, y int, dir int) bool {
	if !m.Valid(x, y) {
		return false
	}
	var bit int = int((1 << (int(dir/2) - 1)) & 0x0f)
	for layer := len(m.Data) - 1; layer >= 0; layer-- {
		// TODO : --layer 확인해야 함.
		var tileId = m.Data[layer][x+y*m.Width]
		if tileId == 0 {
			continue
		}
		var collision = m.Collisions[tileId-1]
		if (collision & bit) != 0 {
			return false
		} else if (collision & 0x0f) == 0x0f {
			return false
		} else if m.Priorities[tileId-1] == 0 {
			return true
		}
	}
	return true
}

func (m *GameMap) Valid(x int, y int) bool {
	return x >= 0 && x < m.Width && y >= 0 && y < m.Height
}
