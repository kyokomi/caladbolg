package main

import (
	"io"

	"github.com/k0kubun/pp"
	"github.com/kyokomi/caladbolg/parse"
	"github.com/kyokomi/caladbolg/slack"
)

type CoverageService struct {
	slack.Client
	DryRun bool
}

func NewCoverageService(s slack.Client) CoverageService {
	return CoverageService{Client: s, DryRun: false}
}

const (
	attachmentColorOK = "#00CC00"
	attachmentColorNG = "#F35A00"
)

func (s CoverageService) Send(r io.Reader) error {
	report, err := parse.ParseReport(r)
	if err != nil {
		return err
	}

	// no test filesは50%以上 or 平均カバレッジが70%以下
	reportColor := attachmentColorNG
	if report.ReportOK() {
		reportColor = attachmentColorOK
	}

	// TODO: なんかいい感じにする
	params := s.NewDefaultPostMessageParams()
	{
		a := s.Client.NewAttachment()
		a.Pretext = "<https://honeybadger.io/path/to/event/|CircleCI Coverage Report> - :sushi:" // TODO: CircleCIのURL branch名
		a.Fallback = "CircleCI Coverage Report - :sushi: https://honeybadger.io/path/to/event/"  // TODO: CircleCIのURL branch名
		a.Text = report.AllText()
		a.AuthorName = "kyokomi"                                                   // TODO: circleCIのトリガーになった人
		a.AuthorLink = "https://github.com/kyokomi"                                // TODO: circleCIのトリガーになった人
		a.AuthorIcon = "https://avatars0.githubusercontent.com/u/1456047?v=3&s=48" // TODO: circleCIのトリガーになった人
		a.Color = reportColor
		{
			af := s.Client.NewAttachmentField()
			af.Title = "内容"
			af.Value = "アレコレをそれこれ変更した" // TODO: pullRequestの内容
			a.Fields = append(a.Fields, af)
			af = s.Client.NewAttachmentField()
			af.Title = "no test files"
			af.Value = report.ReportNoTestFiles()
			af.Short = true
			a.Fields = append(a.Fields, af)
			af = s.Client.NewAttachmentField()
			af.Title = "平均カバレッジ"
			af.Value = report.ReportCoverage()
			af.Short = true
			a.Fields = append(a.Fields, af)
		}
		params.Attachments = append(params.Attachments, a)
	}

	if s.DryRun {
		pp.Println(params)
		return nil
	}
	return s.PostMessage("", params)
}
