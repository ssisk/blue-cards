App = Ember.Application.create();

App.Router.map(function() {
    this.resource('why');
    this.resource('picker');
    this.resource('cards');
});

/*App.IndexRoute = Ember.Route.Extend({
    redirect: function() {
        this.transitionTo('picker');
    }
});*/

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

App.PickerController = Ember.ObjectController.extend({
    generate: function() {
        this.set('cardsToDisplay', App.Card.find());
        this.set('weHaveCards', true);
    },  
    cardsToDisplay: [],
    weHaveCards: false
});

App.Card.FIXTURES = [
    // Base
    {id: 100, name: 'Library'},
    // Seaside
    {id: 200, name: 'Pirate Ship'},
    // Intrigue
    {id: 300, name: 'Baron'},
    {id: 301, name: 'Masquerade'},

];

/* pickCards will be some sort of strategy pattern */
/*

    options looks something like: 
    { 
        noAttack: true,
    }

    I think there are two types of card filters:
    * filters that operate on only one card ("No attack cards")
    * filters that operate on the set of cards ("no unbalanced attack cards")

    the one card filters could probably be run on allCards at the beginning,
    but the multiple card ones will probably need to adopt a "generate 10, 
    check for compliance, then mutate if necessary" strategy
*/
/*
App.pickCards = function (options, allCards) {
    return allCards;
};*/