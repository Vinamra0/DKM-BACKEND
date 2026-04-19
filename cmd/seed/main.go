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
		{Slug: "amoxyclav-625", Name: "Amoxyclav-625", Description: "Amoxycillin 500mg + Clavulanic Acid 125mg tablets for bacterial infection management.", Composition: "Amoxycillin 500mg + Clavulanic Acid 125mg", DosageForm: "Tablet", Packing: "10 tablets", Company: "Cipla Ltd", Category: "Antibiotics", SubCategory: "Penicillin Based", PackageType: "Blister", Tags: []string{"antibiotic", "infection"}, Generics: []string{}, Specifications: map[string]any{"strength": "500/125mg", "route": "Oral", "shelf_life": "24 months", "storage": "15-30°C"}, Image: "https://images.unsplash.com/photo-1584308666744-24d5c474f2ae?auto=format&fit=crop&q=80&w=800", IsActive: true},
		{Slug: "paracetamol-500", Name: "Paracetamol-500", Description: "Paracetamol 500mg tablets for fever reduction and mild-to-moderate pain relief.", Composition: "Paracetamol 500mg", DosageForm: "Tablet", Packing: "15 tablets", Company: "GSK Pharma", Category: "Pain Relief", SubCategory: "Analgesic", PackageType: "Blister", Tags: []string{"pain", "fever"}, Generics: []string{}, Specifications: map[string]any{"strength": "500mg", "route": "Oral", "shelf_life": "36 months", "storage": "20-25°C"}, Image: "https://images.unsplash.com/photo-1584308666744-24d5c474f2ae?auto=format&fit=crop&q=80&w=800", IsActive: true},
		{Slug: "azithromycin-500", Name: "Azithromycin-500", Description: "Azithromycin 500mg tablets commonly prescribed for respiratory and skin infections.", Composition: "Azithromycin 500mg", DosageForm: "Tablet", Packing: "6 tablets", Company: "Pfizer", Category: "Antibiotics", SubCategory: "Macrolide", PackageType: "Blister", Tags: []string{"antibiotic", "infection"}, Generics: []string{}, Specifications: map[string]any{"strength": "500mg", "route": "Oral", "shelf_life": "30 months", "storage": "15-30°C"}, Image: "https://images.unsplash.com/photo-1471864190281-a93a3070b6de?auto=format&fit=crop&q=80&w=800", IsActive: true},
		{Slug: "ibuprofen-400", Name: "Ibuprofen-400", Description: "Ibuprofen 400mg tablets for inflammation, muscle soreness, and joint pain relief.", Composition: "Ibuprofen 400mg", DosageForm: "Tablet", Packing: "10 tablets", Company: "Abbott", Category: "Pain Relief", SubCategory: "NSAID", PackageType: "Blister", Tags: []string{"pain", "inflammation"}, Generics: []string{}, Specifications: map[string]any{"strength": "400mg", "route": "Oral", "shelf_life": "36 months", "storage": "20-25°C"}, Image: "https://images.unsplash.com/photo-1471864190281-a93a3070b6de?auto=format&fit=crop&q=80&w=800", IsActive: true},
		{Slug: "cetirizine-10", Name: "Cetirizine-10", Description: "Cetirizine 10mg anti-allergy tablets for sneezing, itching, and seasonal rhinitis.", Composition: "Cetirizine 10mg", DosageForm: "Tablet", Packing: "10 tablets", Company: "Lupin", Category: "Allergy Care", SubCategory: "Antihistamine", PackageType: "Blister", Tags: []string{"allergy", "antihistamine"}, Generics: []string{}, Specifications: map[string]any{"strength": "10mg", "route": "Oral", "shelf_life": "36 months", "storage": "15-30°C"}, Image: "https://images.unsplash.com/photo-1576602976047-174e57a47881?auto=format&fit=crop&q=80&w=800", IsActive: true},
		{Slug: "pantoprazole-40", Name: "Pantoprazole-40", Description: "Pantoprazole 40mg gastro-resistant tablets for acidity, reflux, and ulcer care.", Composition: "Pantoprazole 40mg", DosageForm: "Tablet", Packing: "10 tablets", Company: "Dr. Reddy's", Category: "Gastro Care", SubCategory: "PPI", PackageType: "Blister", Tags: []string{"gastric", "acidity"}, Generics: []string{}, Specifications: map[string]any{"strength": "40mg", "route": "Oral", "shelf_life": "24 months", "storage": "20-25°C"}, Image: "https://images.unsplash.com/photo-1587854692152-cbe660dbde88?auto=format&fit=crop&q=80&w=800", IsActive: true},
		{Slug: "vitamin-c-1000", Name: "Vitamin-C-1000", Description: "Vitamin C 1000mg effervescent tablets for daily immunity and antioxidant support.", Composition: "Ascorbic Acid 1000mg", DosageForm: "Effervescent Tablet", Packing: "20 tablets", Company: "Vitex", Category: "Supplements", SubCategory: "Vitamins", PackageType: "Tube", Tags: []string{"vitamin", "immunity"}, Generics: []string{}, Specifications: map[string]any{"strength": "1000mg", "route": "Oral", "shelf_life": "24 months", "storage": "15-30°C"}, Image: "https://images.unsplash.com/photo-1616671276441-2f2c277b8bf9?auto=format&fit=crop&q=80&w=800", IsActive: true},
		{Slug: "calcium-d3", Name: "Calcium-D3", Description: "Calcium with Vitamin D3 tablets for bone health, recovery support, and daily supplementation.", Composition: "Calcium Carbonate 500mg + Vitamin D3 250IU", DosageForm: "Tablet", Packing: "30 tablets", Company: "Fortified", Category: "Supplements", SubCategory: "Minerals", PackageType: "Blister", Tags: []string{"calcium", "bone-health"}, Generics: []string{}, Specifications: map[string]any{"calcium": "500mg", "vitamin_d3": "250IU", "route": "Oral", "shelf_life": "36 months", "storage": "20-25°C"}, Image: "https://images.unsplash.com/photo-1616671276441-2f2c277b8bf9?auto=format&fit=crop&q=80&w=800", IsActive: true},
		{Slug: "ors-hydrate", Name: "ORS-Hydrate", Description: "Oral rehydration salts for electrolyte replenishment during dehydration and heat exhaustion.", Composition: "Glucose + Electrolytes", DosageForm: "Powder", Packing: "1L sachet", Company: "Cipla", Category: "Wellness", SubCategory: "Hydration", PackageType: "Sachet", Tags: []string{"hydration", "electrolytes"}, Generics: []string{}, Specifications: map[string]any{"sodium": "75mmol/L", "potassium": "20mmol/L", "route": "Oral", "shelf_life": "36 months", "storage": "15-30°C"}, Image: "https://images.unsplash.com/photo-1607619056574-7b8d3ee536b2?auto=format&fit=crop&q=80&w=800", IsActive: true},
		{Slug: "cough-syrup-dx", Name: "Cough-Syrup-DX", Description: "Dextromethorphan-based cough syrup for dry cough suppression and throat comfort.", Composition: "Dextromethorphan 10mg/5ml", DosageForm: "Syrup", Packing: "100ml", Company: "Pharm-lab", Category: "Respiratory Care", SubCategory: "Cough Suppressant", PackageType: "Bottle", Tags: []string{"cough", "respiratory"}, Generics: []string{}, Specifications: map[string]any{"strength": "10mg/5ml", "route": "Oral", "shelf_life": "24 months", "storage": "20-25°C"}, Image: "https://images.unsplash.com/photo-1584017911766-d451b3d0e843?auto=format&fit=crop&q=80&w=800", IsActive: true},
		{Slug: "diclofenac-gel", Name: "Diclofenac-Gel", Description: "Topical diclofenac gel for back pain, sports strain, and localized inflammation relief.", Composition: "Diclofenac Sodium 1%", DosageForm: "Gel", Packing: "30g", Company: "Reckitt", Category: "Pain Relief", SubCategory: "Topical", PackageType: "Tube", Tags: []string{"pain", "topical"}, Generics: []string{}, Specifications: map[string]any{"strength": "1%", "route": "Topical", "shelf_life": "24 months", "storage": "15-30°C"}, Image: "https://images.unsplash.com/photo-1512069772995-ec65ed45afd6?auto=format&fit=crop&q=80&w=800", IsActive: true},
		{Slug: "montelukast-10", Name: "Montelukast-10", Description: "Montelukast 10mg tablets for allergy-triggered breathing discomfort and seasonal symptoms.", Composition: "Montelukast 10mg", DosageForm: "Tablet", Packing: "10 tablets", Company: "Merck", Category: "Respiratory Care", SubCategory: "Leukotriene Inhibitor", PackageType: "Blister", Tags: []string{"allergy", "respiratory"}, Generics: []string{}, Specifications: map[string]any{"strength": "10mg", "route": "Oral", "shelf_life": "30 months", "storage": "20-25°C"}, Image: "https://images.unsplash.com/photo-1576602976047-174e57a47881?auto=format&fit=crop&q=80&w=800", IsActive: true},
		{Slug: "levocetirizine-5", Name: "Levocetirizine-5", Description: "Levocetirizine 5mg tablets for nighttime allergy symptom control and skin irritation relief.", Composition: "Levocetirizine 5mg", DosageForm: "Tablet", Packing: "10 tablets", Company: "Wockhardt", Category: "Allergy Care", SubCategory: "Antihistamine", PackageType: "Blister", Tags: []string{"allergy", "antihistamine"}, Generics: []string{}, Specifications: map[string]any{"strength": "5mg", "route": "Oral", "shelf_life": "36 months", "storage": "15-30°C"}, Image: "https://images.unsplash.com/photo-1576602976047-174e57a47881?auto=format&fit=crop&q=80&w=800", IsActive: true},
		{Slug: "omeprazole-20", Name: "Omeprazole-20", Description: "Omeprazole 20mg capsules for daily acid control, reflux management, and gastric protection.", Composition: "Omeprazole 20mg", DosageForm: "Capsule", Packing: "10 capsules", Company: "Glaxo", Category: "Gastro Care", SubCategory: "PPI", PackageType: "Blister", Tags: []string{"gastric", "acidity"}, Generics: []string{}, Specifications: map[string]any{"strength": "20mg", "route": "Oral", "shelf_life": "24 months", "storage": "20-25°C"}, Image: "https://images.unsplash.com/photo-1587854692152-cbe660dbde88?auto=format&fit=crop&q=80&w=800", IsActive: true},
		{Slug: "zinc-50", Name: "Zinc-50", Description: "Zinc 50mg tablets to support immunity, recovery, and nutritional balance.", Composition: "Zinc 50mg", DosageForm: "Tablet", Packing: "30 tablets", Company: "Vitex", Category: "Supplements", SubCategory: "Minerals", PackageType: "Blister", Tags: []string{"zinc", "immunity"}, Generics: []string{}, Specifications: map[string]any{"strength": "50mg", "route": "Oral", "shelf_life": "36 months", "storage": "15-30°C"}, Image: "https://images.unsplash.com/photo-1616671276441-2f2c277b8bf9?auto=format&fit=crop&q=80&w=800", IsActive: true},
		{Slug: "protein-plus", Name: "Protein-Plus", Description: "Daily nutrition supplement powder for energy support, recovery, and general wellness.", Composition: "Whey Protein Isolate 20g + Carbs + Vitamins", DosageForm: "Powder", Packing: "500g", Company: "Nutri-health", Category: "Wellness", SubCategory: "Protein Supplement", PackageType: "Jar", Tags: []string{"protein", "nutrition"}, Generics: []string{}, Specifications: map[string]any{"protein": "20g per serving", "serving_size": "25g", "shelf_life": "18 months", "storage": "15-25°C"}, Image: "https://images.unsplash.com/photo-1514996937319-344454492b37?auto=format&fit=crop&q=80&w=800", IsActive: true},
		{Slug: "moxifloxacin-400", Name: "Moxifloxacin-400", Description: "Moxifloxacin 400mg tablets for advanced bacterial infection care under medical supervision.", Composition: "Moxifloxacin 400mg", DosageForm: "Tablet", Packing: "5 tablets", Company: "Bayer", Category: "Antibiotics", SubCategory: "Fluoroquinolone", PackageType: "Blister", Tags: []string{"antibiotic", "infection"}, Generics: []string{}, Specifications: map[string]any{"strength": "400mg", "route": "Oral", "shelf_life": "24 months", "storage": "15-30°C"}, Image: "https://images.unsplash.com/photo-1471864190281-a93a3070b6de?auto=format&fit=crop&q=80&w=800", IsActive: true},
		{Slug: "naproxen-250", Name: "Naproxen-250", Description: "Naproxen 250mg tablets for longer-lasting relief from muscle pain and inflammation.", Composition: "Naproxen 250mg", DosageForm: "Tablet", Packing: "10 tablets", Company: "Sun Pharma", Category: "Pain Relief", SubCategory: "NSAID", PackageType: "Blister", Tags: []string{"pain", "inflammation"}, Generics: []string{}, Specifications: map[string]any{"strength": "250mg", "route": "Oral", "shelf_life": "36 months", "storage": "20-25°C"}, Image: "https://images.unsplash.com/photo-1512069772995-ec65ed45afd6?auto=format&fit=crop&q=80&w=800", IsActive: true},
		{Slug: "saline-nasal-spray", Name: "Saline-Nasal-Spray", Description: "Isotonic saline nasal spray for congestion relief, hydration, and sinus care.", Composition: "Sodium Chloride 0.9%", DosageForm: "Spray", Packing: "15ml", Company: "Nasalmed", Category: "Respiratory Care", SubCategory: "Nasal Spray", PackageType: "Bottle", Tags: []string{"nasal", "congestion"}, Generics: []string{}, Specifications: map[string]any{"strength": "0.9%", "route": "Nasal", "shelf_life": "24 months", "storage": "15-30°C"}, Image: "https://images.unsplash.com/photo-1584017911766-d451b3d0e843?auto=format&fit=crop&q=80&w=800", IsActive: true},
		{Slug: "probiotic-caps", Name: "Probiotic-Caps", Description: "Probiotic capsules formulated to support digestion, gut balance, and recovery after antibiotics.", Composition: "Lactobacillus 10 Billion CFU", DosageForm: "Capsule", Packing: "10 capsules", Company: "Modicare", Category: "Gastro Care", SubCategory: "Probiotic", PackageType: "Blister", Tags: []string{"probiotic", "digestion"}, Generics: []string{}, Specifications: map[string]any{"cfu": "10 Billion per capsule", "strains": "Multiple lactobacillus", "shelf_life": "18 months", "storage": "2-8°C"}, Image: "https://images.unsplash.com/photo-1587854692152-cbe660dbde88?auto=format&fit=crop&q=80&w=800", IsActive: true},
		{Slug: "iron-folate", Name: "Iron-Folate", Description: "Iron and folic acid tablets for deficiency support, energy, and nutritional recovery.", Composition: "Iron 60mg + Folic Acid 500mcg", DosageForm: "Tablet", Packing: "30 tablets", Company: "Ferro Labs", Category: "Supplements", SubCategory: "Iron", PackageType: "Blister", Tags: []string{"iron", "energy"}, Generics: []string{}, Specifications: map[string]any{"iron": "60mg", "folic_acid": "500mcg", "route": "Oral", "shelf_life": "36 months", "storage": "15-30°C"}, Image: "https://images.unsplash.com/photo-1616671276441-2f2c277b8bf9?auto=format&fit=crop&q=80&w=800", IsActive: true},
		{Slug: "hand-sanitizer-500ml", Name: "Hand-Sanitizer-500ml", Description: "Hospital-grade hand sanitizer for hygiene protection in clinics, offices, and homes.", Composition: "Alcohol 70% + Glycerin", DosageForm: "Liquid", Packing: "500ml", Company: "Hygiene Care", Category: "Wellness", SubCategory: "Hygiene", PackageType: "Bottle", Tags: []string{"sanitizer", "hygiene"}, Generics: []string{}, Specifications: map[string]any{"alcohol": "70%", "shelf_life": "30 months", "storage": "15-30°C"}, Image: "https://images.unsplash.com/photo-1584483766114-2cea6facdf57?auto=format&fit=crop&q=80&w=800", IsActive: true},
		{Slug: "vitamin-b-complex", Name: "Vitamin-B-Complex", Description: "Vitamin B complex tablets to support metabolism, nerve health, and daily vitality.", Composition: "B1 + B2 + B3 + B5 + B6 + B12", DosageForm: "Tablet", Packing: "30 tablets", Company: "Vitex", Category: "Supplements", SubCategory: "Vitamins", PackageType: "Blister", Tags: []string{"vitamin-b", "metabolism"}, Generics: []string{}, Specifications: map[string]any{"b_vitamins": "Complete B-complex", "route": "Oral", "shelf_life": "36 months", "storage": "15-30°C"}, Image: "https://images.unsplash.com/photo-1616671276441-2f2c277b8bf9?auto=format&fit=crop&q=80&w=800", IsActive: true},
		{Slug: "anti-fungal-cream", Name: "Anti-Fungal-Cream", Description: "Topical anti-fungal cream for skin irritation, athlete's foot, and fungal rash care.", Composition: "Terbinafine 1%", DosageForm: "Cream", Packing: "15g", Company: "Derma Plus", Category: "Skin Care", SubCategory: "Antifungal", PackageType: "Tube", Tags: []string{"fungal", "skin-care"}, Generics: []string{}, Specifications: map[string]any{"strength": "1%", "route": "Topical", "shelf_life": "24 months", "storage": "15-30°C"}, Image: "https://images.unsplash.com/photo-1550572017-edd951b55104?auto=format&fit=crop&q=80&w=800", IsActive: true},
		{Slug: "aloe-derma-lotion", Name: "Aloe-Derma-Lotion", Description: "Soothing aloe-based lotion for dry skin, irritation relief, and everyday hydration.", Composition: "Aloe Vera 50% + Moisturizing Agents", DosageForm: "Lotion", Packing: "100ml", Company: "Skin Naturals", Category: "Skin Care", SubCategory: "Moisturizer", PackageType: "Bottle", Tags: []string{"aloe", "moisturizer"}, Generics: []string{}, Specifications: map[string]any{"aloe": "50%", "route": "Topical", "shelf_life": "24 months", "storage": "15-30°C"}, Image: "https://images.unsplash.com/photo-1550572017-edd951b55104?auto=format&fit=crop&q=80&w=800", IsActive: true},
		{Slug: "thermo-relief-patch", Name: "Thermo-Relief-Patch", Description: "Heat therapy patch for neck stiffness, joint soreness, and targeted pain comfort.", Composition: "Heat-generating compound", DosageForm: "Patch", Packing: "5 patches", Company: "Heat Therapy Co", Category: "Pain Relief", SubCategory: "Heat Patch", PackageType: "Box", Tags: []string{"heat-therapy", "pain"}, Generics: []string{}, Specifications: map[string]any{"duration": "8 hours", "route": "Topical", "shelf_life": "30 months", "storage": "15-30°C"}, Image: "https://images.unsplash.com/photo-1512069772995-ec65ed45afd6?auto=format&fit=crop&q=80&w=800", IsActive: true},
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
