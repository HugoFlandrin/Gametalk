document.addEventListener('DOMContentLoaded', function () {
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
    });

});


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

let bntModify = document.getElementById("modify")

bntModify.addEventListener('click', function () {
    window.location.href = "/ProfileManager"
})
document.addEventListener("DOMContentLoaded", () => {
    const profileImage = document.querySelector('.cardImg');
    const profileButtons = document.querySelectorAll('.myBtnCategorie');
    const modifyProfileButton = document.getElementById('modify');


    profileImage.addEventListener('mouseover', () => {
        profileImage.style.transition = 'transform 0.5s';
        profileImage.style.transform = 'scale(1.1)';
    });

    profileImage.addEventListener('mouseout', () => {
        profileImage.style.transform = 'scale(1)';
    });

  
    profileButtons.forEach(button => {
        button.addEventListener('mouseover', () => {
            button.style.transition = 'transform 0.3s';
            button.style.transform = 'translateY(-5px)';
        });

        button.addEventListener('mouseout', () => {
            button.style.transform = 'translateY(0)';
        });
    });

   
    modifyProfileButton.addEventListener('mouseover', () => {
        modifyProfileButton.style.transition = 'transform 0.3s';
        modifyProfileButton.style.transform = 'rotate(5deg)';
    });

    modifyProfileButton.addEventListener('mouseout', () => {
        modifyProfileButton.style.transform = 'rotate(0)';
    });
});
