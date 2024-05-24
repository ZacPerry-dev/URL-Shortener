# UrlShortener

UrlShortener is a simple URL shortening service written in Go (Golang). It provides functionality to shorten URLs, redirect to the original URL using the shortened version, and delete URLs from the database. This project was created primarily as a learning exercise to gain more experience with Golang.

## Features

- **Shorten URLs:** Generate a shortened version of any given URL.
- **Redirect:** Navigate to the original URL using the shortened link.
- **Delete URLs:** Remove URLs from the database.
- **Frontend:** Basic frontend built with HTMX and Go templates.
- **Docker Support:** Easily set up and run the application using Docker Compose.
- **Unit Testing:** Included small tests for the handlers.
- **MongoDB:** Utilizes MongoDB for storing URL data.

## Technologies Used

- **Go:** Backend logic and HTTP server using the standard library.
- **HTMX:** For frontend interactivity.
- **Go Templates:** For server-side rendering of HTML.
- **MongoDB:** Database for storing URLs.
- **Docker:** Containerization and orchestration using Docker Compose.

## Getting Started

### Prerequisites

- [Docker](https://www.docker.com/get-started)
- [Docker Compose](https://docs.docker.com/compose/install/)

### Installation

1. Clone the repository:

    ```bash
    git clone https://github.com/yourusername/UrlShortener.git
    cd UrlShortener
    ```

2. Start the application using Docker Compose:

    ```bash
    docker-compose up
    ```

3. The application should now be running at `http://localhost:8080`.

### Stopping the Application

To stop the application, simply run:

```bash
docker-compose down
```

## Acknowledgements
- [Go](https://golang.org/)
- [HTMX](https://htmx.org/)
- [MongoDB](https://www.mongodb.com/)
- [Inspired by Coding Challenges](https://codingchallenges.fyi/challenges/challenge-url-shortener)
