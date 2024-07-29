module github.com/mitchrodrigues/talent-review-backend

go 1.22

require (
	github.com/gertd/go-pluralize v0.2.1
	github.com/golly-go/golly v0.4.1-0.20240704203050-3c4f3a889fca
	github.com/golly-go/plugins/eventsource v0.0.0-20240709080818-046b88b7850a
	github.com/golly-go/plugins/gql v0.0.0-20240708075728-1dadf2cba570
	github.com/golly-go/plugins/orm v0.0.0-20240716054529-38f92aa8128f
	github.com/golly-go/plugins/passport v0.0.0-20240614194328-e5cc8c23ee48
	github.com/google/uuid v1.6.0
	github.com/graphql-go/graphql v0.8.1
	github.com/jinzhu/gorm v1.9.16
	github.com/lestrrat-go/jwx/v2 v2.1.0
	github.com/mailgun/mailgun-go/v4 v4.12.0
	github.com/osteele/liquid v1.4.0
	github.com/sirupsen/logrus v1.9.3
	github.com/spf13/cobra v1.8.0
	github.com/spf13/viper v1.19.0
	github.com/stretchr/testify v1.9.0
	github.com/workos/workos-go/v4 v4.13.0
	gorm.io/driver/sqlite v1.5.6
	gorm.io/gorm v1.25.10
)

require (
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.3.0 // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/fsnotify/fsnotify v1.7.0 // indirect
	github.com/gabriel-vasile/mimetype v1.4.3 // indirect
	github.com/go-chi/chi/v5 v5.0.8 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.22.0 // indirect
	github.com/goccy/go-json v0.10.3 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/pgx/v5 v5.6.0 // indirect
	github.com/jackc/puddle/v2 v2.2.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/lestrrat-go/blackmagic v1.0.2 // indirect
	github.com/lestrrat-go/httpcc v1.0.1 // indirect
	github.com/lestrrat-go/httprc v1.0.5 // indirect
	github.com/lestrrat-go/iter v1.0.2 // indirect
	github.com/lestrrat-go/option v1.0.1 // indirect
	github.com/lib/pq v1.3.0 // indirect
	github.com/magiconair/properties v1.8.7 // indirect
	github.com/mattn/go-sqlite3 v1.14.22 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/osteele/tuesday v1.0.3 // indirect
	github.com/pelletier/go-toml/v2 v2.2.2 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/rogpeppe/go-internal v1.11.0 // indirect
	github.com/sagikazarmark/locafero v0.6.0 // indirect
	github.com/sagikazarmark/slog-shim v0.1.0 // indirect
	github.com/segmentio/asm v1.2.0 // indirect
	github.com/sourcegraph/conc v0.3.0 // indirect
	github.com/spf13/afero v1.11.0 // indirect
	github.com/spf13/cast v1.6.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/crypto v0.24.0 // indirect
	golang.org/x/exp v0.0.0-20240613232115-7f521ea00fb8 // indirect
	golang.org/x/net v0.23.0 // indirect
	golang.org/x/sync v0.7.0 // indirect
	golang.org/x/sys v0.21.0 // indirect
	golang.org/x/term v0.21.0 // indirect
	golang.org/x/text v0.16.0 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	gorm.io/driver/postgres v1.5.9 // indirect
)

// replace github.com/golly-go/plugins/mongo => ../../golly-go/plugins/mongo
// replace github.com/golly-go/plugins/orm => ../../golly-go/plugins/orm
// replace github.com/golly-go/plugins/eventsource => ../../golly-go/plugins/eventsource

// replace github.com/golly-go/plugins/workers => ../../golly-go/plugins/workers
// replace github.com/golly-go/plugins/kafka => ../../golly-go/plugins/kafka
// replace github.com/golly-go/plugins/passport => ../../golly-go/plugins/passport
// replace github.com/golly-go/golly => ../../golly-go/golly
// replace github.com/tmc/langchaingo => github.com/mitchrodrigues/langchaingo v0.0.0-20231125195403-51a3a0a0f54a
