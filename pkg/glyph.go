package msdf

import (
	"encoding/json"
	"fmt"
	"os"
)

type Glyph struct {
	Canvas  *Canvas
	Options *GlyphOptions
}

func newGlyph(c *Canvas, m *GlyphOptions) *Glyph {
	return &Glyph{
		Canvas:  c,
		Options: m,
	}
}

func (g *Glyph) setUV(uv Coords) {
	g.Options.AtlasBounds = uv
}

type Metadata struct {
	Altas   MetadataAtlas   `json:"atlas"`
	Metrics MetadataMetrics `json:"metrics"`
	Glyphs  []GlyphOptions  `json:"glyphs"`
}

type MetadataAtlas struct {
	Type          string  `json:"type"`
	DistanceRange float64 `json:"distanceRange"`
	Size          float64 `json:"size"`
	Width         int     `json:"width"`
	Height        int     `json:"height"`
	YOrigin       string  `json:"yOrigin"`
}

type MetadataMetrics struct {
	EmSize             int `json:"emSize"`
	LineHeight         int `json:"lineHeight"`
	Ascender           int `json:"ascender"`
	Descender          int `json:"descender"`
	UnderlineY         int `json:"underlineY"`
	UnderlineThickness int `json:"underlineThickness"`
}

type GlyphOptions struct {
	Unicode     int     `json:"unicode"`
	Advance     float64 `json:"advance"`
	PlaneBounds Coords  `json:"planeBounds"`
	AtlasBounds Coords  `json:"atlasBounds"`
}

type Coords struct {
	Top    float64 `json:"top"`
	Bottom float64 `json:"bottom"`
	Left   float64 `json:"left"`
	Right  float64 `json:"right"`
}

func (m *Metadata) Save(path string) error {
	file, _ := os.Create(path)
	defer file.Close()
	jsonData, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling to JSON: %v", err)
	}
	_, err = file.WriteString(string(jsonData))

	return err
}
