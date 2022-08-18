db = db.getSiblingDB('data');

db.createCollection('items');

var items_to_insert = 100000;
var items_array = [];

for (var i = 0; i < items_to_insert; i++) {
  items_array.push({
    app: 'vecro-sim',
    id: i,
    value: Math.random()
  });
}

db.items.insertMany(items_array);