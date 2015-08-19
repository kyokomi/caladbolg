package parse

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"strings"
)

type State string

const (
	OK      State = "ok"
	NG      State = "ng"
	UNKNOWN State = "?"
)

type Coverage struct {
	BaseText string
	State    State
	Path     string
	Sec      float64
	Coverage float64
}

func ParseLine(line string) (*Coverage, error) {
	fields := strings.Fields(line)

	if len(fields) < 2 {
		return nil, fmt.Errorf("two fields required, have %d", len(fields))
	}

	c := &Coverage{}
	c.BaseText = line
	c.State = State(fields[0])
	c.Path = fields[1]

	switch c.State {
	case OK:
		fmt.Sscanf(fields[2], "%fs", &c.Sec)
		fmt.Sscanf(fields[4], "%f%%", &c.Coverage)
	case NG:
	case UNKNOWN:
	default:
		return nil, fmt.Errorf(`first field does not start with "Benchmark"`)
	}

	return c, nil
}

type Report struct {
	Coverages []Coverage
}

func ParseReport(r io.Reader) (Report, error) {
	report := Report{}

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		c, err := ParseLine(scanner.Text())
		if err != nil {
			log.Println(err)
			continue
		}
		report.Coverages = append(report.Coverages, *c)
	}
	if err := scanner.Err(); err != nil {
		return report, err
	}
	return report, nil
}

func (r Report) AllText() string {
	var text []string
	for _, report := range r.Coverages {
		text = append(text, report.BaseText)
	}
	return strings.Join(text, "\n")
}

func (r Report) OKCount() int {
	var count int
	for _, report := range r.Coverages {
		if report.State == OK {
			count++
		}
	}
	return count
}

func (r Report) NGCount() int {
	var count int
	for _, report := range r.Coverages {
		if report.State == NG {
			count++
		}
	}
	return count
}

func (r Report) NotTestCount() int {
	var count int
	for _, report := range r.Coverages {
		if report.State == UNKNOWN {
			count++
		}
	}
	return count
}

func (r Report) ReportOK() bool {
	// TODO: configにしたさ
	return r.AVG() >= 70 && r.NGCount() <= len(r.Coverages)/2
}

func (r Report) ReportCoverage() string {
	// TODO: templateにしたさ
	return fmt.Sprintf("%.2f%%", r.AVG())
}

func (r Report) ReportNoTestFiles() string {
	// TODO: templateにしたさ
	return fmt.Sprintf("%d/%d", r.NotTestCount(), len(r.Coverages))
}

func (r Report) AVG() float64 {
	return r.SUM() / float64(len(r.Coverages))
}

func (r Report) SUM() float64 {
	var sumCoverage float64
	for _, report := range r.Coverages {
		if report.State == OK {
			sumCoverage += report.Coverage
		}
	}
	return sumCoverage
}
