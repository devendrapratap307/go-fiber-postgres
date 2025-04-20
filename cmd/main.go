package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-fibre-postgres/models"
	"github.com/go-fibre-postgres/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type Book struct {
	Author    string `json:"author"`
	Title     string `json:"title"`
	Year      int    `json:"year"`
	Publisher string `json:"publisher"`
}

type Repository struct {
	DB *gorm.DB
}

func (r *Repository) CreateBook(context *fiber.Ctx) error {
	book := Book{}
	err := context.BodyParser(&book)
	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{"message": "Invalid input- request failed."})
		return err
	}
	err = r.DB.Create(&book).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "Failed to create book"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{"message": "Book created successfully"})
	return nil
}

func (r *Repository) DeleteBook(context *fiber.Ctx) error {
	bookModel := models.Books{}
	id := context.Params("id")
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{"message": "ID is required"})
		return nil
	}

	if err := r.DB.Delete(bookModel, id).Error; err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "Failed to delete book"})
		return err
	}
	context.Status(http.StatusOK).JSON(&fiber.Map{"message": "book deleted successfully"})
	return nil
}

func (r *Repository) GetBookById(context *fiber.Ctx) error {
	id := context.Params("id")
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{"message": "ID is required"})
		return nil
	}
	bookModel := &models.Books{}
	if err := r.DB.Where("id = ?", id).First(bookModel).Error; err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "Book not found"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{"message": "Book fetched successfully", "data": bookModel})
	return nil

	// var book Book
	// if err := r.DB.First(&book, id).Error; err != nil {
	// 	return context.Status(404).JSON(fiber.Map{"error": "Book not found"})
	// }
	// return context.JSON(book)
}

func (r *Repository) GetBooks(context *fiber.Ctx) error {
	bookModels := &[]models.Books{}
	err := r.DB.Find(bookModels).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "Failed to retrieve books"})
		return err
	}
	context.Status(http.StatusOK).JSON(&fiber.Map{"message": "Books retrieved successfully", "data": bookModels})
	return nil

	// var books []Book
	// if err := r.DB.Find(&books).Error; err != nil {
	// 	return c.Status(500).JSON(fiber.Map{"error": "Failed to retrieve books"})
	// }

	// return c.JSON(books)
}

func (r *Repository) SetUpRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Post("/create_books", r.CreateBook)
	api.Delete("/delete_book/:id", r.DeleteBook)
	api.Get("/get_book/:id", r.GetBookById)
	api.Get("/books", r.GetBooks)

	// v1 := api.Group("/v1")
	// users := v1.Group("/users")
	// app.Get("/users", r.GetUsers)
	// app.Post("/users", r.CreateUser)
	// app.Get("/users/:id", r.GetUserByID)
	// app.Put("/users/:id", r.UpdateUser)
	// app.Delete("/users/:id", r.DeleteUser)
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config := &storage.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASSWORD"),
		User:     os.Getenv("DB_USER"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}

	db, err := storage.NewConnection(config)
	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
	}

	err = models.MigrateBooks(db)
	if err != nil {
		log.Fatal("Error migrating database: ", err)
	}

	r := Repository{
		DB: db,
	}
	app := fiber.New()
	r.SetUpRoutes(app)

	app.Listen(":3000")
	log.Println("Server is running on port 3000")
}
