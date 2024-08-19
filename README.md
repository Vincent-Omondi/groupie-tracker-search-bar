# groupie-tracker

groupie tracker consists of receiving a given API and manipulating the data contained in it to create a site and display the information about music artists, their concert locations, and related dates. The application is designed to provide users with an interactive experience to explore various artists and their related events. Built with Go, JavaScript, HTML, and CSS, and containerized with Docker.

## Project Overview

![](./static/imgs/background.png)

## Prerequisites

- Go programming language
- Git

**Optional**

- Docker

## Getting Started and Installation

1. **Clone the repository:**

   ```bash
   git clone https://learn.zone01kisumu.ke/git/johnodhiambo0/groupie-tracker.git
   cd groupie-tracker
   ```

2. **Install dependencies:**
   Ensure Go modules are enabled using:

   ```bash
   go mod tidy
   ```

3. **Run the application:**

   ```bash
   go run .
   ```

4. **Containerization with Docker (Optional):**
   - Ensure Docker is installed and running.
   - Build the Docker image:
     ```bash
     docker build -t groupie-tracker .
     ```
   - Run the Docker container:
     ```bash
     docker run -p 8080:8080 groupie-tracker
     ```

## Implementation

The application is structured into several packages:

- `api`: handles API requests and data fetching.
- `controllers`: contains the routing logic and serves HTTP handlers.
- `templates`: holds HTML templates for rendering pages.
- `handlers`: processes requests and renders the appropriate responses.

Key features include:

- Error handling for better user experience.
- Dynamic filtering of artists based on search queries.

## Usage Example

Once the application is running, you can access it at `http://localhost:8080/` to view the list of artists. You can view details of a single artist by clicking on their image or navigating to `/artist/{id}`, replacing `{id}` with the artist’s ID e.g `http://localhost:8080/artist/18`.

![](./static/imgs/artist-details.png)

## Contributors

- ![Vincent](https://github.com/Vincent-Omondi)

- ![Hillary](https://github.com/ombima56)

- ![John](https://github.com/johneliud)

## Acknowledgement

- ![Zone01 Kisumu](https://www.zone01kisumu.ke/)

Feel free to contribute to the project by submitting pull requests or raising issues!
