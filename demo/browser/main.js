require.config({
    paths: {
        'deje': 'js-deje/deje'
    },
    shim: {
        'deje/vendor/autobahn': {
            exports: 'ab',
        }
    }
});

var client;
var URL   = "ws://" + window.location.host + "/ws";
var TOPIC = "deje://demo/";

$('.reconnector .url').attr(  'placeholder', URL);
$('.reconnector .topic').attr('placeholder', TOPIC);

function __placeholder(name){
    return function() {
        console.log("Function " + name + " called before it was ready");
    }
}

function loggit(info) {
    $('#log').text( $('#log').text() + info + "\n" );
}

function log_clear() {
    $('#log').text('');
}

$('.reconnector button').click(function() {
    reconnect(
        $('.reconnector .url').val(),
        $('.reconnector .topic').val()
    );
});

function get_msg_input() {
    var data = $('#message-input').val();
    try {
        return JSON.parse(data);
    } catch (e) {
        return loggit("<your browser>: Message is not valid JSON")
    }
}

$('#message-submit').click(function(){
    var data = get_msg_input()
    if (data != undefined) {
        client.publish(data);
    }
});
$('#message-listen').click(function(){
    var data = get_msg_input()
    if (data != undefined) {
        client.cb_managers.msg.run(client.topic, data);
    }
});

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

function render_state() {
    var chooser = $('#timestamps-chooser span');
    chooser.empty();
    for (var t=0; t<client.timestamps.length; t++) {
        if ( t > 0 ) {
            chooser.append(',');
        }
        $('<a/>', {
            href: '#',
            text: JSON.stringify(client.timestamps[t])
        })
        .addClass('event')
        .data('hash', client.timestamps[t])
        .click(goto_event_by_element)
        .appendTo(chooser);
    }
    $('#state-data').text(
        "hash: '" + client.state.hash + "'\n\n"
        + JSON.stringify(client.state.content, null, 4)
    );
}

function render_callbacks() {
    var root = $('.callbacks-selection');
    root.empty();
    var on_change = function() {
        var callback = $(this).data("callback");
        callback.enabled = $(this).is(":checked");
    };
    for (var m in client.cb_managers) {
        var manager = client.cb_managers[m];
        for (var c in manager.callbacks) {
            var callback = manager.callbacks[c];
            var checkbox = $(document.createElement('input'))
                .attr('type', 'checkbox')
                .attr('checked', callback.enabled)
                .data('callback', callback)
                .change(on_change);
            var label = $(document.createElement('label'))
                .append(checkbox)
                .append("on_" + m + " :: " + c);
            root.append(label);
            root.append("<br>");
        }
    }
}

function render_event(ev) {
    var ev_hash = ev.getHash();
    var ev_content = ev.getContent();

    var element = $(document.createElement('div'))
        .addClass('event')
        .data('hash', ev_hash);
    var promote_button = $(document.createElement('button'))
        .text("Promote")
        .click(promote_event_by_element);
    var goto_button = $(document.createElement('button'))
        .text("Goto")
        .click(goto_event_by_element);
    var hash_bar = $(document.createElement('div'))
        .addClass('hash')
        .text(ev_hash)
        .append(promote_button)
        .append(goto_button);
    var content = $(document.createElement('pre'))
        .addClass('content')
        .text( JSON.stringify(ev_content, null, 4) );

    element.append(hash_bar);
    element.append(content);
    return element;
}

function render_events() {
    $(".events .event").not(".input").remove();
    var history = client.getHistory();
    for (var i = 0; i<history.length; i++) {
        var ev = history[i];
        $(".events").prepend(render_event(ev));
    }
}

function promote_event_by_hash(hash) {
    ev = client.getEvent(hash);
    if (ev != undefined) {
        client.promoteEvent(ev);
    } else {
        loggit("<your browser>: No such hash: " + hash);
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
        loggit("<your browser>: No such hash: " + hash);
    }
    return false;
}

save_event = __placeholder('save_event');
require(['deje/event'], function(DejeEvent) {
    save_event = function() {
        var content = {
            "parent": $('.event.input .parent').val() || "",
            "handler": $('.event.input .handler').val() || "",
        };
        try {
            content.args = JSON.parse($('.event.input .args').val());
        } catch (e) {
            return loggit("<your browser>: Invalid JSON for event args")
        }
        client.storeEvent(new DejeEvent(content));
    }
});

reconnect = __placeholder('reconnect');
require(['deje/client'], function(DejeClient) {
    reconnect = function(url, topic) {
        if (client && client.session) {
            client.session.close()
        }
        client = new DejeClient(url || URL, topic || TOPIC, { 'logger': loggit });
        client.connect();
        window.client = client;

        client.cb_managers.store_event.add('render_events', render_events);
        client.cb_managers.goto_event.add('render_state', render_state);
        client.cb_managers.update_ts.add('render_state', render_state);

        render_events();
        render_state();
        render_callbacks();
    }

    $(document).ready(function() {
        reconnect();
    });
});
