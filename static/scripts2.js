document.addEventListener('DOMContentLoaded', function() {
    const searchButton = document.querySelector('.search-button');
    if (searchButton) {
        searchButton.addEventListener('click', performSearch);
    }

    const searchInput = document.getElementById('search-input');
    if (searchInput) {
        searchInput.addEventListener('keydown', function(event) {
            if (event.key === 'Enter') {
                event.preventDefault(); // Prevent form submission
                performSearch();
            }
        });
    }
});

function performSearch() {
    const searchInput = document.getElementById('search-input');
    const searchTerm = searchInput.value.trim();

    if (searchTerm === "") {
        showError("Please enter a search term.");
        return;
    }

    fetch(`/artists?query=${encodeURIComponent(searchTerm)}`)
        .then(response => {
            if (response.ok) {
                return response.text(); 
            } else if (response.status === 404) {
                showError("No artists found matching the search term.");
                return Promise.reject("No results found");
            } else {
                throw new Error("An error occurred while searching.");
            }
        })
        .then(html => {
            const parser = new DOMParser();
            const doc = parser.parseFromString(html, 'text/html');
            const newContent = doc.querySelector('#content-grid');
            const contentGrid = document.getElementById('content-grid');
            if (newContent && contentGrid) {
                contentGrid.innerHTML = newContent.innerHTML;
            }
        })
        .catch(error => {
            console.error('Error:', error);
        });
}

function showError(message) {
    const contentGrid = document.getElementById('content-grid');
    if (contentGrid) {
        contentGrid.innerHTML = `<p class="error-message">${message}</p>`;
    }
}