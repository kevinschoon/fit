$(function() {
    $('#daterange .input-daterange').datepicker({});
    console.log(window.location.pathname);
    console.log(window.location.href);
    var urlParams;
    (window.onpopstate = function() {
        var match,
            pl = /\+/g, // Regex for replacing addition symbol with a space
            search = /([^&=]+)=?([^&]*)/g,
            decode = function(s) {
                return decodeURIComponent(s.replace(pl, " "));
            },
            query = window.location.search.substring(1);

        urlParams = {};
        while (match = search.exec(query))
            urlParams[decode(match[1])] = decode(match[2]);
    })();
    console.log(urlParams)
    for (var i in urlParams) {
        console.log(urlParams[i])
    }
    $("#submit").on('click', function() {
        console.log($("#aggregation").val());
        window.location = "/"
    });
})
