(function (d, axios) {
    "use strict";
    var inputFile = d.querySelector("#inputFile");
    var divNotification = d.querySelector("#alert");

    inputFile.addEventListener("change", addFile);    
    function addFile(e) {
        var file = e.target.files[0]
        if (!file) {
            return
        }
        upload(file);
    }

    function upload(file) {
        var formData = new FormData()
        formData.append("file", file)
        post("/upload", formData)
            .then(onResponse)
            .catch(onResponse);
    }

    function onResponse(response) {
        var className = (response.status !== 400) ? "success" : "error";
        divNotification.innerHTML = response.data;
        divNotification.classList.add(className);
        setTimeout(function () {
            divNotification.classList.remove(className);
            location.href = "/";
        }, 3000);
    }
})(document, axios)