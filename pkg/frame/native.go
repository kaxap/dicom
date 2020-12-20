package frame

import (
	"image"
	"image/color"
)

// NativeFrame represents a native image frame
type NativeFrame struct {
	// Data is a slice of pixels, where each pixel can have multiple values
	Data          []int
	Max           int
	Min           int
	Rows          int
	Cols          int
	BitsPerSample int
}

// IsEncapsulated indicates if the frame is encapsulated or not.
func (n *NativeFrame) IsEncapsulated() bool { return false }

// GetNativeFrame returns a NativeFrame from this frame. If the underlying frame
// is not a NativeFrame, ErrorFrameTypeNotPresent will be returned.
func (n *NativeFrame) GetNativeFrame() (*NativeFrame, error) {
	return n, nil
}

// GetEncapsulatedFrame returns ErrorFrameTypeNotPresent, because this struct
// does not hold encapsulated frame data.
func (n *NativeFrame) GetEncapsulatedFrame() (*EncapsulatedFrame, error) {
	return nil, ErrorFrameTypeNotPresent
}

// GetImage returns an image.Image representation the frame, using default
// processing. This default processing is basic at the moment, and does not
// autoscale pixel values or use window width or level info.
func (n *NativeFrame) GetImage() (image.Image, error) {
	i := image.NewGray16(image.Rect(0, 0, n.Cols, n.Rows))
	step := len(n.Data) / n.Rows / n.Cols

	max := n.Data[0]
	min := n.Data[0]
	// get min and max for normalization
	for j := step; j < len(n.Data); j += step {
		if n.Data[j] > max {
			max = n.Data[j]
		}
		if n.Data[j] < min {
			min = n.Data[j]
		}
	}

	for j := 0; j < len(n.Data); j += step {
		i.SetGray16(j%n.Cols, j/n.Rows,
			color.Gray16{Y: uint16(0xFFFF * (1 - float64(n.Data[j]-min)/float64(max-min)))}) // for now, assume we're not overflowing uint16, assume gray image
	}
	return i, nil
}
