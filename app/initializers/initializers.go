package initializers

import (
	"github.com/golly-go/golly"
	"github.com/golly-go/plugins/orm"
	"github.com/mitchrodrigues/talent-review-backend/app/controllers"
	"github.com/mitchrodrigues/talent-review-backend/app/domains/accounts"
	"github.com/mitchrodrigues/talent-review-backend/app/domains/audits"
	"github.com/mitchrodrigues/talent-review-backend/app/domains/employees"
	"github.com/mitchrodrigues/talent-review-backend/app/domains/reviews"
	"github.com/mitchrodrigues/talent-review-backend/app/utils/esbackend"
	"github.com/mitchrodrigues/talent-review-backend/app/utils/openai"
	"github.com/mitchrodrigues/talent-review-backend/app/utils/workos"
)

// Preboots lists the preboots
var Preboots = []golly.PrebootFunc{
	configPreboot,
}

// Initializers default app initializers - not sure if i like this yet
// id like eto keep the seperated for cleanliness
var Initializers = []golly.GollyAppFunc{
	defaultConfigs,

	workos.Initializer,
	esbackend.Initializer,

	// Todo make this better this sucks
	// the idea of this was to have the configs in advance
	func(app golly.Application) error {
		return orm.Initializer(orm.Config{
			ConnectionName: orm.DefaultConnection,
			Driver:         orm.DriverT(app.Config.GetString("db.driver")),
			Host:           app.Config.GetString("db.host"),
			User:           app.Config.GetString("db.username"),
			Password:       app.Config.GetString("db.password"),
			Database:       app.Config.GetString("db.name"),
			Port:           app.Config.GetInt("db.port"),
		})(app)
	},

	accounts.Initializer,
	employees.Initalizer,
	reviews.Initializer,
	audits.Initialize,

	func(app golly.Application) error {
		return openai.Initailizer(openai.Config{
			Token: app.Config.GetString("openai.token"),
			AIPretendsToBe: "Pretend you are a performance management assistant providing insights and " +
				"recommendations for team growth, performance reviews, and inclusivity.",
			AIScenarioContext: "You are assisting in the Talent Radar application, focusing on team growth " +
				"and performance metrics. Your tasks include generating insights from performance data, offering " +
				"personalized development suggestions, reviewing feedback for inclusivity, predicting future " +
				"performance trends, and ensuring fair and unbiased performance reviews.",
		})
	},

	// kafka.InitializerPublisher,
	controllers.Initializer,
}

func defaultConfigs(app golly.Application) error {
	app.Config.SetDefault("db", map[string]interface{}{
		"host":     "127.0.0.1",
		"port":     5432,
		"username": "app",
		"password": "password",
		"name":     "talent_review",
		"driver":   "postgres",
	})

	return nil
}
