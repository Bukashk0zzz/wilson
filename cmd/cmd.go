package cmd

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/trntv/wilson/pkg/config"
	"github.com/trntv/wilson/pkg/runner"
	"github.com/trntv/wilson/pkg/scheduler"
	"github.com/trntv/wilson/pkg/task"
	"github.com/trntv/wilson/pkg/watch"
	"io/ioutil"
	"os"
	"strings"
)

var debug, silent bool
var cfg *config.Config

var configFile string

var tasks = make(map[string]*task.Task)
var contexts = make(map[string]*runner.ExecutionContext)
var pipelines = make(map[string]*scheduler.Pipeline)
var watchers = make(map[string]*watch.Watcher)

var cancel = make(chan struct{})
var done = make(chan bool)

func NewRootCommand() *cobra.Command {
	err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}

	cmd := &cobra.Command{
		Use:     "wilson",
		Short:   "Wilson the task runner",
		Version: "0.2.1",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if debug {
				log.SetLevel(log.DebugLevel)
			} else {
				log.SetLevel(log.InfoLevel)
			}

			if silent {
				log.SetOutput(ioutil.Discard)
				quiet = true
			}
		},
	}

	cmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "enable debug")
	cmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "config file to use")
	cmd.PersistentFlags().BoolVarP(&quiet, "silent", "q", false, "silence output")

	err = cmd.MarkPersistentFlagFilename("config", "yaml", "yml")
	if err != nil {
		log.Warning(err)
	}

	cmd.AddCommand(NewListCommand())
	cmd.AddCommand(NewRunCommand())
	cmd.AddCommand(NewWatchCommand())
	cmd.AddCommand(NewAutocompleteCommand(cmd))

	return cmd
}

func Execute() error {
	cmd := NewRootCommand()
	return cmd.Execute()
}

func Abort() {
	close(cancel)
	<-done
}

func parseConfigFlag() string {
	for i, arg := range os.Args {
		if strings.HasPrefix(arg, "--config") || strings.HasPrefix(arg, "-c") {
			file := strings.TrimPrefix(arg, "--config")
			file = strings.TrimPrefix(file, "-c")
			file = strings.TrimLeft(file, " =")
			if file != "" {
				return file
			}

			if len(os.Args) >= i+2 {
				return os.Args[i+1]
			}
		}
	}

	return ""
}

func loadConfig() error {
	var err error
	configFile = parseConfigFlag()
	cfg, err = config.Load(configFile)
	if err != nil {
		return err
	}

	for name, def := range cfg.Tasks {
		tasks[name] = task.BuildTask(def)
		tasks[name].Name = name
	}

	for name, def := range cfg.Contexts {
		contexts[name], err = runner.BuildContext(def, &config.Get().WilsonConfig)
		if err != nil {
			return fmt.Errorf("context %s build failed: %v", name, err)
		}
	}

	for name, stages := range cfg.Pipelines {
		pipelines[name], err = scheduler.BuildPipeline(stages, tasks)
		if err != nil {
			return err
		}
	}

	tr := runner.NewTaskRunner(contexts, make([]string, 0), true, false)
	for name, def := range cfg.Watchers {
		watchers[name], err = watch.BuildWatcher(name, def, tasks[def.Task], tr)
		if err != nil {
			return fmt.Errorf("watcher %s build failed: %v", name, err)
		}
	}

	return nil
}
