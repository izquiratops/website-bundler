import htmlPath from './index.html';
import './style.css';

fetch(htmlPath)
    .then(response => response.text())
    .then(html => {
        document.querySelector("body").innerHTML = html;
    });