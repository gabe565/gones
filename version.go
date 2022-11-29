package main

var (
	Version = "next"
	Commit  = ""
)

func buildVersion() string {
	result := Version
	if Commit != "" {
		result += " (" + Commit + ")"
	}
	return result
}
