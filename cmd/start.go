package cmd

import (
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/danesparza/authserver/api"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the API and UI services",
	Long:  `Start the API and UI services`,
	Run:   start,
}

func start(cmd *cobra.Command, args []string) {

	//	Create a router and setup our REST endpoints...
	UIRouter := mux.NewRouter()
	ServiceRouter := mux.NewRouter()

	//	Setup our UI routes
	UIRouter.HandleFunc("/", api.ShowUI)

	//	Setup our Service routes
	ServiceRouter.HandleFunc("/", api.HelloWorld)

	//	Setup the CORS options:
	log.Printf("[INFO] Allowed CORS origins: %s\n", viper.GetString("apiservice.allowed-origins"))

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   strings.Split(viper.GetString("apiservice.allowed-origins"), ","),
		AllowCredentials: true,
	}).Handler(ServiceRouter)

	//	Format the bound interface:
	formattedAPIInterface := viper.GetString("apiservice.bind")
	if formattedAPIInterface == "" {
		formattedAPIInterface = "127.0.0.1"
	}

	formattedUIInterface := viper.GetString("uiservice.bind")
	if formattedUIInterface == "" {
		formattedUIInterface = "127.0.0.1"
	}

	//	Log our TLS key information
	log.Printf("[INFO] API TLS cert: %s\n", viper.GetString("apiservice.tlscert"))
	log.Printf("[INFO] API TLS key: %s\n", viper.GetString("apiservice.tlskey"))
	log.Printf("[INFO] UI TLS cert: %s\n", viper.GetString("uiservice.tlscert"))
	log.Printf("[INFO] UI TLS key: %s\n", viper.GetString("uiservice.tlskey"))

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Printf("[INFO] Starting API service: https://%s:%s\n", formattedAPIInterface, viper.GetString("apiservice.port"))
		log.Printf("[ERROR] %v\n", http.ListenAndServeTLS(viper.GetString("apiservice.bind")+":"+viper.GetString("apiservice.port"), viper.GetString("apiservice.tlscert"), viper.GetString("apiservice.tlskey"), corsHandler))
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Printf("[INFO] Starting UI service: https://%s:%s\n", formattedUIInterface, viper.GetString("uiservice.port"))
		log.Printf("[ERROR] %v\n", http.ListenAndServeTLS(viper.GetString("uiservice.bind")+":"+viper.GetString("uiservice.port"), viper.GetString("uiservice.tlscert"), viper.GetString("uiservice.tlskey"), UIRouter))
	}()

	//	Wait for everything to stop...
	wg.Wait()

}

func init() {
	rootCmd.AddCommand(startCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
