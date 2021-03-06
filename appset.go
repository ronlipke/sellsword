package sellsword

import (
	"fmt"
	"github.com/fatih/color"
	"io/ioutil"
	"os"
	"strings"
)

type AppSet struct {
	Apps []*App
	Home string
}

func NewAppSet(home string) (*AppSet, error) {
	as := new(AppSet)
	if _, err := os.Stat(home); os.IsNotExist(err) {
		return as, err
	}
	as.Home = home
	return as, nil
}

func (as *AppSet) FindApps(appNames ...string) error {
	if _, err := os.Stat(as.Home); os.IsNotExist(err) {
		red := GetTermPrinterF(color.FgRed)
		Logger.Errorln(red("The Home directory that you have specified, %s, does not exist.", as.Home))
	} else {
		if appNames[0] == "all" {
			di, _ := ioutil.ReadDir(as.Home)
			for i := range di {
				if di[i].Name() != "config" {
					name := strings.Split(di[i].Name(), ".ssw")[0]
					a, _ := NewApp(name, as.Home)
					as.Apps = append(as.Apps, a)
				}
			}
		} else {
			for i := range appNames {
				a, _ := NewApp(appNames[i], as.Home)
				as.Apps = append(as.Apps, a)
			}
		}
	}
	return nil
}

func (as *AppSet) ListApps(appNames []string) {
	if len(appNames) == 0 {
		as.FindApps("all")
	} else {
		as.FindApps(appNames...)
	}
	for i := range as.Apps {
		cyan := GetTermPrinter(color.FgCyan)
		red := GetTermPrinter(color.FgRed)
		green := GetTermPrinter(color.FgGreen)
		fmt.Printf("%s:\n", cyan(as.Apps[i].Name))
		current, err := as.Apps[i].Current()
		if err != nil {
			fmt.Printf("%s\n", red("No environment currently in use"))
		} else {
			fmt.Printf("\t%s\t%s\n", green(current.Name), green("CURRENT"))
		}
		envs := as.Apps[i].ListEnvs()
		for i := range envs {
			if envs[i].Name != current.Name {
				fmt.Printf("\t%s\n", envs[i].Name)
			}
		}
	}
}
