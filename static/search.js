document.addEventListener("DOMContentLoaded", () => {
    const searchInput = document.getElementById("search-input");
    const suggestionsList = document.getElementById("suggestions");
    const searchButton = document.getElementById("search-button");
    let cachedSuggestions = [];
  
    // Fetch and cache suggestions on page load
    fetch('/search-suggestions?q=')
      .then(response => response.json())
      .then(suggestions => {
        cachedSuggestions = suggestions;
      })
      .catch(err => {
        console.error("Error fetching initial suggestions:", err);
      });
  
    // Debounce function
    function debounce(func, delay) {
      let timeoutId;
      return function (...args) {
        clearTimeout(timeoutId);
        timeoutId = setTimeout(() => func.apply(this, args), delay);
      };
    }
  
    // Filter suggestions based on input
    function filterSuggestions(query) {
      query = query.toLowerCase();
      return cachedSuggestions.filter(suggestion => 
        suggestion.toLowerCase().includes(query)
      );
    }
  
    // Display suggestions
    function displaySuggestions(suggestions) {
      suggestionsList.innerHTML = "";
      const seenSuggestions = new Set();
      
      suggestions.forEach(suggestion => {
        const nameTypeCombo = suggestion.toLowerCase();
        
        if (!seenSuggestions.has(nameTypeCombo)) {
          seenSuggestions.add(nameTypeCombo);
  
          const li = document.createElement("li");
          li.textContent = suggestion;
          li.addEventListener("click", () => {
            searchInput.value = suggestion.split(" - ")[0];
            performSearch(searchInput.value);
          });
          suggestionsList.appendChild(li);
        }
      });
    }
  
    // Throttled search function
    const throttledSearch = debounce((query) => {
      if (query === "") {
        suggestionsList.innerHTML = "";
        return;
      }
  
      const filteredSuggestions = filterSuggestions(query);
      displaySuggestions(filteredSuggestions);
    }, 100);
  
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
  
    // // Handle form submission
    // searchInput.closest("form").addEventListener("submit", function(e) {
    //   e.preventDefault();
    //   performSearch(searchInput.value.trim());
    // });
  });