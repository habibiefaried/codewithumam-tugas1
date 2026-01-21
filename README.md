# codewithumam-tugas1

A simple Hello World HTTP server built with Go's net/http package.

## Local Development

Run the server locally:

```bash
go run main.go
```

The server will start on port 8080 (or the PORT environment variable if set).

Visit `http://localhost:8080` to see "Hello World!"

## Deploy to Railway using GitHub Actions

### Prerequisites

1. A Railway account (sign up at https://railway.app)
2. A GitHub repository with this code
3. Railway CLI installed (optional, for local setup)

### Setup Steps

1. **Create a Railway Project (IMPORTANT - Do this first!):**
   - Go to https://railway.app and log in
   - Click "New Project"
   - Select "Empty Project" (not GitHub repo - we'll deploy via GitHub Actions)
   - In your new project, click "New Service"
   - Select "Empty Service" or "Dockerfile" (Railway will auto-detect your Go app)
   - Note: You now have a project and service created

2. **Link Your Local Repository (Choose ONE method):**

   **Method A: Using Railway CLI (Recommended)**
   ```bash
   # Make sure you're logged in
   railway login
   
   # Create and link a new project/service (if you don't have one)
   railway init
   # Follow the prompts to create a new project/service or link to existing
   
   # OR if you already created the project via web UI, link to it:
   railway link
   # Select your project and service from the list
   ```

   **Method B: Using Railway Web UI**
   - Go to your Railway project
   - In the service settings, you can link it to your GitHub repo
   - This will automatically create the necessary configuration

3. **Get Railway Token:**
   - Go to Railway Dashboard → Settings → Tokens
   - Click "New Token"
   - Give it a name (e.g., "GitHub Actions Deploy")
   - Copy the token (you'll need it for step 4)

4. **Add Railway Token to GitHub Secrets:**
   - Go to your GitHub repository
   - Navigate to Settings → Secrets and variables → Actions
   - Click "New repository secret"
   - Name: `RAILWAY_TOKEN`
   - Value: Paste your Railway token from step 3
   - Click "Add secret"

5. **Verify Railway Configuration (CRITICAL):**
   - After linking, you should have a `.railway` directory in your project
   - This contains the project and service IDs needed for deployment
   - **IMPORTANT:** You MUST commit this directory to git for GitHub Actions to work
   - Without this directory, Railway CLI will prompt for service selection even with `--service` flag
   - Check if it exists: `ls -la .railway` (should see config.json or similar)
   - If missing, run `railway link` again locally, then commit:
     ```bash
     git add .railway
     git commit -m "Add Railway service link"
     git push
     ```

6. **Deploy:**
   - Push your code (including `.railway` directory) to the `master` branch
   - The GitHub Action will automatically deploy to Railway
   - Or manually trigger via Actions tab → "Deploy to Railway" → "Run workflow"

### Troubleshooting

**Error: "Available options can not be empty"**
- This means you don't have any Railway projects yet
- Solution: Create a project first via Railway web UI (step 1), then try `railway link` again
- Or use `railway init` to create a new project from CLI

**Container stuck on "Starting Container"**
1. **Check the logs:**
   - Go to your Railway project → Service → Deployments
   - Click on the latest deployment
   - Check the "Logs" tab for any errors

2. **Verify the PORT environment variable:**
   - Go to Service → Variables
   - Ensure `PORT` is set (Railway usually sets this automatically)
   - Your app is configured to use `PORT` or default to 8080

3. **Check if the binary is being built correctly:**
   - In the build logs, verify `go build` completed successfully
   - The binary should be named `out` (Railway's default)

4. **Verify the start command:**
   - Service → Settings → "Start Command"
   - Should be `./out` (this is set in `Procfile` and `railway.toml`)

5. **Common fixes:**
   - Make sure your app binds to `0.0.0.0` (your code already does this with `:` prefix)
   - Ensure the app is actually listening (your code does this)
   - Try redeploying after checking logs

6. **Manual verification:**
   ```bash
   # Check if Railway can see your service
   railway status
   
   # View logs
   railway logs
   ```

### Setting Deployment Region

You can configure which region Railway deploys to:

**Via Railway Dashboard (Recommended):**
- Go to your Railway project → Service → Settings
- Scroll down to "Region" section
- Select your preferred region (e.g., `us-west`, `us-east`, `eu-west`, `ap-south`, etc.)
- Save changes

**Via railway.toml (Recommended for Code as Config):**
Add the region to your `railway.toml` file:
```toml
[deploy]
region = "asia-southeast1"
```
This ensures the region is version-controlled and applied consistently.

**Note:** Setting the region in `railway.toml` is the recommended approach as it's version-controlled. The dashboard method also works and persists across deployments.

### Environment Variables

Railway will automatically set the `PORT` environment variable. The server is configured to use it.

### Manual Railway Deployment (Alternative)

If you prefer not to use GitHub Actions:

```bash
# Install Railway CLI
npm i -g @railway/cli

# Login
railway login

# Create and link a project (if not already done)
railway init
# OR link to existing project
railway link

# Deploy
railway up
```

## Endpoints

- `GET /` - Returns "Hello World!"
- `GET /health` - Health check endpoint (returns "OK")