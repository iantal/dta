package service

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/hashicorp/go-hclog"
	"github.com/iantal/dta/internal/utils"
)

// Gradlew handles execution of various gradlew comands
type Gradlew struct {
	log         hclog.Logger
	projectRoot string
}

// NewGradlew created a new Gradlew
func NewGradlew(l hclog.Logger, destPath string) *Gradlew {
	return &Gradlew{l, destPath}
}

// Projects parses the output of `gradlew projects`
func (g *Gradlew) Projects() ([]string, error) {
	if !g.hasGradlew() {
		return nil, fmt.Errorf("Could not execute gradlew")
	}
	g.log.Info("Started parsing projects")
	err, stdout, stderr := utils.CMD("/bin/sh", "gradlew", "projects")
	if err != nil {
		g.log.Error("Error executing [gradlew projects] command")
		return nil, fmt.Errorf("Error executing [gradlew projects] command for %s", g.projectRoot)
	}

	if isFailedCommand(stdout) || len(stderr) != 0 {
		g.log.Error("Received build failure for gradlew projects")
		return nil, fmt.Errorf("Received build failure for [gradlew projects] for %s", g.projectRoot)
	}

	return extractProjects(stdout), nil
}

func isFailedCommand(stdout []string) bool {
	for _, line := range stdout {
		fmt.Println(line)
		if strings.Contains(line, "BUILD SUCCESSFUL") {
			return false
		}
	}
	return true
}

func extractProjects(data []string) []string {
	var res []string
	for _, out := range data {
		p := extractProject(out)
		if p != "" {
			res = append(res, p)
		}
	}
	return res
}

func extractProject(line string) string {
	var re = regexp.MustCompile(`(?m).*Project\s+':(?P<project>.*)'`)

	result := make(map[string]string)
	match := re.FindStringSubmatch(line)
	if match == nil {
		return ""
	}
	for i, name := range re.SubexpNames() {
		if i != 0 && name != "" {
			result[name] = match[i]
		}
	}
	return result["project"]
}


func (g *Gradlew) Dependencies(project string, isSubproject bool) string {
	if isSubproject {
		err, stdout, serr := utils.CMD("/bin/sh", "gradlew", "dependencies")
		if err != nil {
			g.log.Info("Root dep error " + strings.Join(serr, "\n"))

			g.log.Error("Error executing [gradlew dependencies] for " + project)
			return ""
		}
		return strings.Join(stdout, "\n")
	}

	c := project + ":dependencies"
	err, stdout, serr := utils.CMD("/bin/sh", "gradlew", c)
	if err != nil {
		g.log.Info("subproject dep error " + strings.Join(serr, "\n"))

		g.log.Error("Error executing [gradlew dependencies]", "cmd", c)
		return ""
	}
	return strings.Join(stdout, "\n")
}

func (g *Gradlew) hasGradlew() bool {
	os.Chdir(g.projectRoot)
	err, stdout, _ := utils.CMD("ls")
	if err != nil {
		g.log.Error("Error checking presence of gradlew command in", "path", g.projectRoot)
		return false
	}
	for _, out := range stdout {
		if out == "gradlew" {
			return true
		}
	}
	return false
}
