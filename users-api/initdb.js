db.createUser({
  user: 'myuser',
  pwd: 'password',
  roles: [
    {
      role: 'readWrite',
      db: 'users-service-db'
    }
  ]
})