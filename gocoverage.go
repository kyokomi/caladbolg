package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"

	"github.com/nlopes/slack"
)

type CoverageService struct {
	Slack
}

func NewCoverageService(s Slack) CoverageService {
	return CoverageService{Slack: s}
}

const (
	attachmentColorOK = "#00CC00"
	attachmentColorNG = "#F35A00"
)

type coverage struct {
	Text     string
	TestOK   bool
	Path     string
	Sec      float64
	Coverage float64
}

func (c *coverage) Scan(text string) error {
	c.Text = text

	fields := strings.Fields(text)
	i := 0
	c.TestOK = fields[i] == "ok"
	i++
	c.Path = fields[i]
	if c.TestOK {
		i++
		c.Sec, _ = strconv.ParseFloat(strings.Replace(fields[i], "s", "", 1), 32)
		i++
		i++
		c.Coverage, _ = strconv.ParseFloat(strings.Replace(fields[i], "%", "", 1), 32)
	}

	return nil
}

type coverages []coverage

func (cs coverages) AllText() string {
	var text []string
	for _, report := range cs {
		text = append(text, report.Text)
	}
	return strings.Join(text, "\n")
}

func (cs coverages) OKCount() int {
	var count int
	for _, report := range cs {
		if report.TestOK {
			count++
		}
	}
	return count
}

func (cs coverages) NGCount() int {
	var count int
	for _, report := range cs {
		if !report.TestOK {
			count++
		}
	}
	return count
}

func (cs coverages) ReportOK() bool {
	return cs.AVG() >= 70 && cs.NGCount() <= len(cs)/2
}

func (cs coverages) ReportCoverage() string {
	return fmt.Sprintf("%f.2%%", cs.AVG())
}

func (cs coverages) ReportNoTestFiles() string {
	return fmt.Sprintf("%d/%d", cs.NGCount(), len(cs))
}

func (cs coverages) AVG() float64 {
	return cs.SUM() / float64(len(cs))
}

func (cs coverages) SUM() float64 {
	var sumCoverage float64
	for _, report := range cs {
		sumCoverage += report.Coverage
	}
	return sumCoverage
}

func (cs coverages) Scan(r io.Reader) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		c := coverage{}
		if err := c.Scan(scanner.Text()); err != nil {
			log.Println(err)
			continue
		}
		cs = append(cs, c)
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return cs
}

func (c CoverageService) Send(r io.Reader) error {
	cs := coverages{}
	if err := cs.Scan(r); err != nil {
		return err
	}

	// no test filesは50%以上 or 平均カバレッジが70%以下
	reportColor := attachmentColorNG
	if cs.ReportOK() {
		reportColor = attachmentColorOK
	}

	params := c.NewDefaultPostMessageParams()
	params.Attachments = []slack.Attachment{
		{
			Pretext:    "<https://honeybadger.io/path/to/event/|CircleCI Coverage Report> - :sushi:", // TODO: CircleCIのURL branch名
			Fallback:   "CircleCI Coverage Report - :sushi: https://honeybadger.io/path/to/event/",   // TODO: CircleCIのURL branch名
			Text:       cs.AllText(),
			AuthorName: "kyokomi",                                                   // TODO: circleCIのトリガーになった人
			AuthorLink: "https://github.com/kyokomi",                                // TODO: circleCIのトリガーになった人
			AuthorIcon: "https://avatars0.githubusercontent.com/u/1456047?v=3&s=48", // TODO: circleCIのトリガーになった人
			Fields: []slack.AttachmentField{
				{
					Title: "内容",
					Value: "アレコレをそれこれ変更した", // TODO: pullRequestの内容
				},
				{
					Title: "no test files",
					Value: cs.ReportNoTestFiles(),
					Short: true,
				},
				{
					Title: "平均カバレッジ",
					Value: cs.ReportCoverage(),
					Short: true,
				},
			},
			Color: reportColor,
		},
	}

	return c.PostMessage("", params)
}
