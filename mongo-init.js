db = db.getSiblingDB('telemetry_database');

db.createUser({
    user: "root",
    pwd: "mongodb",
    roles: [{ role: "readWrite", db: "telemetry_database" }]
});
