package initializer

const TEMPLATE_FOLDER_PATH string = "template"
const FOLDER_STRUCTURE_PATH string = "template/folder-structure.yaml"
const GENERATED_ROOT_FOLDER = "generated"
const PROJECT_ROOT_FOLDER = "root"
const MAX_NAME_LENGTH = 64

var TEMPLATE_REGISTRY = map[string]string{
	".env":                "env.tmpl",
	"env.go":              "util_env.tmpl",
	"stringutils.go":      "stringutils.tmpl",
	"config_models.go":    "config_models.tmpl",
	"config.go":           "config.tmpl",
	"main.go":             "main.tmpl",
	"database.go":         "database.tmpl",
	"factory.go":          "factory.tmpl",
	"redis.go":            "redis.tmpl",
	"middleware.go":       "middleware.tmpl",
	"auth.go":             "auth.tmpl",
	"response_model.go":   "response_model.tmpl",
	"error_handler.go":    "error_handler.tmpl",
	"success_handler.go":  "success_handler.tmpl",
	"pagination.go":       "pagination.tmpl",
	"server_routes.go":    "server_routes.tmpl",
	".gitignore":          "gitignore.tmpl",
	"README.md":           "readme.tmpl",
	"controller.go":       "controller.tmpl",
	"service.go":          "service.tmpl",
	"repository.go":       "repository.tmpl",
	"model.go":            "model.tmpl",
	"dto.go":              "dto.tmpl",
	"routes.go":           "routes.tmpl",
	"entity.go":           "entity.tmpl",
	"validator.go":        "validator.tmpl",
	"custom_validator.go": "custom_validator.tmpl",
	"jwt.go":              "token.tmpl",
	"go.mod":              "gomod.tmpl",
	"Dockerfile":          "dockerfile.tmpl",
	"docker-compose.yaml": "dockercompose.tmpl",
}
