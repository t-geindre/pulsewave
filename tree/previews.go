package tree

import (
	"fmt"
)

func AttachOscPreviews(tree Node, names ...string) {
	for _, n := range tree.QueryAll(names...) {
		n.AttachPreview(func() (string, string) {
			var wf string
			for _, wn := range n.QueryAll("Waveform", "Type") {
				wfn := wn.(*selectorNode)
				wf = wfn.CurrentOption().Label()
			}

			gains := n.QueryAll("Gain")
			if len(gains) == 1 {
				gain := gains[0].(*sliderNode).Val()
				return fmt.Sprintf("%s | %.2f", wf, gain), ""
			}

			rates := n.QueryAll("Rate")
			if len(rates) == 1 {
				rate := rates[0].(*sliderNode).Val()
				return fmt.Sprintf("%s | %.2f Hz", wf, rate), ""
			}

			return wf, ""
		})
	}
}

func AttachPreviewToParent(tee Node, names ...string) {
	for _, n := range tee.QueryAll(names...) {
		parent := n.Parent()
		parent.AttachPreview(n.Preview)
	}
}
