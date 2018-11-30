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
        preupload(file);
    }

    function preupload(file) {
        var formData = new FormData()
        formData.append("file", file)
        post("HTTP_API/upload", formData)
            .then(onPreResponse)
            .catch(onPreResponse);
    }

    function onPreResponse(response) {
        var file = e.target.files[0]
        if (!file) {
            return
        }
        var presurl = response.data;
        upload(file, presurl)
    }

    function upload(file, presurl) {
        var formData = new FormData()
        formData.append("file", file)
        post(presurl, formData)
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