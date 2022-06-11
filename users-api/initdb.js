db.createUser({
  user: 'myuser',
  pwd: 'password',
  roles: [
    {
      role: 'readWrite',
      db: 'user-service-db'
    }
  ]
})