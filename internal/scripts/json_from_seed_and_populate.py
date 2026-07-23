import os
import json
import boto3
from botocore.config import Config
from pymongo import MongoClient
from dotenv import load_dotenv

# --- 0. LOAD ENVIRONMENT VARIABLES FROM .env FILE ---
load_dotenv()

# --- 1. CONFIGURATION ---
# Cloudflare R2 Config
ACCOUNT_ID = os.getenv('R2_ACCOUNT_ID')
ACCESS_KEY = os.getenv('R2_ACCESS_KEY')
SECRET_KEY = os.getenv('R2_SECRET_KEY')
BUCKET_NAME = os.getenv('R2_BUCKET_NAME')
R2_PUBLIC_BASE_URL = os.getenv('R2_PUBLIC_BASE_URL')

# MongoDB Config
MONGO_URI = os.getenv('MONGO_URI')  # Change this to your MongoDB connection string if needed
DB_NAME = os.getenv('DB_NAME')  # Change this to your database name if needed
COLLECTION_NAME = 'products'

# Local Filepaths
SEED_DIRECTORY = os.getenv('SEED_DIRECTORY', 'C:/Users/salma/Development/Jiyu/Golang/Projects/mkluxe-backend/internal/seed/products')
OUTPUT_JSON_FILE = os.getenv('OUTPUT_JSON_FILE', 'C:/Users/salma/Development/Jiyu/Golang/Projects/mkluxe-backend/internal/seed/products/master_seed.json')

# --- VALIDATION ---
def validate_config():
    """Check that all required environment variables are set"""
    required_vars = ['R2_ACCOUNT_ID', 'R2_ACCESS_KEY', 'R2_SECRET_KEY', 'R2_BUCKET_NAME', 'R2_PUBLIC_BASE_URL', 'MONGO_URI', 'DB_NAME']
    missing_vars = [var for var in required_vars if not os.getenv(var)]
    
    if missing_vars:
        print(f"❌ Error: Missing required environment variables: {', '.join(missing_vars)}")
        print("Please ensure your .env file contains all required variables.")
        exit(1)

# --- 2. INITIALIZE S3 CLIENT ---
s3_client = boto3.client(
    's3',
    endpoint_url=f'https://{ACCOUNT_ID}.r2.cloudflarestorage.com',
    aws_access_key_id=ACCESS_KEY,
    aws_secret_access_key=SECRET_KEY,
    config=Config(signature_version='s3v4')
)

def get_content_type(filename):
    ext = filename.lower()
    if ext.endswith('.png'): return 'image/png'
    if ext.endswith(('.jpg', '.jpeg')): return 'image/jpeg'
    if ext.endswith('.webp'): return 'image/webp'
    if ext.endswith('.gif'): return 'image/gif'
    return 'application/octet-stream'

def stage_1_upload_and_build_json(test_mode):
    print("\n🚀 Stage 1: Uploading images and building unified JSON...")
    
    existing_products = []
    existing_slugs = set()
    
    # SMART CHECK: Read existing master_seed.json to avoid duplicate work
    if os.path.exists(OUTPUT_JSON_FILE):
        print(f"📄 Found existing {os.path.basename(OUTPUT_JSON_FILE)}. Loading memory to skip processed items...")
        try:
            with open(OUTPUT_JSON_FILE, 'r', encoding='utf-8') as f:
                existing_products = json.load(f)
                # Build a set of slugs we already know about
                existing_slugs = {p.get('slug') for p in existing_products if p.get('slug')}
        except Exception as e:
            print(f"⚠️ Warning: Could not read existing {OUTPUT_JSON_FILE}. Starting fresh. Error: {e}")

    # Start our new list with the existing data
    all_products = existing_products.copy()
    processed_count = 0
    skipped_count = 0
    
    # Loop through each category folder
    for category_slug in os.listdir(SEED_DIRECTORY):
        category_path = os.path.join(SEED_DIRECTORY, category_slug)
        if not os.path.isdir(category_path):
            continue
            
        # Loop through each product folder
        for product_slug in os.listdir(category_path):
            # SMART CHECK: If we already processed this slug, skip it entirely
            if product_slug in existing_slugs:
                # print(f"  ⏭️ Skipping: {product_slug} (Already in master_seed.json)")
                skipped_count += 1
                continue

            product_path = os.path.join(category_path, product_slug)
            if not os.path.isdir(product_path):
                continue
                
            json_file_path = os.path.join(product_path, 'data.json')
            if not os.path.exists(json_file_path):
                continue

            print(f"\n📦 Processing NEW product: {product_slug}")
            
            media_list = []
            image_extensions = ('.png', '.jpg', '.jpeg', '.webp', '.gif')
            
            all_files = sorted(os.listdir(product_path))
            image_files = [f for f in all_files if f.lower().endswith(image_extensions)]
            
            # Upload images and build the new media array
            for index, img_filename in enumerate(image_files):
                local_img_path = os.path.join(product_path, img_filename)
                r2_key = f"products/{category_slug}/{product_slug}/{img_filename}"
                
                print(f"  📷 Uploading: {r2_key}")
                
                s3_client.upload_file(
                    local_img_path,
                    BUCKET_NAME,
                    r2_key,
                    ExtraArgs={'ContentType': get_content_type(img_filename)}
                )
                
                public_url = f"{R2_PUBLIC_BASE_URL}/{r2_key}"
                is_primary = True if index == 0 else False
                
                media_list.append({
                    "url": public_url,
                    "alt_text": f"{product_slug.replace('-', ' ').title()} image {index + 1}",
                    "is_primary": is_primary
                })

            # Read the local data.json
            with open(json_file_path, 'r', encoding='utf-8') as f:
                product_data = json.load(f)
            
            # Replace the old media block with the new R2 URLs
            product_data['media'] = media_list
            product_data['category_slug'] = category_slug
            product_data['slug'] = product_slug

            # Add to our master list and memory
            all_products.append(product_data)
            existing_slugs.add(product_slug)
            processed_count += 1
            
            # --- TEST MODE BREAK ---
            if test_mode and processed_count >= 1:
                print("\n🛑 TEST MODE ENABLED: Stopping after 1 new product.")
                break 
        
        # Break the outer loop if test mode is enabled
        if test_mode and processed_count >= 1:
            break

    # Save the updated master list to a single JSON file
    with open(OUTPUT_JSON_FILE, 'w', encoding='utf-8') as f:
        json.dump(all_products, f, indent=4)
        
    print(f"\n✅ Stage 1 Complete!")
    print(f"   - Skipped: {skipped_count} (Already in master JSON)")
    print(f"   - Processed: {processed_count} (New entries)")
    print(f"   - Total saved in {os.path.basename(OUTPUT_JSON_FILE)}: {len(all_products)}")

def stage_2_seed_mongodb():
    print(f"\n🚀 Stage 2: Loading {os.path.basename(OUTPUT_JSON_FILE)} into MongoDB...")
    
    if not os.path.exists(OUTPUT_JSON_FILE):
        print(f"❌ Error: {OUTPUT_JSON_FILE} not found. Run Stage 1 first.")
        return

    # Load the unified JSON file
    with open(OUTPUT_JSON_FILE, 'r', encoding='utf-8') as f:
        all_products = json.load(f)

    # Initialize MongoDB client
    mongo_client = MongoClient(MONGO_URI)
    db = mongo_client[DB_NAME]
    products_collection = db[COLLECTION_NAME]

    # Insert into the database
    synced_count = 0
    for product in all_products:
        product_slug = product.get("slug")
        
        products_collection.update_one(
            {"slug": product_slug},
            {"$set": product},
            upsert=True
        )
        synced_count += 1

    print(f"✅ Stage 2 Complete! {synced_count} products populated in MongoDB.")

if __name__ == "__main__":
    print("========================================")
    print("      MK Luxe Super Seed Script")
    print("========================================\n")
    
    # Validate configuration before proceeding
    validate_config()
    
    # 1. Ask for Test Mode
    test_mode_input = ""
    while test_mode_input not in ['y', 'n']:
        test_mode_input = input("Run in TEST MODE? (Only processes 1 NEW product) [y/n]: ").strip().lower()
    
    is_test_mode = (test_mode_input == 'y')
    
    # 2. Ask for Stage Selection
    stage_selection = ""
    while stage_selection not in ['1', '2', '3']:
        print("\nSelect the stage(s) to run:")
        print("  1 - Stage 1 Only (Upload NEW images to R2 & append to master JSON)")
        print("  2 - Stage 2 Only (Seed MongoDB from master JSON)")
        print("  3 - Both Stages")
        stage_selection = input("Enter your choice [1/2/3]: ").strip()

    # 3. Execute based on choices
    if stage_selection in ['1', '3']:
        stage_1_upload_and_build_json(is_test_mode)
        
    if stage_selection in ['2', '3']:
        stage_2_seed_mongodb()
        
    print("\n🎉 All requested operations complete!\n")
