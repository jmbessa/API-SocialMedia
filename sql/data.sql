INSERT INTO users (nome, nick, email, password)
values
("user 1", "user_1", "user1@gmail.com", "$2a$10$BpvIAF83paKczAQjjwBH8Om1h/0Wrk3arODnlQTHm8lly./IqRaIa"),
("user 2", "user_2", "user2@gmail.com", "$2a$10$BpvIAF83paKczAQjjwBH8Om1h/0Wrk3arODnlQTHm8lly./IqRaIa"),
("user 3", "user_3", "user3@gmail.com", "$2a$10$BpvIAF83paKczAQjjwBH8Om1h/0Wrk3arODnlQTHm8lly./IqRaIa");

INSERT INTO followers(user_id, follower_id)
values
(1, 2)
(3, 1)
(1, 3)

INSERT INTO posts(title, content, authorId)
values 
("User 1 post", "This is the post of the user 1!", 1),
("User 2 post", "This is the post of the user 2!", 2),
("User 3 post", "This is the post of the user 3!", 3);