package main

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"os"
	"os/user"
	"path"
)

var log = logrus.New()

var sswVersion = "0.0.1"

func runShow(args []string, sswHome string) {
	as := new(AppSet)
	as.Home = sswHome
	green := getTermPrinter(color.FgGreen)
	blue := getTermPrinter(color.FgCyan)
	if len(args) == 0 {
		as.findApps("all")
	} else {
		as.findApps(args[0])
	}
	fmt.Println("Environments in use:")
	for i := range as.Apps {
		env, err := as.Apps[i].Current()
		if err != nil {
			fmt.Printf("%s\tno environment currently configured\n", green(as.Apps[i].Name))
		} else {
			fmt.Printf("%s\t%s\n", green(as.Apps[i].Name), blue(env.Name))
		}
	}
}

func runLoad(args []string, as *AppSet) {
	if len(args) == 0 {
		as.findApps("all")
	} else {
		as.findApps(args...)
	}
	for i := range as.Apps {
		as.Apps[i].Load()
	}
}

func getTermPrinter(colorName color.Attribute) func(...interface{}) string {
	newColor := color.New(colorName)
	newColor.EnableColor()
	return newColor.SprintFunc()
}

func getTermPrinterF(colorName color.Attribute) func(string, ...interface{}) string {
	newColor := color.New(colorName)
	newColor.EnableColor()
	return newColor.SprintfFunc()
}

func mkdirP(directories []string) {
	for dir := range directories {
		_, stat_err := os.Stat(directories[dir])
		if stat_err != nil {
			if _, ok := stat_err.(*os.PathError); ok {
				mkdir_err := os.MkdirAll(directories[dir], 0755)
				if mkdir_err != nil {
					log.Error("Received error " + mkdir_err.Error())
					os.Exit(0)
				} else {
					log.Println("Created directory " + directories[dir])
				}
			} else {
				log.Error("Received error " + stat_err.Error())
				os.Exit(0)
			}
		} else {
			log.Println("Directory ", directories[dir], "already exists")
		}
	}
}

func main() {
	log.Out = os.Stderr
	formatter := &logrus.TextFormatter{}
	formatter.ForceColors = true
	log.Formatter = formatter
	log.Level = logrus.InfoLevel
	var Verbose bool
	usr, _ := user.Current()
	SswHome := path.Join(usr.HomeDir, "/.ssw")
	var sswCmd = &cobra.Command{
		Use:   "sellsword",
		Short: "Sellsword is a generic command-line tool for switching between application environments",
		Long:  `Sellsword is a generic command-line tool for switching between application environments`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Please invoke one of Sellsword's subcommands to get started. Type sellsword help for more information")
		},
	}
	sswCmd.PersistentFlags().StringVarP(&SswHome, "ssw-home", "s", SswHome, "Home directory for Sellsword")
	sswCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")

	sswCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		// for some reason I have to look up the verbose flag rather than just access the Verbose var
		v := sswCmd.Flags().Lookup("verbose").Value.String()
		if v == "true" {
			log.Level = logrus.DebugLevel
		}
	}

	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version number of sellsword",
		Long:  `All software has versions. This is Sellsword's`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Sellsword version %s\n", sswVersion)
		},
	}
	sswCmd.AddCommand(versionCmd)

	// var initCmd = &cobra.Command{
	// 	Use:   "init",
	// 	Short: "Initialize Sellsword",
	// 	Long:  `This command creates the Sellsword directory structure and downloads default configurations from git@github.com:bryanwb/sellsword.git`,
	// 	Run: func(cmd *cobra.Command, args []string) {
	// 		usr, _ := user.Current()
	// 		homedir := path.Join(usr.HomeDir, ".ssw/config")
	// 		awsdir := path.Join(usr.HomeDir, ".ssw/aws")
	// 		chefdir := path.Join(usr.HomeDir, ".ssw/chef")
	// 		mkdirP([]string{homedir, awsdir, chefdir})
	// 	},
	// }
	// sswCmd.AddCommand(initCmd)

	var loadCmd = &cobra.Command{
		Use:   "load",
		Short: "Loads current Sellsword configurations",
		Long:  `This command loads all default environment configurations for use by the shell`,
		Run: func(cmd *cobra.Command, args []string) {
			as := new(AppSet)
			as.Home = SswHome
			runLoad(args, as)
		},
	}
	sswCmd.AddCommand(loadCmd)

	var showCmd = &cobra.Command{
		Use:   "show",
		Short: "Show Sellsword environments",
		Long:  `Show current Sellsword environments`,
		Run: func(cmd *cobra.Command, args []string) {
			runShow(args, SswHome)
		},
	}
	sswCmd.AddCommand(showCmd)

	var listCmd = &cobra.Command{
		Use:   "list [env ...]",
		Short: "list available Sellsword environments",
		Long:  `List available Sellsword environments`,
		Run: func(cmd *cobra.Command, args []string) {
			as := new(AppSet)
			as.Home = SswHome
			as.listApps(args)
		},
	}
	sswCmd.AddCommand(listCmd)

	var useCmd = &cobra.Command{
		Use:   "use app env",
		Short: "Load environment and set it as default for application",
		Long:  `Load environment and set it as default for application`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 2 {
				red := getTermPrinter(color.FgRed)
				fmt.Fprintf(os.Stderr, "%s\n", red("Usage: ssw use app_name environment"))
				fmt.Fprintf(os.Stderr, "%s\n",
					red("Execute `ssw list` to show available applications and environments"))
			} else {
				as := new(AppSet)
				as.Home = SswHome
				appName := args[0]
				envName := args[1]
				as.findApps(appName)
				app := as.Apps[0]
				app.MakeCurrent(envName)
			}
		},
	}
	sswCmd.AddCommand(useCmd)

	sswCmd.Execute()

}
