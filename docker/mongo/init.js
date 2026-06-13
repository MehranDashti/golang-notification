db = db.getSiblingDB('golang_notification');

db.createUser({
  user: 'golang_notification_user',
  pwd: 'Mehran123456!',
  roles: [{ role: 'dbOwner', db: 'golang_notification' }]
});