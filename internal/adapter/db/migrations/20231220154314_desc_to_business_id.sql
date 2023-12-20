-- Add a new column 'business_id' in 'exp_desc_to_business' to hold the foreign key reference
ALTER TABLE "exp_desc_to_business"
ADD COLUMN "business_id" bigint;

-- Add a foreign key constraint to the 'business_id' column
ALTER TABLE "exp_desc_to_business"
ADD CONSTRAINT "exp_desc_to_business_business_id_fk"
FOREIGN KEY ("business_id") REFERENCES "exp_businesses"("id")
ON DELETE NO ACTION ON UPDATE NO ACTION;

-- Drop the old foreign key constraint on 'business_name'
ALTER TABLE "exp_desc_to_business" DROP CONSTRAINT "exp_desc_to_business_business_name_exp_businesses_name_fk";

-- Drop the 'business_name' column as it's no longer needed
ALTER TABLE "exp_desc_to_business" DROP COLUMN "business_name";
