package commands

import (
	"fmt"
	"io"
	"text/tabwriter"

	cmds "github.com/ipfs/go-ipfs-cmds"
	"github.com/ipfs/go-ipfs/core/commands/cmdenv"

	"github.com/ipfs/go-ipfs-provider/batched"
)

var statProvideCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Returns statistics about the node's (re)provider system.",
		ShortDescription: `
Returns statistics about the content the node is advertising.

This interface is not stable and may change from release to release.
`,
	},
	Arguments: []cmds.Argument{},
	Options:   []cmds.Option{},
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		nd, err := cmdenv.GetNode(env)
		if err != nil {
			return err
		}

		if !nd.IsOnline {
			return ErrNotOnline
		}

		sys, ok := nd.Provider.(*batched.BatchProvidingSystem)
		if !ok {
			return fmt.Errorf("can only return stats if Experimental.AcceleratedDHTClient is enabled")
		}

		stats, err := sys.Stat(req.Context)
		if err != nil {
			return err
		}

		if err := res.Emit(stats); err != nil {
			return err
		}

		return nil
	},
	Encoders: cmds.EncoderMap{
		cmds.Text: cmds.MakeTypedEncoder(func(req *cmds.Request, w io.Writer, s *batched.BatchedProviderStats) error {
			wtr := tabwriter.NewWriter(w, 0, 0, 1, ' ', 0)
			defer wtr.Flush()

			fmt.Fprintf(wtr, "TotalProvides: %d\n", s.TotalProvides)
			fmt.Fprintf(wtr, "AvgProvideDuration: %v\n", s.AvgProvideDuration)
			fmt.Fprintf(wtr, "LastReprovideDuration: %v\n", s.LastReprovideDuration)
			fmt.Fprintf(wtr, "LastReprovideBatchSize: %d\n", s.LastReprovideBatchSize)
			return nil
		}),
	},
	Type: batched.BatchedProviderStats{},
}
