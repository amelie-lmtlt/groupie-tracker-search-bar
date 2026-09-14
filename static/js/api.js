"use strict";

function getArtistData(callback, pagination, page, search) {

    var request = new XMLHttpRequest();

    var parameters = "";
    var queryParameters = [];

    if (typeof(pagination) === "string") {
        queryParameters.push("n=" + pagination);
    }

    if (typeof(page) === "string") {
        queryParameters.push("page=" + page);
    }

    if (typeof(search) === "string" && search !== "") {
        queryParameters.push("search=" + encodeURIComponent(search));
    }

    if (queryParameters.length > 0) {
        parameters = "?" + queryParameters.join("&");
    }

    request.open("GET", "/artist/" + parameters, true);

    request.onreadystatechange = function () {

        if (request.readyState === 4) {

            if (request.status === 200) {
                callback(request.status, request.responseText);
            } else {
                console.log("Error :", request.status);
            }
        }
    };

    request.send();
}

function getRelationData(callback, id) {
    if (id === undefined) {
        return "error: id was undefined";
    }
    if (type(id) !== Number) {
        return "error: id was not a number";
    }

    var request = new XMLHttpRequest();

    request.open("GET", "/relation/" + id, true);

    request.onreadystatechange = function () {
        if (request.readyState === 4) {
            if (request.status === 200) {
                callback(request.status, request.responseText);
            } else {
                console.log("Error :", request.status);
            }
        }
    };
    request.send();
}
