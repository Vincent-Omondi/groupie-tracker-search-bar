document.addEventListener("DOMContentLoaded", () => {
  const searchInput = document.getElementById("search-input");
  const suggestionsList = document.getElementById("suggestions");
  const searchButton = document.getElementById("search-button");

  // Debounce function
  function debounce(func, delay) {
    let timeoutId;
    return function (...args) {
      clearTimeout(timeoutId);
      timeoutId = setTimeout(() => func.apply(this, args), delay);
    };
  }

  // Throttled search function
  const throttledSearch = debounce((query) => {
    if (query === "") {
      suggestionsList.innerHTML = ""; // Clear suggestions if input is empty
      return;
    }

    fetch(`/search-suggestions?q=${encodeURIComponent(query)}`)
      .then(response => response.json())
      .then(suggestions => {
        suggestionsList.innerHTML = "";
        const seenSuggestions = new Set();  // Track unique name-type combinations
        
        suggestions.forEach(suggestion => {
          const nameTypeCombo = suggestion.toLowerCase();  // Make case-insensitive comparison
          
          if (!seenSuggestions.has(nameTypeCombo)) {  // Only add if not seen
            seenSuggestions.add(nameTypeCombo);

            const li = document.createElement("li");
            li.textContent = suggestion;
            li.addEventListener("click", () => {
              searchInput.value = suggestion.split(" - ")[0]; // Set the input value to the artist/member name
              performSearch(searchInput.value);
            });
            suggestionsList.appendChild(li);
          }
        });
      })
      .catch(err => {
        console.error("Error fetching suggestions:", err);
      });
  }, 500);

  function performSearch(query) {
    window.location.href = `/?query=${encodeURIComponent(query)}`;
  }

  searchInput.addEventListener("input", function() {
    const query = searchInput.value.trim();
    throttledSearch(query);
  });

  searchButton.addEventListener("click", function() {
    performSearch(searchInput.value.trim());
  });

  document.addEventListener("click", function(event) {
    if (!searchInput.contains(event.target) && !suggestionsList.contains(event.target)) {
      suggestionsList.innerHTML = "";
    }
  });

  // Handle form submission
  searchInput.closest("form").addEventListener("submit", function(e) {
    e.preventDefault();
    performSearch(searchInput.value.trim());
  });
});
