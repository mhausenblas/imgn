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
                var img = response.data[i].src;
                var meta = response.data[i].meta;
                buffer += "<a href='" + img + "' title='"+ meta +"'><img class='gimg' src='" + img + "'/></a>";
            }
            divGallery.innerHTML = buffer;
        }
        
    }
})(document, axios)