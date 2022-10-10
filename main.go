package main

import (
	"context"
	"embed"
	"encoding/json"
	authHttp "github.com/deathsgun/art/auth/http"
	"github.com/deathsgun/art/config"
	configHttp "github.com/deathsgun/art/config/http"
	configModel "github.com/deathsgun/art/config/model"
	"github.com/deathsgun/art/di"
	"github.com/deathsgun/art/i18n"
	"github.com/deathsgun/art/provider"
	providerInit "github.com/deathsgun/art/provider/init"
	"github.com/deathsgun/art/untis"
	untisHttp "github.com/deathsgun/art/untis/http"
	"github.com/deathsgun/art/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"io/fs"
	"net/http"
	"os"
	"strings"
)

//go:embed views
var views embed.FS

func main() {
	app := fiber.New(fiber.Config{
		AppName: "Azubi Report Tool",
		Views:   setupEngine(),
	})
	app.Use(logger.New())
	//TODO: CSRF
	//TODO: CORS
	//TODO: Caching
	app.Static("/assets", "./assets")

	setupDatabase()

	di.Set[untis.IUntisService]("untis", untis.NewService())
	di.Set[i18n.ITranslationService]("translation", i18n.New())
	di.Set[config.IConfigService]("configService", config.New())
	di.Set[provider.IProviderService]("providerService", provider.New())

	providerInit.InitializeProvider()

	authHttp.Initialize(app)
	untisHttp.Initialize(app)
	configHttp.Initialize(app)

	err := app.Listen(":3000")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to start webserver on port 3000")
	}
}

func setupEngine() *html.Engine {
	views, err := fs.Sub(views, "views")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get views folder")
	}
	engine := html.NewFileSystem(http.FS(views), ".gohtml")
	if os.Getenv("DEV_MODE") == "true" {
		engine = html.New("./views", ".gohtml")
		engine.Reload(true)
		engine.Debug(true)
	}
	engine.AddFunc("translate", func(acceptHeader string, id string, args ...any) string {
		translationService := di.Instance[i18n.ITranslationService]("translation")
		return translationService.Translate(context.WithValue(context.Background(), i18n.LanguageCtxKey, acceptHeader), id, args)
	})
	engine.AddFunc("contains", func(slice []any, v any) bool {
		for _, b := range slice {
			if v == b {
				return true
			}
		}
		return false
	})
	engine.AddFunc("lowercase", strings.ToLower)
	engine.AddFunc("struct", func(params ...any) interface{} {
		obj := map[string]any{}
		for i, param := range params {
			if i%2 == 0 {
				continue
			}
			obj[params[i-1].(string)] = param
		}
		data, err := json.Marshal(obj)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to convert to struct")
		}
		var result interface{}
		if err = json.Unmarshal(data, &result); err != nil {
			log.Fatal().Err(err).Msg("Failed to convert to struct")
		}
		return result
	})
	engine.AddFunc("hasCapability", func(name string, capability string) bool {
		providerService := di.Instance[provider.IProviderService]("providerService")
		prov, ok := providerService.GetProvider(name)
		if !ok {
			return false
		}
		var c provider.Capability
		switch capability {
		case "configurable":
			c = provider.Configurable
		case "server":
			c = provider.ConfigServer
		case "username":
			c = provider.ConfigServer
		case "password":
			c = provider.ConfigPassword
		case "import":
			c = provider.TypeImport
		case "export":
			c = provider.TypeExport
		default:
			return false
		}
		return utils.Contains(prov.Capabilities(), c)
	})
	return engine
}

func setupDatabase() {
	db, err := gorm.Open(sqlite.Open("out/art.db"), &gorm.Config{PrepareStmt: true})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to open database")
	}
	err = db.AutoMigrate(&configModel.ProviderConfig{})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to auto migrate provider configs")
	}
	di.Set[*gorm.DB]("database", db)
}
