package textgrid
import (
	"fmt"
	"strconv"
	"strings"
	"errors"
	"os"
	"bytes"
	"bufio"
)

type TextGrid struct {
	Name 		string
	MaxTime		float64
	MinTime		float64
	Tiers		[]*IntervalTier
}

func (G *TextGrid) Str () string {
	return fmt.Sprintf("TextGrid <%s, %d Tiers>", G.Name, len(G.Tiers))
}

func (G *TextGrid) GetFirst (tierName string) *IntervalTier {
	for _, I := range G.Tiers {
		if I.Name == tierName {
			return I
		}
	}
	return nil
}

func (G *TextGrid) GetList (tierName string) []*IntervalTier {
	var it []*IntervalTier
	for _, I := range G.Tiers {
		if I.Name == tierName {
			it = append(it, I)
		}
	}
	return it
}

func (G *TextGrid) GetNames () []string {
	var names []string
	for _, I := range G.Tiers {
		names = append(names, I.Name)
	}
	return names
}


func (G *TextGrid) Append (tier *IntervalTier) error {
	if G.MaxTime > 0 && tier.MaxTime > 0 && tier.MaxTime > G.MaxTime {
		return errors.New(tier.Str())
	}
	G.Tiers = append(G.Tiers, tier)
	return nil
}


func (G *TextGrid) Write (filePath string) error {
	if G.Tiers == nil {
		return errors.New("TextGrid nil")
	}

	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	var buf bytes.Buffer
	buf.WriteString("File type = \"ooTextFile\"\n")
	buf.WriteString("Object class = \"TextGrid\"\n")
	buf.WriteString(fmt.Sprintf("xmin = %v\n", G.MinTime))
	buf.WriteString(fmt.Sprintf("xmax = %v\n", G.MaxTime))
	buf.WriteString("tiers? <exists>\n")
	buf.WriteString(fmt.Sprintf("size = %d\n", len(G.Tiers)))
	buf.WriteString("item []:\n")

	for index, tier := range G.Tiers {
		buf.WriteString(fmt.Sprintf("\titem [%d]:\n", index))
		buf.WriteString("\t\tclass = \"IntervalTier\"\n")
		buf.WriteString(fmt.Sprintf("\t\tname = \"%s\"\n", tier.Name))
		buf.WriteString(fmt.Sprintf("\t\txmin = %v\n",tier.MinTime))
		buf.WriteString(fmt.Sprintf("\t\txmax = %v\n",tier.MaxTime))
		buf.WriteString(fmt.Sprintf("\t\tintervals: size = %d\n", len(tier.Intervals)))
		for j, interval := range tier.Intervals {
			buf.WriteString(fmt.Sprintf("\t\t\tintervals [%d]:\n", j))
            buf.WriteString(fmt.Sprintf("\t\t\t\txmin = %v\n", interval.MinTime))
            buf.WriteString(fmt.Sprintf("\t\t\t\txmax = %v\n", interval.MaxTime))
            buf.WriteString(fmt.Sprintf("\t\t\t\ttext = \"%s\"\n", interval.Mark))
		}
	}
	
	f.Write(buf.Bytes())
	return nil
}

func (G *TextGrid) Read (filePath string) error {
	f, err := os.Open(filePath)
	if err != nil {
	 return err
	}
	defer f.Close()
	rd := bufio.NewReader(f)

	line, _ := rd.ReadString('\n') // File type = "ooTextFile"
	line, _ = rd.ReadString('\n') // Object class = "TextGrid"

	line, _ = rd.ReadString('\n') // xmin =
	line = strings.Replace(line, "\n", "", -1)
	G.MinTime, _ = strconv.ParseFloat(strings.Split(line, " = ")[1], 64)

	line, _ = rd.ReadString('\n') // xmax = 
	line = strings.Replace(line, "\n", "", -1)
	G.MaxTime, _ = strconv.ParseFloat(strings.Split(line, " = ")[1], 64)

	line, _ = rd.ReadString('\n') // tiers? <exists>
	line, _ = rd.ReadString('\n') // size = 
	m, _ := strconv.Atoi(strings.Split(line, " = ")[1])

	line, _ = rd.ReadString('\n') // item []:
	for i := 0; i < m; i++ {
		var tier IntervalTier
		line, _ = rd.ReadString('\n') // item [n]: 
		line, _ = rd.ReadString('\n') // class = "IntervalTier"

		line, _ = rd.ReadString('\n') // name = 
		line = strings.Replace(line, "\n", "", -1)
		tier.Name = strings.Split(line, " = ")[1]

		line, _ = rd.ReadString('\n') // xmin = 
		tier.MinTime, _ = strconv.ParseFloat(strings.Split(line, " = ")[1], 64)

		line, _ = rd.ReadString('\n') // xmax = 
		tier.MaxTime, _ = strconv.ParseFloat(strings.Split(line, " = ")[1], 64)

		line, _ = rd.ReadString('\n') // intervals: size = 
		n, _ := strconv.Atoi(strings.Split(line, " = ")[1])
		for j := 0; j < n; j++ {
			var interval Interval
			line, _ = rd.ReadString('\n') // intervals [n]: 
			line, _ = rd.ReadString('\n') // xmin = 
			interval.MinTime, _ = strconv.ParseFloat(strings.Split(line, " = ")[1], 64)

			line, _ = rd.ReadString('\n') // xmax = 
			interval.MaxTime, _ = strconv.ParseFloat(strings.Split(line, " = ")[1], 64)

			line, _ = rd.ReadString('\n') // text = 
			line = strings.Replace(line, "\n", "", -1)
			interval.Mark = strings.Split(line, " = ")[1]

			tier.Intervals = append(tier.Intervals, &interval)
		}
		G.Tiers = append(G.Tiers, &tier)
	}
	return nil
}