define(['jquery'], function($) {

    function Logger(element, filter_input) {
        var self = this;

        this.data = [];
        this.filter_value = '';

        this.element = $(element);
        this.filter_input = $(filter_input).on(
            'change keydown keyup',
            function() { self.set_filter( $(this).val() ); }
        );
    }
    Logger.prototype.display = function() {
        var callback = this.__contains_filter_value.bind(this);
        var filtered = this.data.filter( callback );
        this.element.text( filtered.join("\n" ) );
    }
    Logger.prototype.append = function(info) {
        this.data.push(info);
        this.display();
    }
    Logger.prototype.clear = function() {
        this.data = [];
        this.display();
    }
    Logger.prototype.set_filter = function(value) {
        this.filter_value = value;
        this.display();
    }
    Logger.prototype.__contains_filter_value = function(item) {
        if (item == undefined) {
            return;
        }
        return item.toLowerCase().contains( this.filter_value.toLowerCase() );
    }

    return Logger;
});
