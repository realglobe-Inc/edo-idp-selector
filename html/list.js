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

function query_parse(raw) {
    var queries = {};
    var q = raw.split("&");
    for (var i = 0; i < q.length; i++) {
        var elem = q[i].split("=");

        var key = elem[0];
        var val = elem[1];
        if (val) {
            val = decodeURIComponent(val.replace(/\+/g, " ")).replace(/\n/g, "<br/>");
        }

        queries[key] = val;
    }

    return queries;
}

function list() {
    var ticket = location.hash.substring(1);
    var queries = query_parse(window.location.search.substring(1));

    if (ticket) {
        document.write('<b>ticket:</b> ' + ticket + '<br/>');
    }
    if (queries && Object.keys(queries).length > 0) {
        for (key in queries) {
            document.write('<b>' + key + ':</b> ' + queries[key] + '<br/>');
        }
        document.write('<br/>');
    }

    var apiUri = "/api/info/issuer";
    var selUri = "/select";
    var ticket = location.hash.substring(1)

    var xhr = new XMLHttpRequest();
    xhr.open("GET", apiUri, false);
    xhr.send();
    if (xhr.status === 200) {
        var idps = JSON.parse(xhr.responseText);
        for (var i = 0; i < idps.length; i++ ) {
            document.write('<b><a href="' + selUri + "?ticket=" + encodeURIComponent(ticket) + "&issuer=" + encodeURIComponent(idps[i].issuer) + '">' + idps[i].issuer + "</a></b><br/>");
            for (key in idps[i]) {
                var result = key.match(/^issuer_name#?(.*)?$/)
                if (result) {
                    var prefix = '&nbsp;&nbsp;&nbsp;&nbsp;'
                    if (result[1]) {
                        prefix += "(" + result[1] + ") "
                    }
                    document.write(prefix + '<a href="' + selUri + "?ticket=" + encodeURIComponent(ticket) + "&issuer=" + encodeURIComponent(idps[i].issuer) + '">' + idps[i][key] + "</a><br/>");
                }
            }
            document.write('<br/>');
        }
    } else {
        document.write(xhr.statusText);
    }
}
