package tree

import (
	"fmt"
	"strings"
)

func AttachOscPreviews(tree Node, names ...string) {
	for _, n := range tree.QueryAll(names...) {
		n.AttachPreview(func() string {
			var wf string
			for _, wn := range n.QueryAll("Waveform", "Type") {
				wfn := wn.(*selectorNode)
				wf = wfn.Options()[int(wfn.Val())].Label()
			}
			gain := n.QueryAll("Gain")[0].(*sliderNode).Val()

			return fmt.Sprintf("%s | %.2f", wf, gain)
		})
	}
}

// AttachNameIfSubNodeVal attaches a preview function to all nodes with the given names
// If subNode value is val, attach parent name to current node preview
// If multiple subNodes are found, join names with join
// If none are found with the value, use none
func AttachNameIfSubNodeVal(tree Node, subNode string, val float32, join, none string, names ...string) {
	for _, n := range tree.QueryAll(names...) {
		n.AttachPreview(func() string {
			var ns []string
			for _, sn := range n.QueryAll(subNode) {
				if vn, ok := sn.(ValueNode); ok && vn.Val() == val {
					ns = append(ns, vn.Parent().Label())
				}
			}

			if len(ns) == 0 {
				return none
			}

			return strings.Join(ns, join)
		})
	}
}

func AttachPreviewToParent(tee Node, names ...string) {
	for _, n := range tee.QueryAll(names...) {
		parent := n.Parent()
		parent.AttachPreview(n.Preview)
	}
}
