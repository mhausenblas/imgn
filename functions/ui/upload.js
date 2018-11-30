(function (d, axios) {
    "use strict";
    var inputFile = d.querySelector("#inputFile");
    var divNotification = d.querySelector("#alert");
    var fn = "";

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
        fn = file
        axios.post("HTTP_API/upload", formData)
            .then(onPreResponse)
            .catch(onPreResponse);
    }

    function onPreResponse(response) {
        var presurl = response.data;
        upload(presurl)
    }

    function upload(presurl) {
        axios.put(presurl, fn)
            .then(onResponse)
            .catch(onResponse);
    }

    function onResponse(response) {
        var className = (response.status !== 400) ? "success" : "error";
        divNotification.innerHTML = "Successfully uploaded " + fn.name + " to gallery!"
        divNotification.classList.add(className);
        setTimeout(function () {
            divNotification.classList.remove(className);
            location.href = "/";
        }, 3000);
    }
})(document, axios)