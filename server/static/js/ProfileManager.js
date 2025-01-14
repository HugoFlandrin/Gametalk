const menuLogo = document.getElementById("menu");
const sidebarMenu = document.getElementById("sidebarMenu");

menuLogo.addEventListener('click', function () {
    console.log("click");
    sidebarMenu.classList.toggle('hidden');
    sidebarMenu.classList.toggle('expanded');
});

document.addEventListener('click', function (event) {
    if (!event.target.closest('.sidebar') && !event.target.matches('.logoHeaderMenu')) {
        sidebarMenu.classList.add('hidden');
        sidebarMenu.classList.remove('expanded');
    }
})


document.addEventListener('DOMContentLoaded', function () {
    const menuLogo = document.getElementById("profile");
    const sidebarMenu = document.getElementById("sidebarMenu-2");

    menuLogo.addEventListener('click', function () {
        console.log("click");
        sidebarMenu.classList.toggle('hidden');
        sidebarMenu.classList.toggle('expanded');
    });

    document.addEventListener('click', function (event) {
        if (!event.target.closest('.sidebar-2') && !event.target.matches('.logoHeader')) {
            sidebarMenu.classList.add('hidden');
            sidebarMenu.classList.remove('expanded');
        }
    });
});


//searchBar
document.addEventListener('DOMContentLoaded', () => {
    const searchIcon = document.getElementById('search-menu');
    const searchBar = document.getElementById('searchHeader');
    const searchInput = document.getElementById('inputSearch');

    searchIcon.addEventListener('click', (event) => {
        event.stopPropagation();
        searchBar.classList.toggle('visible');
        searchBar.classList.toggle('hidden');
        if (searchBar.classList.contains('visible')) {
            searchInput.focus();
        }
    });

    document.addEventListener('click', (event) => {
        if (!searchBar.contains(event.target) && !searchIcon.contains(event.target)) {
            searchBar.classList.add('hidden');
            searchBar.classList.remove('visible');
        }
    });

    searchBar.addEventListener('click', (event) => {
        event.stopPropagation();
    });
});
