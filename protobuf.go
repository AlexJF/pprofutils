package pprofutils

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/google/pprof/profile"
)

type Protobuf struct {
	SampleTypes bool
}

func (p Protobuf) Convert(protobuf *profile.Profile, text io.Writer) error {
	w := bufio.NewWriter(text)
	if p.SampleTypes {
		var sampleTypes []string
		for _, sampleType := range protobuf.SampleType {
			sampleTypes = append(sampleTypes, sampleType.Type+"/"+sampleType.Unit)
		}
		w.WriteString(strings.Join(sampleTypes, " ") + "\n")
	}
	for _, sample := range protobuf.Sample {
		var frames []string
		for i := range sample.Location {
			loc := sample.Location[len(sample.Location)-i-1]
			for j := range loc.Line {
				line := loc.Line[len(loc.Line)-j-1]
				frames = append(frames, line.Function.Name)
			}
		}
		var values []string
		for _, val := range sample.Value {
			values = append(values, fmt.Sprintf("%d", val))
			if !p.SampleTypes {
				break
			}
		}
		fmt.Fprintf(
			w,
			"%s %s\n",
			strings.Join(frames, ";"),
			strings.Join(values, " "),
		)
	}
	return w.Flush()
}
