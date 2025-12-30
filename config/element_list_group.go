package config

type Grouped struct {
	HasError bool
	Elements []FormElement
}

type Groups []Grouped

func SplitGroup(elements []FormElement) Groups {
	result := Groups{}
	t := 0
	g := Grouped{}
	var lastIsSet bool
	for idx, ele := range elements {
		if idx == 0 {
			if !g.HasError {
				if he, ok := ele.(HasError); ok {
					g.HasError = he.HasError()
				}
			}
			g.Elements = append(g.Elements, ele)
			t += ele.Cols()
			continue
		}
		isSet := ele.ElementType() == "fieldset" || ele.ElementType() == "langset"
		cols := ele.Cols()
		if isSet || (isSet != lastIsSet) || cols == 0 || t+cols > 12 {
			result = append(result, g)
			g = Grouped{}
			t = 0
		}
		lastIsSet = isSet
		if !g.HasError {
			if he, ok := ele.(HasError); ok {
				g.HasError = he.HasError()
			}
		}
		g.Elements = append(g.Elements, ele)
		t += ele.Cols()
	}
	if len(g.Elements) > 0 {
		result = append(result, g)
	}
	return result
}
