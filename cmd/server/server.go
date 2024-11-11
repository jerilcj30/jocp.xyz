package main

import (
	"database/sql"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/csrf"
	_ "github.com/jackc/pgx/v4/stdlib"
	controllers "github.com/jerilcj30/jocp/controllers"
	"github.com/jerilcj30/jocp/migrations"
	"github.com/jerilcj30/jocp/models"
	"github.com/jerilcj30/jocp/templates"
	"github.com/jerilcj30/jocp/views"
	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"
)

func main() {

	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}

	// *********************** set up the database *******************************

	db, err := openDB()
	if err != nil {
		panic(err)
	}
	// We also defer a call to db.Close(), so that the connection pool is closed
	// before the main() function exits.
	defer db.Close()

	err = migrate(db, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	// ************************* set up services *******************************

	userService := models.UserService{
		DB: db,
	}
	sessionService := models.SesssionService{
		DB: db,
	}
	leadService := models.LeadsService{
		DB: db,
	}

	// ************************* set up middleware *******************************

	//csrf middleware - using gorilla/csrf package
	csrfKey := os.Getenv("CSRF_KEY")

	// Retrieve the environment variable as a string
	csrfSecureStr := os.Getenv("CSRF_SECURE")

	// Convert the string to a bool
	csrfSecure, err := strconv.ParseBool(csrfSecureStr)
	if err != nil {
		log.Fatalf("Invalid value for CSRF_SECURE: %v", err)
	}

	csrfMw := csrf.Protect(
		[]byte(csrfKey),
		// Todo: Fix this before deploying to production. It requires https
		csrf.Secure(csrfSecure))

	// ************************ set up default controllers ********************************

	userC := controllers.UserHandler{
		UserService:    &userService,
		SessionService: &sessionService,
	}
	homeC := controllers.StaticHandler{}

	leadC := controllers.LeadsHandler{
		LeadsService: &leadService,
	}

	// ************************ injecting default user directly ********************************

	_, err = userC.UserService.Create("jerilcj3@gmail.com", "MM!!cr0507t")
	if err != nil {
		log.Printf("User already exists")
	}

	// ************************ set up homeServices controllers ********************************

	homeServicesC := controllers.StaticHandler{}

	umw := controllers.UserMiddleware{
		SessionService: &sessionService,
	}

	// ************************* Set up default templates *******************************

	userC.Templates.SignIn = views.Must(views.ParseFS(
		templates.FS,
		"layout.gohtml",
		"signin.gohtml",
	))

	homeC.Templates.New = views.Must(views.ParseFS(
		templates.FS,
		"layout.gohtml",
		"home.gohtml",
	))

	// ************************* Set up homeServices templates *******************************
	homeServicesC.Templates.New = views.Must(views.ParseFS(
		templates.FS,
		"layout.gohtml",
		"homeservices/home.gohtml",
	))

	// ************************* set up default routes **********************************

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(csrfMw)
	r.Use(umw.SetUser)

	r.Get("/", homeC.New)
	r.Get("/signin", userC.SignIn)
	// The below uses chi subrouter
	r.Route("/users/me", func(r chi.Router) {
		r.Use(umw.RequireUser) // require user to sign in to access the below pages
		r.Get("/", userC.CurrentUser)
		/*
			/users/me/hello
			r.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, "Hello")
			})
		*/
	})
	// r.Get("/users/me", userC.CurrentUser)

	r.Post("/signin", userC.ProcessSignIn)
	r.Post("/signout", userC.ProcessSignOut)

	r.Post("/create_lead", leadC.CreateLead)

	assetsHandler := http.FileServer(http.Dir("assets"))
	r.Get("/assets/*", http.StripPrefix("/assets", assetsHandler).ServeHTTP)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})

	// ************************* set up home services routes **********************************
	r.Get("/homeservices", homeServicesC.New)

	// ***************** start the server *******************************************

	fmt.Println("Starting the server on :3000...")
	http.ListenAndServe(os.Getenv("SERVER_ADDRESS"), r)
}

func openDB() (*sql.DB, error) {

	db, err := sql.Open("pgx", os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func migrate(db *sql.DB, migrationFS fs.FS, dir string) error {

	// implements embedded sql migrations
	goose.SetBaseFS(migrationFS)

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("migrate: %w", err)
	}

	if err := goose.Up(db, dir); err != nil {
		return fmt.Errorf("migrate: %w", err)
	}

	return nil
}
