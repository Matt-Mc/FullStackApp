# FullStackApp

This project is a full-stack application that fetches data from the Open Parliament API and displays it in a web interface.

## Project Structure

- **Backend**: Contains the Go server that fetches data from the Open Parliament API and serves it to the frontend.
- **Frontend**: Contains the web interface for displaying the data.

## Running the Application

### Backend

1. Navigate to the `Backend` directory:
   ```bash
   cd Backend
   ```

2. Run the Go server:
   ```bash
   go run main.go
   ```

### Frontend

1. Navigate to the `Frontend` directory:
   ```bash
   cd Frontend
   ```

2. Install dependencies:
   ```bash
   npm install
   ```

3. Start the frontend development server:
   ```bash
   npm run dev
   ```

## Features

- Fetches and displays Canadian parliamentary bills.
- Fetches and displays information about Members of Parliament (MPs).
- Displays a total tally of MPs by party.
- Visualizes MP distribution by party using a pie chart with party-specific colors.

## Technologies Used

- **Backend**: Go
- **Frontend**: React with TypeScript

## License

This project is licensed under the MIT License. 