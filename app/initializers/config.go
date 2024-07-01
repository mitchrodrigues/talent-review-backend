package initializers

import (
	"github.com/golly-go/golly"
)

func configPreboot() error {
	golly.SetName("talent-review")
	golly.SetVersion(0, 0, 1, "")

	// append our config initailizer to the first config slot always
	// you can manually do this in config.Initializers but for saftey i
	// chose todo this here to guarantee its always close to the first one
	golly.RegisterInitializerEx(true, configInitializer)

	return nil
}

// ConfigInitializer initializes various config options
func configInitializer(a golly.Application) error {
	a.Config.SetDefault("bind", "localhost:9008")

	// Cors config
	a.Config.SetDefault("cors", map[string]interface{}{
		"origins": []string{
			"http://localhost:*",
			"http://127.0.0.1:*",
			"http://192.168.0.80:*",
			"https://app.talent-radar.io",
			"https://api.talent-radar.io",
		},
		"headers": []string{
			"Accept",
			"Authorization",
			"Content-Type",
			"X-CSRF-Token",
			"Link",
		},
		"methods": []string{
			"GET",
			"POST",
			"PUT",
			"DELETE",
			"OPTIONS",
			"HEAD",
			"PATCH",
		},
	})
	return nil
}
