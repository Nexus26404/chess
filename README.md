# Five in a Row Game

A real-time multiplayer Five in a Row (Gomoku) game implemented using Go, Gin, WebSocket, and Vue.js.

## Features

- Real-time gameplay using WebSocket
- Multiplayer support
- Simple and intuitive UI
- Game state synchronization
- Move validation and win detection

## Technology Stack

### Backend
- Go (Golang)
- Gin Web Framework
- Gorilla WebSocket

### Frontend
- Vue.js
- HTML5 Canvas
- WebSocket client

## Setup Instructions

### Prerequisites
- Go 1.16 or higher
- Node.js 14 or higher
- npm or yarn

### Backend Setup
1. Clone the repository
```bash
git clone [your-repository-url]
cd chess
```

2. Install Go dependencies
```bash
go mod tidy
```

3. Run the server
```bash
go run main.go
```

### Frontend Setup
1. Navigate to the frontend directory
```bash
cd frontend
```

2. Install dependencies
```bash
npm install
# or
yarn install
```

3. Run the development server
```bash
npm run serve
# or
yarn serve
```

## Usage

1. Open your browser and navigate to `http://localhost:8080`
2. Create a new game or join an existing one
3. Share the game URL with your opponent
4. Play the game by taking turns placing pieces on the board

## Contributing

Feel free to submit issues and enhancement requests.

## License

MIT License

Copyright (c) 2025 ZZZ

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
