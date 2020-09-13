CREATE TABLE users(
  id SERIAL PRIMARY KEY,
  username VARCHAR UNIQUE NOT NULL,
  password TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE user_profiles(
  id SERIAL PRIMARY KEY,
  first_name VARCHAR NULL,
  last_name VARCHAR NULL,
  age SMALLINT NULL,
  email VARCHAR NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  user_id INT NOT NULL,
  
  FOREIGN KEY(user_id) REFERENCES users(id)
);

CREATE TABLE user_images(
  id SERIAL PRIMARY KEY,
  image_path VARCHAR NOT NULL,
  user_id INT NOT NULL,

  FOREIGN KEY(user_id) REFERENCES users(id)
);