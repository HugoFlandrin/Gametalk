document.addEventListener('DOMContentLoaded', function () {
    const cardBox = document.getElementById("cardBox");
    const cardBoxMini = document.querySelector(".card-box-mini");
    const images = [
        './static/img/valo.png',
        './static/img/GTA6.png',
        './static/img/LOL.png',
        './static/img/DBZ.png',
    ];
    const contentData = [
        {
            title: "Tournament",
            text: "Esports tournaments are video game competitions in which players or teams compete for prizes. They range from small local events to major international competitions.",
            buttonText: "Discover",
            route: "/TournoiPage"
        },
        {
            title: "ComingSoon",
            text: "Find all the video game releases for December 2025, including Grand Theft Auto VI, Pokémon Legends Z-A and Marvel 1943: Rise of Hydra.",
            buttonText: "Discover",
            route: "/AvenirPage"
        },
        {
            title: "Games",
            text: "Video games are interactive digital entertainments in which players play on different devices such as consoles, PCs and cell phones. They cover a wide range of genres, including action, adventure, strategy and role-playing.",
            buttonText: "Discover",
            route: "/ThemeSelection"
        },
        {
            title: "Wiki",
            text: "Discovered the most historic games like dragon ball Z, Inazuma Eleven, Mario kart",
            buttonText: "Discover",
            route: "/WikiPage"
        }
    ];
    let currentIndex = 0;

    function changeBackground(direction) {
        if (direction === 'next') {
            currentIndex = (currentIndex + 1) % images.length;
        } else if (direction === 'prev') {
            currentIndex = (currentIndex - 1 + images.length) % images.length;
        }

        const newImage = images[currentIndex];
        const newData = contentData[currentIndex];

        cardBox.style.transition = 'background-image 0.5s ease-in-out';
        cardBox.style.backgroundImage = `url("${newImage}")`;


        const buttons = cardBoxMini.querySelectorAll(".myBtn");
        buttons[0].textContent = newData.title;
        buttons[1].textContent = newData.buttonText;
        cardBoxMini.querySelector(".text").textContent = newData.text;

        buttons[0].addEventListener('click', () => {
            window.location.href = newData.route;
        });
        buttons[1].addEventListener('click', () => {
            window.location.href = newData.route;
        });
    }

    cardBox.addEventListener('click', function (event) {
        const rect = cardBox.getBoundingClientRect();
        const posX = event.clientX - rect.left;
        const width = rect.width;
        if (posX > width / 2) {
            changeBackground('next');
        } else {
            changeBackground('prev');
        }
    });


    changeBackground('next');

    const menuLogo = document.getElementById("menu");
    const sidebarMenu = document.getElementById("sidebarMenu");

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
});


document.addEventListener('DOMContentLoaded', function () {
    const game = document.getElementById("Game");
    game.addEventListener('click', function () {
        if (getCookie("session_id") != null && getCookie("session_id") != 0){
            window.location.href = '/ThemeSelection'
        } else {
            window.location.href = '/Login'
        }
    })

    const rumor = document.getElementById("rumor");
    rumor.addEventListener('click', function () {
        if (getCookie("session_id") != null && getCookie("session_id") != 0){
            window.location.href = '/Rumor'
        } else {
            window.location.href = '/Login'
        }
    })

    const message = document.getElementById("Message");
    message.addEventListener('click', function () {
        if (getCookie("session_id") != null && getCookie("session_id") != 0){
            window.location.href = '/Message'
        } else {
            window.location.href = '/Login'
        }
    })
});