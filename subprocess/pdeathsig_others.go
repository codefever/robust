// +build !linux

package subprocess

import "log"

func maybeSetUpPdeathsig(_ interface{}) {
}

func init() {
	log.Printf("WARNING: cannot use Pdeathsig in this OS.")
}
