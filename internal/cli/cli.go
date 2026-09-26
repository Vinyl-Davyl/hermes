package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/Vinyl-Davyl/hermes/internal/agents"
	"github.com/Vinyl-Davyl/hermes/internal/doctor"
	"github.com/Vinyl-Davyl/hermes/internal/export"
	"github.com/Vinyl-Davyl/hermes/internal/handoff"
	"github.com/Vinyl-Davyl/hermes/internal/importpkg"
	"github.com/Vinyl-Davyl/hermes/internal/manifest"
	"github.com/Vinyl-Davyl/hermes/internal/pack"
	"github.com/Vinyl-Davyl/hermes/internal/version"
	"github.com/Vinyl-Davyl/hermes/web"
)

func New() *cobra.Command {
	root := &cobra.Command{
		Use:   "hermes",
		Short: "Hand work from one coding agent to the next",
		Long: `Hermes is a local CLI. Pack the git state (and an optional session),
then resume in Cursor, Claude Code, Codex, or another agent.

  hermes handoff              pack + resume recipe (default: Cursor)
  hermes handoff cursor       same, explicit destination
  hermes handoff claude cursor
  hermes resume               reuse the latest pack (no glob needed)
  hermes list
  hermes doctor

No MCP. No plugin. No cloud. You do not need a global install —
./bin/hermes works. make install puts hermes on your PATH.`,
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	root.AddCommand(
		handoffCmd(),
		resumeCmd(),
		toCmd(),
		fromCmd(),
		listCmd(),
		doctorCmd(),
		validateCmd(),
		landingCmd(),
		versionCmd(),
		exportCmd(),
		importCmd(),
	)
	return root
}

func agentIDs() string {
	return strings.Join(agents.IDs(), "|")
}

func handoffCmd() *cobra.Command {
	var (
		message      string
		from         string
		to           string
		session      string
		sessionID    string
		includeFiles bool
		output       string
		zip          bool
		decisions    string
		copyPrompt   bool
		openAgent    bool
		noOpen       bool
	)
	cmd := &cobra.Command{
		Use:   "handoff [to] | [from] [to]",
		Short: "Pack this repo and prepare the next agent",
		Long: `One name is WHERE YOU ARE GOING. Two names are FROM then TO.

  You are in Cursor, limit hit, continue in Claude:
    hermes handoff claude
    hermes handoff cursor claude
    hermes handoff --from cursor --to claude

  You are in Antigravity, continue in Cursor:
    hermes handoff antigravity cursor

  Several Cursor chats? List them, then pick one:
    hermes list --agent cursor --here
    hermes handoff cursor claude --id e1884976

Do not write ./handoff-* in zsh. Use hermes resume.`,
		Args: cobra.MaximumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			src, dest, err := parseHandoffArgs(args, from, to)
			if err != nil {
				return err
			}
			res, err := handoff.Run(handoff.Options{
				Message:      message,
				From:         src,
				To:           dest,
				SessionPath:  session,
				SessionID:    sessionID,
				IncludeFiles: includeFiles,
				Output:       output,
				Zip:          zip,
				Decisions:    decisions,
				Copy:         copyPrompt,
				Open:         openAgent,
				NoOpen:       noOpen,
			})
			if err != nil {
				return err
			}
			printHandoff(cmd, res)
			return nil
		},
	}
	cmd.Flags().StringVarP(&message, "message", "m", "", "what you were doing")
	cmd.Flags().StringVar(&from, "from", "", "source agent: auto|none|"+agentIDs())
	cmd.Flags().StringVar(&to, "to", "", "destination agent (default cursor)")
	cmd.Flags().StringVar(&session, "session", "", "explicit session file")
	cmd.Flags().StringVar(&sessionID, "id", "", "session id from hermes list (prefix is enough)")
	cmd.Flags().BoolVar(&includeFiles, "include-files", false, "copy changed files into the pack")
	cmd.Flags().StringVarP(&output, "output", "o", "", "output directory")
	cmd.Flags().BoolVar(&zip, "zip", false, "also write a .zip")
	cmd.Flags().StringVar(&decisions, "decisions", "", "notes the next agent must not redo")
	cmd.Flags().BoolVar(&copyPrompt, "copy", false, "copy PROMPT.md to the clipboard")
	cmd.Flags().BoolVar(&openAgent, "open", false, "start the destination CLI (default when it is on PATH)")
	cmd.Flags().BoolVar(&noOpen, "no-open", false, "do not start the next agent; only write the pack")
	return cmd
}

func parseHandoffArgs(args []string, fromFlag, toFlag string) (from, to string, err error) {
	from = fromFlag
	to = toFlag
	switch len(args) {
	case 0:
		// flags only
	case 1:
		if to == "" {
			to = args[0]
		} else if from == "" {
			from = args[0]
		}
	case 2:
		if from == "" {
			from = args[0]
		}
		if to == "" {
			to = args[1]
		}
	}
	if from == "" {
		from = "auto"
	}
	if to == "" {
		to = "cursor"
	}
	if from != "auto" && from != "none" {
		if _, ok := agents.ByID(from); !ok {
			return "", "", fmt.Errorf("unknown source %q (try: auto, none, %s)", from, agentIDs())
		}
	}
	if _, ok := agents.ByID(to); !ok {
		return "", "", fmt.Errorf("unknown destination %q (try: %s)", to, agentIDs())
	}
	return from, to, nil
}

func printHandoff(cmd *cobra.Command, res handoff.Result) {
	out := cmd.OutOrStdout()
	fmt.Fprintf(out, "packed  %s\n", res.Dir)
	if res.Zip != "" {
		fmt.Fprintf(out, "zip     %s\n", res.Zip)
	}
	if res.FromAgent != "" {
		fmt.Fprintf(out, "from    %s", res.FromAgent)
		if res.SessionTitle != "" {
			fmt.Fprintf(out, "  (%s)", res.SessionTitle)
		}
		if res.SessionID != "" {
			fmt.Fprintf(out, "  id=%s", res.SessionID)
		}
		fmt.Fprintln(out)
	} else {
		fmt.Fprintf(out, "from    git only\n")
	}
	fmt.Fprintf(out, "to      %s\n", res.ToAgent)
	if res.Copied {
		fmt.Fprintf(out, "copied  PROMPT.md → clipboard\n")
	}
	if res.Opened {
		fmt.Fprintf(out, "opened  %s (seeded from the pack)\n", res.ToAgent)
	}
	for _, w := range res.Warnings {
		fmt.Fprintf(cmd.ErrOrStderr(), "warn    %s\n", w)
	}
	fmt.Fprintln(out)
	fmt.Fprintln(out, res.Instructions)
	fmt.Fprintf(out, "\nnext    hermes resume %s\n", res.ToAgent)
}

func resumeCmd() *cobra.Command {
	var (
		target string
		print  bool
	)
	cmd := &cobra.Command{
		Use:   "resume [to]",
		Short: "Reuse the latest pack (no ./handoff-* glob)",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 1 && target == "" {
				target = args[0]
			}
			if target == "" {
				target = "cursor"
			}
			res, err := handoff.Resume("", target)
			if err != nil {
				return err
			}
			if print {
				fmt.Fprintln(cmd.OutOrStdout(), res.Prompt)
				return nil
			}
			fmt.Fprintf(cmd.OutOrStdout(), "pack    %s\n\n", res.PackDir)
			fmt.Fprintln(cmd.OutOrStdout(), res.Instructions)
			return nil
		},
	}
	cmd.Flags().StringVar(&target, "to", "", "destination agent")
	cmd.Flags().BoolVar(&print, "print", false, "print PROMPT.md")
	return cmd
}

func toCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "to <agent>",
		Short: "Same as hermes handoff <agent>",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			res, err := handoff.Run(handoff.Options{From: "auto", To: args[0]})
			if err != nil {
				return err
			}
			printHandoff(cmd, res)
			return nil
		},
	}
	return cmd
}

func fromCmd() *cobra.Command {
	var message string
	cmd := &cobra.Command{
		Use:   "from <agent>",
		Short: "Pack a session from one agent (destination still Cursor unless you handoff)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			res, err := handoff.Run(handoff.Options{From: args[0], To: "cursor", Message: message})
			if err != nil {
				return err
			}
			printHandoff(cmd, res)
			return nil
		},
	}
	cmd.Flags().StringVarP(&message, "message", "m", "", "what you were doing")
	return cmd
}

func exportCmd() *cobra.Command {
	var (
		message      string
		from         string
		session      string
		includeFiles bool
		output       string
		zip          bool
		decisions    string
	)
	cmd := &cobra.Command{
		Use:    "export",
		Short:  "Write a pack only (prefer: hermes handoff)",
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			res, err := export.Run(export.Options{
				Message:       message,
				From:          from,
				SessionPath:   session,
				IncludeFiles:  includeFiles,
				Output:        output,
				Zip:           zip,
				Decisions:     decisions,
				HermesVersion: version.String,
			})
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "packed  %s\n", res.Dir)
			if res.Zip != "" {
				fmt.Fprintf(cmd.OutOrStdout(), "zip     %s\n", res.Zip)
			}
			if res.FromAgent != "" {
				fmt.Fprintf(cmd.OutOrStdout(), "from    %s\n", res.FromAgent)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "branch  %s  (%d files)\n", res.Git.Branch, len(res.Git.ChangedFiles))
			for _, w := range res.Warnings {
				fmt.Fprintf(cmd.ErrOrStderr(), "warn    %s\n", w)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "\nnext    hermes resume cursor\n")
			return nil
		},
	}
	cmd.Flags().StringVarP(&message, "message", "m", "", "what you were doing")
	cmd.Flags().StringVar(&from, "from", "auto", "agent: auto|none|"+agentIDs())
	cmd.Flags().StringVar(&session, "session", "", "explicit session file")
	cmd.Flags().BoolVar(&includeFiles, "include-files", false, "copy changed files into the pack")
	cmd.Flags().StringVarP(&output, "output", "o", "", "output directory")
	cmd.Flags().BoolVar(&zip, "zip", false, "also write a .zip")
	cmd.Flags().StringVar(&decisions, "decisions", "", "notes the next agent must not redo")
	return cmd
}

func importCmd() *cobra.Command {
	var (
		target string
		print  bool
	)
	cmd := &cobra.Command{
		Use:    "import [pack-dir]",
		Short:  "Prepare a pack (prefer: hermes resume)",
		Hidden: true,
		Args:   cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := ""
			if len(args) == 1 {
				path = args[0]
			}
			resolved, err := pack.Resolve(path)
			if err != nil {
				if os.IsNotExist(err) {
					return fmt.Errorf("no handoff pack here. Run: hermes handoff %s", target)
				}
				return err
			}
			res, err := importpkg.Run(resolved, target)
			if err != nil {
				return err
			}
			if print {
				fmt.Fprintln(cmd.OutOrStdout(), res.Prompt)
				return nil
			}
			fmt.Fprintln(cmd.OutOrStdout(), res.Instructions)
			fmt.Fprintln(cmd.OutOrStdout(), "PROMPT.md refreshed in the pack.")
			return nil
		},
	}
	cmd.Flags().StringVar(&target, "to", "cursor", "destination: "+agentIDs())
	cmd.Flags().BoolVar(&print, "print", false, "print PROMPT.md")
	return cmd
}

func listCmd() *cobra.Command {
	var (
		agent string
		query string
		here  bool
	)
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Sessions Hermes can see on this machine",
		RunE: func(cmd *cobra.Command, args []string) error {
			var sessions []agents.Session
			if agent != "" {
				ag, ok := agents.ByID(agent)
				if !ok {
					return fmt.Errorf("unknown agent %q (try: %s)", agent, agentIDs())
				}
				var err error
				sessions, err = ag.Discover()
				if err != nil {
					return err
				}
			} else {
				var err error
				sessions, err = agents.DiscoverAll()
				if err != nil {
					return err
				}
			}
			cwd, _ := os.Getwd()
			filterCwd := ""
			if here {
				filterCwd = cwd
			}
			sessions = agents.Filter(sessions, query, filterCwd)
			if len(sessions) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "no sessions found")
				fmt.Fprintln(cmd.OutOrStdout(), "that is fine — hermes handoff --from none still packs git")
				return nil
			}
			current := ""
			for _, s := range sessions {
				if s.Agent != current {
					current = s.Agent
					fmt.Fprintf(cmd.OutOrStdout(), "\n# %s\n\n", current)
				}
				title := s.Title
				if title == "" {
					title = s.ID
				}
				fmt.Fprintf(cmd.OutOrStdout(), "  %s  %s  %s\n    %s\n",
					s.Modified.Format("2006-01-02 15:04"), s.ID, title, s.Path)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&agent, "agent", "", "filter: "+agentIDs())
	cmd.Flags().StringVarP(&query, "query", "q", "", "search title, id, or path")
	cmd.Flags().BoolVar(&here, "here", false, "prefer sessions that look like this directory")
	return cmd
}

func doctorCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "doctor",
		Short: "Check git and which agents Hermes can see",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprint(cmd.OutOrStdout(), doctor.Format(doctor.Run("")))
			return nil
		},
	}
}

func validateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "validate [pack-dir]",
		Short: "Check a handoff pack",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := ""
			if len(args) == 1 {
				path = args[0]
			}
			resolved, err := pack.Resolve(path)
			if err != nil {
				if os.IsNotExist(err) {
					return fmt.Errorf("no handoff pack here. Run: hermes handoff")
				}
				return err
			}
			m, err := manifest.Read(resolved)
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "ok  hermes pack v%s  branch=%s  from=%s\n", m.Version, m.Branch, m.FromAgent)
			return nil
		},
	}
}

func landingCmd() *cobra.Command {
	var addr string
	cmd := &cobra.Command{
		Use:     "site",
		Aliases: []string{"landing"},
		Short:   "Serve the local landing page",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintf(cmd.OutOrStdout(), "hermes site  http://%s\n", addr)
			return web.Serve(addr)
		},
	}
	cmd.Flags().StringVar(&addr, "addr", "127.0.0.1:8787", "listen address")
	return cmd
}

func versionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintf(cmd.OutOrStdout(), "hermes %s\n", version.String)
		},
	}
}

func Main() {
	if err := New().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "hermes:", err)
		os.Exit(1)
	}
}
