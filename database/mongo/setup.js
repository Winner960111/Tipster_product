use tipster;

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
            user: "root",
            pwd: "pass.123",
            roles: [{role: "readWrite", db: "tipster"}]
        }
    ),
    () => db.grantRolesToUser("root", [{role: "readWrite", db: "tipster"}])
);
