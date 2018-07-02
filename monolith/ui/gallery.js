(function (d, axios) {
    "use strict";
    var divGallery = d.querySelector("#gallery");

    fillGallery();

    function fillGallery(){
        get("/gallery")
            .then(onRenderGallery)
            .catch(onRenderGallery);
    }

    function onRenderGallery(response) {
        if (response.status == 200){
            var buffer = "";
            for (let i = 0; i < response.data.length; i++) {
                var img = response.data[i];
                buffer += "<a href='" + img + "'><img class='gimg' src='" + img + "'/></a>";
            }
            divGallery.innerHTML = buffer;
        }
        
    }
})(document, axios)