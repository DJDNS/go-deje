define(['jquery'], function($) {
    function Injector(root, connector) {
        var self = this;
        this.root = $(root);
        this.connector = connector;

        this.root.find('.msg-submit').click(function() {
            var data = self.get_msg_input();
            if (data != undefined) {
                self.connector.client.publish(data);
            }
        });
        this.root.find('.msg-listen').click(function() {
            var data = self.get_msg_input();
            var client = self.connector.client;
            if (data != undefined) {
                client.cb_managers.msg.run(client.topic, data);
            }
        });
    }
    Injector.prototype.get_msg_input = function() {
        var data = this.root.find('.msg-input').val();
        try {
            return JSON.parse(data);
        } catch (e) {
            return logger.append("<your browser>: Message is not valid JSON")
        }
    }

    return Injector;
});
