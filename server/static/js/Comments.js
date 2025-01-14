function getCookie(name) {
    const regex = new RegExp(`(?:^|;\\s*)${name}=([^;]*)`);
    const match = regex.exec(document.cookie);
    return match ? match[1].replace(/\s+/g, '') : null;
}

function displayUser(id) {
    const encodedUser = encodeURIComponent(id);
    window.location.href = `/User?user=${encodedUser}`;
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
                if (like.post_id === parseInt(id_post)) {
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
                if (dislike.post_id === parseInt(id_post)) {
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

async function setup() {
    try {
        const response = await fetch(`http://localhost:8081/api/getposts`, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json'
            }
        });

        if (!response.ok) {
            throw new Error('Network response was not ok');
        }

        const data = await response.json();

        for (const comment of data) {
            if (comment.topic_id == topicId) {
                const userid = comment.created_by;
                const userResponse = await fetch(`http://localhost:8081/api/getuser?id=${userid}`, {
                    method: 'GET',
                    headers: {
                        'Content-Type': 'application/json'
                    }
                });

                if (!userResponse.ok) {
                    throw new Error('Network response was not ok');
                }

                const userdata = await userResponse.json();

                const commentsSection = document.getElementById('commentsSection');
                const newComment = document.createElement('div');
                newComment.classList.add('comment');
                newComment.id = comment.id;
                const { likecount: likes, dislikecount: dislikes } = await getLikes(newComment.id);

                const parsedCommentContent = parseCommentContent(comment.body);

                newComment.innerHTML = `
                    <div class="comment-header">
                        <img src="http://localhost:8081/${userdata.profile_picture.String}" alt="User Avatar" class="avatar">
                        <div class="user-info">
                            <span class="username">${userdata.username}</span>
                            <span class="time">${comment.CreatedAtFormatted}</span>
                        </div>
                    </div>
                    <div class="comment-content">${parsedCommentContent}</div>
                    <div class="comment-footer">
                        <div>
                            <span class="like-counter">${likes}</span>
                            <img class="likeAndUnLike-like-btn" src="../static/img/like.png" alt="like" onclick="handleLikeDislike(${newComment.id}, 'like')">
                            <span class="unlike-counter">${dislikes}</span>
                            <img class="likeAndUnLike-unlike-btn" src="../static/img/unLike.png" alt="Unlike" onclick="handleLikeDislike(${newComment.id}, 'dislike')">
                        </div>
                        ${userdata.id == getCookie("session_id") ? `
                        <button class="modif-comment-btn" onclick="editComment(this)">✏️</button>
                        <button class="delete-comment-btn" onclick="deleteComment(this)">🗑️</button>
                        ` : ''}
                    </div>
                `;

                commentsSection.prepend(newComment);
                newComment.querySelector('.avatar').addEventListener('click', () => displayUser(userdata.id));
                updateLikeDislikeCounts(newComment.id);
            }
        }
    } catch (error) {
        console.error('Error:', error);
    }
}

const queryString = window.location.search;
const urlParams = new URLSearchParams(queryString);
const topicId = urlParams.get('id');
setup()

document.addEventListener('DOMContentLoaded', async function () {
    const title = document.querySelector(".title-top-topic");
    if (title) {
        const response = await fetch('http://localhost:8081/api/gettopics', {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json'
            },
        });
        const data = await response.json();
        const topicId = new URLSearchParams(window.location.search).get('id');
        data.forEach(element => {
            if (element.id == topicId) {
                title.innerHTML = element.body;
            }
        });
    }

    const menuLogo = document.getElementById("menu");
    const sidebarMenu = document.getElementById("sidebarMenu");
    if (menuLogo && sidebarMenu) {
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
    }
});

document.addEventListener('DOMContentLoaded', function () {
    const menuLogo = document.getElementById("profile");
    const sidebarMenu = document.getElementById("sidebarMenu-2");
    if (menuLogo && sidebarMenu) {
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
    }
});

document.addEventListener('DOMContentLoaded', () => {
    const searchIcon = document.getElementById('search-menu');
    const searchBar = document.getElementById('searchHeader');
    const searchInput = document.getElementById('inputSearch');

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


document.addEventListener('DOMContentLoaded', () => {
    const searchIcon = document.getElementById('search-menu');
    const searchBar = document.getElementById('searchHeader');
    const searchInput = document.getElementById('inputSearch');

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
});


let commentId = 1;

async function searchGif(event) {
    const query = event.target.value;
    if (query.length < 3) return;

    const apiKey = 'JU77cxigeGaROzWJq9goayq09PxIwb62';
    const response = await fetch(`https://api.giphy.com/v1/gifs/search?api_key=${apiKey}&q=${query}&limit=10`);
    const data = await response.json();
    const gifResults = document.getElementById('gifResults');

    gifResults.innerHTML = '';
    data.data.forEach(gif => {
        const img = document.createElement('img');
        img.src = gif.images.fixed_height.url;
        img.alt = gif.title;
        img.onclick = () => addGifToComment(gif.images.fixed_height.url);
        gifResults.appendChild(img);
    });

    gifResults.style.display = 'block';
}

function toggleGifSearch() {
    const gifSearch = document.getElementById('gifSearch');
    gifSearch.style.display = gifSearch.style.display === 'none' ? 'block' : 'none';
    gifSearch.focus();
}

function addGifToComment(gifUrl) {
    const commentInput = document.getElementById('commentInput');
    commentInput.value += ` ![GIF](${gifUrl}) `;
    document.getElementById('gifResults').style.display = 'none';
    document.getElementById('gifSearch').value = '';
    updateCharacterCount();
}

function toggleImageUpload() {
    const fileInput = document.getElementById('imageUpload');
    fileInput.click();
}

function handleImageUpload() {
    const fileInput = document.getElementById('imageUpload');
    const file = fileInput.files[0];
    if (file) {
        const reader = new FileReader();
        reader.onload = function (e) {
            const commentInput = document.getElementById('commentInput');
            commentInput.value += ` ![Image](${e.target.result}) `;
            updateCharacterCount();
        };
        reader.readAsDataURL(file);
    }
}

function clearComment() {
    document.getElementById('commentInput').value = '';
    updateCharacterCount();
}

async function postComment() {
    const user_id = getCookie("session_id")
    const commentInput = document.getElementById('commentInput').value;
    const postData = {
        body: commentInput,
        created_by: parseInt(user_id),
        topic_id: parseInt(topicId),
    }
    const post = await fetch('http://localhost:8081/api/post', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(postData)
    });
    if (!post.ok) {
        throw new Error('Network response was not ok');
    }
    const comment = await post.json();
    if (commentInput.trim() !== '') {
        const response = await fetch(`http://localhost:8081/api/getuser?id=${user_id}`, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json'
            }
        });
        if (!response.ok) {
            throw new Error('Network response was not ok');
        }
        const data = await response.json();
        const commentsSection = document.getElementById('commentsSection');
        const newComment = document.createElement('div');
        newComment.classList.add('comment');
        newComment.id = comment.post_id;
        const { likecount: likes,  dislikecount: dislikes } = await getLikes(newComment.id);

        newComment.innerHTML = `
            <div class="comment-header">
                <img src="http://localhost:8081/${data.profile_picture.String}" alt="User Avatar" class="avatar">
                <div class="user-info">
                    <span class="username">${data.username}</span>
                    <span class="time">${new Date().toLocaleString()}</span>
                </div>
            </div>
            <div class="comment-content">${parseCommentContent(commentInput)}</div>
            <div class="comment-footer">
                <div>
                    <span class="like-counter">${likes}</span> 
                    <img class="likeAndUnLike-like-btn" src="../static/img/like.png" alt="like" onclick="handleLikeDislike(${newComment.id}, 'like')">
                    <span class="unlike-counter">${dislikes}</span>  
                    <img class="likeAndUnLike-unlike-btn" src="../static/img/unLike.png" alt="Unlike" onclick="handleLikeDislike(${newComment.id}, 'dislike')">
                </div>
                <button class="modif-comment-btn" onclick="editComment(this)">✏️</button>
                <button class="delete-comment-btn" onclick="deleteComment(this)">🗑️</button>
            </div>
        `;
        commentsSection.prepend(newComment);
        updateLikeDislikeCounts(newComment.id);
        newComment.querySelector('.avatar').addEventListener('click', () => displayUser(data.id));
        clearComment()
    }
}

function parseCommentContent(content) {
    const imageRegex = /!\[Image\]\((.*?)\)/g;
    const gifRegex = /!\[GIF\]\((.*?)\)/g;

    let parsedContent = content.replace(imageRegex, '<img src="$1" style="max-width: 15%;" />');
    parsedContent = parsedContent.replace(gifRegex, '<img src="$1" style="max-width: 15%;" />');

    return parsedContent;
}

function deleteComment(button) {
    const comment = button.closest('.comment');
    console.log(comment)
    fetch(`http://localhost:8081/api/delete_post?id=${comment.id}&user_id=${getCookie("session_id")}`, {
        method: 'DELETE'
    })
    .then(response => response.json())
    .then(data => {
        if (comment) {
            comment.remove();
        }
    })
    .catch(error => console.error('Error:', error));
}

function editComment(button) {
    const comment = button.closest('.comment');
    console.log(comment)
    const postID = comment.id;
    const newText = prompt('Modifier le contenu du post:');
    if (newText !== null) {
        updateComment(postID, newText, comment);
    }
}

function updateCharacterCount() {
    const commentInput = document.getElementById('commentInput').value;
    const characterCount = commentInput.length;
    document.getElementById('characterCount').innerText = `${characterCount} /10000`;
}

function updateComment(postID, newText, comment){
    const postData = {
        id: parseInt(postID, 10),
        body: newText
    };

    fetch('http://localhost:8081/api/update_post', {
        method: 'PUT',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(postData)
    })
        .then(response => response.json())
        .then(data => {
            if (data.status === 'success') {
                comment.querySelector('.comment-content').innerText = newText;
            }
        })
        .catch(error => console.error('Error:', error));
}

async function likeComment(commentId) {
    try {
        const user_id = getCookie("session_id");
        const response = await fetch('http://localhost:8081/api/like', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ post_id: commentId, user_id: parseInt(user_id) })
        });
        if (!response.ok) {
            throw new Error('Network response was not ok');
        }
        updateLikeDislikeCounts(commentId);
    } catch (error) {
        console.error('Error:', error);
    }
}

async function unlikeComment(commentId) {
    try {
        const user_id = getCookie("session_id");
        const response = await fetch('http://localhost:8081/api/unlike', {
            method: 'DELETE',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ post_id: commentId, user_id: parseInt(user_id), topic_id: 0 })
        });
        if (!response.ok) {
            throw new Error('Network response was not ok');
        }
        updateLikeDislikeCounts(commentId);
    } catch (error) {
        console.error('Error:', error);
    }
}

async function dislikeComment(commentId) {
    try {
        const user_id = getCookie("session_id");
        const response = await fetch('http://localhost:8081/api/dislike', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ post_id: commentId, user_id: parseInt(user_id) })
        });
        if (!response.ok) {
            throw new Error('Network response was not ok');
        }
        updateLikeDislikeCounts(commentId);
    } catch (error) {
        console.error('Error:', error);
    }
}

async function undislikeComment(commentId) {
    try {
        const user_id = getCookie("session_id");
        const response = await fetch('http://localhost:8081/api/undislike', {
            method: 'DELETE',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ post_id: commentId, user_id: parseInt(user_id), topic_id: 0 })
        });
        if (!response.ok) {
            throw new Error('Network response was not ok');
        }
        updateLikeDislikeCounts(commentId);
    } catch (error) {
        console.error('Error:', error);
    }
}

async function handleLikeDislike(commentId, action) {
    const { user_like, user_dislike } = await getLikes(commentId);
    if (action === 'like') {
        if (user_like) {
            await unlikeComment(commentId);
        } else {
            if (user_dislike) await undislikeComment(commentId);
            await likeComment(commentId);
        }
    } else if (action === 'dislike') {
        if (user_dislike) {
            await undislikeComment(commentId);
        } else {
            if (user_like) await unlikeComment(commentId);
            await dislikeComment(commentId);
        }
    }
    updateLikeDislikeCounts(commentId);
}

async function likeComment(commentId) {
    const user_id = getCookie("session_id");
    const response = await fetch('http://localhost:8081/api/like', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ post_id: commentId, user_id: parseInt(user_id) })
    });
    if (!response.ok) {
        throw new Error('Network response was not ok');
    }
}

async function unlikeComment(commentId) {
    const user_id = getCookie("session_id");
    const response = await fetch('http://localhost:8081/api/unlike', {
        method: 'DELETE',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ post_id: commentId, user_id: parseInt(user_id) })
    });
    if (!response.ok) {
        throw new Error('Network response was not ok');
    }
}

async function dislikeComment(commentId) {
    const user_id = getCookie("session_id");
    const response = await fetch('http://localhost:8081/api/dislike', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ post_id: commentId, user_id: parseInt(user_id) })
    });
    if (!response.ok) {
        throw new Error('Network response was not ok');
    }
}

async function undislikeComment(commentId) {
    const user_id = getCookie("session_id");
    const response = await fetch('http://localhost:8081/api/undislike', {
        method: 'DELETE',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ post_id: commentId, user_id: parseInt(user_id) })
    });
    if (!response.ok) {
        throw new Error('Network response was not ok');
    }
}

async function updateLikeDislikeCounts(commentId) {
    const { likecount, dislikecount, user_like, user_dislike } = await getLikes(commentId);
    const commentElement = document.getElementById(commentId);
    if (commentElement) {
        const likeCounter = commentElement.querySelector('.like-counter');
        const dislikeCounter = commentElement.querySelector('.unlike-counter');
        const likeBtn = commentElement.querySelector('.likeAndUnLike-like-btn');
        const dislikeBtn = commentElement.querySelector('.likeAndUnLike-unlike-btn');

        likeCounter.innerText = likecount;
        dislikeCounter.innerText = dislikecount;

        likeBtn.classList.toggle('liked', user_like);
        dislikeBtn.classList.toggle('disliked', user_dislike);
    }
}

document.addEventListener("DOMContentLoaded", function() {
    // Animation pour les boutons
    const buttons = document.querySelectorAll('.comment-actions button');
    buttons.forEach(button => {
        button.addEventListener('mouseenter', () => {
            button.style.transform = 'scale(1.1)';
            button.style.transition = 'transform 0.2s';
        });
        button.addEventListener('mouseleave', () => {
            button.style.transform = 'scale(1)';
            button.style.transition = 'transform 0.2s';
        });
    });

    // Animation pour les commentaires
    const comments = document.querySelectorAll('.comment');
    comments.forEach((comment, index) => {
        comment.style.opacity = '0';
        comment.style.transition = `opacity 0.5s ease ${(index + 1) * 0.2}s`;
    });

    setTimeout(() => {
        comments.forEach(comment => {
            comment.style.opacity = '1';
        });
    }, 100);

    // Animation pour la zone de texte
    const textarea = document.getElementById('commentInput');
    textarea.addEventListener('focus', () => {
        textarea.style.boxShadow = '0 0 10px #ff6347';
        textarea.style.transition = 'box-shadow 0.3s';
    });
    textarea.addEventListener('blur', () => {
        textarea.style.boxShadow = 'none';
        textarea.style.transition = 'box-shadow 0.3s';
    });
});