document.addEventListener('DOMContentLoaded', function() {
    // Add event listener for search button
    const searchButton = document.querySelector('.search-button');
    if (searchButton) {
        searchButton.addEventListener('click', performSearch);
    }

    // Add event listener for Enter key in the search input
    const searchInput = document.getElementById('search-input');
    if (searchInput) {
        searchInput.addEventListener('keydown', function(event) {
            if (event.key === 'Enter') {
                performSearch();
            }
        });
    }

    // Add event listeners for tabs
    const tabs = document.querySelectorAll('.tab');
    tabs.forEach(tab => {
        tab.addEventListener('click', function() {
            switchTab(this);
        });
    });

    // Add event listener for "Artists" tab to reset the search
    const artistsButton = document.getElementById('artists-btn');
    if (artistsButton) {
        artistsButton.addEventListener('click', resetArtists);
    }

    // Save the original artist list
    const contentGrid = document.getElementById('content-grid');
    if (contentGrid) {
        contentGrid.setAttribute('data-original-content', contentGrid.innerHTML);
    }
});

function performSearch() {
    const searchInput = document.getElementById('search-input');
    const searchTerm = searchInput.value.toLowerCase();
    const artistCards = document.querySelectorAll('.content-card');

    artistCards.forEach(card => {
        const artistName = card.querySelector('.content-title').textContent.toLowerCase();
        if (artistName.includes(searchTerm)) {
            card.style.display = 'block';
        } else {
            card.style.display = 'none';
        }
    });
}

function switchTab(clickedTab) {
    const tabs = document.querySelectorAll('.tab');
    tabs.forEach(tab => {
        tab.classList.remove('active');
    });
    clickedTab.classList.add('active');
    // Add logic here to show/hide content based on the selected tab
}

function resetArtists() {
    const searchInput = document.getElementById('search-input');
    const contentGrid = document.getElementById('content-grid');

    if (searchInput && contentGrid) {
        // Clear the search input
        searchInput.value = '';

        // Reset the artist list to the original content
        const originalContent = contentGrid.getAttribute('data-original-content');
        contentGrid.innerHTML = originalContent;
    }
}

document.addEventListener("DOMContentLoaded", function() {
    var topbar = document.getElementById("topbar");
    var mainContent = document.getElementById("main-content");

    window.addEventListener("scroll", function() {
        if (window.scrollY > mainContent.offsetTop - topbar.offsetHeight) {
            topbar.classList.add("scrolled");
        } else {
            topbar.classList.remove("scrolled");
        }
    });
});

