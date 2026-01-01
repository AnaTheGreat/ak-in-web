#!/usr/bin/env python3

with open('main.go', 'r') as f:
    content = f.read()

# We need to replace lines 64-73 (approximately)
# Find the initDB function and fix it
import re

# Pattern to find the connection string lines
pattern = r'(\s*)connStr := fmt\.Stringf\(\s*"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",\s*getEnv\("DB_HOST", "localhost"\),\s*getEnv\("DB_PORT", "5432"\),\s*getEnv\("DB_USER", "postgres"\),\s*getEnv\("DB_PASSWORD", ""\),\s*getEnv\("DB_NAME", "akinweb"\)'

# Replacement with conditional password logic
replacement = '''\tdbHost := getEnv("DB_HOST", "localhost")
\tdbPort := getEnv("DB_PORT", "5432")
\tdbUser := getEnv("DB_USER", "postgres")
\tdbPassword := getEnv("DB_PASSWORD", "")
\tdbName := getEnv("DB_NAME", "akinweb")
	
\t// Build connection string - only include password if it's not empty
\tvar connStr string
\tif dbPassword == "" {
\t\tconnStr = fmt.Sprintf(
\t\t\t"host=%s port=%s user=%s dbname=%s sslmode=disable",
\t\t\tdbHost, dbPort, dbUser, dbName)
\t} else {
\t\tconnStr = fmt.Sprintf(
\t\t\t"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
\t\t\tdbHost, dbPort, dbUser, dbPassword, dbName)
\t}'''

# Do the replacement
new_content = re.sub(pattern, replacement, content, flags=re.DOTALL)

# Write back
with open('main.go', 'w') as f:
    f.write(new_content)

print("Fixed database connection string in main.go")
