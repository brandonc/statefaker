package statefaker

// Options holds configuration options for generating fake state
type Options struct {
	NumOutputs          int
	NumResources        int
	MultiInstanceChance int // percentage chance (0-100) that a resource has multiple instances
	MultiInstanceMin    int // minimum number of instances for multi-instance resources
	MultiInstanceMax    int // maximum number of instances for multi-instance resources
	ModuleChance        int // percentage chance (0-100) that a resource appears within a module
}

// Option is a function type for configuring Options
type Option func(*Options)

// DefaultOptions returns the default configuration
func DefaultOptions() Options {
	return Options{
		NumOutputs:          3,
		NumResources:        3,
		MultiInstanceChance: 10, // 10% chance
		MultiInstanceMin:    3,
		MultiInstanceMax:    50,
		ModuleChance:        70, // 70% chance
	}
}

// WithOutputs sets the number of outputs to generate
func WithOutputs(count int) Option {
	return func(opts *Options) {
		opts.NumOutputs = count
	}
}

// WithResources sets the number of resources to generate
func WithResources(count int) Option {
	return func(opts *Options) {
		opts.NumResources = count
	}
}

// WithMultiInstanceChance sets the percentage chance (0-100) that a resource has multiple instances
func WithMultiInstanceChance(percentage int) Option {
	return func(opts *Options) {
		if percentage < 0 {
			percentage = 0
		}
		if percentage > 100 {
			percentage = 100
		}
		opts.MultiInstanceChance = percentage
	}
}

// WithMultiInstanceMin sets the minimum number of instances for multi-instance resources
func WithMultiInstanceMin(min int) Option {
	return func(opts *Options) {
		if min < 1 {
			min = 1
		}
		opts.MultiInstanceMin = min
	}
}

// WithMultiInstanceMax sets the maximum number of instances for multi-instance resources
func WithMultiInstanceMax(max int) Option {
	return func(opts *Options) {
		if max < 1 {
			max = 1
		}
		opts.MultiInstanceMax = max
	}
}

// WithModuleChance sets the percentage chance (0-100) that a resource appears within a module
func WithModuleChance(percentage int) Option {
	return func(opts *Options) {
		if percentage < 0 {
			percentage = 0
		}
		if percentage > 100 {
			percentage = 100
		}
		opts.ModuleChance = percentage
	}
}

// ApplyOptions applies the given options to the base configuration
func ApplyOptions(opts ...Option) Options {
	options := Options{}
	for _, opt := range opts {
		opt(&options)
	}

	// Ensure min <= max for multi-instance configuration
	if options.MultiInstanceMin > options.MultiInstanceMax {
		options.MultiInstanceMax = options.MultiInstanceMin
	}

	return options
}
