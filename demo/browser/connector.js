define(['jquery', 'deje/client'], function($, DejeClient) {
    function Connector(url, topic, logger) {
        var self = this;
        this.url    = url;
        this.topic  = topic;
        this.logger = logger;
        this.log_func = function(item) {
            self.logger.append(item);
        }

        this.client = undefined;
    }
    Connector.prototype.setup_interface = function(root) {
        var self = this;
        root = $(root);
        root.find('.url').attr('placeholder', this.url);
        root.find('.topic').attr('placeholder', this.topic);
        root.find('button').click(function(){
            self.reconnect(
                root.find('.url').val(),
                root.find('.topic').val()
            );
        });
        return this;
    }
    Connector.prototype.reconnect = function(url, topic) {
        var client = this.client;
        if (client && client.session) {
            client.session.close()
        }

        this.url = url || this.url;
        this.topic = topic || this.topic;

        client = new DejeClient(this.url, this.topic, { 'logger': this.log_func });
        client.connect();
        this.client = client;

        client.cb_managers.store_event.add('render_events', this.render_events.bind(this));
        client.cb_managers.goto_event.add('render_state', this.render_state.bind(this));
        client.cb_managers.update_ts.add('render_state', this.render_state.bind(this));

        this.render_events();
        this.render_state();
        this.render_callbacks();
    }

    Connector.prototype.render_state = function() {
        var chooser = $('#timestamps-chooser span');
        chooser.empty();
        for (var t=0; t<this.client.timestamps.length; t++) {
            if ( t > 0 ) {
                chooser.append(',');
            }
            $('<a/>', {
                href: '#',
                text: JSON.stringify(this.client.timestamps[t])
            })
            .addClass('event')
            .data('hash', this.client.timestamps[t])
            .click(goto_event_by_element)
            .appendTo(chooser);
        }
        $('#state-data').text(
            "hash: '" + this.client.state.hash + "'\n\n"
            + JSON.stringify(this.client.state.content, null, 4)
        );
    }
    Connector.prototype.render_callbacks = function() {
        var root = $('.callbacks-selection');
        root.empty();
        var on_change = function() {
            var callback = $(this).data("callback");
            callback.enabled = $(this).is(":checked");
        };
        for (var m in this.client.cb_managers) {
            var manager = this.client.cb_managers[m];
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
    Connector.prototype.__render_event = function(ev) {
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
    Connector.prototype.render_events = function() {
        $(".events .event").not(".input").remove();
        var history = this.client.getHistory();
        for (var i = 0; i<history.length; i++) {
            var ev = history[i];
            $(".events").prepend(this.__render_event(ev));
        }
    }

    return Connector;
});
