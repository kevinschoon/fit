$(function() {
    function update(uri, key, value) {
        /*Thanks http://stackoverflow.com/questions/5999118/add-or-update-query-string-parameter#6021027*/
        if (typeof value != 'undefined') {
            var re = new RegExp("([?&])" + key + "=.*?(&|$)", "i");
            var separator = uri.indexOf('?') !== -1 ? "&" : "?";
            if (uri.match(re)) {
                return uri.replace(re, '$1' + key + "=" + value + '$2');
            } else {
                return uri + separator + key + "=" + value;
            }
        }
        return uri
    }
    $("#submit").on('click', function() {
        var target = window.location.pathname
        var Q = window.location.search
        Q = update(Q, "X", $("#x").val())
        Q = update(Q, "Y", $("#y").val())
        Q = update(Q, "aggr", $("#aggr").val())
        Q = update(Q, "fn", $("#fn").val())
        Q = update(Q, "start", $("#start").val())
        Q = update(Q, "end", $("#end").val())
        window.location = window.location.pathname + Q
    });
})
