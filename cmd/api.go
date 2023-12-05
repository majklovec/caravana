package cmd

import (
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/Masterminds/sprig"
	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

type deployment struct {
	DOMAIN   string `yaml:"DOMAIN"`
	TEMPLATE string `yaml:"TEMPLATE"`
}

// apiCmd represents the api command
var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		gin.SetMode(gin.ReleaseMode)
		r := gin.Default()

		r.ForwardedByClientIP = true
		r.SetTrustedProxies([]string{"127.0.0.1"})

		tpl := template.Must(template.New("base").Funcs(sprig.FuncMap()).ParseGlob("views/*.tpl"))

		r.SetHTMLTemplate(tpl)

		// r.LoadHTMLGlob("views/*")
		r.StaticFile("/", "./public/index.html")
		r.GET("/templates", func(context *gin.Context) {
			color.NoColor = true
			body, err := listDirectories(templatesDir, 2)
			if err != nil {
				context.AbortWithError(http.StatusBadRequest, err)
				return
			}
			context.HTML(http.StatusOK, "items.tpl", gin.H{
				"templates": body,
			})

		})
		r.POST("/deploy", func(context *gin.Context) {
			var body deployment

			if err := context.BindJSON(&body); err != nil {
				context.AbortWithError(http.StatusBadRequest, err)
				return
			}

			saveConfig(filepath.Join(configDir, body.DOMAIN+".yaml"), body)

			if err := processConfig(body.DOMAIN); err != nil {
				context.AbortWithError(http.StatusBadRequest, err)
				return
			}

			context.JSON(http.StatusAccepted, &body)
		})
		r.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
		})
		r.Run()
	},
}

func init() {
	rootCmd.AddCommand(apiCmd)

}
