// Burger menus
document.addEventListener('DOMContentLoaded', function () {
    const hamburger = document.querySelector("#hamburger");
    const menu = document.querySelector('#menu');

    if (screen.width < 1024){
        menu.classList.remove('flex')
        menu.classList.add('hidden')
    }

    hamburger.addEventListener('click', () => {
        // for (var i=1; i<hamburger.length;i++){
        //     console.log(hamburger[i])
        //     hamburger[i].classList.toggle('hidden')
        //     hamburger[i].classList.toggle('block')
        // }
        
        menu.classList.toggle('flex')
        menu.classList.toggle('hidden')
    });
});

window.addEventListener('resize', function (event) {
    const menu = document.querySelector('#menu');
    const width = event.target.screen.width
    
    if (width >= 1024){
        menu.classList.remove('hidden')
        menu.classList.add('flex')
    }
});