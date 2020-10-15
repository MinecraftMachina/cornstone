module cornstone

go 1.15

replace github.com/dghubble/sling v1.3.0 => github.com/ViRb3/sling v1.3.0-fork

replace github.com/cavaliercoder/grab v2.0.0+incompatible => github.com/ViRb3/grab v1.0.0-new

require (
	github.com/cavaliercoder/grab v2.0.0+incompatible
	github.com/dghubble/sling v1.3.0
	github.com/h2non/filetype v1.1.0
	github.com/mholt/archiver/v3 v3.3.2
	github.com/pkg/errors v0.9.1
	github.com/schollz/progressbar/v3 v3.6.0
	github.com/spf13/cobra v1.0.0
	github.com/spf13/viper v1.7.1
)
