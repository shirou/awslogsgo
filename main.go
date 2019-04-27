package main

import (
	"fmt"
	"os"
	"time"

	"github.com/urfave/cli"
)

const VERSION = "0.0.4"

const callDeadLine = 10 * time.Second

type Args struct {
	NoColor bool
	Start   string
	End     string
}

var awsCommandFlags = []cli.Flag{
	cli.StringFlag{
		Name:   "aws-access-key-id",
		Usage:  "aws access key id",
		EnvVar: "AWS_ACCESS_KEY_ID",
	},
	cli.StringFlag{
		Name:   "aws-secret-access-key",
		Usage:  "aws secret access key",
		EnvVar: "AWS_SECRET_ACCESS_KEY",
	},
	cli.StringFlag{
		Name:   "aws-session-token",
		Usage:  "aws session token",
		EnvVar: "AWS_SESSION_TOKEN",
	},
	cli.StringFlag{
		Name:   "profile",
		Usage:  "aws profile",
		EnvVar: "AWS_PROFILE",
	},
	cli.StringFlag{
		Name:   "aws-region",
		Usage:  "aws region",
		EnvVar: "AWS_REGION",
	},
}

var commandGet = cli.Command{
	Name:      "get",
	Usage:     "get log stream",
	ArgsUsage: "log_group_name log_stream_name",
	Action:    runGet,
	Flags: append(awsCommandFlags, []cli.Flag{
		cli.StringFlag{
			Name:  "f,filter-pattern",
			Usage: "A valid CloudWatch Logs filter `PATTERN` to use for filtering the response. If not provided, all the events are matched."},
		cli.BoolFlag{Name: "w,watch", Usage: "Query for new log lines constantly"},
		cli.BoolFlag{Name: "G,no-group", Usage: "Do not display group name"},
		cli.BoolFlag{Name: "S,no-stream", Usage: "Do not display stream name"},
		cli.BoolFlag{
			Name:  "timestamp",
			Usage: "Add creation timestamp to the output"},
		cli.BoolFlag{
			Name:  "ingestion-time",
			Usage: "Add ingestion time to the output"},
		cli.StringFlag{
			Name:  "s,start",
			Value: "5m",
			Usage: "`START` time"},
		cli.StringFlag{Name: "e,end", Usage: "`END` time"},
		cli.BoolFlag{Name: "no-color", Usage: "Do not color output"},
		//		cli.StringFlag{Name: "q,query", Usage: "JMESPath `QUERY` to use in filtering the response data"},
	}...),
}

var commandListGroups = cli.Command{
	Name:   "groups",
	Usage:  "list log groups",
	Action: runListGroups,
	Flags: append(awsCommandFlags, []cli.Flag{
		cli.StringFlag{
			Name:  "p,log-group-prefix",
			Value: "/",
			Usage: "List only groups matching the `PREFIX`"},
	}...),
}

var commandListStreams = cli.Command{
	Name:      "streams",
	Usage:     "list log stream",
	ArgsUsage: "log_group_name",
	Action:    runListStreams,
	Flags: append(awsCommandFlags, []cli.Flag{
		cli.StringFlag{
			Name:  "s,start",
			Value: "1h",
			Usage: "`START` time"},
		cli.StringFlag{Name: "e,end", Usage: "`END` time"},
		cli.StringFlag{
			Name:  "p,log-stream-prefix",
			Usage: "List only stream matching the prefix"},
	}...),
}

func main() {
	app := cli.NewApp()
	app.Name = "awslogsgo"
	app.Usage = "AWS Cloudwatch Logs Reader"
	app.Version = VERSION
	app.Commands = []cli.Command{
		commandGet,
		commandListGroups,
		commandListStreams,
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func runGet(c *cli.Context) error {
	start, err := parseTime(c.String("start"))
	if err != nil {
		return err
	}
	end, err := parseTime(c.String("end"))
	if err != nil {
		return err
	}
	filter := c.String("f")

	group := c.Args().Get(0)
	streamPrefix := c.Args().Get(1)

	ac, err := awsConfig(awsConfigArgs(c))
	if err != nil {
		return err
	}

	r, err := NewCloudwatchLogsReader(ac, group, streamPrefix, filter, start, end)
	if err != nil {
		return err
	}

	stream, err := r.Stream(c.Bool("w"))
	if err != nil {
		return err
	}

	w, err := NewStdoutWriter(
		c.Bool("no-color"),
		c.Bool("G"),
		c.Bool("S"),
		c.Bool("timestamp"),
		c.Bool("ingestion-time"),
	)
	if err != nil {
		return err
	}

	return w.Write(stream)
}

func runListGroups(c *cli.Context) error {
	prefix := c.String("p")

	ac, err := awsConfig(awsConfigArgs(c))
	if err != nil {
		return err
	}

	return ListGroup(ac, prefix)
}

func runListStreams(c *cli.Context) error {
	group := c.Args().Get(0)
	prefix := c.String("p")
	start, err := parseTime(c.String("start"))
	if err != nil {
		return err
	}
	end, err := parseTime(c.String("end"))
	if err != nil {
		return err
	}

	ac, err := awsConfig(awsConfigArgs(c))
	if err != nil {
		return err
	}

	st, err := ListStreams(ac, group, prefix, start, end)
	if err != nil {
		return err
	}
	for _, s := range st {
		fmt.Println(s)
	}
	return nil
}
