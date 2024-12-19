package config

type Choice struct {
	Group   string   `json:"group"`
	Option  []string `json:"option"` //["value","text"]
	Checked bool     `json:"checked"`
}

func (c *Choice) Clone() *Choice {
	r := &Choice{
		Group:   c.Group,
		Option:  make([]string, len(c.Option)),
		Checked: c.Checked,
	}
	copy(r.Option, c.Option)
	return r
}

func (c *Choice) Merge(source *Choice) *Choice {
	var found bool
	for _, v := range source.Option {
		if len(v) == 0 {
			continue
		}
		for _, v2 := range c.Option {
			if len(v2) == 0 {
				continue
			}
			if v[0] == v2[0] {
				found = true
				break
			}
		}
		if !found {
			c.Option = append(c.Option, v)
		} else {
			found = false
		}
	}
	return c
}
