export function setCookie(name, value, days) {
    let expires = "";
    if (days) {
        let date = new Date();
        date.setTime(date.getTime() + (days * 24 * 60 * 60 * 1000));
        expires = "; expires=" + date.toUTCString();
    }
    document.cookie = name + "=" + (value || "") + expires + "; path=/";
}

export function getCookie(name) {
    const regex = new RegExp(`(?:^|;\\s*)${name}=([^;]*)`);
    const match = regex.exec(document.cookie);
    return match ? match[1].replace(/\s+/g, '') : null;
}


