document.addEventListener('DOMContentLoaded', function() {
    // Add event listener for search button
    const searchButton = document.querySelector('.search-button');
    if (searchButton) {
        searchButton.addEventListener('click', performSearch);
    }

    // Add event listeners for tabs
    const tabs = document.querySelectorAll('.tab');
    tabs.forEach(tab => {
        tab.addEventListener('click', function() {
            switchTab(this);
        });
    });
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