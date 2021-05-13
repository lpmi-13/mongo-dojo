config="{
  _id: 'dojo', members: [
  {_id: 0, host: '192.168.42.100:27017'},
  {_id: 1, host: '192.168.42.101:27017'},
  {_id: 2, host: '192.168.42.102:27017'}
  ]
}"

mongo localhost:27017 --eval "rs.initiate($config)"
