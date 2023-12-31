package options

// Validate checks Options and return a slice of found errs.
func (o *Options) Validate() []error {
	var errs []error

	errs = append(errs, o.Log.Validate()...)
	errs = append(errs, o.TCP.Validate()...)
	errs = append(errs, o.UnixSocket.Validate()...)
	errs = append(errs, o.ServerRunOptions.Validate()...)

	return errs
}
