document.addEventListener('DOMContentLoaded', () => {
    const tournaments = document.querySelectorAll('.tournament');
    tournaments.forEach((tournament, index) => {
        tournament.style.opacity = 0;
        tournament.style.transform = 'translateY(20px)';
        setTimeout(() => {
            tournament.style.transition = 'opacity 0.5s ease, transform 0.5s ease';
            tournament.style.opacity = 1;
            tournament.style.transform = 'translateY(0)';
        }, index * 100);
    });
});


function displayUser(id) {
    const encodedUser = encodeURIComponent(id);
    window.location.href = `/User?user=${encodedUser}`;
}

function getCookie(name) {
    const regex = new RegExp(`(?:^|;\\s*)${name}=([^;]*)`);
    const match = regex.exec(document.cookie);
    return match ? match[1].replace(/\s+/g, '') : null;
}

async function getLikes(id_post) {
    let likecount = 0;
    let user_like = false
    let dislikecount = 0;
    let user_dislike = false
    try {
        const response = await fetch('http://localhost:8081/api/getLikes', {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json'
            }
        });
        if (!response.ok) {
            throw new Error('Network response was not ok');
        }
        const data = await response.json();
        if (data != null) {
            data.forEach(like => {
                if (like.topic_id === id_post) {
                    if (like.user_id == parseInt(getCookie("session_id"))) {
                        user_like = true;
                    }
                    likecount++;
                }
            });
        }
    } catch (error) {
        console.error('Error:', error);
    }
    try {
        const response = await fetch('http://localhost:8081/api/getUnLikes', {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json'
            }
        });
        if (!response.ok) {
            throw new Error('Network response was not ok');
        }
        const data = await response.json();
        if (data != null) {
            data.forEach(dislike => {
                if (dislike.topic_id === id_post) {
                    if (dislike.user_id == parseInt(getCookie("session_id"))) {
                        user_dislike = true;
                    }
                    dislikecount++;
                }
            });
        }
    } catch (error) {
        console.error('Error:', error);
    }
    return {likecount, dislikecount, user_like, user_dislike};
}

async function getUser(user_id) {
    try {
        const response = await fetch('http://localhost:8081/api/getuser?id=' + encodeURI(user_id), {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json'
            }
        });
        if (!response.ok) {
            throw new Error('Network response was not ok');
        }
        const data = await response.json();
        return data
    } catch (error) {
        console.error('Error:', error);
    }
}


document.addEventListener('DOMContentLoaded', () => {
    let userID = getCookie("session_id");
    userID = parseInt(userID);
    let count = 0

    fetch('http://localhost:8081/api/gettopics', {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json'
        }
    }).then(response => {
        if (!response.ok) {
            throw new Error('Network response was not ok');
        }
        return response.json();
    }).then(data => {
            const cardContainer = document.querySelector('.cardContainer');
            if (cardContainer) {
                data.forEach(async post => {
                    if (post.category == 4){
                        count++
                        let newCard;
                        if (userID == post.created_by) {
                            newCard = await createCardElement(post.id, post.body, post.created_by, post.title, post.CategoryName, post.CreatedAtFormatted );
                        } else {
                            newCard = await createCardElementOther(post.id, post.body, post.created_by, post.title, post.CategoryName, post.CreatedAtFormatted );
                        }
                        cardContainer.appendChild(newCard);
                    }
                });
            } else {
                console.error("Card container element not found");
            }
            if (count == 0){
                let errorWrite = document.createElement("p")
                errorWrite.innerHTML = "No one post Topics"
                errorWrite.className = "title-top"
                cardContainer.appendChild(errorWrite)
            }
    }).catch(error => console.error('Error:', error));

    setupMenuToggle('menu', 'sidebarMenu', '.sidebar', '.logoHeaderMenu');
    setupMenuToggle('profile', 'sidebarMenu-2', '.sidebar-2', '.logoHeader');
    setupSearchToggle('search-menu', 'searchHeader', 'inputSearch');
    setupResponsePopup();
    const cards = document.querySelectorAll('.card');
    if (cards.length) {
        cards.forEach(addCardFunctionality);
    } else {
        console.error("No cards found");
    }
});


function setupMenuToggle(menuId, sidebarId, sidebarClass, logoClass) {
    const menuLogo = document.getElementById(menuId);
    const sidebarMenu = document.getElementById(sidebarId);

    if (menuLogo && sidebarMenu) {
        menuLogo.addEventListener('click', () => {
            sidebarMenu.classList.toggle('hidden');
            sidebarMenu.classList.toggle('expanded');
        });

        document.addEventListener('click', (event) => {
            if (!event.target.closest(sidebarClass) && !event.target.matches(logoClass)) {
                sidebarMenu.classList.add('hidden');
                sidebarMenu.classList.remove('expanded');
            }
        });
    }
}

function setupSearchToggle(searchIconId, searchBarId, searchInputId) {
    const searchIcon = document.getElementById(searchIconId);
    const searchBar = document.getElementById(searchBarId);
    const searchInput = document.getElementById(searchInputId);

    if (searchIcon && searchBar && searchInput) {
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
    }
}

function setupResponsePopup() {
    const cardContainer = document.querySelector('.cardContainer');
    const responseButton = document.getElementById('response');
    const popup = document.getElementById('popup');
    const popupClose = document.getElementById('popup-close');
    const postButton = document.getElementById('post-button');
    const responseInput = document.getElementById('response-input');
    const responseInputTitle = document.getElementById('response-input-title');

    let isEditing = false;
    let currentPostID = null;
    let currentCard = null;

    if (responseButton && popup && popupClose && postButton && responseInput) {
        responseButton.addEventListener('click', () => {
            popup.style.display = 'block';
            postButton.innerText = 'Send';
            isEditing = false;
            currentPostID = null;
            currentCard = null;
        });

        popupClose.addEventListener('click', () => {
            popup.style.display = 'none';
        });

        window.addEventListener('click', (event) => {
            if (event.target === popup) {
                popup.style.display = 'none';
            }
        });

        postButton.addEventListener('click', async () => {
            const userInput = responseInput.value.trim();
            const userInputTitle = responseInputTitle.value.trim();
            const categoryChoose = document.getElementById("category-select").value
            let userID = parseInt(getCookie("session_id"), 10);
            if (userInput) {
                if (isEditing) {
                    updatePost(currentPostID, userInput, currentCard);
                } else {
                    const postData = {
                        title: userInputTitle,
                        body: userInput,
                        created_by: userID,
                        category: parseInt(categoryChoose)
                    };

                    try {
                        const response = await fetch('http://localhost:8081/api/topic', {
                            method: 'POST',
                            headers: {
                                'Content-Type': 'application/json'
                            },
                            body: JSON.stringify(postData)
                        });

                        if (!response.ok) {
                            throw new Error('Network response was not ok');
                        }

                        const data = await response.json();

                        if (data.status === 'success') {
                            const newCard = await createCardElement(data.post_id, userInput, userID, userInputTitle, data.category, data.time);
                            cardContainer.appendChild(newCard);
                            responseInput.value = '';
                            popup.style.display = 'none';
                        } else {
                            alert('Failed to create post');
                        }
                    } catch (error) {
                        console.error('Error:', error);
                    }
                }
            }
        });

        responseInput.addEventListener('keypress', (event) => {
            if (event.key === 'Enter') {
                event.preventDefault();
                postButton.click();
            }
        });
    }

    function addEditFunctionality(editBtn) {
        editBtn.addEventListener('click', () => {
            const postID = editBtn.dataset.postId;
            currentPostID = postID;
            currentCard = editBtn.closest('.card');
            const currentText = currentCard.querySelector('.textCardContainer').innerText;
            responseInput.value = currentText;
            postButton.innerText = 'Modifier';
            isEditing = true;
            popup.style.display = 'block';
        });
    }

    document.querySelectorAll('.edit-btn').forEach(addEditFunctionality);
}


async function createCardElement(postID, userInput, userID, title, category, time) {
    const { likecount: likes, user_like: user_like, dislikecount: dislikes, user_dislike: user_dislike } = await getLikes(postID);
    const user = await getUser(userID)
    const newCard = document.createElement('div');
    newCard.className = 'card';
    newCard.dataset.postId = postID;
    newCard.dataset.userId = userID;
    newCard.innerHTML = `
        <div class="content">
            <div class="cardHeader">
                <img class="cardProfileImg" src="http://localhost:8081/${user.profile_picture.String}" alt="profile">
                <div class="cardHeaderText">
                    <p class="NameUser">${user.username}</p>
                    <p class="DatePost">${time}</p>
                </div>
            </div>
            <div class="cardBody">
                <p class="cardTitle">${title}</p>
                <p class="textCardContainer">${userInput}</p>
                <span class="categoryTag">${category}</span>
            </div>
            <div class="cardFooter">
                <div>
                    <span class="like-counter">${likes}</span> 
                    <img class="likeAndUnLike like-btn" src="../static/img/like.png" alt="like">
                    <span class="unlike-counter">${dislikes}</span>  
                    <img class="likeAndUnLike unlike-btn" src="../static/img/unLike.png" alt="Unlike">
                    <img class="response" id="addComments" src="../static/img/discuter.png" alt="response">Reply
                </div>
                <div>
                    <span class="delete-btn">&#10060;</span>
                    <img class="response edit-btn" data-post-id="${postID}" src="../static/img/edit.png" alt="edit">Modify
                </div>
            </div>
        </div>
    `;
    addCardFunctionality(newCard, true, user_like, user_dislike);
    newCard.querySelector('.cardProfileImg').addEventListener('click', () => displayUser(user.id));
    return newCard;
}

async function createCardElementOther(postID, userInput, userID, title, category, time) {
    const { likecount: likes, user_like: user_like, dislikecount: dislikes, user_dislike: user_dislike } = await getLikes(postID);
    const user = await getUser(userID)
    const newCard = document.createElement('div');
    newCard.className = 'card';
    newCard.dataset.postId = postID;
    newCard.dataset.userId = userID;
    newCard.innerHTML = `
        <div class="content">
            <div class="cardHeader">
                <img class="cardProfileImg" src="http://localhost:8081/${user.profile_picture.String}" alt="profile">
                <div class="cardHeaderText">
                    <p class="NameUser">${user.username}</p>
                    <p class="DatePost">${time}</p>
                </div>
            </div>
            <div class="cardBody">
                <p class="cardTitle">${title}</p>
                <p class="textCardContainer">${userInput}</p>
                <span class="categoryTag">${category}</span>
            </div>            
            <div class="cardFooter">
                <div>
                    <span class="like-counter">${likes}</span> 
                    <img class="likeAndUnLike like-btn" src="../static/img/like.png" alt="like">
                    <span class="unlike-counter">${dislikes}</span>  
                    <img class="likeAndUnLike unlike-btn" src="../static/img/unLike.png" alt="Unlike">
                    <img class="response" id="addComments" src="../static/img/discuter.png" alt="response">Reply
                </div>
            </div>
        </div>
    `;
    addCardFunctionality(newCard, false, user_like, user_dislike);
    newCard.querySelector('.cardProfileImg').addEventListener('click', () => displayUser(user.id));
    return newCard;
}


function addCardFunctionality(card, state, user_like, user_dislike) {
    const likeBtn = card.querySelector('.like-btn');
    const unlikeBtn = card.querySelector('.unlike-btn');
    const likeCounter = card.querySelector('.like-counter');
    const unlikeCounter = card.querySelector('.unlike-counter');
    const comments = card.querySelector("#addComments")

    let likeCount = parseInt(likeCounter.textContent, 10);
    let unlikeCount = parseInt(unlikeCounter.textContent, 10);
    let likeClicked = user_like;
    let unlikeClicked = user_dislike;
    
    const userId = getCookie("session_id");
    const postId = card.dataset.postId;
    
    comments.addEventListener('click', function () {
        window.location.href = `/Comments?id=${postId}`
    })  

    function updatePostText() {
        updatePost(postId, card.querySelector('.textCardContainer').textContent, card);
    }

    async function likePost(userId, postId) {

        const postData = {
            user_id: parseInt(userId),
            post_id: 0,
            topic_id: parseInt(postId),
        };

        let state = false
        try {
            const response = await fetch('http://localhost:8081/api/getLikes', {
                method: 'GET',
                headers: {
                    'Content-Type': 'application/json'
                }
            });
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
            const data = await response.json();
            if (data != null) {
                data.forEach(like => {
                    if (like.user_id === getCookie("session_id")){
                        state = true
                    }
                });
            }
        } catch (error) {
            console.error('Error:', error);
        }

        if (state == false){
            const response = await fetch('http://localhost:8081/api/like',{
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(postData)
            });
            const data = await response.json();
            if (data.status !== 'success') {
                alert('Failed to like post');
            }
        }
    }


    async function unlikePost() {

        const postData = {
            user_id: parseInt(userId),
            post_id: 0,
            topic_id: parseInt(postId)
        };

        const response = await fetch('http://localhost:8081/api/unlike', {
            method: 'DELETE',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(postData)
        });
        const data = await response.json();
        if (data.status !== 'success') {
            alert('Failed to unlike post');
        }
    }

    async function dislikePost() {

        const postData = {
            user_id: parseInt(userId),
            post_id: 0,
            topic_id: parseInt(postId)
        };


        const response = await fetch('http://localhost:8081/api/dislike', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(postData)
        });
        const data = await response.json();
        if (data.status !== 'success') {
            alert('Failed to dislike post');
        }
    }

    async function undislikePost() {

        const postData = {
            user_id: parseInt(userId),
            post_id: 0,
            topic_id: parseInt(postId)
        };


        const response = await fetch('http://localhost:8081/api/undislike', {
            method: 'DELETE',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(postData)
        });
        const data = await response.json();
        if (data.status !== 'success') {
            alert('Failed to undislike post');
        }
    }

    if (likeBtn) {
        likeBtn.addEventListener('click', async () => {
            if (!likeClicked) {
                likeCount++;
                likeClicked = true;
                if (unlikeClicked) {
                    unlikeCount--;
                    unlikeClicked = false;
                    await undislikePost();
                }
                await likePost(userId, postId);
            } else {
                likeCount--;
                likeClicked = false;
                await unlikePost();
            }
            updatePostText();
            likeCounter.innerText = likeCount;
            unlikeCounter.innerText = unlikeCount;
        });
    }

    if (unlikeBtn) {
        unlikeBtn.addEventListener('click', async () => {
            if (!unlikeClicked) {
                unlikeCount++;
                unlikeClicked = true;
                if (likeClicked) {
                    likeCount--;
                    likeClicked = false;
                    await unlikePost();
                }
                await dislikePost();
            } else {
                unlikeCount--;
                unlikeClicked = false;
                await undislikePost();
            }
            updatePostText();
            likeCounter.innerText = likeCount;
            unlikeCounter.innerText = unlikeCount;
        });
    }

    if (state){
        const deleteBtn = card.querySelector('.delete-btn');
        const editBtn = card.querySelector('.edit-btn');     

        if (deleteBtn) {
            deleteBtn.addEventListener('click', () => {
                const postID = card.dataset.postId;
                fetch(`http://localhost:8081/api/delete_topic?id=${postID}&user_id=${card.dataset.userId}`, {
                    method: 'DELETE'
                })
                    .then(response => response.json())
                    .then(data => {
                        if (data.status === 'success') {
                            card.remove();
                        } else {
                            alert('Failed to delete post');
                        }
                    })
                    .catch(error => console.error('Error:', error));
            });
        }

        if (editBtn) {
            editBtn.addEventListener('click', () => {
                const postID = card.dataset.postId;
                const currentText = card.querySelector('.textCardContainer').innerText;
                const newText = prompt('Modifier le contenu du post:', currentText);
                if (newText !== null) {
                    updatePost(postID, newText, card);
                }
            });
        }
    }
}

function updatePost(postID, newText, card) {
    const postData = {
        id: parseInt(postID, 10),
        body: newText
    };

    fetch('http://localhost:8081/api/update_topic', {
        method: 'PUT',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(postData)
    })
        .then(response => response.json())
        .then(data => {
            if (data.status === 'success') {
                card.querySelector('.textCardContainer').innerText = newText;
            } else {
                alert('Failed to update post');
            }
        })
        .catch(error => console.error('Error:', error));
}

document.addEventListener("DOMContentLoaded", function() {
    const posts = document.querySelectorAll('.card');
    const observerOptions = {
        threshold: 0.1
    };

    const observer = new IntersectionObserver((entries, observer) => {
        entries.forEach(entry => {
            if (entry.isIntersecting) {
                entry.target.classList.add('fadeInUp');
                observer.unobserve(entry.target);
            }
        });
    }, observerOptions);

    posts.forEach(post => {
        observer.observe(post);
    });
});

const buttons = document.querySelectorAll('.likeAndUnLike, .response');
buttons.forEach(button => {
    button.addEventListener('mouseover', () => {
        button.style.transform = 'scale(1.2)';
    });
    button.addEventListener('mouseout', () => {
        button.style.transform = 'scale(1)';
    });
});

const addPostButton = document.querySelector('.addPost');
addPostButton.addEventListener('click', () => {
    const popup = document.getElementById('popup');
    popup.style.display = 'block';
});

const popupCloseButton = document.getElementById('popup-close');
popupCloseButton.addEventListener('click', () => {
    const popup = document.getElementById('popup');
    popup.style.display = 'none';
});

const postButton = document.getElementById('post-button');
postButton.addEventListener('mouseover', () => {
    postButton.style.transform = 'scale(1.05)';
});
postButton.addEventListener('mouseout', () => {
    postButton.style.transform = 'scale(1)';
});

