package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/joho/godotenv"

	"dkmbackend/internal/config"
	"dkmbackend/internal/db"
	"dkmbackend/internal/models"
	"dkmbackend/internal/repository/mongoimpl"
)

func main() {
	_ = godotenv.Load()
	cfg := config.Load()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := db.Connect(ctx, cfg.MongoURI)
	if err != nil {
		log.Fatalf("mongo connect: %v", err)
	}
	defer client.Disconnect(context.Background())

	database := client.Database(cfg.MongoDB)

	blogRepo := mongoimpl.NewBlogRepository(database)
	prodRepo := mongoimpl.NewProductRepository(database)
	careerRepo := mongoimpl.NewCareerRepository(database)

	// sample blogs
	blogs := []models.Blog{
		{Title: "Understanding Generic Medicines", Excerpt: "Generic medicines are just as effective as brand-name medicines but cost significantly less. Learn more about how they work.", Content: "Full content here...", Author: "Dr. Sharma", Date: "2023-10-15", Image: "https://images.unsplash.com/photo-1584308666744-24d5c474f2ae?auto=format&fit=crop&q=80&w=800", Category: "Education", Slug: "understanding-generic-medicines"},
		{Title: "The Importance of Vaccination", Excerpt: "Vaccines are crucial for preventing serious diseases. Discover why staying up-to-date with vaccinations is vital for public health.", Content: "Full content here...", Author: "Nurse Rina", Date: "2023-11-02", Image: "https://images.unsplash.com/photo-1633613286991-611fe299c4be?auto=format&fit=crop&q=80&w=800", Category: "Health Awareness", Slug: "importance-of-vaccination"},
		{Title: "Healthy Living Tips for Winter", Excerpt: "Winter brings specific health challenges. Here are some tips to stay healthy and active during the colder months.", Content: "Full content here...", Author: "Wellness Team", Date: "2023-12-01", Image: "https://images.unsplash.com/photo-1516481265257-97e5f4bc50d5?auto=format&fit=crop&q=80&w=800", Category: "Lifestyle", Slug: "healthy-living-winter"},
	}

	for _, b := range blogs {
		if err := blogRepo.Create(context.Background(), &b); err != nil {
			fmt.Printf("blog insert error: %v\n", err)
		} else {
			fmt.Printf("inserted blog: %s\n", b.Title)
		}
	}

	// sample products
	products := []models.Product{
		{Name: "Amoxyclav-625", Description: "Amoxycillin 500mg + Clavulanic Acid 125mg", Category: "Antibiotics", Image: "https://images.unsplash.com/photo-1584308666744-24d5c474f2ae?auto=format&fit=crop&q=80&w=800", Price: 0.0},
		{Name: "Paracetamol-500", Description: "Paracetamol 500mg", Category: "Pain Relief", Image: "https://images.unsplash.com/photo-1584308666744-24d5c474f2ae?auto=format&fit=crop&q=80&w=800", Price: 0.0},
	}
	for _, p := range products {
		if err := prodRepo.Create(context.Background(), &p); err != nil {
			fmt.Printf("product insert error: %v\n", err)
		} else {
			fmt.Printf("inserted product: %s\n", p.Name)
		}
	}

	// sample careers
	careers := []models.Career{
		{Title: "Sales Executive", Location: "Mumbai", Type: "Full-Time", Description: "Responsible for pharma sales.", Active: true},
		{Title: "Quality Assurance", Location: "Delhi", Type: "Full-Time", Description: "QA role.", Active: true},
	}
	for _, c := range careers {
		if err := careerRepo.Create(context.Background(), &c); err != nil {
			fmt.Printf("career insert error: %v\n", err)
		} else {
			fmt.Printf("inserted career: %s\n", c.Title)
		}
	}

	fmt.Println("seeding complete")
}
