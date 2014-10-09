require.config({
    paths: {
        'deje': 'js-deje/deje',
        'jquery': 'https://code.jquery.com/jquery-2.1.1.min'
    },
    shim: {
        'deje/vendor/autobahn': {
            exports: 'ab',
        }
    }
});

require(['jquery', 'deje/event', 'log', 'connector', 'injector'],
        function($, DejeEvent, Logger, Connector, Injector) {
    var client;
    var URL   = "ws://" + window.location.host + "/ws";
    var TOPIC = "deje://demo/";

    var logger = new Logger('#log', '#log_filter');
    var connector = new Connector(URL, TOPIC, logger)
        .setup_interface('.reconnector');
    var injector = new Injector('.submitter', connector);

    $('.js-edit-timestamps').click(function(){
        var per_ts_prefix = '\n        ';
        $('#message-input').val('{\n' +
            '    "instructions": "Edit as desired, then click the *Listen to the voices* button below.",\n' +
            '    "type": "02-publish-timestamps",\n' +
            '    "timestamps": [' + ((client.timestamps.length > 0) ? per_ts_prefix : '') +
             client.timestamps.map(JSON.stringify).join(',' + per_ts_prefix) +
             ((client.timestamps.length > 0) ? '\n    ' : '') + ']\n}'
        );
        return false;
    });
    $('.js-clear-timestamps').click(function(){
        client.setTimestamps([]);
        return false;
    });

    function promote_event_by_hash(hash) {
        ev = client.getEvent(hash);
        if (ev != undefined) {
            client.promoteEvent(ev);
        } else {
            logger.append("<your browser>: No such hash: " + hash);
        }
        return false;
    }

    function promote_event_by_element() {
        hash = $(this).closest('.event').data('hash');
        return promote_event_by_hash(hash);
    }

    function goto_event_by_element() {
        hash = $(this).closest('.event').data('hash');

        ev = client.getEvent(hash);
        if (ev != undefined) {
            client.applyEvent(ev);
            render_state();
        } else {
            logger.append("<your browser>: No such hash: " + hash);
        }
        return false;
    }

    save_event = function() {
        var content = {
            "parent": $('.event.input .parent').val() || "",
            "handler": $('.event.input .handler').val() || "",
        };
        try {
            content.args = JSON.parse($('.event.input .args').val());
        } catch (e) {
            return logger.append("<your browser>: Invalid JSON for event args")
        }
        client.storeEvent(new DejeEvent(content));
    }


    $(document).ready(function() {
        connector.reconnect();
    });

});
