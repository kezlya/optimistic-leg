package main

func main() {
	StartServer()
}

func whatToDo(hive *Hive) map[int]BotOder {
	food := hive.Map.food()
	home := hive.Map.home(hive.Username)

	actions := make(map[int]BotOder)

	for id, ant := range hive.Ants {
		antPoint := point{y: ant.Y, x: ant.X}

		homeDir, isHome := antPoint.isHomeAround(hive.Map, hive.Username)
		if isHome && ant.Payload > 0 {
			actions[id] = BotOder{Unload, homeDir}
			continue
		}

		if mealDir, isMeal := antPoint.isMealAround(hive.Map); isMeal {
			if ant.Health < 9 {
				actions[id] = BotOder{Eat, mealDir}
				continue
			}
			if ant.Payload < 9 && !(isHome && homeDir == mealDir) {
				actions[id] = BotOder{Load, mealDir}
				continue
			}
		}

		if ant.Payload < 9 {
			actions[id] = BotOder{Move, antPoint.towardsFood(food)}
			continue
		}

		actions[id] = BotOder{Move, antPoint.towardsHome(home)}
	}

	return actions
}

func (m *Map) food() map[point]int {
	food := make(map[point]int)
	for y, row := range m.Cells {
		for x, c := range row {
			if c.Food > 0 {
				food[point{y: y, x: x}] = c.Food
			}
		}
	}
	return food
}

func (m *Map) home(username string) []point {
	home := make([]point, 0)
	for y, row := range m.Cells {
		for x, c := range row {
			if c.Hive == username {
				home = append(home, point{y: y, x: x})
			}
		}
	}
	return home
}

func (p *point) towardsFood(f map[point]int) Dir {
	effort := 10000000
	dir := Up
	for to := range f {
		ticks := p.distance(to)
		if ticks < effort {
			dir = p.move(to)
			effort = ticks
		}
	}
	return dir
}

func (p *point) towardsHome(h []point) Dir {
	effort := 10000000
	dir := Up
	for _, to := range h {
		ticks := p.distance(to)
		if ticks < effort {
			dir = p.move(to)
			effort = ticks
		}
	}
	return dir
}

func (p *point) isMealAround(world *Map) (d Dir, y bool) {
	if p.y > 0 && world.Cells[p.y-1][p.x].Food > 0 {
		return Up, true
	}
	if p.y < world.Height-1 && world.Cells[p.y+1][p.x].Food > 0 {
		return Down, true
	}
	if p.x < world.Width-1 && world.Cells[p.y][p.x+1].Food > 0 {
		return Right, true
	}
	if p.x > 0 && world.Cells[p.y][p.x-1].Food > 0 {
		return Left, true
	}
	return
}

func (p *point) isHomeAround(world *Map, username string) (d Dir, y bool) {
	if p.y > 0 && world.Cells[p.y-1][p.x].Hive == username {
		return Up, true
	}
	if p.y < world.Height-1 && world.Cells[p.y+1][p.x].Hive == username {
		return Down, true
	}
	if p.x < world.Width-1 && world.Cells[p.y][p.x+1].Hive == username {
		return Right, true
	}
	if p.x > 0 && world.Cells[p.y][p.x-1].Hive == username {
		return Left, true
	}
	return
}

func (from *point) distance(to point) int {
	w, h := 0, 0
	if from.x > to.x {
		w = from.x - to.x
	} else {
		w = to.x - from.x
	}
	if from.y > to.y {
		h = from.y - to.y
	} else {
		h = to.y - from.y
	}
	return w + h
}

func (from *point) move(to point) Dir {
	if from.x > to.x {
		return Left
	}
	if from.x < to.x {
		return Right
	}
	if from.y > to.y {
		return Up
	}
	// from.y < to.y
	return Down
}
