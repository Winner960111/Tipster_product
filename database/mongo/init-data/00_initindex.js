use tipster;

initIndexes = {
    start: function () {

        db.getCollection("users").createIndex(
            { 'username': 1 }, 
            { 
                'name': "idx_username"
            }
        );
        db.getCollection("users").createIndex(
            { 'email': 1 }, 
            { 
                'name': "idx_email_unique", 
                'unique': true
            }
        );
        db.getCollection("users").createIndex(
            { 'tags': 1 },  // Multikey index for array field
            { 
                'name': "idx_user_tags"
            }
        );
        db.getCollection("users").createIndex(
            { 'createdAt': -1 }, 
            { 
                'name': "idx_user_createdAt"
            }
        );
        db.getCollection("users").createIndex(
            { 'updatedAt': -1 }, 
            { 
                'name': "idx_user_updatedAt"
            }
        );

        // Tips Collection Indexes
        db.getCollection("tips").createIndex(
            { 'tipsterId': 1 }, 
            { 
                'name': "idx_tipsterId"
            }
        );
        db.getCollection("tips").createIndex(
            { 'tags': 1 },  // Multikey index for array field
            { 
                'name': "idx_tip_tags"
            }
        );
        db.getCollection("tips").createIndex(
            { 'createdAt': -1 }, 
            { 
                'name': "idx_tip_createdAt"
            }
        );
        db.getCollection("tips").createIndex(
            { 'updatedAt': -1 }, 
            { 
                'name': "idx_tip_updatedAt"
            }
        );

        // Comments Collection Indexes
        db.getCollection("comments").createIndex(
            { 'tipId': 1 }, 
            { 
                'name': "idx_comment_tipId"
            }
        );
        db.getCollection("comments").createIndex(
            { 'userId': 1 }, 
            { 
                'name': "idx_comment_userId"
            }
        );
        db.getCollection("comments").createIndex(
            { 'parentId': 1 }, 
            { 
                'name': "idx_comment_parentId"
            }
        );
        db.getCollection("comments").createIndex(
            { 'createdAt': -1 }, 
            { 
                'name': "idx_comment_createdAt"
            }
        );
        db.getCollection("comments").createIndex(
            { 'updatedAt': -1 }, 
            { 
                'name': "idx_comment_updatedAt"
            }
        );
    }
};

initIndexes.start();
