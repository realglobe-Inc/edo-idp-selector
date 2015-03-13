// Copyright 2015 realglobe, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

function list() {
    var apiUri = "/list";
    var rediUri = "/redirect";

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
