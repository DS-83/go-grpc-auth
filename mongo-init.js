db = db.getSiblingDB("admin")

db.createUser(
    {
        user: "admin",
        pwd: "admin", //passwordPrompt(), // or cleartext password
      roles: [
          { role: "userAdminAnyDatabase", db: "admin" },
          { role: "readWriteAnyDatabase", db: "admin" }
        ]
    }
    )
    
db = db.getSiblingDB("photogramm")
    // db.dropUser("dev")
db.createUser({
    user: 'dev',
    pwd: 'SuperSecretPassword',
    roles: [
        {
            role: 'dbAdmin',
            db: 'photogramm',
        },
    ],
});


db.users.createIndex( { username: 1 } )

db.createCollection('revokedTokens');

db.adminCommand( { shutdown: 1 } )
