import htmlPath from './index.html';

fetch(htmlPath)
    .then(response => response.text())
    .then(html => {
        document.querySelector("body").innerHTML = html;
    });