"use strict";

function getArtistData(callback, pagination, page) {

    var request = new XMLHttpRequest();

    var parameters = "";
    if (pagination !== undefined || page !== undefined) {
        parameters += "?"
        parameters += typeof(pagination) !== "string" ? "" : "n=" + pagination + "&";
        parameters += typeof(page) !== "string" ? "" : "page=" + page;
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
