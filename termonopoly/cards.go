package termonopoly

type ChanceCard struct {
	text   string
	action ChanceAction
}

type CommunityCard struct {
	text   string
	action ChanceAction
}

type ChanceAction interface {
	OnDraw(player *Player)
}

type Advance struct {
	dest string
	pass bool
}

func (c *Advance) OnDraw(player *Player) {
	player.MoveTo(c.dest, c.pass)
}

type Back struct {
	dist int
}

func (c *Back) OnDraw(player *Player) {
	player.MoveDist(40-c.dist, false)
}

type Sum struct {
	rollCount int
}

func (c *Sum) OnDraw(player *Player) {
	total := 0
	for i := 0; i < c.rollCount; i++ {
		l, _ := RollDice()

		total += l
	}

	player.Spend(total)
}

type Near struct {
	nodeType string
}

func (c *Near) OnDraw(player *Player) {
	player.MoveUntil(c.nodeType, true)
}

type Trans struct {
	amount    int
	pay       bool
	toPlayers bool
}

func (c *Trans) OnDraw(player *Player) {

	if c.pay {
		if c.toPlayers {
			for _, pl := range Players {
				if !pl.IsBankrupt() && pl != player {
					player.Spend(c.amount)
					pl.Receive(c.amount)
				}
			}
		} else {
			player.Spend(c.amount)
		}
	} else {

		if c.toPlayers {
			for _, pl := range Players {
				if !pl.IsBankrupt() && pl != player {
					pl.Spend(c.amount)
					player.Receive(c.amount)
				}
			}
		} else {
			player.Receive(c.amount)
		}
	}
}

type JailFree struct {
}

func (c *JailFree) OnDraw(player *Player) {
	player.JailFree++
}

type Repair struct {
	house int
	hotel int
}

func (c *Repair) OnDraw(player *Player) {
	total := 0

	for _, colour := range player.FullSets {
		for _, prop := range player.Properties[colour] {
			if prop.Houses == 5 {
				total += c.hotel
			} else {
				total += c.house * prop.Houses
			}
		}
	}

	player.Spend(total)
}
