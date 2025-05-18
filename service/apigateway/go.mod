module github.com/MamangRust/payment-gateway-monolith-grpc/service/apigateway

go 1.23.4

require (
	github.com/MamangRust/payment-gateway-monolith-grpc/pkg v0.0.0-00010101000000-000000000000
	github.com/MamangRust/payment-gateway-monolith-grpc/shared v0.0.0-00010101000000-000000000000
	github.com/golang-jwt/jwt/v5 v5.2.2
	github.com/labstack/echo-jwt/v4 v4.3.1
	github.com/labstack/echo/v4 v4.13.3
	github.com/spf13/viper v1.20.1
	github.com/swaggo/echo-swagger v1.4.1
	go.uber.org/zap v1.27.0
	golang.org/x/time v0.11.0
	google.golang.org/grpc v1.72.0
)

require (
	github.com/KyleBanks/depth v1.2.1 // indirect
	github.com/fsnotify/fsnotify v1.9.0 // indirect
	github.com/gabriel-vasile/mimetype v1.4.9 // indirect
	github.com/ghodss/yaml v1.0.0 // indirect
	github.com/go-openapi/jsonpointer v0.21.1 // indirect
	github.com/go-openapi/jsonreference v0.21.0 // indirect
	github.com/go-openapi/spec v0.21.0 // indirect
	github.com/go-openapi/swag v0.23.1 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.26.0 // indirect
	github.com/go-viper/mapstructure/v2 v2.2.1 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/labstack/gommon v0.4.2 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/mailru/easyjson v0.9.0 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/pelletier/go-toml/v2 v2.2.4 // indirect
	github.com/sagikazarmark/locafero v0.9.0 // indirect
	github.com/sourcegraph/conc v0.3.0 // indirect
	github.com/spf13/afero v1.14.0 // indirect
	github.com/spf13/cast v1.8.0 // indirect
	github.com/spf13/pflag v1.0.6 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	github.com/swaggo/files/v2 v2.0.2 // indirect
	github.com/swaggo/swag v1.16.4
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.2 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/crypto v0.38.0 // indirect
	golang.org/x/net v0.40.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.25.0 // indirect
	golang.org/x/tools v0.33.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250505200425-f936aa4a68b2 // indirect
	google.golang.org/protobuf v1.36.6
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/MamangRust/payment-gateway-monolith-grpc/shared => ../../shared

replace github.com/MamangRust/payment-gateway-monolith-grpc/pkg => ../../pkg
