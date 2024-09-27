document.addEventListener('DOMContentLoaded', function() {
    // Add event listeners for tabs
    const tabs = document.querySelectorAll('.tab');
    tabs.forEach(tab => {
        tab.addEventListener('click', function() {
            switchTab(this);
        });
    });
});

function switchTab(clickedTab) {
    const tabs = document.querySelectorAll('.tab');
    tabs.forEach(tab => {
        tab.classList.remove('active');
    });
    clickedTab.classList.add('active');
    // Add logic here to show/hide content based on the selected tab
}

document.addEventListener('DOMContentLoaded', function() {
    // Add event listeners for tabs
    const tabs = document.querySelectorAll('.tab');
    tabs.forEach(tab => {
        tab.addEventListener('click', function() {
            switchTab(this);
        });
    });

    // Add scroll event listener
    var topbar = document.getElementById("topbar");
    var mainContent = document.getElementById("main-content");
    window.addEventListener('scroll', function() {
        const topbar = document.querySelector('.topbar');
        const mainContent = document.querySelector('.main-content');
        if (mainContent) {
            const mainContentTop = mainContent.offsetTop;
            if (window.scrollY >= mainContentTop) {
                topbar.classList.add('scrolled');
            } else {
                topbar.classList.remove('scrolled');
            }
        }
    });
});

function openTab(evt, tabName) {
    var i, tabcontent, tablinks;
    tabcontent = document.getElementsByClassName("tab-content");
    for (i = 0; i < tabcontent.length; i++) {
        tabcontent[i].style.display = "none";
    }
    tablinks = document.getElementsByClassName("tab");
    for (i = 0; i < tablinks.length; i++) {
        tablinks[i].className = tablinks[i].className.replace(" active", "");
    }
    document.getElementById(tabName).style.display = "block";
    evt.currentTarget.className += " active";
}