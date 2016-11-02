$(function() {
    /*
      var rules_basic = {
          condition: 'AND',
          rules: [{
              id: 'price',
              operator: 'less',
              value: 10.25
          }, {
              condition: 'OR',
              rules: [{
                  id: 'category',
                  operator: 'equal',
                  value: 2
              }, {
                  id: 'category',
                  operator: 'equal',
                  value: 1
              }]
          }]
      };

    */

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

    $('#builder').queryBuilder({
        plugins: ['bt-tooltip-errors'],
        filters: [{
            id: 'name',
            label: 'Dataset',
            type: 'string'
        }, {
            id: 'column',
            label: 'Column',
            type: "string"
        }],
    });

    $('#submit').on('click', function() {
        var result = $('#builder').queryBuilder('getSQL', 'question_mark');
      console.log(JSON.stringify(result, null, 2))
        if (result.sql.length) {
            //alert(result.sql + '\n\n' + JSON.stringify(result.params, null, 2));
            var Q = window.location.search
            Q = update(Q, "x", JSON.stringify(result, null, 2))
            //window.location = window.location.pathname + Q
        }
    });

})
