package vfl

// StandardMetrics defines the standard spacing and size values used in VFL.
// These values follow Apple's Human Interface Guidelines for standard spacing.
// Users can override these values by providing custom metrics in ParseOptions.
var StandardMetrics = map[string]float64{
	// Standard spacing values
	"standard":        8.0,  // Default spacing between views
	"aqua-space":      8.0,  // macOS Aqua guidelines standard space
	"aqua-space-small": 4.0,  // macOS Aqua guidelines small space
	"aqua-space-large": 20.0, // macOS Aqua guidelines large space

	// Common spacing presets
	"default-spacing": 8.0,
	"small-spacing":   4.0,
	"medium-spacing":  8.0,
	"large-spacing":   20.0,
	"xl-spacing":      32.0,

	// Standard sizes
	"button-height":   30.0,
	"textfield-height": 22.0,
	"toolbar-height":  40.0,
	"sidebar-width":   200.0,
	"inspector-width": 250.0,

	// iOS specific
	"ios-margin":      16.0,
	"ios-spacing":     8.0,
	"ios-button-height": 44.0,

	// Material Design inspired
	"material-unit":   8.0,
	"material-margin": 16.0,
	"material-gutter": 24.0,
}

// GetMetric returns a standard metric value by name.
// It returns the metric value and true if found, or 0 and false if not found.
// This function is safe for concurrent use.
func GetMetric(name string) (float64, bool) {
	value, exists := StandardMetrics[name]
	return value, exists
}

// RegisterMetric allows registration of custom metrics.
// Custom metrics override standard metrics with the same name.
// This function is NOT safe for concurrent use.
func RegisterMetric(name string, value float64) {
	StandardMetrics[name] = value
}