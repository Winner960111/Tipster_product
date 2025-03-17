use transaction;

errorHandler = function (action, errorAction) {
    try {
        action();
    } catch (ex) {
        print(ex);
        if (errorAction != null) errorAction();
    }
};

errorHandler(
    () => db.createUser(
        {
            user: "support",
            pwd: "pass.123",
            roles: [{role: "readWrite", db: "transaction"}]
        }
    ),
    () => db.grantRolesToUser("support", [{role: "readWrite", db: "transaction"}])
);
