use transaction;

initIndexes = {
    start: function () {
        db.getCollection("transaction.test").createIndex(
            { 'Telegram.PlatformIdentity': 1}, 
            { 
                'name': "tgIX", 
                'unique': true, 
                'partialFilterExpression': { "Telegram.PlatformIdentity": { '$exists': true } },
                'collation': { 'locale': 'en', 'strength': 2 }
            });
        db.getCollection("transaction.test").createIndex(
            { 'Google.PlatformIdentity': 1}, 
            { 
                'name': "gooIX", 
                'unique': true, 
                'partialFilterExpression': { "Google.PlatformIdentity": { '$exists': true } },
                'collation': { 'locale': 'en', 'strength': 2 }
            });
    }
};

initIndexes.start();
