package main

import (
	"io"

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

func (c CoverageService) Send(r io.Reader) error {
	// TODO: inputを元にparamsを作る
	text := `
?       _/home/ubuntu/.go_project/src/github.com/hogehogee/fuga    [no test files]
?       _/home/ubuntu/.go_project/src/github.com/hogehogee/fuga/ab    [no test files]
?       _/home/ubuntu/.go_project/src/github.com/hogehogee/fuga/appirater    [no test files]
?       _/home/ubuntu/.go_project/src/github.com/hogehogee/fuga/channels    [no test files]
?       _/home/ubuntu/.go_project/src/github.com/hogehogee/fuga/internal    [no test files]
ok      _/home/ubuntu/.go_project/src/github.com/hogehogee/fuga/legacy    0.062s    coverage: 37.0% of statements
ok      _/home/ubuntu/.go_project/src/github.com/hogehogee/fuga/legacy/models    0.011s    coverage: 0.3% of statements
?       _/home/ubuntu/.go_project/src/github.com/hogehogee/fuga/middleware    [no test files]
ok      _/home/ubuntu/.go_project/src/github.com/hogehogee/fuga/news    0.019s    coverage: 12.3% of statements
ok      _/home/ubuntu/.go_project/src/github.com/hogehogee/fuga/notifications    0.023s    coverage: 37.2% of statements
ok      _/home/ubuntu/.go_project/src/github.com/hogehogee/fuga/papers    0.043s    coverage: 81.8% of statements
?       _/home/ubuntu/.go_project/src/github.com/hogehogee/fuga/platform    [no test files]
?       _/home/ubuntu/.go_project/src/github.com/hogehogee/fuga/tabs    [no test files]
?       _/home/ubuntu/.go_project/src/github.com/hogehogee/fuga/testhelper    [no test files]
ok      _/home/ubuntu/.go_project/src/github.com/hogehogee/fuga/trending    0.020s    coverage: 100.0% of statements
?       _/home/ubuntu/.go_project/src/github.com/hogehogee/fuga/urls    [no test files]
ok      _/home/ubuntu/.go_project/src/github.com/hogehogee/fuga/users    0.018s    coverage: 17.6% of statements
ok      _/home/ubuntu/.go_project/src/github.com/hogehogee/fuga/weather    0.010s    coverage: 16.4% of statements
	`

	params := c.NewDefaultPostMessageParams()
	params.Attachments = []slack.Attachment{
		{
			Pretext:    "<https://honeybadger.io/path/to/event/|CircleCI Coverage Report> - :sushi:", // TODO: CircleCIのURL branch名
			Fallback:   "CircleCI Coverage Report - :sushi: https://honeybadger.io/path/to/event/",   // TODO: CircleCIのURL branch名
			Text:       text,                                                                         // TODO: go test coverage内容
			AuthorName: "kyokomi",                                                                    // TODO: circleCIのトリガーになった人
			AuthorLink: "https://github.com/kyokomi",                                                 // TODO: circleCIのトリガーになった人
			AuthorIcon: "https://avatars0.githubusercontent.com/u/1456047?v=3&s=48",                  // TODO: circleCIのトリガーになった人
			Fields: []slack.AttachmentField{
				{
					Title: "内容",
					Value: "アレコレをそれこれ変更した", // TODO: pullRequestの内容
				},
				{
					Title: "no test files",
					Value: "5/10", // TODO: no test filesのcount
					Short: true,
				},
				{
					Title: "平均カバレッジ",
					Value: "30%", // TODO: coverageの平均 no test filesは0%として計算する
					Short: true,
				},
			},
			Color: attachmentColorNG, // TODO: no test filesは50%以上 or 平均カバレッジが70%以下
		},
	}

	return c.PostMessage("", params)
}
