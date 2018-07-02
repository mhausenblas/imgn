"use strict";

function post(url, data) {
    return axios.post(url, data)
        .then(function (response) {
            return response;
        }).catch(function (error) {
            return error.response;
        });
}