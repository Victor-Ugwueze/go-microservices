db.createUser({
  user: 'myuser',
  pwd: 'password',
  roles: [
    {
      role: 'readWrite',
      db: 'orders-service-db'
    }
  ]
})