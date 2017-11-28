package components

import "time"

type HeadHeight struct {
	GameComponent
}

func (c *FreeMove) Update(elapsed time.Duration) {
	temp := c.Transform().Pos()
	temp[1] = 1.8
	c.Transform().SetPos(temp)
}