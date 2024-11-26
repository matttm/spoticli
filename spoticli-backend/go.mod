module github.com/matttm/spoticli/spoticli-backend

go 1.23.1

require (
	github.com/aws/aws-sdk-go-v2 v1.32.0
	github.com/aws/aws-sdk-go-v2/config v1.27.41
	github.com/aws/aws-sdk-go-v2/service/s3 v1.65.0
	github.com/go-sql-driver/mysql v1.8.1
	github.com/gorilla/mux v1.8.1
	github.com/matttm/spoticli/spoticli-models v0.0.0
	github.com/rs/cors v1.11.1
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.6.6 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.17.39 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.16.15 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.19 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.19 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.1 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.3.19 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.12.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.4.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.12.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.18.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.24.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.28.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.32.0 // indirect
	github.com/aws/smithy-go v1.22.0 // indirect
	github.com/coder/flog v1.1.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/fatih/color v1.13.0 // indirect
	github.com/gookit/color v1.5.4 // indirect
	github.com/gookit/goutil v0.6.17 // indirect
	github.com/gookit/gsr v0.1.0 // indirect
	github.com/gookit/slog v0.5.7 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/testify v1.9.0 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/xo/terminfo v0.0.0-20220910002029-abceb7e1c41e // indirect
	go.uber.org/mock v0.5.0 // indirect
	golang.org/x/sync v0.8.0 // indirect
	golang.org/x/sys v0.25.0 // indirect
	golang.org/x/text v0.18.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/matttm/spoticli/spoticli-models => ../spoticli-models
