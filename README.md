# Repository Management Echo

## Description
A Go-based application that manages repository pulls using the Echo framework. It provides an API to pull the latest changes from a specified branch of a repository.

## Installation
1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd <repository-directory>
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Create a `.env` file in the root directory with the following variables:
   ```
   REPO_MANAGEMENT_API_KEY=<your_api_key>
   PORT=8080
   ```

## Usage
To run the application, use:
```bash
   go run main.go
```

The application will start on the specified port (default is 8080).

## API Endpoints
- `POST /pull-repo`: Pulls the latest changes from the specified repository and branch. Requires an API key.

## Environment Variables
- `REPO_MANAGEMENT_API_KEY`: API key for authentication.
- `PORT`: Port number for the server to listen on (default is 8080).

## License
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
