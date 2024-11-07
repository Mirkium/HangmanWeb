const navBar = document.querySelector(".header");

let prevSroll = window.scrollY;

window.addEventListener("scroll", function() {
    let currScroll = window.scrollY;

    if (currScroll > prevSroll) {
        navBar.style.transform = `translateY(-105%)`;
    } else {
        navBar.style.transform = `translateY(0%)`;
    }
})