module registration

go 1.19

require github.com/aws/aws-lambda-go v1.34.1

require data v1.0.0

require (
	github.com/aws/aws-sdk-go v1.44.75 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
)

replace data v1.0.0 => ../data
