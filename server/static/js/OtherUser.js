// OtherUser.js

document.addEventListener('DOMContentLoaded', function () {
    fetch('http://localhost:8081/api/getusers', {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json'
        }
    })
    .then(response => {
        if (!response.ok) {
            throw new Error('Network response was not ok');
        }
        return response.json();
    })
    .then(users => {
        const id = new URLSearchParams(window.location.search).get('user') || '';
        const user = users.find(user => user.id === parseInt(id));
        if (user) {
            setUserElements(user);
        } else {
            console.error('User not found');
        }
    })
    .catch(error => {
        console.error('Error:', error);
    });
});

async function setUserElements(user) {
    const profilePicture = document.getElementById("profileImg");
    const username = document.getElementById("username");
    const creationDate = document.getElementById("creationDate");
    const biography = document.getElementById("biography");
    const postCount = document.getElementById("postCount");
    const commentCount = document.getElementById("commentCount");

    if (!user.profile_picture.String) {
        profilePicture.src = "/static/img/agent.png";
    } else {
        profilePicture.src = "http://localhost:8081/" + user.profile_picture.String;
    }

    username.innerHTML = user.username;

    creationDate.innerHTML = user.created_at ? "Account created on: " + user.created_at : "";
    biography.innerHTML = user.biography.String || "No biography.";

    try {
        let posts = await getUserPosts(user);
        postCount.innerHTML = posts.length + " publications";
        displayPosts(posts);
    } catch (error) {
        console.error('Error fetching user posts:', error);
    }

    try {
        let comments = await getUserComments(user);
        commentCount.innerHTML = comments.length + "comments";
    } catch (error) {
        console.error('Error fetching user comments:', error);
    }

    document.getElementById('postsButton').addEventListener('click', async () => {
        try {
            let posts = await getUserPosts(user);
            displayPosts(posts);
        } catch (error) {
            console.error('Error displaying posts:', error);
        }
    });

    document.getElementById('commentsButton').addEventListener('click', async () => {
        try {
            let comments = await getUserComments(user);
            displayComments(comments);
        } catch (error) {
            console.error('Error displaying comments:', error);
        }
    });
}

async function getUserPosts(user) {
    const response = await fetch('http://localhost:8081/api/gettopics', {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json'
        }
    });

    if (!response.ok) {
        throw new Error('Network response was not ok');
    }

    const data = await response.json();
    return filterPosts(data, user);
}

async function getUserComments(user) {
    const response = await fetch('http://localhost:8081/api/getposts', {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json'
        }
    });

    if (!response.ok) {
        throw new Error('Network response was not ok');
    }

    const data = await response.json();
    return filterPosts(data, user);
}

function filterPosts(posts, user) {
    return posts.filter(post => post.created_by === user.id);
}       

function displayPosts(posts) {
    const postsContainer = document.getElementById('postsContainer');
    postsContainer.innerHTML = '';

    posts.forEach(post => {
        const postElement = document.createElement('div');
        postElement.classList.add('post');
        postElement.innerHTML = `
            <h3>${post.title}</h3>
            <p>${post.body}</p>
            <p>Posted on: ${new Date(post.created_at).toLocaleString()}</p>
        `;
        postsContainer.appendChild(postElement);
    });
}

function displayComments(comments) {
    const postsContainer = document.getElementById('postsContainer');
    postsContainer.innerHTML = '';

    comments.forEach(comment => {
        const commentElement = document.createElement('div');
        commentElement.classList.add('comment');
        commentElement.innerHTML = `
            <p>${comment.body}</p>
            <p>Commented on: ${new Date(comment.created_at).toLocaleString()}</p>
        `;
        postsContainer.appendChild(commentElement);
    });
}
