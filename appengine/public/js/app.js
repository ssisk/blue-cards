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
    //id: DS.attr('number'),
    set: DS.attr('number'),
    name: DS.attr('string'),
    type: DS.attr('string'),
    cost: DS.attr('string') // costs can be either coins ($4) or potions, thus string

    
});

/*

what needs to happen:
1. when the user clicks generate, we call the server
    localhost:8080/get10
    {
  "NoAttack": true,
  "NoAttackWithoutDefense": false,
  "SetsAvailable": [1, 5],
  "ForbiddenCards": [6, 7]
}
2. when the data is received, change the model

*/

App.PickerController = Ember.ObjectController.extend({
    generate: function() {
        that = this;
        jQuery.post(
            "/get10",
            JSON.stringify({
                                   "NoAttack": true,
                                    "NoAttackWithoutDefense": false,
                                    "SetsAvailable": [1, 5],
                                    "ForbiddenCards": [6, 7]
                                }),
            function(data, textStatus, jqXHR) {that.cardsReceived(data, textStatus, jqXHR)}
        )
    },  
    cardsReceived: function(data, textStatus, jqXHR) {
        // data is of form: [{"Id":165,"Set":7,"Name":"Squire","Cost":"$2","Type":"Action"},]
        cards = JSON.parse(data);
        console.log(cards)
        this.set('cardsToDisplay', cards);
        this.set('weHaveCards', true);
    },
    cardsToDisplay: [],
    weHaveCards: false
});

App.Card.FIXTURES = [
    // Base
    {id: 100, set: 1, name: 'Library', type: 'Action', cost: '$5'},
    // Seaside
    {id: 200, set: 2, name: 'Pirate Ship', type: 'Action', cost: '$4'},
    // Intrigue
    {id: 300, set: 3, name: 'Baron', type: 'Action', cost: '$4'},
    {id: 301, set: 3, name: 'Masquerade', type: 'Action', cost: '$3'},

];

/* pickCards will be some sort of strategy pattern */

/*
App.pickCards = function (options, allCards) {
    return allCards;
};*/