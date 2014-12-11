function list() {
    var apiUri = "/list";
    var rediUri = "/redirect"

    var prefix = rediUri;
    if (window.location.search.length > 0) {
        prefix += window.location.search + "&";
    } else {
        prefix += "?";
    }

    var xhr = new XMLHttpRequest();
    xhr.open("GET", apiUri, false);
    xhr.send();
    if (xhr.status === 200) {
        var idps = JSON.parse(xhr.responseText);
        idps.sort(
            function(a, b) {
                if (a.name < b.name) {
                    return -1;
                } else if (a.name > b.name) {
                    return 1;
                } else {
                    return 0;
                }
            }
        );
        for (var i = 0; i < idps.length; i++ ) {
            document.write('<a href="' + prefix + "idp=" + encodeURIComponent(idps[i].id) + '">' + idps[i].name + "</a><br/>");
        }
    } else {
        document.write(xhr.statusText);
    }
}
