package build

import (
	"github.com/hpcng/warewulf/internal/pkg/assets"
	"github.com/hpcng/warewulf/internal/pkg/overlay"
	"github.com/hpcng/warewulf/internal/pkg/wwlog"
	"github.com/spf13/cobra"
	"os"
)

func CobraRunE(cmd *cobra.Command, args []string) error {


	var updateNodes []assets.NodeInfo

	if len(args) > 0 && BuildAll == false {
		nodes, err := assets.FindAllNodes()
		if err != nil {
			wwlog.Printf(wwlog.ERROR, "Cloud not get nodeList: %s\n", err)
			os.Exit(1)
		}

		for _, node := range nodes {
			if SystemOverlay == true && node.SystemOverlay == args[0] {
				updateNodes = append(updateNodes, node)
			} else if node.RuntimeOverlay == args[0] {
				updateNodes = append(updateNodes, node)
			}
		}
	} else {
		var err error
		updateNodes, err = assets.FindAllNodes()
		if err != nil {
			wwlog.Printf(wwlog.ERROR, "Cloud not get nodeList: %s\n", err)
			os.Exit(1)
		}
	}

	wwlog.Printf(wwlog.DEBUG, "Checking on system overlay update\n")
	if SystemOverlay == true || BuildAll == true {
		wwlog.Printf(wwlog.INFO, "Updating System Overlays...\n")
		err := overlay.SystemBuild(updateNodes, true)
		if err != nil {
			wwlog.Printf(wwlog.WARN, "Some system overlays failed to be generated\n")
		}
	}

	wwlog.Printf(wwlog.DEBUG, "Checking on system overlay update\n")
	if SystemOverlay == false || BuildAll == true {
		wwlog.Printf(wwlog.INFO, "Updating Runtime Overlays...\n")
		err := overlay.RuntimeBuild(updateNodes, true)
		if err != nil {
			wwlog.Printf(wwlog.WARN, "Some runtime overlays failed to be generated\n")
		}
	}

	return nil
}