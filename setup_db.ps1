# PostgreSQL setup automation script

# Configuration
$DB_USER = "adityatomar"
$DB_PASSWORD = "password"
$DB_NAME = "assignment"
$POSTGRES_SUPERUSER = "postgres"

Write-Host "🚀 Starting database setup..."

# Check if PostgreSQL is installed
try {
    $psqlVersion = psql --version
    Write-Host "✅ PostgreSQL found: $psqlVersion"
} catch {
    Write-Host "❌ PostgreSQL not found. Please install PostgreSQL and add it to your PATH"
    exit 1
}

# Create SQL commands file
$sqlCommands = @"
-- Create user if not exists
DO `$`$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = '$DB_USER') THEN
        CREATE USER $DB_USER WITH PASSWORD '$DB_PASSWORD';
    END IF;
END
`$`$;

-- Create database if not exists
SELECT 'CREATE DATABASE $DB_NAME'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = '$DB_NAME')\gexec

-- Grant privileges
GRANT ALL PRIVILEGES ON DATABASE $DB_NAME TO $DB_USER;

-- Connect to the database and set up permissions
\c $DB_NAME

-- Grant schema privileges
GRANT ALL ON SCHEMA public TO $DB_USER;
"@

# Save SQL commands to temporary file
$tempFile = New-TemporaryFile
$sqlCommands | Out-File -FilePath $tempFile -Encoding UTF8

Write-Host "📦 Creating database and user..."

# Run SQL commands
try {
    Get-Content $tempFile | psql -U $POSTGRES_SUPERUSER
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✅ Database and user created successfully"
    } else {
        Write-Host "❌ Error creating database and user"
        exit 1
    }
} catch {
    Write-Host "❌ Error: $_"
    exit 1
} finally {
    Remove-Item $tempFile
}

# Create or update .env file
$envContent = @"
DB_HOST=localhost
DB_PORT=5432
DB_USER=$DB_USER
DB_PASSWORD=$DB_PASSWORD
DB_NAME=$DB_NAME
SERVER_PORT=8080
"@

$envContent | Out-File -FilePath ".env" -Encoding UTF8
Write-Host "✅ Created .env file with database configuration"

Write-Host "`n🎉 Setup complete! You can now run the application with:"
Write-Host "go run main.go"

# Test database connection
Write-Host "`n🔍 Testing database connection..."
$env:PGPASSWORD = $DB_PASSWORD
$testConnection = psql -U $DB_USER -d $DB_NAME -h localhost -c "\conninfo"
if ($LASTEXITCODE -eq 0) {
    Write-Host "✅ Database connection successful!"
} else {
    Write-Host "⚠️ Could not verify database connection. You may need to:"
    Write-Host "  1. Check if PostgreSQL service is running"
    Write-Host "  2. Verify pg_hba.conf allows password authentication"
    Write-Host "  3. Restart PostgreSQL service if you made configuration changes"
}