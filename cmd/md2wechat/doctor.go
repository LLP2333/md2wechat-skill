package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/geekjourneyx/md2wechat-skill/internal/config"
	"github.com/geekjourneyx/md2wechat-skill/internal/doctor"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Diagnose local md2wechat readiness",
	Long: `Diagnose local md2wechat readiness without calling remote APIs.

doctor checks config loading, default conversion mode, API key presence,
theme compatibility, layout catalog availability, and WeChat draft credential
presence. It does not validate live authentication, upload images, or create drafts.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		report := runDoctor()
		if jsonOutput {
			responseSuccessWith(codeDoctorCompleted, "Doctor report completed", report)
			return nil
		}
		printDoctorReport(report)
		return nil
	},
}

func runDoctor() doctor.Report {
	config.SetQuiet(true)
	loaded, err := config.Load()
	if err != nil {
		return doctor.LoadError(err)
	}
	return doctor.RunLocal(loaded)
}

func printDoctorReport(report doctor.Report) {
	_, _ = fmt.Fprintf(os.Stdout, "md2wechat doctor\n\n")
	_, _ = fmt.Fprintf(os.Stdout, "Overall: %s\n", report.Overall)
	_, _ = fmt.Fprintf(os.Stdout, "Live checks: %v\n\n", report.Live)
	for _, check := range report.Checks {
		_, _ = fmt.Fprintf(os.Stdout, "%s %s - %s\n", strings.ToUpper(check.Status), check.ID, check.Message)
		for _, action := range check.NextActions {
			_, _ = fmt.Fprintf(os.Stdout, "  - %s\n", action)
		}
	}
	_, _ = fmt.Fprintf(os.Stdout, "\nReadiness:\n")
	_, _ = fmt.Fprintf(os.Stdout, "  format_local: %v\n", report.Readiness.FormatLocal)
	_, _ = fmt.Fprintf(os.Stdout, "  format_api: %v\n", report.Readiness.FormatAPI)
	_, _ = fmt.Fprintf(os.Stdout, "  advanced_layout: %v\n", report.Readiness.AdvancedLayout)
	_, _ = fmt.Fprintf(os.Stdout, "  draft: %v\n", report.Readiness.Draft)
}
