# Groupie Tracker Search Bar

## Overview

The Groupie Tracker Search Bar is designed to provide a powerful search functionality for tracking artists, bands, and related information on a web interface. This project aims to deliver an intuitive user experience for searching specific text inputs such as artist names, members, locations, album dates, and more. The search functionality is responsive, case-insensitive, and includes typing suggestions, ensuring a seamless and efficient user interaction.

## Features

### Search Functionality

The search bar allows users to perform searches on multiple attributes of an artist or band. The key features include:

- **Search Attributes:**
  - Artist/Band Name: Find artists or bands by name.
  - Members: Look up band members.
  - Locations: Search based on performance or consert locations.
  - First Album Date: Discover bands by the release date of their first album.
  - Creation Date: Filter based on the year the band was formed.

- **Case-Insensitive Search:**
  - The search input is treated as case-insensitive, ensuring that users do not have to worry about capitalization when typing in their query.

- **Real-Time Typing Suggestions:**
  - As the user types, the search bar provides instant suggestions, improving user experience and helping them find the right artist or band quickly.
  - Suggestions are pulled from cached data to minimize delay and improve performance.

- **Suggestion Categorization:**
  - Suggestions are categorized to differentiate between various attributes. For example, searching for "phil" may return both Phil Collins - member and Phil Collins - artist/band, clearly indicating whether the result refers to a band name or a member.

- **Debounced Input for Optimized Performance:**
  - The search input is debounced, reducing the number of calls to the backend and improving performance by limiting how frequently suggestions are updated during typing.

- **Keyboard Navigation:**
  - Users can navigate through suggestions using the arrow keys and select a suggestion by pressing the enter key, allowing for seamless keyboard interactions.

### Search Workflow

1. **Fetching Initial Suggestions:**
   - When the page loads, initial suggestions are fetched from `/search-suggestions?q=` and cached to improve performance during user input.

2. **Debounced Search:**
   - When the user types into the search bar, the input is processed with a debounce function to ensure that backend requests are limited while still providing responsive search suggestions.

3. **Filtering Suggestions:**
   - Suggestions are filtered in real-time based on the current input, and displayed in a dropdown. Only unique results are shown to avoid redundancy.

4. **Performing a Search:**
   - Once the user selects a suggestion or presses the enter key, the search is executed, and the results are displayed accordingly.

## Example

Imagine you have created a card system to display the band data. The user can directly search for the band or member they want to see. For example:

1. A user types "Phil" in the search bar.
2. The search bar suggests options like:
   - Phil Collins - member
   - Phil Collins - artist/band
3. The user selects Phil Collins - artist/band and the page is redirected to show details about Phil Collins as an artist.

Here's a visual representation of the search functionality:

![Search Bar Example](/static/img/search.png)

This screenshot demonstrates the real-time suggestions and categorization of results as the user types "Phil" into the search bar.

## Installation and Setup

To run this project locally:

1. Clone the repository:
   ```bash
   git clone https://github.com/Vincent-Omondi/groupie-tracker-search-bar.git
   cd groupie-tracker-search-bar
   ```

2. Run the application:
   ```bash
   go run .
   ```

3. Open your browser and navigate to `http://localhost:8080`.

## Project Structure

- **Controllers:**
  - `handlers.go`: Manages requests, handles artist data, and filters search results.
  - `routes.go`: Registers routes for the application, including the `/artists` and `/search-suggestions` endpoints.
- **API:**
  - `api.go`: Defines the logic to fetch artist, location, and relation data from external APIs.
- **Static:**
  - `search.js`: Implements the search bar functionality, including debounced input, real-time suggestions, and search execution.
- **Templates:**
  - `artists.html`: Renders the search results on the frontend.

  ## Contributors

This project exists thanks to all the people who contribute:

- [Vincent Omondi](https://github.com/Vincent-Omondi)
- [Hillary Ombima](https://github.com/ombima56)
- [John Eluid](https://github.com/johneliud)

## Contributing

We welcome contributions to the Groupie Tracker Search Bar project! If you'd like to contribute, please follow these steps:

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

Please make sure to update tests as appropriate and adhere to the project's coding standards.
