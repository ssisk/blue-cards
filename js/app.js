App = Ember.Application.create();

App.Router.map(function() {
    this.resource('why');
    this.resource('picker');
    this.resource('cards');
});

App.CardsRoute = Ember.Route.extend({
    model: function() {
	return App.Card.find();
    }
})

App.Store = DS.Store.extend({
    adapter: 'DS.FixtureAdapter'
});

App.Card = DS.Model.extend({
    name: DS.attr('string'),
//    set: DS.attr('number')
});

App.Card.FIXTURES = [
    // Base
    {id: 1, name: 'Library'},
    // Seaside
    {id: 2, name: 'Pirate Ship'},
    // Intrigue
    {id: 3, name: 'Baron'},

];